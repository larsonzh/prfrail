// B2 in-container paid probe client (Node, no dependencies).
// Talks to GitHub/Copilot exclusively through the enforcement proxy via a manual CONNECT tunnel,
// so the container's only egress path is the allowlist proxy. Prints key=value lines to stdout and
// appends them to result.txt next to this file.
'use strict';
const http = require('http');
const tls = require('tls');
const fs = require('fs');
const path = require('path');

const PROXY_HOST = '127.0.0.1';
const PROXY_PORT = Number(process.env.B2_PROXY_PORT || '39877');
const REPORT = path.join(__dirname, 'node-result.txt');
const PAID = process.env.B2_PAID === '1';
const GH_TOKEN = process.env.COPILOT_GITHUB_TOKEN || '';

function safe(value) {
  let out = '';
  for (const ch of String(value)) {
    const code = ch.codePointAt(0);
    if (code >= 32 && code < 127) out += ch;
    else if (code >= 128) out += '\\u' + code.toString(16).padStart(4, '0');
    else out += '\\x' + code.toString(16).padStart(2, '0');
  }
  return out;
}
function addLine(key, value) {
  const line = key + '=' + safe(value);
  console.log(line);
  fs.appendFileSync(REPORT, line + '\n');
}

function connectTunnel(target) {
  return new Promise((resolve, reject) => {
    const req = http.request({ host: PROXY_HOST, port: PROXY_PORT, method: 'CONNECT', path: target + ':443' });
    req.on('connect', (res, socket) => {
      if (res.statusCode !== 200) { reject(new Error('proxy CONNECT status ' + res.statusCode)); return; }
      resolve(socket);
    });
    req.on('error', reject);
    req.end();
  });
}

async function httpsJson(host, urlPath, method, headers, body) {
  const socket = await connectTunnel(host);
  const secure = tls.connect({ socket, servername: host });
  await new Promise((resolve, reject) => {
    secure.once('secureConnect', resolve);
    secure.once('error', reject);
  });
  const payload = body ? Buffer.from(body, 'utf8') : null;
  const lines = [method + ' ' + urlPath + ' HTTP/1.1', 'Host: ' + host];
  for (const key of Object.keys(headers)) lines.push(key + ': ' + headers[key]);
  if (payload) {
    lines.push('Content-Type: application/json');
    lines.push('Content-Length: ' + payload.length);
  }
  lines.push('Connection: close', '', '');
  secure.write(lines.join('\r\n'));
  if (payload) secure.write(payload);
  const chunks = [];
  await new Promise((resolve) => {
    secure.on('data', (chunk) => chunks.push(chunk));
    secure.on('end', resolve);
    secure.on('close', resolve);
    secure.on('error', resolve);
  });
  const raw = Buffer.concat(chunks).toString('utf8');
  const split = raw.indexOf('\r\n\r\n');
  const head = split >= 0 ? raw.slice(0, split) : raw;
  const rest = split >= 0 ? raw.slice(split + 4) : '';
  const statusMatch = head.match(/^HTTP\/1\.1\s+(\d{3})/);
  const status = statusMatch ? Number(statusMatch[1]) : 0;
  // chunked transfer decoding (GitHub answers with chunked encoding)
  let bodyText = rest;
  if (/transfer-encoding:\s*chunked/i.test(head)) {
    let out = '';
    let index = 0;
    while (index < rest.length) {
      const lineEnd = rest.indexOf('\r\n', index);
      if (lineEnd < 0) break;
      const size = parseInt(rest.slice(index, lineEnd), 16);
      if (!size) break;
      out += rest.slice(lineEnd + 2, lineEnd + 2 + size);
      index = lineEnd + 2 + size + 2;
    }
    bodyText = out;
  }
  const requestId = (head.match(/x-request-id:\s*([^\r\n]+)/i) || [])[1] || '';
  return { status, head, body: bodyText, requestId };
}

async function main() {
  try { fs.rmSync(REPORT, { force: true }); } catch (error) { /* ignore */ }
  addLine('client', 'node ' + process.version);
  addLine('proxy', PROXY_HOST + ':' + PROXY_PORT);
  addLine('cwd', process.cwd());
  addLine('githubTokenPresent', String(Boolean(GH_TOKEN)));
  addLine('githubTokenLength', String(GH_TOKEN.length));
  addLine('paid', String(PAID));

  const baseHeaders = {
    Accept: 'application/json',
    'User-Agent': 'GitHubCopilotChat/0.26.7',
    'Editor-Version': 'vscode/1.99.0',
    'Editor-Plugin-Version': 'copilot-chat/0.26.7',
    'Copilot-Integration-Id': 'vscode-chat'
  };

  let copilotToken = process.env.B2_COPILOT_TOKEN || '';
  if (copilotToken) {
    addLine('copilotTokenSource', 'gh-exchange-injected');
    addLine('copilotTokenPresent', 'true');
  }
  const attempts = copilotToken ? [] : [
    { label: 'cli-ua', headers: { 'User-Agent': 'GitHubCopilotCLI/1.0.83', Accept: 'application/json' }, scheme: 'token' },
    { label: 'gh-api-version', headers: { 'User-Agent': 'GitHubCopilotCLI/1.0.83', Accept: 'application/vnd.github+json', 'X-GitHub-Api-Version': '2022-11-28' }, scheme: 'token' },
    { label: 'vscode-profile', headers: baseHeaders, scheme: 'token' },
    { label: 'curl-ua', headers: { 'User-Agent': 'curl/8.4.0', Accept: 'application/json' }, scheme: 'token' }
  ];
  for (const attempt of attempts) {
    try {
      const tokenHeaders = Object.assign({}, attempt.headers, { Authorization: attempt.scheme + ' ' + GH_TOKEN });
      const response = await httpsJson('api.github.com', '/copilot_internal/v2/token', 'GET', tokenHeaders, null);
      addLine('tokenStatus_' + attempt.label, String(response.status));
      if (response.status !== 200) addLine('tokenBodyHead_' + attempt.label, response.body.slice(0, 120));
      if (response.status === 200) {
        const parsed = JSON.parse(response.body);
        copilotToken = parsed.token || '';
        addLine('copilotTokenPresent', String(Boolean(copilotToken)));
        addLine('tokenProfile', attempt.label);
        break;
      }
    } catch (error) {
      addLine('tokenError_' + attempt.label, error.message);
    }
  }
  if (!copilotToken) { addLine('aborted', 'no copilot token'); return; }

  const apiHeaders = Object.assign({}, baseHeaders, { Authorization: 'Bearer ' + copilotToken });
  try {
    const models = await httpsJson('api.individual.githubcopilot.com', '/models', 'GET', apiHeaders, null);
    addLine('modelsStatus', String(models.status));
    addLine('modelsBodyHead', models.body.slice(0, 200));
  } catch (error) {
    addLine('modelsError', error.message);
  }
  if (PAID) {
    const payload = JSON.stringify({
      messages: [{ role: 'user', content: 'Reply with exactly PONG' }],
      max_tokens: 16,
      temperature: 0
    });
    try {
      const chat = await httpsJson('api.individual.githubcopilot.com', '/chat/completions', 'POST', apiHeaders, payload);
      addLine('chatStatus', String(chat.status));
      addLine('chatRequestId', chat.requestId);
      addLine('chatBodyHead', chat.body.slice(0, 1200));
      addLine('pongSeen', String(/PONG/.test(chat.body)));
    } catch (error) {
      addLine('chatError', error.message);
    }
  }
  addLine('done', 'true');
}

main().catch((error) => { addLine('fatal', error && error.message ? error.message : String(error)); });
