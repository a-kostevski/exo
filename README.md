# Productivity Tools

A collection of command-line tools to enhance personal knowledge management and digital workflow.

## Tools

### Zet

A CLI tool for managing a Zettelkasten note-taking system. Creates and manages markdown notes with automatic daily logs and templating support.

#### Features

- Create new notes with automatic linking
- Daily note management
- Template support for notes and daily entries
- Automatic cross-linking between notes
- Configurable editor integration

#### Usage

```bash
# Create a new note
zet new "Note Title"

# Create a daily note
zet new --day

# Create note without opening editor
zet new "Note Title" --no-open
```

#### Environment Variables

- `LIFE`: Base directory for daily notes and templates
- `ZETDIR`: Directory for Zettelkasten notes
- `EDITOR`: Preferred text editor

## Installation

```bash
# Install using go
go install github.com/your-username/productivity-tools/cmd/zet@latest

# Or build from source
git clone https://github.com/your-username/productivity-tools.git
cd productivity-tools
go install ./cmd/zet
```

## Configuration

By default, the configuration file is searched in the following locations:
- `$XDG_CONFIG_HOME/zet/zet.yaml`
- `$HOME/zet.yaml`

Example configuration:
```yaml
# zet.yaml
templates:
  path: ~/.config/zet/templates
```

## Directory Structure

```
.
├── cmd/          # Command-line tools
│   └── zet/      # Zettelkasten note manager
├── internal/     # Internal packages
├── pkg/         # Public packages
└── docs/        # Documentation
```

## Development

Requirements:
- Go 1.21 or higher
- Git

```bash
# Clone repository
git clone https://github.com/your-username/productivity-tools.git

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build ./cmd/zet
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Planned Features

- [ ] Project management tool
- [ ] Reference management system
- [ ] Enhanced note searching and linking
- [ ] Export capabilities
- [ ] Integration with external tools

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
