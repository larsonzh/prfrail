// Command enforcement-proxy is the B2 external-enforcement allowlist proxy for the ProofRail
// AgentRunner proof: it listens on a fixed loopback endpoint and is the only egress path a
// zero-capability AppContainer is allowed to reach (once the container SID holds a loopback
// exemption). It never decodes TLS: HTTPS is forwarded as a CONNECT tunnel, plain HTTP is
// forwarded with a Host allowlist. Everything is deny-by-default and logged as JSONL evidence.
//
// This is a tool (tools/agent-probe), not production runtime code.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	defaultListen      = "127.0.0.1:39877"
	defaultAllow       = "api.githubcopilot.com,api.github.com,github.com,githubusercontent.com,copilot-proxy.githubusercontent.com"
	defaultConnectDial = 15 * time.Second
	defaultIdleTimeout = 5 * time.Minute
	defaultHTTPTimeout = 2 * time.Minute
)

// config carries every runtime knob. The listen address is never changed silently: when the
// endpoint is occupied the process fails fast instead of picking another port.
type config struct {
	listen      string
	allowRules  []string
	upstream    string
	logPath     string
	observe     bool
	idleTimeout time.Duration
}

// allowlist matches hosts (port-insensitive). A rule matches the host itself and any subdomain.
type allowlist struct {
	rules []string
}

func newAllowlist(spec string) (*allowlist, error) {
	list := &allowlist{}
	for _, raw := range strings.Split(spec, ",") {
		rule := strings.ToLower(strings.TrimSpace(raw))
		rule = strings.TrimPrefix(rule, ".")
		if rule == "" {
			continue
		}
		// A rule is a bare host suffix. Anything else (scheme, path, port, wildcard) would be
		// accepted here but never match a host, so it is rejected instead of failing silently.
		if strings.ContainsAny(rule, " /\\:?@#*[]") {
			return nil, fmt.Errorf("invalid allowlist rule %q: use a bare host suffix such as example.com", raw)
		}
		if strings.HasPrefix(rule, "-") || strings.HasSuffix(rule, "-") {
			return nil, fmt.Errorf("invalid allowlist rule %q: host labels cannot start or end with a hyphen", raw)
		}
		list.rules = append(list.rules, rule)
	}
	if len(list.rules) == 0 {
		return nil, errors.New("allowlist must contain at least one rule")
	}
	return list, nil
}

func (a *allowlist) allowed(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	host = strings.TrimSuffix(host, ".")
	if host == "" {
		return false
	}
	for _, rule := range a.rules {
		if host == rule || strings.HasSuffix(host, "."+rule) {
			return true
		}
	}
	return false
}

// eventLogger writes one JSON object per line; the file is the proof evidence for networkControl.
type eventLogger struct {
	mu          sync.Mutex
	w           *os.File
	failureOnce sync.Once
	failure     error
}

func newEventLogger(path string) (*eventLogger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &eventLogger{w: file}, nil
}

// log writes one audit record and reports whether the trail is intact. The enforcement claim rests on
// the audit trail, so a caller must treat a write failure as fatal for the request (fail closed) rather
// than letting it pass unlogged.
func (l *eventLogger) log(record map[string]any) error {
	if l == nil || l.w == nil {
		return errors.New("audit log is not available")
	}
	if _, ok := record["ts"]; !ok {
		record["ts"] = time.Now().UTC().Format(time.RFC3339Nano)
	}
	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("audit record could not be serialized: %w", err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, err := l.w.Write(append(payload, '\n')); err != nil {
		l.fail(err)
		return err
	}
	if err := l.w.Sync(); err != nil {
		l.fail(err)
		return err
	}
	return nil
}

// fail records the first audit failure once and surfaces it on stderr, because the operator must know
// that decisions from this point on are being denied rather than allowed.
func (l *eventLogger) fail(err error) {
	l.failureOnce.Do(func() {
		l.failure = err
		fmt.Fprintf(os.Stderr, "enforcement-proxy: AUDIT TRAIL BROKEN (%v): requests are denied until the log is writable again\n", err)
	})
}

func (l *eventLogger) err() error {
	if l == nil {
		return errors.New("audit log is not available")
	}
	return l.failure
}

func (l *eventLogger) close() {
	if l == nil || l.w == nil {
		return
	}
	_ = l.w.Close()
}

type proxy struct {
	cfg       config
	allow     *allowlist
	logger    *eventLogger
	transport *http.Transport
}

func newProxy(cfg config, allow *allowlist, logger *eventLogger) *proxy {
	pr := &proxy{cfg: cfg, allow: allow, logger: logger}
	pr.transport = &http.Transport{
		DialContext:           (&net.Dialer{Timeout: defaultConnectDial}).DialContext,
		MaxIdleConns:          32,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
		ResponseHeaderTimeout: defaultHTTPTimeout,
	}
	if cfg.upstream != "" {
		upstreamURL, err := url.Parse(cfg.upstream)
		if err == nil {
			pr.transport.Proxy = http.ProxyURL(upstreamURL)
		}
	}
	return pr
}

func (p *proxy) decide(host string) (bool, string) {
	if p.cfg.observe {
		return true, "observe"
	}
	if p.allow.allowed(host) {
		return true, "allowlisted"
	}
	return false, "not-on-allowlist"
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleConnect(w, r)
		return
	}
	p.handlePlainHTTP(w, r)
}

func splitHostPortOrHost(value string) (string, string) {
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return value, ""
	}
	return host, port
}

func (p *proxy) handleConnect(w http.ResponseWriter, r *http.Request) {
	host, port := splitHostPortOrHost(r.Host)
	allowed, reason := p.decide(host)
	started := time.Now()
	if !allowed {
		_ = p.logger.log(map[string]any{
			"event": "connect", "decision": "deny", "reason": reason,
			"host": host, "port": port, "client": r.RemoteAddr,
		})
		http.Error(w, "blocked by ProofRail enforcement proxy", http.StatusForbidden)
		return
	}
	if err := p.logger.err(); err != nil {
		http.Error(w, "blocked: enforcement audit trail is not writable", http.StatusForbidden)
		return
	}
	// Pre-flight audit write: an allow decision is only valid once it is on disk, so a broken audit
	// trail turns into a denial instead of an unlogged tunnel.
	if err := p.logger.log(map[string]any{
		"event": "connect", "phase": "authorized", "decision": "allow", "reason": reason,
		"host": host, "port": port, "client": r.RemoteAddr, "upstream": p.cfg.upstream,
	}); err != nil {
		http.Error(w, "blocked: enforcement audit trail is not writable", http.StatusForbidden)
		return
	}

	target := r.Host
	if port == "" {
		target = net.JoinHostPort(host, "443")
	}

	clientConn, err := hijack(w)
	if err != nil {
		_ = p.logger.log(map[string]any{"event": "connect", "decision": "error", "reason": err.Error(), "host": host, "port": port})
		return
	}
	defer clientConn.Close()

	targetConn, err := p.dialTunnel(target, host, port)
	if err != nil {
		_ = p.logger.log(map[string]any{
			"event": "connect", "decision": "deny", "reason": "dial-failed: " + err.Error(),
			"host": host, "port": port, "client": r.RemoteAddr, "upstream": p.cfg.upstream,
		})
		_, _ = clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\nContent-Length: 0\r\n\r\n"))
		return
	}
	defer targetConn.Close()

	if _, err := clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		return
	}
	if err := p.logger.log(map[string]any{
		"event": "connect", "phase": "established", "decision": "allow", "reason": reason,
		"host": host, "port": port, "client": r.RemoteAddr, "upstream": p.cfg.upstream,
	}); err != nil {
		// The tunnel is already established, so the only honest option is to close it immediately: an
		// unlogged session must not continue.
		return
	}

	var up, down int64
	done := make(chan struct{}, 2)
	go func() { up = copyWithIdle(targetConn, clientConn, p.cfg.idleTimeout); done <- struct{}{} }()
	go func() { down = copyWithIdle(clientConn, targetConn, p.cfg.idleTimeout); done <- struct{}{} }()
	<-done
	<-done

	_ = p.logger.log(map[string]any{
		"event": "connect", "phase": "closed", "decision": "allow", "reason": reason,
		"host": host, "port": port, "client": r.RemoteAddr, "upstream": p.cfg.upstream,
		"bytesToTarget": up, "bytesToClient": down,
		"durationMs": time.Since(started).Milliseconds(),
	})
}

func (p *proxy) dialTunnel(target, host, port string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: defaultConnectDial}
	if p.cfg.upstream == "" {
		return dialer.Dial("tcp", target)
	}
	upstreamURL, err := url.Parse(p.cfg.upstream)
	if err != nil {
		return nil, err
	}
	if _, _, err := net.SplitHostPort(upstreamURL.Host); err != nil {
		upstreamURL.Host = net.JoinHostPort(upstreamURL.Host, "8080")
	}
	conn, err := dialer.Dial("tcp", upstreamURL.Host)
	if err != nil {
		return nil, err
	}
	request := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\nProxy-Connection: keep-alive\r\n\r\n", target, target)
	if _, err := conn.Write([]byte(request)); err != nil {
		conn.Close()
		return nil, err
	}
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		conn.Close()
		return nil, err
	}
	fields := strings.Fields(line)
	if len(fields) < 2 || !strings.HasPrefix(fields[1], "2") {
		conn.Close()
		return nil, fmt.Errorf("upstream %s refused CONNECT: %s", upstreamURL.Host, strings.TrimSpace(line))
	}
	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return nil, err
		}
		if header == "\r\n" || header == "\n" {
			break
		}
	}
	if reader.Buffered() > 0 {
		conn.Close()
		return nil, errors.New("upstream sent tunnel bytes before the client handshake completed")
	}
	_ = host
	_ = port
	return conn, nil
}

func (p *proxy) handlePlainHTTP(w http.ResponseWriter, r *http.Request) {
	if !r.URL.IsAbs() || r.URL.Host == "" {
		// Audit gap avoided: malformed requests are recorded before they are rejected.
		_ = p.logger.log(map[string]any{
			"event": "http", "decision": "deny", "reason": "malformed-request: absolute-form request URI required",
			"host": r.Host, "port": "", "method": r.Method, "client": r.RemoteAddr,
		})
		http.Error(w, "absolute-form request URI required (this is a forward proxy)", http.StatusBadRequest)
		return
	}
	host := r.URL.Hostname()
	port := r.URL.Port()
	if port == "" {
		if r.URL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	allowed, reason := p.decide(host)
	started := time.Now()
	if !allowed {
		_ = p.logger.log(map[string]any{
			"event": "http", "decision": "deny", "reason": reason,
			"host": host, "port": port, "method": r.Method, "client": r.RemoteAddr,
		})
		http.Error(w, "blocked by ProofRail enforcement proxy", http.StatusForbidden)
		return
	}
	if err := p.logger.err(); err != nil {
		http.Error(w, "blocked: enforcement audit trail is not writable", http.StatusForbidden)
		return
	}
	// Pre-flight audit write (see handleConnect): no allow decision without a durable record.
	if err := p.logger.log(map[string]any{
		"event": "http", "phase": "authorized", "decision": "allow", "reason": reason,
		"host": host, "port": port, "method": r.Method, "client": r.RemoteAddr,
	}); err != nil {
		http.Error(w, "blocked: enforcement audit trail is not writable", http.StatusForbidden)
		return
	}
	outbound := r.Clone(r.Context())
	outbound.RequestURI = ""
	removeHopByHop(outbound.Header)
	response, err := p.transport.RoundTrip(outbound)
	if err != nil {
		_ = p.logger.log(map[string]any{
			"event": "http", "decision": "deny", "reason": "roundtrip-failed: " + err.Error(),
			"host": host, "port": port, "method": r.Method, "client": r.RemoteAddr, "upstream": p.cfg.upstream,
		})
		http.Error(w, "upstream request failed", http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	removeHopByHop(response.Header)
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(response.StatusCode)
	written, _ := io.Copy(w, response.Body)
	_ = p.logger.log(map[string]any{
		"event": "http", "decision": "allow", "reason": reason,
		"host": host, "port": port, "method": r.Method, "client": r.RemoteAddr, "upstream": p.cfg.upstream,
		"status": response.StatusCode, "bytesToClient": written,
		"durationMs": time.Since(started).Milliseconds(),
	})
}

var hopByHop = []string{
	"Connection", "Proxy-Connection", "Keep-Alive", "Proxy-Authenticate",
	"Proxy-Authorization", "Te", "Trailer", "Transfer-Encoding", "Upgrade",
}

func removeHopByHop(header http.Header) {
	for _, key := range hopByHop {
		header.Del(key)
	}
}

func hijack(w http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "connection cannot be hijacked", http.StatusInternalServerError)
		return nil, errors.New("response writer does not support hijacking")
	}
	conn, _, err := hijacker.Hijack()
	return conn, err
}

// copyWithIdle copies until EOF, refreshing a read deadline so a stalled peer cannot pin the tunnel.
func copyWithIdle(dst net.Conn, src net.Conn, idle time.Duration) int64 {
	if idle <= 0 {
		idle = defaultIdleTimeout
	}
	buf := make([]byte, 32*1024)
	var total int64
	for {
		_ = src.SetReadDeadline(time.Now().Add(idle))
		n, readErr := src.Read(buf)
		if n > 0 {
			_ = dst.SetWriteDeadline(time.Now().Add(idle))
			if _, writeErr := dst.Write(buf[:n]); writeErr != nil {
				return total
			}
			total += int64(n)
		}
		if readErr != nil {
			return total
		}
	}
}

func run(cfg config, stdout io.Writer) error {
	if err := validateListenAddress(cfg.listen); err != nil {
		return err
	}
	allow, err := newAllowlist(strings.Join(cfg.allowRules, ","))
	if err != nil {
		return err
	}
	logger, err := newEventLogger(cfg.logPath)
	if err != nil {
		return err
	}
	defer logger.close()

	listener, err := net.Listen("tcp", cfg.listen)
	if err != nil {
		_ = logger.log(map[string]any{
			"event": "start", "decision": "error", "reason": "listen-failed: " + err.Error(),
			"listen": cfg.listen,
		})
		return fmt.Errorf("listen %s failed (the port is never switched silently): %w", cfg.listen, err)
	}
	defer listener.Close()

	mode := "enforce"
	if cfg.observe {
		mode = "observe"
	}
	fmt.Fprintf(stdout, "enforcement-proxy mode=%s listen=%s upstream=%q allow=%s log=%s\n",
		mode, cfg.listen, cfg.upstream, strings.Join(allow.rules, ","), cfg.logPath)
	if err := logger.log(map[string]any{
		"event": "start", "decision": "ok", "mode": mode, "listen": cfg.listen,
		"allow": allow.rules, "upstream": cfg.upstream,
	}); err != nil {
		return fmt.Errorf("audit trail is not writable, refusing to serve: %w", err)
	}

	server := &http.Server{
		Handler:           newProxy(cfg, allow, logger),
		ReadHeaderTimeout: 15 * time.Second,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		_ = logger.log(map[string]any{"event": "stop", "decision": "ok", "reason": "signal"})
		_ = server.Close()
	}()
	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	_ = logger.log(map[string]any{"event": "stop", "decision": "ok"})
	return nil
}

func main() {
	var (
		listen   = flag.String("listen", defaultListen, "fixed loopback endpoint; never changed automatically")
		allow    = flag.String("allow", defaultAllow, "comma-separated host allowlist (matches host and subdomains)")
		upstream = flag.String("upstream", "", "optional upstream HTTP proxy for the host-side egress path")
		logPath  = flag.String("log", "enforcement-proxy.jsonl", "JSONL audit log path")
		observe  = flag.Bool("observe", false, "discovery mode: allow everything but log every target (never proof evidence)")
		idle     = flag.Duration("idle-timeout", defaultIdleTimeout, "idle timeout for tunneled connections")
	)
	flag.Parse()

	cfg := config{
		listen:      *listen,
		allowRules:  strings.Split(*allow, ","),
		upstream:    *upstream,
		logPath:     *logPath,
		observe:     *observe,
		idleTimeout: *idle,
	}
	if _, port, err := net.SplitHostPort(cfg.listen); err != nil || port == "" {
		fmt.Fprintf(os.Stderr, "listen address must include a port: %s\n", cfg.listen)
		os.Exit(2)
	}
	if _, err := strconv.Atoi(mustPort(cfg.listen)); err != nil {
		fmt.Fprintf(os.Stderr, "listen port must be numeric: %s\n", cfg.listen)
		os.Exit(2)
	}
	if cfg.upstream != "" {
		if _, err := url.Parse(cfg.upstream); err != nil {
			fmt.Fprintf(os.Stderr, "invalid upstream proxy URL %q: %v\n", cfg.upstream, err)
			os.Exit(2)
		}
	}
	if err := run(cfg, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "enforcement-proxy:", err)
		os.Exit(1)
	}
}

func mustPort(listen string) string {
	_, port, err := net.SplitHostPort(listen)
	if err != nil {
		return ""
	}
	return port
}

// validateListenAddress keeps the enforcement point on the loopback interface. Exposing it on a
// routable address would turn the allowlist (and, in observe mode, the whole upstream) into an open
// relay for anything on the network, which would invalidate the B2 boundary claims.
func validateListenAddress(listen string) error {
	host, port, err := net.SplitHostPort(listen)
	if err != nil || port == "" {
		return fmt.Errorf("listen address must be host:port, got %q", listen)
	}
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("listen port must be numeric, got %q", listen)
	}
	switch strings.ToLower(host) {
	case "127.0.0.1", "::1", "localhost":
		return nil
	default:
		return fmt.Errorf("listen address must be a loopback address (127.0.0.1, ::1 or localhost), got %q: the enforcement proxy is never exposed to the network", listen)
	}
}
