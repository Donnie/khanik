# Khanik - SSH Surang Manager

Khanik is a Go-based tool for managing SSH tunnels (referred to as "surangs"). It provides a daemon that automatically starts, monitors, and restarts SSH tunnels based on a configuration file.

## Features

- Daemon-based management of SSH tunnels
- Automatic restart of failed tunnels
- IP verification for each tunnel
- Easy-to-use CLI commands

## Installation

To download Khanik (a single executable) do:

```
gh release download 0.0.4 -R Donnie/khanik
chmod +x khanik-macos-arm64
mv khanik-macos-arm64 /usr/local/bin/khanik
which khanik
khanik version
```

## Configuration

Create a `config.yaml` file in a directory. Here's an example configuration:

```yaml
surangs:
  home:
    command: "user@host1.example.com"
    expect_ip: "203.0.113.1"
    port: 8080
  office:
    command: "user@host2.example.com"
    expect_ip: "203.0.113.2"
    port: 8081
```

## Usage

Khanik provides the following commands:

- `khanik start`: Start the surang manager daemon
- `khanik stop`: Stop the surang manager daemon
- `khanik list`: List all configured surangs and their status
- `khanik version`: Display the version of Khanik

### Examples

Start the daemon:
```
khanik start
```

List all surangs:
```
khanik list
```

Stop the daemon:
```
khanik stop
```

## How It Works

Khanik uses SSH's built-in SOCKS proxy functionality to create tunnels. The daemon periodically checks each configured tunnel to ensure it's running and functioning correctly. If a tunnel fails, Khanik automatically attempts to restart it.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)
