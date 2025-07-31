# Keylight

**Keylight** is a lightweight global keyboard event hook and overlay system written entirely in Go for Windows, without any C++ runtime dependencies. It listens to global key events, allows dynamic control of event forwarding to an on-screen overlay, and can be run silently in the background.

---

## Features

- Pure Go implementation (no C++ runtime required)
- Global keyboard hook using [`moutend/go-hook`](https://github.com/moutend/go-hook)
- Transparent, always-on-top overlay rendered with raw Win32 API calls (no external GUI libraries)
- Dual-channel event handling
  - Internal channel for logic and control
  - External channel for overlay updates
- Toggle behavior using control key + hotkeys:
  - `CTRL + 1`: Enable overlay event forwarding
  - `CTRL + 2`: Disable event forwarding
  - `CTRL + 3`: Exit the application
- Hidden console window using `-ldflags="-H=windowsgui"`

---

## Installation

### 1. Clone

```bash
git clone https://github.com/renniemaharaj/keylight
cd keylight
```

### 2. Build

```bash
go build -ldflags="-H=windowsgui" -o keylight.exe ./cmd
```
