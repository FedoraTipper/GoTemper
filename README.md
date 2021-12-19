# GoTemper
Golang reader implementation for TEMPer USB with InfluxDB support

# About

A lot of credit goes to urwen's temper project (https://github.com/urwen/temper), all inner workings of this project were taken from urwen's code (Albeit rewritten in Go).
That said I needed something that had multiple ways of posting stats, and in a more complete Golang solution. Thus, I decided to rewrite urwen's implementation with some key additional features added on.

# Configuration
All configuration files are to be structured as YAML files. An example configuration file can be found in `./configs/`.

The configuration path is hard set for now. (To lazy to add custom flag)

`config.yml` file will be searched for in these directories:

```
./
./configs/
$HOME/<program>/
$HOME/.config/<program-name>/

e.g. ./config/gotemper/config.yml
```

## Logging
List of available logging levels:
- debug
- info (default)
- warning
- error
- fatal

# TODO & Additional features:
- Add humidity logic (I don't have a temper device that has this feature) 
  - `PostStats` in `cmd/gotemper/main.go` will need to reworked to include dynamic payload construction
- Add other support for other temper devices 
  - Create a driver implementation in `internal/temper/<temperDeviceName>.go`
  - Add a driver, Vendor and Product IDs to `temperDevicesMap` in `internal/temper/constants.go` (Bear in mind the read offset)
- Add tty/serial support (https://github.com/urwen/temper/blob/master/temper.py#L220)