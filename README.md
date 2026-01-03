# WAF Rate Limit Tester

Tool for testing NGWAF block & rate limiting rules.

## Features

- **Traffic Modes**: Normal traffic vs. Attack traffic (SQL injection, XSS, path traversal, command injection)
- **Test Modes**: Baseline (evenly distributed) vs. Burst (concentrated spike)
- **Real-time Streaming**: Live progress updates via Server-Sent Events (SSE)
- **Export Results**: JSON and CSV formats
- **Detailed Statistics**: Response times, percentiles (p50/p95/p99), rate limit detection
- **User Agent Control**: Legitimate browser vs. Scanner signatures

## Quick Start

### 1. Download

git clone/download the repo

### 2. Run

cd waf-tester

```
chmod +x waf-tester
./waf-tester
```

### 3. Open Browser

Navigate to `http://localhost:8080`

### 4. Start Testing

Configure and run your first test!

**macOS binary**
MAC users: That's it! No installation, no dependencies.

**Linux/Windows binary**
Build from source (see below).

### Build from Source

**Prerequisites:**

- Go 1.16+ ([Download](https://golang.org/dl/))
- Node.js 14+ ([Download](https://nodejs.org/))

**Build:**

```bash
# Install dependencies (one-time)
make install

# Build single binary
make build

# Or use the build script
./build.sh
```

**Run:**

```bash
./waf-tester
```

Output:

```
WAF Tester started
Open in browser: http://localhost:8080
Press Ctrl+C to stop
```

## Configuration

### Port Configuration

Set custom port with environment variable:

```bash
PORT=9000 ./waf-tester
```

Or create a `.env` file:

```env
PORT=8080
```

### Attack Payloads

Place payload files in `payloads/` directory:

- `sql-injection.txt`
- `xss.txt`
- `path-traversal.txt`
- `command-injection.txt`
- `scanner-user-agents.txt`
- `legitimate-user-agents.txt`
  The binary will load these automatically if present.

## API Endpoints

- `POST /api/test/stream` - Start streaming test with real-time updates
- `POST /api/test/stop` - Cancel a running test
- `POST /api/test/export` - Export results to JSON/CSV
- `GET /api/exports/list` - List all exported files

## Building

### Single Binary

```bash
make build
```

Produces: `./waf-tester`

### Cross-Compile for All Platforms

```bash
make build-all
```

Produces binaries in `./dist/`:

- `waf-tester-darwin-amd64` (macOS Intel)
- `waf-tester-darwin-arm64` (macOS Apple Silicon)
- `waf-tester-linux-amd64` (Linux)
- `waf-tester-windows-amd64.exe` (Windows)

### Development Mode

```bash
make dev
```

Runs separate backend and frontend servers for development.

## Make Commands

```bash
make              # Build single binary
make install      # Install dependencies
make build        # Build single binary
make frontend     # Build frontend only
make backend      # Build backend only
make build-all    # Cross-compile for all platforms
make clean        # Remove build artifacts
make dev          # Run development servers
make help         # Show help
```

## Usage Examples

1. **Test baseline traffic**: Configure 100 requests over 10 seconds with normal traffic
2. **Test burst traffic**: Configure 500 requests in burst mode to simulate traffic spike
3. **Test with attacks**: Enable attack mode to inject SQL/XSS payloads
4. **Export results**: Click Export JSON/CSV after test completion

## Troubleshooting

**Port already in use:**

```bash
PORT=9000 ./waf-tester
```

**Binary won't run (macOS):**

```bash
chmod +x waf-tester
# If still blocked: System Settings → Privacy & Security → Allow
```

**Need to rebuild:**

```bash
make clean
make build
```

## Safety Reminders

- Only test against URLs you own or have explicit permission to test
- Respect rate limits (max 10,000 requests per test)
- Attack mode is for authorized security testing only
- Some targets may temporarily block your IP if you exceed rate limits

## License

MIT

# waf-rule-stresser
