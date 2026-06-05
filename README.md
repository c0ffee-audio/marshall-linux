# marshall-linux

Linux application to control Marshall Bluetooth devices - ANC mode, EQ presets, and battery level - over BLE/GATT.

Available as a desktop GUI (Wails + Svelte) and a CLI.

## Features

- **Noise Control** - switch between Off, ANC, and Transparency
- **Equalizer** - 11 presets (Flat, Rock, Metal, Pop, Hip-Hop, Electronic, Jazz, Bass Boost, Mid Boost, Treble Boost, Workout)
- **Battery** - live battery level
- **Device info** - model name and firmware version

## Supported devices

Any device using the Zound Industries BLE protocol (`-1337-1dea-feed-c0ffee70c0de` UUID suffix):

- Marshall Motif II ANC
- Marshall Major IV, Monitor III ANC, and others

> See [marshall-protocol](https://github.com/bbuddha/marshall-protocol) for the full reverse-engineered protocol documentation.

## Requirements

- Linux with BlueZ 5.50+
- `webkit2gtk-4.1` (GUI only)
- Go 1.21+
- [Wails v2](https://wails.io) (GUI only)

```
sudo pacman -S webkit2gtk-4.1   # Arch
sudo apt install libwebkit2gtk-4.1-dev  # Debian/Ubuntu
```

## Build

### GUI

```bash
wails build -tags webkit2_41
./build/bin/marshall-linux
```

### CLI

```bash
go build ./cmd/marshall/
./marshall
```

## CLI usage

```
marshall <name|address> <command> [args]

Commands:
  info                          Model, firmware, battery, ANC mode
  anc [off|anc|transparency]    Get or set ANC mode
  eq <preset>                   Set EQ preset
  battery                       Battery level
  scan                          List all GATT characteristics

Examples:
  marshall "MOTIF II A.N.C." info
  marshall "MOTIF II A.N.C." anc transparency
  marshall "MOTIF II A.N.C." eq rock
```

## How it works

The device exposes two separate Bluetooth transports:

- **BR/EDR** (public address) - A2DP audio
- **BLE** (random address) - GATT control

On first connection, BlueZ pairs the BLE device using Just Works (`NoInputNoOutput`). Subsequent connections reuse the stored bond.

## License

MIT
