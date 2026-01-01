# LinkSleuth

LinkSleuth is a fast, reliable, and extendable URL discovery and analysis tool written in Go. It allows security researchers and developers to discover endpoints, analyze HTTP responses, and detect sensitive exposure through a modular and multi-threaded approach.

## Features

- **Multi-threaded URL Discovery**: Fast scanning using worker pools.
- **HTTP Response Analysis**: Automatic categorization of status codes (2xx, 3xx, 4xx, 5xx).
- **Security Checks**: Detection of sensitive endpoints (admin, login, backup, config).
- **Flexible Reporting**: Export results to JSON, CSV, and interactive HTML formats.
- **User-Agent Rotation**: Built-in support for rotating user-agents to avoid rate-limiting.
- **Modular Architecture**: Easy to extend with custom analyzer or reporter modules.

## Installation

```bash
go install github.com/ismailtsdln/linksluth@latest
```

Or build from source:

```bash
git clone https://github.com/ismailtsdln/LinkSleuth.git
cd LinkSleuth
go build -o linksluth main.go
```

## Usage

### 1. Scan a Target

Scan a domain using a wordlist and save results to JSON.

```bash
./linksluth scan -u https://example.com -w wordlist.txt -o results.json
```

### 2. Analyze Results

Analyze previously saved JSON results in the console.

```bash
./linksluth analyze -i results.json
```

### 3. Generate Reports

Generate an interactive HTML report from JSON results.

```bash
./linksluth report -i results.json -o report.html
```

## Flags

- `-u, --url` : Target URL (required for scan)
- `-w, --wordlist` : Optional directory/file wordlist
- `-t, --threads` : Concurrent threads (default 10)
- `-r, --retry` : Retry failed requests (default 2)
- `-a, --agent` : Custom User-Agent
- `-o, --output` : Output file path
- `-v, --verbose` : Verbose logging

## License

MIT
