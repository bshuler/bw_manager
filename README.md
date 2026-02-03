# BoxWallet Manager

The quickest way to get [BoxWallet](https://github.com/richardltc/boxwallet) up and running.

BoxWallet Manager automatically downloads and installs the latest version of BoxWallet, creating everything you need to run it in your browser.

---

## Features

- Automatic detection of latest BoxWallet release
- Cross-platform support (Linux, macOS, Windows)
- Creates organized version-specific directories
- Generates platform-appropriate run scripts
- Zero configuration required

---

## Installation

### Download

Download the appropriate binary for your system from the [Releases](https://github.com/bshuler/bw_manager/releases/latest) page:

| Operating System | Architecture | Download |
|------------------|--------------|----------|
| Linux | AMD64 (Intel/AMD) | `bw_manager_lin_amd64` |
| Linux | ARM64 (Raspberry Pi, etc.) | `bw_manager_lin_arm64` |
| macOS | Intel | `bw_manager_mac_amd64` |
| macOS | Apple Silicon (M1/M2/M3) | `bw_manager_mac_arm64` |
| Windows | AMD64 (Intel/AMD) | `bw_manager_win.exe` |

---

## Platform-Specific Instructions

### Linux

```bash
# Download (choose the correct architecture)
wget https://github.com/bshuler/bw_manager/releases/latest/download/bw_manager_lin_amd64

# Make executable
chmod +x bw_manager_lin_amd64

# Run
./bw_manager_lin_amd64
```

### macOS

macOS includes Gatekeeper security that blocks unsigned applications. You'll need to remove the quarantine attribute before running.

#### Intel Mac

```bash
# Download
curl -LO https://github.com/bshuler/bw_manager/releases/latest/download/bw_manager_mac_amd64

# Remove macOS quarantine attribute
xattr -d com.apple.quarantine bw_manager_mac_amd64

# Make executable
chmod +x bw_manager_mac_amd64

# Run
./bw_manager_mac_amd64
```

#### Apple Silicon (M1/M2/M3)

```bash
# Download
curl -LO https://github.com/bshuler/bw_manager/releases/latest/download/bw_manager_mac_arm64

# Remove macOS quarantine attribute
xattr -d com.apple.quarantine bw_manager_mac_arm64

# Make executable
chmod +x bw_manager_mac_arm64

# Run
./bw_manager_mac_arm64
```

> **Note:** If you see "cannot be opened because the developer cannot be verified", run the `xattr` command above to remove the quarantine flag.

### Windows

1. Download `bw_manager_win.exe` from the [Releases](https://github.com/bshuler/bw_manager/releases/latest) page
2. Open the folder containing the downloaded file
3. Double-click `bw_manager_win.exe` to run

> **Note:** Windows SmartScreen may show a warning. Click "More info" then "Run anyway" to proceed.

---

## Usage

1. Run `bw_manager` in the directory where you want to install BoxWallet
2. The manager will:
   - Detect the latest BoxWallet version
   - Create a version-specific subfolder (e.g., `v0.0.5`)
   - Download and extract BoxWallet
   - Generate the appropriate run script for your platform
3. Navigate to the created directory
4. Run the startup script:
   - **Linux/macOS:** `./run_boxwallet.sh`
   - **Windows:** `run_boxwallet.bat`
5. Open your browser to `http://localhost:4000`

---

## Verifying Downloads

Each release includes a `checksums.txt` file containing SHA256 hashes. Verify your download:

```bash
# Linux/macOS
sha256sum bw_manager_lin_amd64
# Compare output with checksums.txt

# Windows (PowerShell)
Get-FileHash bw_manager_win.exe -Algorithm SHA256
```

---

## Building from Source

Requires Go 1.24 or later.

```bash
# Clone the repository
git clone https://github.com/bshuler/bw_manager.git
cd bw_manager

# Build for your current platform
go build -ldflags="-s -w" -o bw_manager

# Or use the build script for all platforms
./build.sh
```

---

## License

See [LICENSE](LICENSE) for details.
