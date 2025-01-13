# deeploy

A next-gen deployment platform for terminal-loving devs. Ship apps with a sleek TUI or CLI - built for those who live and breathe the command line.

## Features

- Modern Terminal UI AND CLI commands
- Docker-based deployments
- Open source and self-hosted
- Built with Go

## Why deeploy?

- First deployment platform with a proper Terminal UI (TUI first approach)
- Not just another web dashboard (but we'll have that too)
- Simple, fast, and dev-friendly
- Built with Go, not PHP or JS

## Quick Start

```bash
# Install server
curl -fsSL https://deeploy.sh/install.sh | sh

# Install CLI tool
brew install deeploy-cli
# or
curl -fsSL deeploy.sh/install-cli | sh

# Connect to your server
deeploy connect http://your-server:8090
```

## Usage

```bash
# Interactive TUI mode
deeploy

# Or use CLI commands
deeploy deploy myapp
deeploy logs myapp
deeploy status
```

## Inspiration

Deeploy draws inspiration from popular open source deployment platforms:

- [Coolify](https://coolify.io)
- [Dokploy](https://dokploy.com)

## Current Status

This is a pre-alpha release. The platform is under active development with upcoming features including:

- User Authentication
- Project Management
- Container Deployments
- Domain Management
- Templates
- And more!

## Requirements

- Linux server (Ubuntu 24.04 recommended)
- Docker

## Built With

- Go
- Bubbletea (TUI)
- SQLite
- Templ + templUI (Web UI coming soon)
- Docker

## Contributing

We welcome contributions from the community! Whether it's adding new features, improving existing ones, or enhancing documentation, your input is valuable. Please check our [contributing guidelines](CONTRIBUTING.md) for more information on how to get involved.

## License

Deeploy is open-source software licensed under the [MIT license](LICENSE).

## Support

For support, questions, or discussions, please [open an issue](https://github.com/axzilla/deeploy/issues) on our GitHub repository or [visit our community (GitHub Discussions)](https://github.com/axzilla/deeploy/discussions).

---

Built with ❤️ by the dev community, for the dev community.
