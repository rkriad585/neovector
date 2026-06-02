# Windows Setup

## Prerequisites

- Windows 10 or later
- PowerShell 5.1+ or PowerShell 7+

## Installation

### Option 1: Install script (recommended)

```powershell
irm https://raw.githubusercontent.com/rkriad585/neovector/main/installer.ps1 | iex
```

Restart your terminal, then:

```powershell
neovector --help
```

### Option 2: Go install

```powershell
go install github.com/rkriad585/neovector@latest
```

### Option 3: Download binary

Download from [releases](https://github.com/rkriad585/neovector/releases).
Place `neovector-windows-amd64.exe` in a directory on your `PATH`
and rename it to `neovector.exe`.

## Build from source

```powershell
git clone https://github.com/rkriad585/neovector.git
cd neovector
go build -o neovector.exe .
```

## Uninstall

```powershell
irm https://raw.githubusercontent.com/rkriad585/neovector/main/installer.ps1 | iex -- --selfuninstall
```
