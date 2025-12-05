# UniFi Protect View Switcher

Let's create a go cli application for switching between different camera views
in UniFi Protect.

## Features

- List available Viewports
- Switch between Viewports by name or ID
- List available camera views
- Switch to a specified camera view by name or ID
- Authenticate with UniFi Protect using an API token
- Read configuration from a file in XDG config directory
- Config will have fields for UniFi Protect URL and API token
- Support for command line arguments to override config file settings
- Error handling for network issues and invalid inputs
- Logging for debugging purposes
- Unit tests for core functionality
- Documentation on how to use the CLI tool
- Use existing Taskfile for build and run tasks
- Use goreleaser for packaging the application
- Follow best practices for Go project structure and code organization
- Ensure compatibility with the latest version of UniFi Protect API
- Provide examples of usage in the README file
- Implement a help command to display usage information
