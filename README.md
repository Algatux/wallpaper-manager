# Wallpaper Manager Daemon

A daemon application that automates wallpaper rotation for Hyprland window manager. The daemon acts as a wrapper around **awww** (a fast and efficient wallpaper daemon for Wayland compositors), enabling automated scheduling and multi-monitor wallpaper management.

## Overview

**wallpaper-manager** is a lightweight background service designed to:

- Automatically rotate wallpapers at configurable intervals
- Manage independent wallpaper pools for multiple monitors
- Seamlessly handle wallpaper changes across Hyprland displays
- Run efficiently in the background with minimal overhead

The daemon reads wallpapers from configured directories, discovers supported image formats, and orchestrates changes through the awww wallpaper daemon.

## Usage

### Basic Invocation

```bash
./wallpaper-manager [options] <config.json>
```

### Command-Line Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-i` | int | `3600` | Interval in seconds between automatic wallpaper changes |
| `-l` | string | `"/tmp/wallpaper_cycle.lock"` | Path to the lock file used to synchronize wallpaper cycling |
| `-d` | bool | `true` | Enable debug mode to print verbose diagnostic messages |

### Positional Arguments

| Argument | Description |
|----------|-------------|
| `<config.json>` | Path to the configuration file in JSON format |

### Configuration Format

Create a JSON file defining monitors and their wallpaper directories:

```json
{
  "monitors": [
    {
      "name": "HDMI-1",
      "path": "/home/user/.wallpapers/main",
      "flags": ""
    },
    {
      "name": "DP-1",
      "path": "/home/user/.wallpapers/secondary",
      "flags": ""
    }
  ]
}
```

### Usage Examples

**Rotate wallpapers every 2 hours (7200 seconds) on two monitors:**

```bash
./wallpaper-manager -i 7200 config.json
```

**Run with custom lock file:**

```bash
./wallpaper-manager -l /var/run/wallpaper_cycle.lock
```

**Use default settings:**

```bash
./wallpaper-manager config.json
```

**Display help and version information:**

```bash
./wallpaper-manager -help
```

## How It Works

1. **Initialization** — The daemon parses command-line arguments and loads the monitor configuration
2. **Wallpaper Discovery** — For each configured monitor, the daemon scans the specified directory and discovers supported image files
3. **Rotation Loop** — At the specified interval, the daemon cycles through available wallpapers and applies them to each monitor via awww
4. **Lock File** — A lock file ensures synchronized wallpaper cycling across multiple processes or system instances

## Supported Image Formats

The daemon supports common image formats including:
- PNG (.png)
- JPEG (.jpg, .jpeg)
- WebP (.webp)

## Building

```bash
go build -o wallpaper-manager ./cmd/daemon
```

## Requirements

- **Go 1.19+** for building
- **awww** daemon running (for actual wallpaper application)
- **Hyprland** window manager (or other Wayland compositors supported by awww)
