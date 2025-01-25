# Exo - A Note-Taking System

Exo is a command-line tool for organizing your thoughts, ideas, and knowledge through various note types and organizational structures.

## Features

- **Daily Notes**: Create and manage daily notes for journaling and task tracking
- **Zettel Notes**: Create and organize atomic notes using the Zettelkasten method
- **Ideas**: Capture and develop ideas quickly
- **Templates**: Customize note templates to match your workflow
- **Configuration**: Flexible configuration options for paths and settings

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/a-kostevski/exo.git
   cd exo
   ```

2. Build the project:
   ```bash
   make build
   ```

3. Initialize Exo:
   ```bash
   exo init
   ```

## Usage

### Daily Notes

Create or open today's daily note:
```bash
exo day
```

### Zettel Notes

Create a new Zettel note:
```bash
exo zet "Your note title"
```

### Ideas

Create a new idea note:
```bash
exo idea "Your idea title"
```

### Templates

List available templates:
```bash
exo templates
```

Install default templates:
```bash
exo templates install
```

### Configuration

List all configuration settings:
```bash
exo config
```

Get a specific setting:
```bash
exo config get data_home
```

Set a configuration value:
```bash
exo config set editor "code -w"
```

## Directory Structure

- `cmd/`: Command-line interface implementation
  - `notes/`: Note-related commands (day, zet, idea)
  - `system/`: System commands (init, config, templates)
- `internal/`: Internal packages
  - `config/`: Configuration management
  - `logger/`: Logging functionality
  - `note/`: Note management and types
  - `templates/`: Template management
  - `utils/`: Utility functions

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
