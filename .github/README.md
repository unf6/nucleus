# nucleus-cli

`nucleus-cli` is the command-line interface for managing and running **Nucleus Shell**, a modern Wayland shell built on **QuickShell** for **Hyprland**.

This repository contains the **cli** for **nucleus-shell**

---

## Features

* Start and stop Nucleus Shell
* Reload running QuickShell instances
* Foreground (debug) and daemonized modes
* Plugin management namespace (WIP)

---

## Requirements

* Linux
* Hyprland
* QuickShell (available in `$PATH`)
* Go ≥ 1.25

---

## Installation

### Build from source

```sh
git clone https://github.com/unf6/nucleus.git nucleus-cli
cd nucleus-cli
go mod tidy
go build
```

(Optional)

```sh
sudo mv nucleus /usr/local/bin/
```

---

## Usage

```sh
nucleus <command> [flags]
```

---

## Commands

### `nucleus run`

Start the Nucleus Shell using QuickShell.

```sh
nucleus run
```

#### Flags

| Flag       | Short | Description                                        |
| ---------- | ----- | -------------------------------------------------- |
| `--reload` | `-r`  | Kill existing QuickShell instances before starting |
| `--debug`  | `-d`  | Run in foreground (no daemon)                      |

#### Examples

Run normally (daemonized):

```sh
nucleus run
```

Run in foreground:

```sh
nucleus run --debug
```

Restart shell cleanly:

```sh
nucleus run --reload
```

---

### `nucleus plugins` (WIP)

Plugin management command namespace.

```sh
nucleus plugins
```

> Subcommands are not implemented yet.

Planned plugin sources:

* Official: `https://github.com/xZepyx/nucleus-plugins.git`

---

## Configuration Paths

The CLI expects Nucleus Shell to be installed at:

```text
~/.config/quickshell/nucleus-shell/
```

Required file:

```text
shell.qml
```

The shell is considered **installed** if this file exists:

```text
~/.config/quickshell/nucleus-shell/shell.qml
```

If missing, `nucleus run` will refuse to start.

---
## Related Repositories

* nucleus-shell: https://github.com/xZepyx/nucleus-shell
* nucleus-plugins: https://github.com/xZepyx/nucleus-plugins
* nucleus-colorschemes: https://github.com/xZepyx/nucleus-colorschemes
* nucleus(cli): https://github.com/unf6/nucleus (this repo)
* zenith(ai backend): https://github.com/xZepyx/zenith
 
---

## Process Behavior

* Launches `quickshell` with `--no-duplicate`
* Daemonized by default
* Foreground mode when `--debug` is set
* Reload uses `pkill -f quickshell`
* Graceful shutdown on `SIGINT` / `SIGTERM`

---

## Internal Layout

```text
cmd/
 ├─ root.go        # Root Cobra command
 ├─ run.go         # `nucleus run`
 └─ plugins/       # Plugin namespace (WIP)

internal/
 ├─ config/        # Config paths, install checks
 └─ shell/         # QuickShell process control
```

---

## Logging

* Uses `github.com/charmbracelet/log`
* Logs to stdout/stderr
* Debug output enabled via `--debug`

---

## Status

* CLI core: usable
* Plugin system: WIP
* Installer command: planned

---

## License

MIT
