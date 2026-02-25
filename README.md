# Notify Record

The notify-record is a simple tool that listens for D-Bus notifications on the local session bus and prints them in a human-readable format. It can be configured to filter notifications based on keywords.
It prints formatted notifications to stdout, which can be piped to a file or another program for further processing.

## Configuration

The configuration file follows the XDG Base Directory specification:

- `$XDG_CONFIG_HOME/notify-record/config.yml` (if `$XDG_CONFIG_HOME` is set)
- `~/.config/notify-record/config.yml` (default)

To set up your configuration:

```bash
mkdir -p ~/.config/notify-record
cp example/config.yml ~/.config/notify-record/config.yml
```

## Example

Terminal1:

```
$ go run main.go 
```

Terminal2:

```
$ notify-send 'Hello world!' 'This is an example notification.'
```

then,

Terminal1:

```
Hello world!: This is an example notification.
```