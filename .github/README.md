<div align="center">

# ✦ Nucleus ✦

<p>
  <img src="https://img.shields.io/github/last-commit/unf6/nucleus?style=for-the-badge&color=8ad7eb&logo=git&logoColor=D9E0EE&labelColor=1E202B" alt="Last Commit" />
  <img src="https://img.shields.io/github/stars/unf6/nucleus?style=for-the-badge&logo=andela&color=86dbd7&logoColor=D9E0EE&labelColor=1E202B" alt="Stars" />
  <img src="https://img.shields.io/github/repo-size/unf6/nucleus?color=86dbce&label=SIZE&logo=protondrive&style=for-the-badge&logoColor=D9E0EE&labelColor=1E202B" alt="Repo Size" />
  &nbsp;
  <img src="https://img.shields.io/badge/Maintenance-Active%20-6BCB77?style=for-the-badge&logo=vercel&logoColor=D9E0EE&labelColor=1E202B" alt="Maintenance" />
</p>

</div>

---
<h2 align="center">✦ Overview ✦ </h2>

#### A blazingly fast CLI for managing and supercharging nucleus-shell.

---

## Features

* Start and stop Nucleus Shell
* Reload running QuickShell instances
* Foreground (debug) and daemonized modes
* Plugin management namespace
* Install and Update nucleus-shell
* Manage quickshell ipc's

---


> [!IMPORTANT]
> * You can also join the [discord server](https://discord.gg/SvQMhuMXXa) for help.
> * **Before reporting an issue:**
  If you encounter a problem in the current release, please first test against the latest source code by cloning the repository (`git clone ...`). This ensures you are not reporting an issue that has already been fixed.
  Only open an issue if the problem is still reproducible on the latest source.

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

(Required)

```sh
sudo mv nucleus /usr/local/bin/
```

### Installing Via Aur

```sh
yay -S nucleus-cli
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

```Nucleus is a beautiful, customizable shell built to get things done

Usage:
  nucleus [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  install     Install or Update Nucleus Shell
  ipc         Interact with the shell via IPC
  plugins     Manage Nucleus Shell plugins
  run         Start The Nucleus Shell

Flags:
  -h, --help   help for nucleus

Use "nucleus [command] --help" for more information about a command.
```

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

### `nucleus plugins`

Plugin management command namespace.

```sh
nucleus plugins
```

Plugin sources:

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
