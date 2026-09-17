package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestAllowlistMatching(t *testing.T) {
	list, err := newAllowlist(" api.githubcopilot.com ,github.com,.githubusercontent.com ")
	if err != nil {
		t.Fatalf("newAllowlist: %v", err)
	}
	cases := []struct {
		host string
		want bool
	}{
		{"api.githubcopilot.com", true},
		{"API.GitHubCopilot.com", true},
		{"api.githubcopilot.com.", true},
		{"github.com", true},
		{"raw.github.com", true},
		{"objects.githubusercontent.com", true},
		{"notgithub.com", false},
		{"github.com.evil.example", false},
		{"example.com", false},
		{"", false},
	}
	for _, item := range cases {
		if got := list.allowed(item.host); got != item.want {
			t.Errorf("allowed(%q) = %v, want %v", item.host, got, item.want)
		}
	}
	if _, err := newAllowlist(""); err == nil {
		t.Fatal("empty allowlist must be rejected")
	}
	if _, err := newAllowlist("bad rule/with/slash"); err == nil {
		t.Fatal("allowlist rule with separators must be rejected")
	}
	for _, spec := range []string{"*", "*.github.com", "https://github.com", "github.com:443", "github.com/path", "api@github.com", "-github.com"} {
		if _, err := newAllowlist(spec); err == nil {
			t.Errorf("allowlist rule %q must be rejected instead of silently never matching", spec)
		}
	}
}

func TestListenAddressMustBeLoopback(t *testing.T) {
	for _, listen := range []string{"127.0.0.1:39877", "localhost:39877", "[::1]:39877"} {
		if err := validateListenAddress(listen); err != nil {
			t.Errorf("validateListenAddress(%q) = %v, want nil", listen, err)
		}
	}
	for _, listen := range []string{"0.0.0.0:39877", ":39877", "192.168.1.10:39877", "example.com:39877", "127.0.0.1", "127.0.0.1:http"} {
		if err := validateListenAddress(listen); err == nil {
			t.Errorf("validateListenAddress(%q) must reject non-loopback or malformed input", listen)
		}
	}
	out := &strings.Builder{}
	if err := run(config{listen: "0.0.0.0:39877", allowRules: []string{"github.com"}, idleTimeout: time.Minute}, out); err == nil {
		t.Fatal("run must refuse a non-loopback listen address")
	}
}

func TestAuditTrailFailureDeniesRequests(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "proxy.jsonl")
	logger, err := newEventLogger(logPath)
	if err != nil {
		t.Fatalf("newEventLogger: %v", err)
	}
	allow, err := newAllowlist("127.0.0.1")
	if err != nil {
		t.Fatalf("newAllowlist: %v", err)
	}
	server := httptest.NewServer(newProxy(config{idleTimeout: time.Minute}, allow, logger))
	defer server.Close()
	logger.close() // the audit trail becomes unusable

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:1/denied", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	transport := &http.Transport{Proxy: http.ProxyURL(mustParseURL(t, server.URL))}
	client := &http.Client{Transport: transport, Timeout: 5 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("request through proxy: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 when the audit trail is unusable", response.StatusCode)
	}
}

func TestConnectAuditTrailFailureDeniesAndDoesNotTunnel(t *testing.T) {
	target, connections := startTCPEcho(t)
	logPath := filepath.Join(t.TempDir(), "proxy.jsonl")
	logger, err := newEventLogger(logPath)
	if err != nil {
		t.Fatalf("newEventLogger: %v", err)
	}
	allow, err := newAllowlist("127.0.0.1")
	if err != nil {
		t.Fatalf("newAllowlist: %v", err)
	}
	server := httptest.NewServer(newProxy(config{idleTimeout: time.Minute}, allow, logger))
	defer server.Close()

	// With a working audit trail the tunnel is established, so the target observes a connection.
	first := connectThroughProxy(t, server.URL, target)
	if !strings.Contains(first, "200") {
		t.Fatalf("status line = %q, want a 200 before the log is broken", first)
	}
	if connections.Load() == 0 {
		t.Fatal("target must have been reached while the audit trail was healthy")
	}

	logger.close() // the audit trail becomes unusable
	before := connections.Load()
	status := connectThroughProxy(t, server.URL, target)
	if !strings.Contains(status, "403") {
		t.Fatalf("status line = %q, want 403 when the audit trail is unusable", status)
	}
	if connections.Load() != before {
		t.Fatal("no tunnel may be opened once the audit trail is broken")
	}
}

func TestMalformedRequestIsAudited(t *testing.T) {
	server, logPath := newTestProxy(t, config{}, "127.0.0.1")
	defer server.Close()
	req, err := http.NewRequest(http.MethodGet, server.URL+"/origin-form", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.StatusCode)
	}
	record := waitForRecord(t, logPath, func(entry map[string]any) bool {
		reason, _ := entry["reason"].(string)
		return strings.HasPrefix(reason, "malformed-request")
	})
	if record["decision"] != "deny" {
		t.Fatalf("decision = %v, want deny", record["decision"])
	}
}

func TestPlainHTTPRecordsPort(t *testing.T) {
	server, logPath := newTestProxy(t, config{}, "127.0.0.1")
	defer server.Close()
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/port-test", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	transport := &http.Transport{Proxy: http.ProxyURL(mustParseURL(t, server.URL))}
	client := &http.Client{Transport: transport, Timeout: 5 * time.Second}
	response, err := client.Do(req)
	if err == nil {
		_ = response.Body.Close()
	}
	record := waitForRecord(t, logPath, func(entry map[string]any) bool {
		return entry["event"] == "http" && entry["port"] == "8080"
	})
	if record["host"] != "127.0.0.1" {
		t.Fatalf("host = %v, want 127.0.0.1", record["host"])
	}
}

func TestDialFailureIsAuditedAsDeny(t *testing.T) {
	server, logPath := newTestProxy(t, config{idleTimeout: time.Minute}, "127.0.0.1")
	defer server.Close()
	// port 1 on loopback is allowlisted but refuses connections, so the tunnel dial fails
	status := connectThroughProxy(t, server.URL, "127.0.0.1:1")
	if !strings.Contains(status, "502") {
		t.Fatalf("status line = %q, want a 502 from the proxy", status)
	}
	record := waitForRecord(t, logPath, func(entry map[string]any) bool {
		reason, _ := entry["reason"].(string)
		return strings.HasPrefix(reason, "dial-failed")
	})
	if record["decision"] != "deny" {
		t.Fatalf("decision = %v, want deny", record["decision"])
	}
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("url.Parse(%q): %v", raw, err)
	}
	return parsed
}

func newTestProxy(t *testing.T, cfg config, allowSpec string) (*httptest.Server, string) {
	t.Helper()
	logPath := filepath.Join(t.TempDir(), "proxy.jsonl")
	logger, err := newEventLogger(logPath)
	if err != nil {
		t.Fatalf("newEventLogger: %v", err)
	}
	list, err := newAllowlist(allowSpec)
	if err != nil {
		t.Fatalf("newAllowlist: %v", err)
	}
	cfg.logPath = logPath
	server := httptest.NewServer(newProxy(cfg, list, logger))
	t.Cleanup(func() {
		server.Close()
		logger.close()
	})
	return server, logPath
}

func readLog(t *testing.T, path string) []map[string]any {
	t.Helper()
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	var records []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(payload)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("log line is not JSON: %v (%q)", err, line)
		}
		records = append(records, record)
	}
	return records
}

// waitForRecord polls the log: the tunnel "established" record is written by the handler goroutine,
// which can land marginally after the client handshake completes.
func waitForRecord(t *testing.T, path string, predicate func(map[string]any) bool) map[string]any {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		for _, record := range readLog(t, path) {
			if predicate(record) {
				return record
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("no matching log record within the deadline; log=%+v", readLog(t, path))
	return nil
}

func startTCPEcho(t *testing.T) (string, *atomic.Int32) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	accepted := &atomic.Int32{}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			accepted.Add(1)
			go func() {
				defer conn.Close()
				_, _ = io.Copy(conn, conn)
			}()
		}
	}()
	t.Cleanup(func() { _ = listener.Close() })
	return listener.Addr().String(), accepted
}

func connectThroughProxy(t *testing.T, proxyServer, target string) string {
	t.Helper()
	proxyURL, err := url.Parse(proxyServer)
	if err != nil {
		t.Fatalf("parse proxy url: %v", err)
	}
	conn, err := net.Dial("tcp", proxyURL.Host)
	if err != nil {
		t.Fatalf("dial proxy: %v", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	request := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", target, target)
	if _, err := conn.Write([]byte(request)); err != nil {
		t.Fatalf("write CONNECT: %v", err)
	}
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}

func TestConnectTunnelAllowedReachesTarget(t *testing.T) {
	target, accepted := startTCPEcho(t)
	server, logPath := newTestProxy(t, config{idleTimeout: time.Minute}, "127.0.0.1")

	status := connectThroughProxy(t, server.URL, target)
	if !strings.Contains(status, "200") {
		t.Fatalf("want 200 handshake, got %q", status)
	}
	if accepted.Load() != 1 {
		t.Fatalf("target accept count = %d, want 1", accepted.Load())
	}
	record := waitForRecord(t, logPath, func(r map[string]any) bool {
		return r["event"] == "connect" && r["decision"] == "allow" && r["phase"] == "established"
	})
	if record["host"] != "127.0.0.1" {
		t.Fatalf("log host = %v, want 127.0.0.1", record["host"])
	}
}

func TestConnectTunnelDeniedByAllowlist(t *testing.T) {
	target, accepted := startTCPEcho(t)
	server, logPath := newTestProxy(t, config{}, "api.githubcopilot.com")

	status := connectThroughProxy(t, server.URL, target)
	if !strings.Contains(status, "403") {
		t.Fatalf("want 403 refusal, got %q", status)
	}
	if accepted.Load() != 0 {
		t.Fatalf("denied CONNECT must not dial the target, accepted=%d", accepted.Load())
	}
	records := readLog(t, logPath)
	if len(records) != 1 || records[0]["decision"] != "deny" || records[0]["reason"] != "not-on-allowlist" {
		t.Fatalf("unexpected deny log: %+v", records)
	}
}

func TestPlainHTTPForwardedAndLogged(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	}))
	defer upstream.Close()

	server, logPath := newTestProxy(t, config{}, "127.0.0.1")
	proxyURL, _ := url.Parse(server.URL)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	response, err := client.Get(upstream.URL)
	if err != nil {
		t.Fatalf("GET through proxy: %v", err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if string(body) != "pong" {
		t.Fatalf("body = %q, want pong", body)
	}
	records := readLog(t, logPath)
	if len(records) < 2 {
		t.Fatalf("want an authorized record plus a completion record, got %+v", records)
	}
	if records[0]["decision"] != "allow" || records[0]["phase"] != "authorized" {
		t.Fatalf("first record = %+v, want the pre-flight authorized record", records[0])
	}
	completion := records[len(records)-1]
	if completion["event"] != "http" || completion["decision"] != "allow" || completion["bytesToClient"] == nil {
		t.Fatalf("last record = %+v, want the completion record with byte counts", completion)
	}
}

func TestPlainHTTPDeniedByAllowlist(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("target must not be reached for a denied host")
	}))
	defer upstream.Close()

	server, logPath := newTestProxy(t, config{}, "github.com")
	proxyURL, _ := url.Parse(server.URL)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	response, err := client.Get(upstream.URL)
	if err != nil {
		t.Fatalf("GET through proxy: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", response.StatusCode)
	}
	records := readLog(t, logPath)
	if len(records) != 1 || records[0]["decision"] != "deny" {
		t.Fatalf("unexpected deny log: %+v", records)
	}
}

func TestConnectUsesUpstreamProxy(t *testing.T) {
	target, accepted := startTCPEcho(t)

	var upstreamHits atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodConnect {
			http.Error(w, "CONNECT only", http.StatusMethodNotAllowed)
			return
		}
		upstreamHits.Add(1)
		backend, err := net.Dial("tcp", r.Host)
		if err != nil {
			http.Error(w, "dial failed", http.StatusBadGateway)
			return
		}
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "no hijack", http.StatusInternalServerError)
			backend.Close()
			return
		}
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			backend.Close()
			return
		}
		defer clientConn.Close()
		defer backend.Close()
		_, _ = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		go func() { _, _ = io.Copy(backend, clientConn) }()
		_, _ = io.Copy(clientConn, backend)
	}))
	defer upstream.Close()

	server, logPath := newTestProxy(t, config{upstream: upstream.URL, idleTimeout: time.Minute}, "127.0.0.1")
	status := connectThroughProxy(t, server.URL, target)
	if !strings.Contains(status, "200") {
		t.Fatalf("want 200, got %q", status)
	}
	if upstreamHits.Load() != 1 {
		t.Fatalf("upstream hits = %d, want 1", upstreamHits.Load())
	}
	if accepted.Load() != 1 {
		t.Fatalf("target accepts = %d, want 1", accepted.Load())
	}
	record := waitForRecord(t, logPath, func(r map[string]any) bool {
		return r["event"] == "connect" && r["decision"] == "allow"
	})
	if record["upstream"] != upstream.URL {
		t.Fatalf("log upstream = %v, want %v", record["upstream"], upstream.URL)
	}
}

func TestObserveModeAllowsUnknownHostAndLogsIt(t *testing.T) {
	target, accepted := startTCPEcho(t)
	server, logPath := newTestProxy(t, config{observe: true, idleTimeout: time.Minute}, "api.githubcopilot.com")

	status := connectThroughProxy(t, server.URL, target)
	if !strings.Contains(status, "200") {
		t.Fatalf("observe mode must permit the tunnel, got %q", status)
	}
	if accepted.Load() != 1 {
		t.Fatalf("target accepts = %d, want 1", accepted.Load())
	}
	record := waitForRecord(t, logPath, func(r map[string]any) bool {
		return r["event"] == "connect" && r["reason"] == "observe"
	})
	if record["decision"] != "allow" {
		t.Fatalf("observe mode must be recorded as allow: %+v", record)
	}
}

func TestRunFailsFastOnOccupiedPort(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	logPath := filepath.Join(t.TempDir(), "occupied.jsonl")
	cfg := config{
		listen:      listener.Addr().String(),
		allowRules:  []string{"api.githubcopilot.com"},
		logPath:     logPath,
		idleTimeout: time.Minute,
	}
	var output strings.Builder
	err = run(cfg, &output)
	if err == nil {
		t.Fatal("run must fail when the fixed port is occupied (no silent port switch)")
	}
	if !strings.Contains(err.Error(), "never switched silently") {
		t.Fatalf("error must state the no-silent-switch rule, got %v", err)
	}
	records := readLog(t, logPath)
	if len(records) != 1 || records[0]["event"] != "start" || records[0]["decision"] != "error" {
		t.Fatalf("occupied-port attempt must be logged: %+v", records)
	}
}

func TestTunnelWorksAcrossRealTLSThroughProxy(t *testing.T) {
	target := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "secure-pong")
	}))
	defer target.Close()

	server, logPath := newTestProxy(t, config{idleTimeout: time.Minute}, "127.0.0.1")
	proxyURL, _ := url.Parse(server.URL)
	transport := &http.Transport{
		Proxy:             http.ProxyURL(proxyURL),
		TLSClientConfig:   target.Client().Transport.(*http.Transport).TLSClientConfig,
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}
	response, err := client.Get(target.URL)
	if err != nil {
		t.Fatalf("HTTPS through proxy: %v", err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if string(body) != "secure-pong" {
		t.Fatalf("body = %q, want secure-pong", body)
	}
	record := waitForRecord(t, logPath, func(r map[string]any) bool {
		return r["event"] == "connect" && r["decision"] == "allow"
	})
	if record["port"] == "" || record["port"] == nil {
		t.Fatalf("log must record the target port: %+v", record)
	}
}
