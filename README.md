# UniFi Protect Control

An interactive TUI (Terminal User Interface) application for managing viewports
and controlling PTZ cameras in UniFi Protect.

## Features

- 🖥️ **Interactive TUI** - Easy-to-use terminal interface with keyboard
  navigation
- 📹 **Viewport Management** - Switch viewports between different liveviews
- 🎥 **PTZ Camera Control** - Move cameras to home position and presets (0-9)
- 🔐 **Secure Authentication** - API token-based authentication
- ⚙️ **Simple Configuration** - Single YAML file configuration
- 🚀 **CLI Mode Available** - Scriptable commands for automation
- ✅ **Well-tested** - Comprehensive unit test coverage

## Quick Start

### 1. Get an API Token

1. Log into your UniFi Protect web interface
2. Navigate to **Settings → Control Plane → Integrations → Your API Keys**
3. Generate an API token and copy it

### 2. Setup

```bash
# Clone and build
git clone https://github.com/methridge/protect.git
cd protect
task build
# or: go build -o bin/protect .

# Create config file
mkdir -p ~/.config/protect
cp config.yaml.example ~/.config/protect/config.yaml
```

### 3. Configure

Edit `~/.config/protect/config.yaml`:

```yaml
protect_url: https://192.168.1.100 # Your UniFi Protect URL
api_token: your-api-token-here # Your API token
log_level: none # Optional: none, debug, info, warn, error
```

### 4. Run

```bash
# Launch the TUI
./bin/protect

# Or use CLI commands for scripting
./bin/protect viewport list
./bin/protect camera goto "Front Door" -- -1
```

## Using the TUI

When you run `protect` without arguments, an interactive menu appears:

```text
UniFi Protect Control

Select an option:

  > Manage Viewports
    Control PTZ Cameras

↑/↓: navigate • enter: select • q: quit
```

### Navigation

- **↑/↓** or **j/k** - Navigate through options
- **Enter** or **Space** - Select current option
- **Esc** or **Backspace** - Go back to previous screen
- **q** or **Ctrl+C** - Quit application

### Viewport Management

1. Select "Manage Viewports" from main menu
2. Choose a viewport from the list
3. Select the liveview you want to switch to
4. See confirmation message

### PTZ Camera Control

1. Select "Control PTZ Cameras" from main menu
2. Choose a camera from the list
3. Select a preset:
   - **Home (-1)** - Return to home position
   - **Preset 0-9** - Move to saved preset positions
4. See confirmation message

## Installation

### Using Homebrew (macOS/Linux)

```bash
brew install methridge/tap/protect
```

### From Release

Download the latest release from the
[releases page](https://github.com/methridge/protect/releases) and extract it to
your PATH.

### From Source

```bash
git clone https://github.com/methridge/protect.git
cd protect
task build
```

The binary will be available in `bin/protect`.

## Configuration

The application reads configuration from a `config.yaml` file. It searches the
following locations in order:

1. `$XDG_CONFIG_HOME/protect/config.yaml` (if XDG_CONFIG_HOME is set)
2. `~/.config/protect/config.yaml` (recommended on Linux/macOS)
3. `~/Library/Application Support/protect/config.yaml` (macOS standard location)

**Note:** On macOS, `~/.config/protect/` is recommended for consistency across
platforms.

### Example Configuration

Create the file `~/.config/protect/config.yaml`:

```yaml
protect_url: https://protect.example.com
api_token: your-api-token-here
log_level: info
```

### Configuration Options

| Option        | Description                                    | Required | Default |
| ------------- | ---------------------------------------------- | -------- | ------- |
| `protect_url` | UniFi Protect server URL                       | Yes      | -       |
| `api_token`   | API authentication token                       | Yes      | -       |
| `log_level`   | Logging level (none, debug, info, warn, error) | No       | none    |

### Environment Variables

Configuration can be overridden using environment variables with the `PROTECT_`
prefix:

```bash
export PROTECT_PROTECT_URL=https://protect.example.com
export PROTECT_API_TOKEN=your-api-token
export PROTECT_LOG_LEVEL=debug
```

### Command-Line Flags

Configuration can also be overridden using command-line flags:

```bash
protect --url https://protect.example.com --token your-api-token --log-level debug
```

## Examples

### Interactive TUI Workflow

The TUI is the easiest way to interact with UniFi Protect:

1. **Launch:** Run `protect` from terminal
2. **Navigate:** Use arrow keys or `j`/`k` to browse
3. **Select:** Press `Enter` on any viewport, camera, or liveview
4. **Quick actions:**
   - Press `q` at any time to quit
   - Press `Esc` to go back to main menu
   - PTZ presets: Select camera → choose preset → press `Enter`

### Scripting with CLI

For automation and integration:

```bash
# Morning routine: Switch to driveway view
protect viewport switch Tower Driveway

# Security patrol: Cycle through camera presets
for i in 1 2 3; do
  protect camera goto "Front Door" $i
  sleep 10
done

# Integration with cron
# Switch to "All Cameras" at 10 PM
0 22 * * * /usr/local/bin/protect viewport switch Tower "All Cameras"
```

## CLI Reference

All TUI features are available via CLI for scripting and automation.

### Global Flags

```bash
-h, --help              Show help
-l, --log-level string  Log level (none, debug, info, warn, error)
-t, --token string      API token
-u, --url string        UniFi Protect URL
```

### Commands

```bash
# Viewports
protect viewport list                           # List all viewports
protect viewport switch <viewport> <liveview>   # Switch viewport to liveview

# Liveviews
protect liveview list                           # List all liveviews
protect liveview switch <viewport> <liveview>   # Alternative switch syntax

# PTZ Cameras
protect camera list                             # List all PTZ cameras
protect camera goto <camera> <preset>           # Move camera to preset (-1 to 9)
protect camera goto "Front Door" -- -1          # Home position (use -- before -1)
```

## Troubleshooting

### Authentication Issues

- Verify your API token in `~/.config/protect/config.yaml`
- Ensure your token has necessary permissions
- Try accessing the URL in a browser first

### Connection Errors

- Verify URL includes protocol (`https://`)
- Check firewall rules and network connectivity
- Use `--log-level debug` for detailed output

### Configuration Not Loading

- Check file exists at `~/.config/protect/config.yaml`
- Verify YAML syntax is correct
- Try overriding with `--url` and `--token` flags

## Development

### Prerequisites

- Go 1.21 or later
- Task (taskfile.dev)

### Building and Testing

```bash
task build           # Build binary
task run             # Run directly
task test            # Run all tests
task test-coverage   # Run tests with coverage
task lint            # Run linter
```

### Project Structure

```text
protect/
├── cmd/                    # Command definitions (root, viewport, liveview, camera)
├── internal/
│   ├── client/            # UniFi Protect API client
│   ├── config/            # Configuration management
│   ├── logger/            # Logging utilities
│   └── tui/               # Terminal UI (Bubble Tea)
├── main.go                # Application entry point
├── Taskfile.yaml          # Task automation
└── .goreleaser.yaml       # Release configuration
```

## API Compatibility

Tested with UniFi Protect v3.x on UDM Pro, UDM SE, UNVR, etc.

**Note:** UniFi Protect API is unofficial and may change without notice.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see the LICENSE file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Zap](https://github.com/uber-go/zap) - Structured logging

## Support

If you encounter any issues or have questions:

- Open an issue on [GitHub](https://github.com/methridge/protect/issues)
- Check existing issues for solutions
- Include debug logs when reporting issues (`--log-level debug`)

---

**Disclaimer:** This is an unofficial tool and is not affiliated with or
endorsed by Ubiquiti Inc.
