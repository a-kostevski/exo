# Exo - Personal Knowledge Management System

Exo is a command-line tool for managing personal knowledge through various note types and organizational structures.

## Features

- **Multiple Note Types**:
  - Zettelkasten notes for atomic, interconnected ideas
  - Daily notes for journaling and daily tracking
  - Automatic linking between notes
  - Smart organization and cross-referencing

- **Example Daily Workflow**:
  - Morning review and planning
  - Quick note capture during the day with command `exo zet "My Note Title"`
  - Evening reflection and summary
  - Automatic linking of new notes to daily entries facilitating note review and information retention.

- **Note Organization**:
  - Automatic file organization based on note type
  - Consistent naming conventions
  - Template-based note creation
  - Configurable directory structure

- **Editor Integration**:
  - Opens notes in your preferred editor

## Installation

```bash
go install github.com/a-kostevski/exo@latest
```

## Quick Start

1. Initialize Exo:
   ```bash
   exo init
   ```
   This creates:
   - Configuration file at `~/.config/exo/config.yaml`
   - Data directory at `~/.local/share/exo`
   - Default templates

   To use a custom data directory:
   ```bash
   # Set custom data directory
   export EXO_DATA_HOME=~/Documents/notes
   
   # Initialize with custom location
   exo init
   
   # Or force reinitialization if already initialized
   exo init --force
   ```

2. Start your day:
   ```bash
   # Open today's daily note for morning review
   exo day
   
   # Create notes throughout the day
   exo zet "Interesting Thought"  # automatically linked in daily note
   
   # Review your day's notes in the evening
   exo day
   ```

3. View your configuration:
   ```bash
   exo config
   ```

## Configuration

Exo uses a YAML configuration file and environment variables to customize its behavior. The configuration is loaded in the following order:

1. Environment variables
2. Configuration file
3. Default values

### Environment Variables

- `EXO_DATA_HOME`: Base directory for storing all notes and data files
  ```bash
  # Default: $XDG_DATA_HOME/exo or ~/.local/share/exo
  export EXO_DATA_HOME=~/Documents/notes
  ```

- `XDG_DATA_HOME`: Base directory for XDG data (used if `EXO_DATA_HOME` is not set)
  ```bash
  # Default: ~/.local/share
  export XDG_DATA_HOME=~/.local/share
  ```

- `XDG_CACHE_HOME`: Base directory for cache files (used for logs)
  ```bash
  # Default: ~/.cache
  export XDG_CACHE_HOME=~/.cache
  ```

### Configuration File

Example configuration:

```yaml
# Editor command to open notes
editor: "vim"  # default: vim

# Directory paths (all support ~ expansion)
data_home: "~/Documents/notes"     # Base directory for all notes
template_dir: "~/.config/exo/templates"  # Note templates
periodic_dir: "~/Documents/notes/periodic"  # Periodic notes
zettel_dir: "~/Documents/notes/zettel"     # Zettelkasten notes

# Logging configuration
log:
  level: "info"    # debug, info, warn, error
  format: "text"   # text or json
  output: "both"   # stdout, stderr, file, or both
  file: "~/.cache/exo/exo.log"  # log file path
```

## Commands

### Note Creation

```bash
# Create a Zettelkasten note
exo zet "My Note Title"

# Create a daily note (today)
exo day
```

### Configuration Management

```bash
# View all settings
exo config

# Get a setting
exo config get [key]

# Set a setting
exo config set [key] [value]

# Initialize configuration
exo init
exo init --force  # force reinitialization
```

### Template Management

```bash
# List available templates
exo templates
```

## Project Structure

```
.
├── cmd/                    # Command implementations
│   ├── config.go          # Configuration commands
│   ├── day.go             # Daily note command
│   ├── init.go            # Initialization command
│   ├── root.go            # Root command
│   ├── templates.go       # Template commands
│   └── zet.go             # Zettelkasten command
├── internal/
│   ├── config/            # Configuration management
│   ├── core/              # Core utilities
│   ├── fs/                # File system operations
│   ├── logger/            # Logging functionality
│   ├── note/              # Note management
│   │   ├── api/          # Note interfaces
│   │   ├── base/         # Base implementations
│   │   ├── factory/      # Note factories
│   │   └── types/        # Note type implementations
│   └── templates/         # Template management
└── main.go                # Application entry point
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Daily Workflow

Here's an example of how to use exo for your daily knowledge management:

### Morning Review (15min) ☀️
1. Run `exo day` to open today's daily note
2. Review yesterday's key points
3. Set learning objectives for today
4. Plan your active learning sessions

### During the Day 📚
1. When you encounter something interesting or have a thought:
   - Run `exo zet "Your Note Title"` to create a new Zettel note
   - The note will automatically be linked in today's daily note
   - Write down your thoughts, insights, or learnings

### Evening Summary (10min) 🌙
1. Open today's daily note with `exo day`
2. Review all the notes you created (linked at the bottom)
3. Summarize key learnings
4. Fill in your daily metrics
5. Plan tomorrow's focus areas

This workflow helps you:
- Start each day with intention
- Capture thoughts and learnings as they happen
- Build a connected knowledge base
- Reflect and improve daily
