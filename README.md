<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Axon - AI-Driven Adventure Game](#axon---ai-driven-adventure-game)
  - [Features](#features)
    - [Core Gameplay](#core-gameplay)
    - [Technical Features](#technical-features)
  - [Architecture](#architecture)
    - [AI Model Strategy](#ai-model-strategy)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Building from Source](#building-from-source)
    - [Setting Up API Keys](#setting-up-api-keys)
  - [Usage](#usage)
    - [Starting Your Adventure](#starting-your-adventure)
    - [Game Commands](#game-commands)
    - [Navigation](#navigation)
  - [Configuration](#configuration)
    - [Save Files](#save-files)
  - [Development](#development)
    - [Technology Stack](#technology-stack)
    - [Design Principles](#design-principles)
    - [Running Tests](#running-tests)
    - [Test Coverage](#test-coverage)
    - [Project Structure](#project-structure)
  - [Contributing](#contributing)
  - [API Costs](#api-costs)
  - [Troubleshooting](#troubleshooting)
    - [Common Issues](#common-issues)
    - [Error Messages](#error-messages)
  - [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Axon - AI-Driven Adventure Game

Axon is a player-driven, prompt-based adventure game where AI models serve as your world builder, storyteller, and dynamic rule-setter. Built entirely as a Terminal User Interface (TUI) using Go and Bubble Tea, Axon runs on virtually any operating system with a terminal.

## Features

### Core Gameplay
- **Player-Driven World Generation**: Describe your ideal adventure setting and watch AI bring it to life
- **Dynamic Storytelling**: AI responds to your actions with contextual, engaging narratives
- **Emergent Gameplay**: Every decision shapes your unique adventure through AI-driven consequences
- **Interactive Inventory System**: Collect and manage items throughout your journey
- **Action Suggestions**: AI provides contextual suggestions to guide your adventure
- **Save/Load System**: Preserve your progress and return to your adventures anytime

### Technical Features
- **Cross-Platform**: Runs on Linux, macOS, Windows, and other Unix-like systems
- **Universal Terminal Compatibility**: Works with minimal terminals and UNIX System V
- **Automatic Terminal Detection**: Adapts interface based on terminal capabilities
- **Monochrome Design**: Clean, terminal-friendly black and white interface
- **Modular Architecture**: Clean separation of concerns for easy maintenance and extension
- **Multiple AI Providers**: Support for OpenRouter and Gemini APIs
- **Intelligent Model Selection**: Different AI models optimized for specific tasks
- **Comprehensive Testing**: High test coverage ensuring reliability

## Architecture

Axon follows a modular design pattern:

```
├── main.go                 # Entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── ai/                # AI client and model selection
│   ├── game/              # Game engine, state, and UI model
│   ├── storage/           # Save/load functionality
│   └── ui/                # Terminal UI styles and components
```

### AI Model Strategy

Axon uses different AI models optimized for specific tasks:
- **World Building**: Claude 3.5 Sonnet (for rich, detailed world creation)
- **Storytelling**: GPT-4o (for engaging narrative responses)
- **Rule Setting**: GPT-4o Mini (for game mechanics and suggestions)
- **Dialogue**: Claude 3 Haiku (for character interactions)

## Installation

### Prerequisites
- Go 1.23 or later
- Terminal with at least 80x24 character display
- OpenRouter API key (recommended) or Gemini API key

### Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd axon

# Install dependencies
go mod download

# Build the game
go build -o axon .

# Run the game
./axon
```

### Setting Up API Keys

Axon requires AI API access. Set up your credentials:

```bash
# For OpenRouter (recommended)
export OPENROUTER_API_KEY="your_openrouter_key_here"

# Or for Gemini
export GEMINI_API_KEY="your_gemini_key_here"

# Run the game
./axon
```

API keys can also be configured in the settings menu or by editing `~/.axon/config.json`.

## Usage

### Starting Your Adventure

1. **Launch Axon**: Run `./axon` in your terminal
2. **Choose "New Game"**: Select option 1 from the main menu
3. **Describe Your World**: Enter a creative description of your desired adventure setting
   - Example: "A cyberpunk city in 2077 where hackers fight against corporate oppression"
   - Example: "A medieval fantasy kingdom threatened by an ancient dragon"
   - Example: "A generation ship traveling to a distant star"
4. **Start Playing**: Type actions and watch your story unfold

### Game Commands

- **Any text**: Describe your action (e.g., "look around", "talk to the guard", "pick up the sword")
- **inventory** or **inv**: Check your items
- **stats**: View your character statistics
- **save [name]**: Save your game (e.g., "save my_adventure")
- **load [name]**: Load a saved game
- **help**: Display available commands
- **q** or **Ctrl+C**: Quit the game

### Navigation

- **↑/↓ Arrow Keys**: Scroll through game history
- **Enter**: Submit your input
- **Backspace**: Edit your current input

## Configuration

Axon creates a configuration file at `~/.axon/config.json` with the following structure:

```json
{
  "terminal": {
    "width": 80,
    "height": 24,
    "color_enabled": false
  },
  "ai": {
    "openrouter_api_key": "",
    "gemini_api_key": "",
    "default_model": "openai/gpt-4o-mini"
  },
  "game": {
    "history_limit": 1000,
    "save_dir": "/home/user/.axon/saves"
  }
}
```

### Save Files

Game saves are stored as JSON files in `~/.axon/saves/`. Each save contains:
- Complete world state and description
- Player character and inventory
- Full conversation history
- Game metadata and timestamps

### Terminal Compatibility

Axon automatically detects and adapts to your terminal type:

- **Modern Terminals**: Full features with color, mouse, and advanced formatting
- **Minimal Terminals** (`dumb`, CI/CD): Plain text mode with essential functionality
- **UNIX System V** (`vt100`, `vt220`): Compatible with legacy UNIX systems

Use the included detection utility:
```bash
# Check your terminal capabilities
./terminal-info

# Test with different terminal types
TERM=dumb ./axon          # Minimal mode
TERM=vt100 ./axon         # System V mode
```

See [TERMINAL_COMPATIBILITY.md](TERMINAL_COMPATIBILITY.md) for detailed information.

## Development

### Technology Stack

- **Language**: Go (with CGO=0 for maximum compatibility)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **AI APIs**: OpenRouter and Gemini
- **Storage**: JSON files
- **Testing**: Go's built-in testing framework

### Design Principles

1. **Extreme Minimalism**: Simple, clean interface focused on the story
2. **Maximum Modularity**: Clear separation between game logic, AI, storage, and UI
3. **Cross-Platform Compatibility**: Runs anywhere Go can compile
4. **Player Agency**: Players drive the narrative through their descriptions and actions

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific module tests
go test ./internal/game/
```

### Test Coverage

- **UI Module**: 100% coverage
- **Storage Module**: 82.4% coverage  
- **AI Module**: 75.6% coverage
- **Game Module**: 55.9% coverage
- **Config Module**: 47.4% coverage

### Project Structure

```
axon/
├── main.go                     # Application entry point
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── README.md                   # This file
└── internal/                   # Private application code
    ├── config/                 # Configuration management
    │   ├── config.go          # Config structures and loading
    │   └── config_test.go     # Configuration tests
    ├── ai/                     # AI client implementation
    │   ├── client.go          # API clients for OpenRouter/Gemini
    │   └── client_test.go     # AI client tests
    ├── game/                   # Core game logic
    │   ├── engine.go          # Game engine and AI integration
    │   ├── model.go           # Bubble Tea model (main UI)
    │   ├── state.go           # Game state management
    │   ├── engine_test.go     # Engine tests
    │   ├── model_test.go      # UI model tests
    │   └── state_test.go      # State management tests
    ├── storage/                # Save/load functionality
    │   ├── storage.go         # File-based storage implementation
    │   └── storage_test.go    # Storage tests
    └── ui/                     # User interface components
        ├── styles.go          # Monochrome styling definitions
        └── styles_test.go     # UI styling tests
```

## Contributing

Axon follows conventional commit messages and maintains high test coverage. When contributing:

1. Write tests for new functionality
2. Maintain the minimalist design philosophy
3. Ensure cross-platform compatibility
4. Follow Go best practices and idioms
5. Update documentation for user-facing changes

## API Costs

Axon is designed to be cost-effective:
- Uses efficient models appropriate for each task
- Limits context length to control token usage
- Provides fallback responses when API calls fail
- Allows gameplay without constant API calls through local commands

## Troubleshooting

### Common Issues

**Game won't start**: Ensure your terminal is at least 80x24 characters

**AI responses fail**: Check your API key configuration and internet connection

**Save files missing**: Verify write permissions to `~/.axon/saves/`

**Display issues**: Try adjusting terminal size or font settings

### Error Messages

Axon provides clear error messages for common issues:
- Missing or invalid API keys
- Network connectivity problems
- File permission errors
- Invalid save file formats

## License

This project follows an open-source approach prioritizing player experience and developer accessibility.

---

**Start your adventure today!** Create worlds limited only by your imagination, and let AI bring your stories to life in the terminal.

