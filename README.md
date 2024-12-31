# Proxy Manager

A simple command-line tool to manage system-wide proxy settings on Linux systems. This tool helps you quickly enable or disable proxy configurations in both `/etc/environment` and `.zshrc` files.

## Installation

1. Clone or download the source code, then navigate to the project directory:
```bash
cd proxy-manager
```

2. Modify the proxy settings in `proxy-manager.go`:
```go
const (
    proxyIP   = "your_proxy_ip"
    proxyPort = "your_proxy_port"
)
```

3. Build the program:
```bash
go build proxy-manager.go
```

4. Install it system-wide (optional but recommended):
```bash
sudo cp proxy-manager /usr/local/bin/
```

## Usage

If you installed it system-wide:
```bash
sudo proxy-manager -aktif
sudo proxy-manager -mati
```

If you're running it from the project directory:
```bash
sudo ./proxy-manager -aktif
sudo ./proxy-manager -mati
```

## What It Does

When enabling proxy (`-aktif`):
- Adds proxy settings to `/etc/environment`
- Adds export commands to your `.zshrc`

When disabling proxy (`-mati`):
- Removes proxy settings from `/etc/environment`
- Removes proxy-related export commands from `.zshrc`

## Requirements

- Linux operating system
- Go compiler (for building)
- sudo privileges (for running)
- ZSH shell with `.zshrc` configuration file

## Note

After enabling or disabling the proxy, you may need to:
1. Restart your terminal
2. Source your `.zshrc`:
```bash
source ~/.zshrc
```

## Troubleshooting

If you encounter permission errors:
- Make sure you're running the command with `sudo`
- Verify that the executable has proper permissions:
```bash
sudo chmod +x /usr/local/bin/px 
```
