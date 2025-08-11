# Style Server

A simple MCP (Model Context Protocol) server that applies Google style guide best practices for Python and Go code.

## Features

- **Python Style Formatting**:
  - Fix indentation to 4 spaces
  - Add proper spacing around operators
  - Convert function names to snake_case
  - Check for line length violations (80 chars)
  - Suggest adding docstrings

- **Go Style Formatting**:
  - Fix spacing around keywords
  - Check for proper comments on exported functions
  - Suggest running gofmt

## Usage

The server provides a `format-code` tool that takes:
- `language`: "python" or "go"
- `code`: The code to format
- `filePath` (optional): File path for context

## Installation

```bash
go build -o style-server ./cmd/style-server/
```

## Example Usage

The server can be used as an MCP server in VS Code or other MCP-compatible clients.

### Configuration

Add to your MCP client configuration:

```json
{
  "servers": {
    "style-server": {
      "type": "stdio",
      "command": "/path/to/style-server"
    }
  }
}
```

## Google Style Guides

This server applies rules based on:
- [Google Python Style Guide](https://google.github.io/styleguide/pyguide.html)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)

## Development

The server is built using the [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk).
