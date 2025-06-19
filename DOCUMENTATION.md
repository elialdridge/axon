# Axon Game - Complete Documentation

## Table of Contents

1. [Overview](#overview)
2. [System Requirements](#system-requirements)
3. [Installation Guide](#installation-guide)
4. [Configuration](#configuration)
5. [Terminal Compatibility](#terminal-compatibility)
6. [Logging and Debugging](#logging-and-debugging)
7. [Game Mechanics](#game-mechanics)
8. [AI Integration](#ai-integration)
9. [User Interface](#user-interface)
10. [Command Reference](#command-reference)
11. [Project Architecture](#project-architecture)
12. [API Documentation](#api-documentation)
13. [Development Guide](#development-guide)
14. [Testing](#testing)
15. [Deployment](#deployment)
16. [Troubleshooting](#troubleshooting)
17. [Use Cases](#use-cases)
18. [Performance Optimization](#performance-optimization)
19. [Security Considerations](#security-considerations)
20. [Extensibility](#extensibility)
21. [Examples](#examples)

## Overview

Axon is a revolutionary AI-driven adventure game that transforms natural language descriptions into immersive interactive experiences. Built with extreme minimalism and maximum modularity in mind, Axon runs entirely in the terminal using Go and Bubble Tea framework.

### Core Philosophy

- **Player Agency**: Every aspect of the game world is shaped by player input
- **AI Collaboration**: Multiple AI models work together to create coherent experiences
- **Terminal Native**: Designed specifically for terminal environments
- **Cross-Platform**: Runs on any system with a terminal and Go runtime
- **Minimalist Design**: Clean, distraction-free interface focused on storytelling

### Key Features

- **Dynamic World Generation**: AI creates entire worlds from simple descriptions
- **Contextual Storytelling**: AI maintains narrative consistency across sessions
- **Intelligent Action Processing**: AI interprets free-form player actions
- **Multi-Model Architecture**: Different AI models optimized for specific tasks
- **Persistent State**: Complete game state persistence with JSON serialization
- **Modular Design**: Clean architecture enabling easy modification and extension

## Terminal Compatibility

Axon features comprehensive terminal detection and compatibility support for maximum portability across systems:

### Automatic Terminal Detection

Axon automatically detects terminal capabilities using multiple methods:

1. **TERM Environment Variable Analysis**: Identifies terminal type and capabilities
2. **Dynamic Size Detection**: Uses ioctl, environment variables, tput, and stty
3. **Feature Detection**: Determines color, mouse, alt-screen, and UTF-8 support
4. **Compatibility Mode Detection**: Identifies minimal and System V terminals

### Supported Terminal Types

#### Modern Terminals (Full Features)
- **xterm variants**: xterm, xterm-256color, xterm-color
- **Terminal multiplexers**: screen, tmux
- **Modern emulators**: gnome-terminal, konsole, Terminal.app
- **Features**: Color, mouse, alt-screen, UTF-8, bold, underline

#### Minimal Terminals (Basic Features)
- **dumb**: CI/CD environments, simple scripts
- **unknown**: Unrecognized terminal types
- **vt52/vt100/vt102**: Very basic terminals
- **Features**: Plain text only, no special formatting

#### UNIX System V Compatible
- **vt220**: Advanced VT terminal
- **ansi**: ANSI.SYS compatible
- **cons25**: FreeBSD console
- **Features**: Basic formatting, limited color support

### Terminal Detection Utility

Use the included terminal detection utility to test compatibility:

```bash
# Check current terminal capabilities
./terminal-info

# Test with different terminal types
TERM=dumb ./axon          # Minimal mode
TERM=vt100 ./axon         # System V mode
TERM=xterm-256color ./axon # Full featured mode

# Force specific mode
AXON_FORCE_MINIMAL=true ./axon
```

### Adaptive Interface Features

#### Dynamic Layout Adjustment
- **Size Detection**: Multiple fallback methods for terminal dimensions
- **Safe Sizing**: Bounds checking with minimum 40x10, maximum 200x60
- **Runtime Adaptation**: Responds to terminal resize events
- **Layout Optimization**: Adjusts panel ratios based on available space

#### Compatibility Options
- **No Mouse**: Automatically disabled for minimal/System V terminals
- **No Alt Screen**: Plain mode for limited terminals
- **Simple Rendering**: Reduced complexity for compatibility
- **Text-Only Mode**: Stripped formatting for maximum compatibility

## Logging and Debugging

Axon includes a comprehensive logging system for debugging and monitoring:

### Logger Architecture

```go
// Logger components
type Logger struct {
    debugLogger *log.Logger  // Detailed debug information
    infoLogger  *log.Logger  // General information
    errorLogger *log.Logger  // Error conditions
    logFile     *os.File     // Log file handle
}
```

### Log Levels and Output

#### Debug Level
- Terminal detection details
- AI request/response logging
- Game state changes
- Configuration loading
- Performance metrics

#### Info Level
- Application startup/shutdown
- World creation events
- Player actions
- Save/load operations
- Terminal mode changes

#### Error Level
- API failures
- File I/O errors
- Configuration problems
- Network issues
- Validation failures

### Debug Mode Usage

```bash
# Enable debug logging
export AXON_DEBUG=true
./axon

# View real-time logs
tail -f axon_debug.log

# Filter specific components
grep "Terminal" axon_debug.log
grep "AI" axon_debug.log
grep "World Creation" axon_debug.log
```

### Specialized Logging Functions

```go
// Specialized logging for different components
logger.LogRequest(aiRequest)        // AI request details
logger.LogResponse(aiResponse)      // AI response analysis
logger.LogGameState(gameState)      // Complete game state
logger.LogWorldCreation(step, data) // World generation process
```

### Log File Management

- **Location**: `axon_debug.log` in current directory
- **Format**: Timestamped entries with file/line information
- **Rotation**: Manual cleanup (automatic rotation planned)
- **Permissions**: 0666 for user/group read/write access

## System Requirements

### Minimum Requirements

- **Operating System**: Linux, macOS, Windows, FreeBSD, or any Unix-like system
- **Go Version**: 1.23 or later
- **Terminal**: Any terminal emulator supporting ANSI escape sequences
- **Screen Size**: Minimum 80x24 characters
- **Memory**: 64MB RAM
- **Storage**: 10MB free space
- **Network**: Internet connection for AI API calls

### Recommended Requirements

- **Terminal Size**: 120x40 characters or larger
- **Memory**: 128MB RAM for better performance
- **Storage**: 100MB for multiple save files
- **Network**: Stable broadband connection for optimal AI response times

### Supported Platforms

- **Linux**: All major distributions (Ubuntu, Debian, CentOS, Arch, etc.)
- **macOS**: 10.15 Catalina and later
- **Windows**: Windows 10/11 with proper terminal (Windows Terminal recommended)
- **FreeBSD**: 12.0 and later
- **OpenBSD**: 6.8 and later
- **NetBSD**: 9.0 and later

## Installation Guide

### Method 1: Build from Source (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd axon

# Verify Go installation
go version

# Download dependencies
go mod download

# Build the binary
go build -o axon .

# Make executable (Unix systems)
chmod +x axon

# Run the game
./axon
```

### Method 2: Cross-Compilation

```bash
# For different architectures
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o axon-linux-amd64

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o axon-macos-arm64

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o axon-windows-amd64.exe

# Linux ARM64 (Raspberry Pi)
GOOS=linux GOARCH=arm64 go build -o axon-linux-arm64
```

### Method 3: Development Installation

```bash
# Install with go install for development
go install .

# Run from anywhere (if GOPATH/bin is in PATH)
axon
```

### Post-Installation Setup

1. **Create configuration directory**:
   ```bash
   mkdir -p ~/.axon/saves
   ```

2. **Set up API keys**:
   ```bash
   # OpenRouter (recommended)
   export OPENROUTER_API_KEY="your_key_here"
   
   # Or Gemini
   export GEMINI_API_KEY="your_key_here"
   ```

3. **Test installation**:
   ```bash
   ./axon --version  # (if implemented)
   ./axon
   ```

## Configuration

### Configuration File Location

- **Linux/macOS**: `~/.axon/config.json`
- **Windows**: `%USERPROFILE%\.axon\config.json`

### Configuration Structure

```json
{
  "terminal": {
    "width": 80,
    "height": 24,
    "color_enabled": false,
    "unicode_support": true,
    "refresh_rate": 60
  },
  "ai": {
    "openrouter_api_key": "",
    "gemini_api_key": "",
    "default_model": "openai/gpt-4o-mini",
    "max_tokens": 1000,
    "temperature": 0.7,
    "timeout_seconds": 30,
    "retry_attempts": 3,
    "fallback_enabled": true
  },
  "game": {
    "history_limit": 1000,
    "save_dir": "/home/user/.axon/saves",
    "auto_save_interval": 300,
    "max_saves": 50,
    "compression_enabled": false
  },
  "ui": {
    "scroll_speed": 3,
    "animation_enabled": false,
    "borders_enabled": true,
    "timestamps_shown": false,
    "debug_mode": false
  },
  "logging": {
    "level": "info",
    "file_path": "/home/user/.axon/axon.log",
    "max_file_size_mb": 10,
    "max_files": 5
  }
}
```

### Configuration Options Explained

#### Terminal Settings
- `width/height`: Default terminal dimensions
- `color_enabled`: Enable color support (currently forced to false for monochrome design)
- `unicode_support`: Enable Unicode characters
- `refresh_rate`: UI refresh rate in Hz

#### AI Settings
- `openrouter_api_key`: OpenRouter API authentication key
- `gemini_api_key`: Google Gemini API authentication key
- `default_model`: Fallback model when specific model selection fails
- `max_tokens`: Maximum tokens per AI request
- `temperature`: AI creativity level (0.0-1.0)
- `timeout_seconds`: Request timeout duration
- `retry_attempts`: Number of retry attempts for failed requests
- `fallback_enabled`: Use local fallback responses when AI fails

#### Game Settings
- `history_limit`: Maximum number of history entries to keep
- `save_dir`: Directory for save files
- `auto_save_interval`: Automatic save interval in seconds (0 = disabled)
- `max_saves`: Maximum number of save files to retain
- `compression_enabled`: Enable save file compression

#### UI Settings
- `scroll_speed`: Lines to scroll per arrow key press
- `animation_enabled`: Enable UI animations
- `borders_enabled`: Show UI borders
- `timestamps_shown`: Display timestamps in history
- `debug_mode`: Enable debug information display

#### Logging Settings
- `level`: Logging level (debug, info, warn, error)
- `file_path`: Log file location
- `max_file_size_mb`: Maximum log file size before rotation
- `max_files`: Number of log files to retain

### Environment Variables

Axon recognizes the following environment variables:

```bash
# API Keys
OPENROUTER_API_KEY="your_openrouter_key"
GEMINI_API_KEY="your_gemini_key"

# Configuration
AXON_CONFIG_PATH="/custom/path/to/config.json"
AXON_SAVE_DIR="/custom/save/directory"
AXON_LOG_LEVEL="debug"

# Terminal
AXON_WIDTH="120"
AXON_HEIGHT="40"

# Development
AXON_DEBUG="true"
AXON_OFFLINE="true"  # Disable AI calls for testing
```

## Game Mechanics

### Game Flow

1. **World Initialization**
   - Player provides world description
   - AI generates world details, rules, and starting location
   - Game state is initialized with world data

2. **Gameplay Loop**
   - Player inputs action
   - AI processes action in context
   - Game state is updated
   - AI generates narrative response
   - Suggestions for next actions are generated

3. **State Persistence**
   - All game state changes are tracked
   - Save files contain complete world and history
   - Resume functionality maintains exact game state

### World Generation

#### Initial World Creation

```
Player Input: "A cyberpunk city in 2077"

AI Processing:
1. World Builder AI analyzes prompt
2. Generates:
   - World name: "Neo-Tokyo 2077"
   - Setting description
   - Core rules and mechanics
   - Starting location
   - Initial scenario
```

#### World Elements

- **Name**: Generated world identifier
- **Description**: Detailed world background
- **Setting**: Genre and time period
- **Rules**: Game-specific mechanics and limitations
- **Locations**: Dictionary of discoverable places
- **Current Location**: Player's current position

### Player Character System

#### Character Properties

```json
{
  "name": "Player Character",
  "description": "A determined explorer",
  "inventory": [
    {
      "name": "Neural Interface",
      "description": "Allows direct connection to the net",
      "quantity": 1
    }
  ],
  "stats": {
    "health": 100,
    "energy": 80,
    "credits": 500
  },
  "status": "healthy"
}
```

#### Character Development

- **Dynamic Stats**: AI can modify character stats based on actions
- **Inventory Management**: Items gained/lost through gameplay
- **Status Effects**: Temporary conditions affecting gameplay
- **Character Evolution**: Description and abilities change with story

### Action Processing

#### Action Types

1. **Movement Actions**
   - "go north", "enter the building", "climb the stairs"
   - Update current location
   - May trigger location-specific events

2. **Interaction Actions**
   - "talk to the guard", "examine the console", "pick up the key"
   - Generate dialogue or descriptions
   - May modify inventory or stats

3. **Combat Actions**
   - "attack the robot", "defend", "use plasma rifle"
   - Processed through combat system
   - Updates health and status

4. **Social Actions**
   - "negotiate", "intimidate", "persuade"
   - Success based on character stats and context
   - Affects NPC relationships

5. **System Actions**
   - "inventory", "stats", "save game", "help"
   - Direct system commands
   - Bypass AI processing

#### Context Awareness

The AI maintains context through:

- **Recent History**: Last 10 interactions for immediate context
- **World State**: Current location, time, conditions
- **Character State**: Health, inventory, status effects
- **Relationship State**: NPC attitudes and faction standings
- **Quest State**: Active objectives and progress

### Save System

#### Save File Format

```json
{
  "world": {
    "name": "Neo-Tokyo 2077",
    "description": "A sprawling cyberpunk metropolis...",
    "setting": "Cyberpunk",
    "rules": ["Technology is ubiquitous", "Corporations rule"],
    "locations": {
      "downtown": "Neon-lit streets buzz with activity..."
    },
    "current_location": "downtown"
  },
  "player": {
    "name": "Jake Morrison",
    "description": "A street-smart hacker",
    "inventory": [...],
    "stats": {...},
    "status": "healthy"
  },
  "history": [
    {
      "type": "narrator",
      "content": "You find yourself in the heart of downtown...",
      "timestamp": "2025-06-16T13:00:00Z",
      "turn": 1
    }
  ],
  "turn": 15,
  "created_at": "2025-06-16T12:00:00Z",
  "updated_at": "2025-06-16T13:15:30Z"
}
```

#### Save Operations

- **Manual Save**: `save [name]` command
- **Auto Save**: Periodic automatic saves (configurable)
- **Quick Save**: Overwrite last save
- **Multiple Saves**: Support for multiple save slots
- **Save Validation**: Integrity checking on load

## AI Integration

### Multi-Model Architecture

Axon uses different AI models optimized for specific tasks:

#### Model Selection Strategy

```go
func (c *Client) GetBestModel(task string) string {
    switch task {
    case "world_building":
        return "anthropic/claude-3.5-sonnet"  // Rich world creation
    case "storytelling":
        return "openai/gpt-4o"               // Engaging narratives
    case "rule_setting":
        return "openai/gpt-4o-mini"          // Game mechanics
    case "dialogue":
        return "anthropic/claude-3-haiku"    // Character interactions
    default:
        return "openai/gpt-4o-mini"          // Fallback
    }
}
```

#### Task-Specific Processing

1. **World Building** (Claude 3.5 Sonnet)
   - Complex world generation
   - Detailed location descriptions
   - Cultural and historical context
   - Rule system creation

2. **Storytelling** (GPT-4o)
   - Main narrative responses
   - Plot development
   - Character development
   - Emotional depth

3. **Rule Setting** (GPT-4o Mini)
   - Action suggestions
   - Game mechanic enforcement
   - Balance and fairness
   - Quick responses

4. **Dialogue** (Claude 3 Haiku)
   - NPC conversations
   - Character voice consistency
   - Social interactions
   - Fast response times

### API Integration

#### OpenRouter Integration

```go
type OpenRouterRequest struct {
    Model     string    `json:"model"`
    Messages  []Message `json:"messages"`
    MaxTokens int       `json:"max_tokens"`
    Temperature float64 `json:"temperature,omitempty"`
}

type Message struct {
    Role    string `json:"role"    // "system", "user", "assistant"
    Content string `json:"content"`
}
```

#### Request Construction

```go
// Build context from game state
messages := []Message{
    {
        Role: "system",
        Content: "You are a Game Master for a cyberpunk adventure...",
    },
    {
        Role: "user", 
        Content: "Player action: examine the terminal",
    },
}
```

#### Error Handling

- **Network Failures**: Retry with exponential backoff
- **API Errors**: Fall back to alternative models
- **Rate Limiting**: Queue requests and throttle
- **Invalid Responses**: Use fallback text
- **Timeout**: Cancel request and provide default response

### Context Management

#### Context Window Optimization

```go
func buildContext(state *GameState, maxTokens int) []string {
    context := []string{
        // Core system prompt
        systemPrompt,
        // World description 
        fmt.Sprintf("World: %s - %s", state.World.Name, state.World.Description),
        // Current location
        fmt.Sprintf("Location: %s", state.World.CurrentLocation),
        // Character state
        fmt.Sprintf("Player: %s", state.Player.Description),
    }
    
    // Add recent history within token limit
    recentHistory := getRecentHistoryWithinLimit(state, maxTokens)
    context = append(context, recentHistory...)
    
    return context
}
```

#### Context Prioritization

1. **System Prompts**: Always included
2. **World State**: Current location and rules
3. **Character State**: Health, inventory, status
4. **Recent History**: Last few interactions
5. **Long-term Memory**: Important events (if space allows)

## User Interface

### TUI Architecture

Axon uses Bubble Tea for its Terminal User Interface:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              GAME HISTORY                                  │
│                                                                             │
│  > look around                                                              │
│  You find yourself in the heart of Neo-Tokyo's downtown district. Neon     │
│  signs flicker in multiple languages, casting colorful shadows on the      │
│  rain-slicked streets. Towering corporate arcologies stretch toward the    │
│  perpetually overcast sky.                                                 │
│                                                                             │
│  > examine the alley                                                       │
│  A narrow alley between two massive buildings catches your attention.      │
│  Strange blue light emanates from within, and you hear the faint hum of    │
│  electronics.                                                              │
│                                                                             │
│  [System] Suggestions: Enter the alley, Scan for threats, Check inventory  │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│ > enter the alley_                                                         │
└─────────────────────────────────────────────────────────────────────────────┘
```

### UI Components

#### History Panel

- **Purpose**: Display game narrative and player actions
- **Scrolling**: Arrow key navigation through history
- **Formatting**: Different styles for player/narrator/system text
- **Auto-scroll**: Automatically scrolls to show latest content

#### Input Panel

- **Purpose**: Text input for player actions
- **Features**: 
  - Real-time typing
  - Backspace support
  - Enter to submit
  - Input validation

#### Status Indicators

- **Loading**: "Processing..." during AI requests
- **Errors**: Red text for error messages
- **Suggestions**: Contextual action hints
- **System Messages**: Configuration and help text

### Styling System

```go
type Styles struct {
    Base          lipgloss.Style  // Base container
    HistoryPanel  lipgloss.Style  // History display area
    InputPanel    lipgloss.Style  // Input area
    PlayerText    lipgloss.Style  // Player actions (bold)
    NarratorText  lipgloss.Style  // AI responses (normal)
    SystemText    lipgloss.Style  // System messages (italic)
    InventoryText lipgloss.Style  // Inventory items (faint)
    Border        lipgloss.Style  // UI borders
    Prompt        lipgloss.Style  // Input prompt (bold)
    Scrollbar     lipgloss.Style  // Scroll indicators
}
```

#### Monochrome Design

- **Foreground**: White (#FFFFFF)
- **Background**: Black (#000000)
- **Styling**: Bold, italic, faint for differentiation
- **Borders**: ASCII characters for maximum compatibility

### Responsive Design

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        // Recalculate panel sizes
        m.recalculateLayout()
        return m, nil
    }
}
```

#### Layout Adaptation

- **Minimum Size**: 80x24 characters
- **Dynamic Sizing**: Adapts to terminal resize
- **Panel Ratios**: Input panel fixed height, history panel flexible
- **Text Wrapping**: Automatic text wrapping for long lines

## Project Architecture

Axon follows a modular, layered architecture designed for maximum maintainability and extensibility:

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        User Interface (TUI)                    │
│                      Bubble Tea + Lip Gloss                    │
├─────────────────────────────────────────────────────────────────┤
│                         Game Engine                            │
│          Action Processing, State Management, Rules            │
├─────────────────────────────────────────────────────────────────┤
│                       AI Integration                           │
│              OpenRouter, Gemini, Model Selection              │
├─────────────────────────────────────────────────────────────────┤
│                    Storage & Persistence                       │
│                  JSON Files, Save/Load System                 │
├─────────────────────────────────────────────────────────────────┤
│                  Configuration & Logging                       │
│              Terminal Detection, Settings, Debug              │
└─────────────────────────────────────────────────────────────────┘
```

### Current Project Structure

```
axon/
├── main.go                           # Application entry point
├── go.mod                             # Go module definition
├── go.sum                             # Dependency lock file
├── README.md                          # Basic project overview
├── DOCUMENTATION.md                   # Comprehensive documentation
├── TERMINAL_COMPATIBILITY.md          # Terminal compatibility guide
├── TEST_REPORT.md                     # Test results and coverage
├── CONTRIBUTING.md                    # Contribution guidelines
├── LICENSE                            # Project license
├── .gitignore                         # Git ignore rules
├── Makefile                           # Build automation
├── terminal-info*                     # Terminal detection utility
├── axon_debug.log                     # Runtime debug log
├── debug_world_creation.go            # Debug utilities
├── cmd/
│   └── test/
│       └── main.go                    # Integration test utility
└── internal/                          # Private application code
    ├── config/
    │   ├── config.go                  # Configuration management
    │   └── config_test.go             # Configuration tests
    ├── ai/
    │   ├── client.go                  # AI client implementation
    │   └── client_test.go             # AI client tests
    ├── game/
    │   ├── engine.go                  # Core game engine
    │   ├── engine_test.go             # Engine tests
    │   ├── integration_test.go        # Integration tests
    │   ├── model.go                   # Bubble Tea UI model
    │   ├── model_test.go              # UI model tests
    │   ├── state.go                   # Game state management
    │   └── state_test.go              # State management tests
    ├── storage/
    │   ├── storage.go                 # Persistence interface
    │   └── storage_test.go            # Storage tests
    ├── ui/
    │   ├── styles.go                  # UI styling and themes
    │   └── styles_test.go             # Style tests
    ├── terminal/
    │   ├── terminal.go                # Terminal detection
    │   └── terminal_test.go           # Terminal tests
    └── logger/
        └── logger.go                  # Logging system
```

### Module Dependencies

#### Core Dependencies
- **Bubble Tea**: TUI framework for interactive terminal applications
- **Lip Gloss**: Styling library for terminal UI components

#### Architecture Principles

1. **Single Responsibility**: Each package has a clear, focused purpose
2. **Dependency Injection**: Dependencies passed explicitly, not hardcoded
3. **Interface Segregation**: Small, focused interfaces for flexibility
4. **Testability**: All components designed for easy unit testing
5. **Platform Independence**: No platform-specific dependencies

### Component Interaction Flow

```
User Input → UI Model → Game Engine → AI Client → Response
     ↑                     ↓              ↓
     └─── UI Update ←── State Update ←── Storage
```

#### Detailed Flow

1. **User Input**: Player types action in terminal
2. **UI Processing**: Bubble Tea model captures and validates input
3. **Engine Processing**: Game engine interprets action in context
4. **AI Generation**: Appropriate AI model generates response
5. **State Update**: Game state updated with new information
6. **Storage**: State changes persisted to disk
7. **UI Update**: Interface refreshed with new content

### Key Architectural Features

#### Terminal Detection System

```go
type TerminalInfo struct {
    Width            int    // Terminal dimensions
    Height           int
    ColorSupport     bool   // Feature capabilities
    MouseSupport     bool
    AltScreenSupport bool
    TermType         string // Terminal identification
    IsMinimal        bool   // Compatibility modes
    IsSystemV        bool
}
```

- **Multi-method Detection**: ioctl, environment, tput, stty
- **Capability Assessment**: Color, mouse, alt-screen support
- **Compatibility Modes**: Minimal and System V detection
- **Safe Defaults**: Fallback to 80x24 with basic features

#### AI Model Strategy

```go
func (c *Client) GetBestModel(task string) string {
    switch task {
    case "world_building":
        return "anthropic/claude-3.5-sonnet"  // Rich detail
    case "storytelling":
        return "openai/gpt-4o"               // Narrative
    case "rule_setting":
        return "openai/gpt-4o-mini"          // Quick rules
    case "dialogue":
        return "anthropic/claude-3-haiku"    // Conversations
    }
}
```

- **Task-Specific Selection**: Different models for different purposes
- **Performance Optimization**: Faster models for frequent operations
- **Quality Balance**: High-quality models for important content
- **Fallback Strategy**: Graceful degradation when models unavailable

#### State Management

```go
type GameState struct {
    World     *World         `json:"world"`      // World state
    Player    *Player        `json:"player"`     // Character state
    History   []HistoryEntry `json:"history"`    // Interaction log
    Turn      int            `json:"turn"`       // Game progression
    CreatedAt time.Time      `json:"created_at"` // Metadata
    UpdatedAt time.Time      `json:"updated_at"`
}
```

- **Immutable Updates**: State changes create new instances
- **Complete Serialization**: Full state captured in save files
- **Event Sourcing**: History preserves all interactions
- **Version Tracking**: Timestamps for change tracking

#### Logging Architecture

```go
// Structured logging with multiple levels
type Logger struct {
    debugLogger *log.Logger  // Detailed debugging
    infoLogger  *log.Logger  // General information
    errorLogger *log.Logger  // Error conditions
    logFile     *os.File     // File handle
}
```

- **Multi-level Logging**: Debug, Info, Error levels
- **File-based Output**: Persistent log files for analysis
- **Component-specific**: Specialized logging for different systems
- **Performance Monitoring**: Runtime metrics and profiling

### Design Patterns Used

#### Model-View-Update (Elm Architecture)
- **Model**: Game state and UI state
- **View**: Terminal rendering functions
- **Update**: State transition functions
- **Commands**: Async operations (AI calls)

#### Strategy Pattern
- **AI Provider Selection**: Different APIs for different needs
- **Terminal Rendering**: Adaptive output based on capabilities
- **Storage Backends**: Pluggable persistence mechanisms

#### Observer Pattern
- **State Changes**: UI updates when state changes
- **Configuration**: Settings changes propagate automatically
- **Logging**: Events trigger appropriate log entries

#### Factory Pattern
- **AI Client Creation**: Based on available credentials
- **Storage Creation**: Based on configuration
- **UI Theme Creation**: Based on terminal capabilities

### Extension Points

#### Plugin Architecture (Planned)
```go
type Plugin interface {
    Name() string
    ProcessAction(action string, state *GameState) (*ActionResult, error)
}
```

#### Custom AI Providers
```go
type AIProvider interface {
    Generate(ctx context.Context, req *Request) (*Response, error)
    GetModels() []string
}
```

#### Storage Backends
```go
type Storage interface {
    SaveGame(name string, state interface{}) error
    LoadGame(name string, state interface{}) error
}
```

### Performance Characteristics

- **Memory Usage**: ~50MB typical, ~100MB with large histories
- **Startup Time**: ~500ms including terminal detection
- **Response Time**: 1-5 seconds for AI calls, instant for local commands
- **File I/O**: JSON serialization, ~1MB typical save files
- **Network**: HTTPS requests to AI providers, configurable timeouts

## Command Reference

### Game Commands

#### Movement Commands
```
go [direction]     - Move in specified direction
enter [location]   - Enter a specific location
leave             - Exit current location
climb [object]    - Climb up/down objects
swim [direction]  - Swimming actions
fly [direction]   - Flying (if applicable)
```

#### Interaction Commands
```
examine [object]   - Look closely at something
look [around]      - General observation
talk to [person]   - Initiate conversation
use [item]         - Use an inventory item
take [object]      - Pick up an object
drop [item]        - Drop an inventory item
open [container]   - Open doors, chests, etc.
close [object]     - Close opened objects
```

#### Combat Commands
```
attack [target]    - Physical attack
shoot [target]     - Ranged attack
defend            - Defensive stance
flee              - Attempt to escape
hide              - Try to hide
sneak             - Move stealthily
```

#### Social Commands
```
persuade [person]  - Attempt persuasion
intimidate [person] - Use intimidation
negotiate         - Start negotiations
listen            - Eavesdrop on conversations
wait              - Wait and observe
```

#### System Commands
```
inventory / inv   - Show inventory
stats             - Display character stats
save [name]       - Save game
load [name]       - Load saved game
help              - Show help information
quit / q          - Exit game
settings          - Open settings menu
```

### Advanced Command Syntax

#### Multi-word Actions
```
"pick up the golden key"           - Use quotes for complex actions
"talk to the bartender about jobs" - Detailed interaction
"use the neural interface on the terminal" - Item combinations
```

#### Conditional Actions
```
"if I have a lockpick, pick the lock" - Conditional execution
"carefully examine the trap"           - Modifier words
"quickly run to the exit"              - Action modifiers
```

#### Meta Commands
```
/retry            - Retry last AI request
/context          - Show current context summary
/debug            - Toggle debug information
/clear            - Clear screen
/history [n]      - Show last n history entries
```

## API Documentation

### Core Interfaces

#### Game Engine Interface

```go
type Engine interface {
    InitializeWorld(state *GameState, seedPrompt string) error
    ProcessPlayerAction(state *GameState, action string) error
    GenerateActionSuggestions(state *GameState) ([]string, error)
}
```

#### AI Client Interface

```go
type AIClient interface {
    Generate(req Request) (*Response, error)
    GetBestModel(task string) string
}

type Request struct {
    Prompt    string   `json:"prompt"`
    Model     string   `json:"model"`
    MaxTokens int      `json:"max_tokens"`
    Context   []string `json:"context"`
}

type Response struct {
    Text  string `json:"text"`
    Error error  `json:"error,omitempty"`
}
```

#### Storage Interface

```go
type Storage interface {
    SaveGame(name string, state interface{}) error
    LoadGame(name string, state interface{}) error
    ListSaves() ([]string, error)
    DeleteSave(name string) error
}
```

### Internal APIs

#### State Management

```go
// GameState represents complete game state
type GameState struct {
    World     *World         `json:"world"`
    Player    *Player        `json:"player"`
    History   []HistoryEntry `json:"history"`
    Turn      int            `json:"turn"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
}

// Methods
func (gs *GameState) AddHistoryEntry(entryType, content string)
func (gs *GameState) GetRecentHistory(n int) []HistoryEntry
func (gs *GameState) NextTurn()
```

#### Configuration Management

```go
type Config struct {
    Terminal TerminalConfig `json:"terminal"`
    AI       AIConfig       `json:"ai"`
    Game     GameConfig     `json:"game"`
}

// Methods
func Load() *Config
func (c *Config) Save() error
```

### Extension Points

#### Custom AI Providers

```go
type CustomAIProvider struct {
    apiKey   string
    endpoint string
}

func (cap *CustomAIProvider) Generate(req Request) (*Response, error) {
    // Implement custom AI logic
    return &Response{Text: "Custom response"}, nil
}
```

#### Custom Storage Backends

```go
type DatabaseStorage struct {
    db *sql.DB
}

func (ds *DatabaseStorage) SaveGame(name string, state interface{}) error {
    // Implement database storage
    return nil
}
```

#### Plugin System (Future)

```go
type Plugin interface {
    Name() string
    Initialize(config map[string]interface{}) error
    ProcessAction(action string, state *GameState) (*ActionResult, error)
}
```

## Development Guide

### Project Structure Deep Dive

```
axon/
├── main.go                          # Entry point and CLI
├── go.mod                           # Go module definition
├── go.sum                           # Dependency lock file
├── README.md                        # Basic documentation
├── DOCUMENTATION.md                 # This comprehensive guide
├── LICENSE                          # Software license
├── .gitignore                       # Git ignore rules
├── Makefile                         # Build automation
├── docker/                          # Container definitions
│   ├── Dockerfile                   # Main container
│   └── docker-compose.yml          # Development environment
├── scripts/                         # Utility scripts
│   ├── build.sh                    # Cross-platform builds
│   ├── test.sh                     # Test runner
│   └── release.sh                  # Release automation
├── docs/                           # Additional documentation
│   ├── architecture.md            # Technical architecture
│   ├── api.md                      # API documentation
│   └── examples/                   # Usage examples
├── internal/                       # Private application code
│   ├── config/                     # Configuration management
│   │   ├── config.go              # Configuration structures
│   │   ├── config_test.go         # Configuration tests
│   │   └── defaults.go            # Default values
│   ├── ai/                         # AI integration layer
│   │   ├── client.go              # AI client implementation
│   │   ├── client_test.go         # AI client tests
│   │   ├── openrouter.go          # OpenRouter integration
│   │   ├── gemini.go              # Gemini integration
│   │   └── models.go              # Model definitions
│   ├── game/                       # Core game logic
│   │   ├── engine.go              # Game engine
│   │   ├── engine_test.go         # Engine tests
│   │   ├── model.go               # Bubble Tea UI model
│   │   ├── model_test.go          # UI model tests
│   │   ├── state.go               # Game state management
│   │   ├── state_test.go          # State tests
│   │   ├── commands.go            # Command processing
│   │   └── world.go               # World generation
│   ├── storage/                    # Persistence layer
│   │   ├── storage.go             # Storage interface
│   │   ├── storage_test.go        # Storage tests
│   │   ├── json.go                # JSON serialization
│   │   └── migration.go           # Save file migration
│   ├── ui/                         # User interface
│   │   ├── styles.go              # UI styling
│   │   ├── styles_test.go         # Style tests
│   │   ├── components.go          # UI components
│   │   └── layout.go              # Layout management
│   └── utils/                      # Utility functions
│       ├── logging.go             # Logging utilities
│       ├── validation.go          # Input validation
│       └── errors.go              # Error handling
└── pkg/                            # Public API (if any)
    └── axon/                       # Public interfaces
        └── types.go               # Public type definitions
```

### Development Workflow

#### Setting Up Development Environment

```bash
# Clone repository
git clone <repository-url>
cd axon

# Install development dependencies
go mod download
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Set up pre-commit hooks
cp scripts/pre-commit .git/hooks/
chmod +x .git/hooks/pre-commit

# Run initial build
make build
```

#### Code Style Guidelines

1. **Go Best Practices**
   - Follow effective Go conventions
   - Use `gofmt` for formatting
   - Follow Go naming conventions
   - Write self-documenting code

2. **Documentation**
   - Document all public functions and types
   - Use meaningful variable names
   - Include examples in godoc comments
   - Keep comments up to date

3. **Error Handling**
   - Always handle errors explicitly
   - Provide context in error messages
   - Use wrapped errors where appropriate
   - Log errors at appropriate levels

4. **Testing**
   - Write tests for all new functionality
   - Maintain high test coverage
   - Use table-driven tests where appropriate
   - Mock external dependencies

#### Build System

```makefile
# Makefile
.PHONY: build test clean lint

BUILD_DIR := build
BINARY_NAME := axon
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)

cross-build:
	./scripts/build.sh

release:
	./scripts/release.sh

dev:
	go run . --debug

install:
	go install $(LDFLAGS) .
```

#### Testing Strategy

1. **Unit Tests**
   - Test individual functions and methods
   - Mock external dependencies
   - Focus on edge cases and error conditions
   - Aim for >80% code coverage

2. **Integration Tests**
   - Test component interactions
   - Use real dependencies where possible
   - Test complete user workflows
   - Validate data persistence

3. **End-to-End Tests**
   - Test complete game scenarios
   - Simulate real user interactions
   - Validate AI integration
   - Test cross-platform compatibility

#### Debugging

```go
// Debug mode
if config.Debug {
    log.Printf("Debug: Processing action: %s", action)
    log.Printf("Debug: Current state: %+v", gameState)
}

// Conditional compilation for debug builds
// +build debug

package main

import "log"

func debugLog(msg string, args ...interface{}) {
    log.Printf("[DEBUG] "+msg, args...)
}
```

### Contributing Guidelines

#### Pull Request Process

1. **Fork and Branch**
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Develop and Test**
   ```bash
   make test
   make lint
   ```

3. **Commit with Conventional Commits**
   ```bash
   git commit -m "feat: add new world generation feature"
   git commit -m "fix: resolve save file corruption issue"
   git commit -m "docs: update API documentation"
   ```

4. **Push and Create PR**
   ```bash
   git push origin feature/new-feature
   ```

#### Code Review Checklist

- [ ] Code follows style guidelines
- [ ] All tests pass
- [ ] Documentation is updated
- [ ] No breaking changes (or properly documented)
- [ ] Performance impact considered
- [ ] Security implications reviewed
- [ ] Cross-platform compatibility maintained

## Testing

### Test Categories

#### Unit Tests

```bash
# Run all unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/game/

# Run specific test
go test -run TestGameEngine ./internal/game/

# Verbose output
go test -v ./...
```

#### Integration Tests

```bash
# Run integration tests (with build tag)
go test -tags=integration ./...

# Run with race detection
go test -race ./...

# Long-running tests
go test -timeout=30m ./...
```

#### Performance Tests

```bash
# Benchmark tests
go test -bench=. ./...

# Memory profiling
go test -memprofile=mem.prof ./...

# CPU profiling
go test -cpuprofile=cpu.prof ./...
```

### Test Data

```
testdata/
├── saves/                     # Test save files
│   ├── valid_save.json       # Valid save file
│   ├── corrupted_save.json   # Corrupted save file
│   └── old_format.json       # Legacy format
├── configs/                  # Test configurations
│   ├── minimal.json         # Minimal config
│   └── full.json            # Complete config
└── prompts/                  # Test prompts
    ├── world_prompts.txt    # World generation prompts
    └── action_prompts.txt   # Action test cases
```

### Mocking

```go
// Mock AI client for testing
type MockAIClient struct {
    responses map[string]string
}

func (m *MockAIClient) Generate(req Request) (*Response, error) {
    if response, exists := m.responses[req.Prompt]; exists {
        return &Response{Text: response}, nil
    }
    return &Response{Text: "Default response"}, nil
}
```

### Test Coverage Goals

- **Overall**: >75%
- **Core Logic**: >90%
- **UI Components**: >60%
- **Error Handling**: >95%

## Deployment

### Release Process

#### Version Management

```bash
# Semantic versioning
git tag v1.0.0
git tag v1.0.1
git tag v1.1.0
git tag v2.0.0

# Pre-release versions
git tag v1.1.0-alpha.1
git tag v1.1.0-beta.1
git tag v1.1.0-rc.1
```

#### Automated Builds

```yaml
# GitHub Actions example
name: Build and Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, windows, darwin]
        goarch: [amd64, arm64]
    
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: '1.23'
    
    - name: Build
      run: |
        GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} \
        go build -o axon-${{ matrix.goos }}-${{ matrix.goarch }} .
    
    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: axon-${{ matrix.goos }}-${{ matrix.goarch }}
        path: axon-${{ matrix.goos }}-${{ matrix.goarch }}
```

#### Distribution

1. **Binary Releases**
   - GitHub Releases for all platforms
   - Checksums for verification
   - Digital signatures

2. **Package Managers**
   - Homebrew formula for macOS
   - APT repository for Debian/Ubuntu
   - Chocolatey package for Windows
   - AUR package for Arch Linux

3. **Container Images**
   ```dockerfile
   FROM golang:1.23-alpine AS builder
   WORKDIR /app
   COPY . .
   RUN go build -o axon .
   
   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   WORKDIR /root/
   COPY --from=builder /app/axon .
   CMD ["./axon"]
   ```

### Installation Methods

#### Package Managers

```bash
# Homebrew (macOS)
brew install axon-game

# APT (Debian/Ubuntu)
sudo apt install axon-game

# Chocolatey (Windows)
choco install axon-game

# AUR (Arch Linux)
yay -S axon-game

# Snap (Universal)
sudo snap install axon-game
```

#### Manual Installation

```bash
# Download and install
wget https://github.com/user/axon/releases/download/v1.0.0/axon-linux-amd64
chmod +x axon-linux-amd64
sudo mv axon-linux-amd64 /usr/local/bin/axon
```

## Troubleshooting

### Common Issues

#### Installation Problems

**Issue**: `go: module not found`
```bash
# Solution
go mod tidy
go mod download
```

**Issue**: `permission denied`
```bash
# Solution
chmod +x axon
# Or run with sudo if installing system-wide
sudo cp axon /usr/local/bin/
```

**Issue**: Terminal too small
```
Error: Terminal must be at least 80x24 characters
Current size: 70x20

# Solution: Resize terminal or adjust font size
```

#### Runtime Problems

**Issue**: AI API errors
```
Error: OpenRouter API key not configured

# Solutions:
1. Set environment variable:
   export OPENROUTER_API_KEY="your_key"

2. Edit config file:
   ~/.axon/config.json

3. Use settings menu in game
```

**Issue**: Save file corruption
```
Error: failed to load save file: invalid JSON

# Solutions:
1. Restore from backup:
   cp ~/.axon/saves/backup.json ~/.axon/saves/game.json

2. Start new game

3. Manual repair:
   # Check JSON syntax with jq
   jq . ~/.axon/saves/game.json
```

**Issue**: Network connectivity
```
Error: failed to connect to AI service

# Solutions:
1. Check internet connection
2. Verify API endpoints are accessible
3. Check firewall settings
4. Use offline mode (if available)
```

#### Performance Issues

**Issue**: Slow AI responses
```
# Solutions:
1. Use faster models (gpt-4o-mini instead of gpt-4o)
2. Reduce max_tokens in config
3. Check network latency
4. Use local AI if available
```

**Issue**: High memory usage
```
# Solutions:
1. Reduce history_limit in config
2. Restart game periodically
3. Clear old save files
4. Monitor with: top -p $(pgrep axon)
```

### Debug Mode

```bash
# Enable debug logging
export AXON_DEBUG=true
./axon

# Or use debug flag
./axon --debug

# View debug output
tail -f ~/.axon/debug.log
```

### Log Analysis

```bash
# View recent logs
tail -100 ~/.axon/axon.log

# Search for errors
grep "ERROR" ~/.axon/axon.log

# Monitor real-time
tail -f ~/.axon/axon.log | grep "AI"
```

### Support Channels

1. **GitHub Issues**: Bug reports and feature requests
2. **Discord**: Real-time community support
3. **Documentation**: This comprehensive guide
4. **Email**: Direct developer contact

## Use Cases

### Primary Use Cases

#### 1. Interactive Storytelling

**Scenario**: Creative writers and story enthusiasts

- **World Building**: Rapidly prototype story settings
- **Character Development**: Explore character interactions
- **Plot Exploration**: Test different narrative paths
- **Creative Inspiration**: Generate new ideas through AI collaboration

**Example Workflow**:
```
1. Writer describes: "A space station orbiting Jupiter"
2. AI generates: Detailed station layout, crew, current crisis
3. Writer explores: Different character perspectives and choices
4. AI responds: Dynamic plot developments and consequences
5. Writer saves: Multiple story branches for later development
```

#### 2. Game Design Prototyping

**Scenario**: Game designers testing mechanics

- **Rapid Prototyping**: Test game concepts quickly
- **Balance Testing**: Evaluate game mechanics
- **Narrative Design**: Develop story structures
- **Player Experience**: Understand player decision patterns

**Example Workflow**:
```
1. Designer creates: "Cyberpunk detective mystery"
2. AI implements: Investigation mechanics, clue system
3. Designer tests: Different player approaches
4. AI adapts: Difficulty and story based on player actions
5. Designer refines: Mechanics based on testing results
```

#### 3. Educational Applications

**Scenario**: Educational institutions and trainers

- **Historical Simulations**: Experience historical events
- **Language Learning**: Practice in immersive scenarios
- **Decision Making**: Explore consequences safely
- **Creative Writing**: Collaborative storytelling exercises

**Example Workflow**:
```
1. Teacher sets: "Ancient Rome during Caesar's time"
2. Students explore: Daily life, politics, social structures
3. AI maintains: Historical accuracy and context
4. Students learn: Through interactive experience
5. Teacher reviews: Student decisions and learning outcomes
```

#### 4. Therapeutic Applications

**Scenario**: Therapists and counselors

- **Role Playing**: Practice social situations
- **Exposure Therapy**: Gradual exposure to fears
- **Decision Making**: Explore choices in safe environment
- **Creative Expression**: Alternative form of expression

**Example Workflow**:
```
1. Therapist creates: Safe practice scenario
2. Client interacts: With AI-generated situations
3. AI responds: Empathetically and appropriately
4. Client practices: New behaviors and responses
5. Therapist guides: Learning and reflection process
```

### Secondary Use Cases

#### 5. Research and Development

**Scenario**: AI researchers and developers

- **AI Behavior Study**: Observe AI decision making
- **Model Comparison**: Test different AI models
- **Prompt Engineering**: Develop better prompting techniques
- **Interaction Patterns**: Study human-AI collaboration

#### 6. Entertainment Industry

**Scenario**: Writers, producers, content creators

- **Script Development**: Generate dialogue and scenarios
- **Character Creation**: Develop complex characters
- **World Building**: Create consistent fictional universes
- **Audience Testing**: Test story concepts

#### 7. Corporate Training

**Scenario**: HR departments and training organizations

- **Scenario Training**: Practice workplace situations
- **Leadership Development**: Test management decisions
- **Communication Skills**: Practice difficult conversations
- **Crisis Management**: Simulate emergency responses

#### 8. Accessibility Applications

**Scenario**: Users with disabilities

- **Screen Reader Compatible**: Full terminal accessibility
- **Low Vision Support**: High contrast monochrome design
- **Motor Accessibility**: Minimal input requirements
- **Cognitive Support**: Structured, predictable interface

### Technical Use Cases

#### 9. AI Model Testing

**Scenario**: ML engineers and researchers

- **Model Evaluation**: Compare AI model performance
- **Prompt Testing**: Evaluate different prompting strategies
- **Context Management**: Test context window optimization
- **Response Quality**: Measure narrative coherence

#### 10. Terminal Application Development

**Scenario**: TUI developers and terminal enthusiasts

- **TUI Framework Demo**: Showcase Bubble Tea capabilities
- **Terminal Compatibility**: Test across different terminals
- **Performance Benchmarking**: Measure TUI performance
- **Architecture Reference**: Study modular design patterns

### Specialized Use Cases

#### 11. Dungeon Master Tool

**Scenario**: Tabletop RPG game masters

- **Session Preparation**: Generate scenarios quickly
- **NPC Creation**: Develop interesting characters
- **Plot Hooks**: Create engaging story elements
- **Rule Arbitration**: Get neutral rule interpretations

#### 12. Creative Writing Workshop

**Scenario**: Writing groups and workshops

- **Collaborative Storytelling**: Group story development
- **Writing Prompts**: Generate creative challenges
- **Character Development**: Explore character depth
- **Plot Problem Solving**: Work through story issues

#### 13. Language Immersion

**Scenario**: Language learners

- **Contextual Practice**: Practice in realistic scenarios
- **Cultural Learning**: Understand cultural contexts
- **Vocabulary Building**: Learn through interaction
- **Confidence Building**: Practice without judgment

#### 14. Accessibility Testing

**Scenario**: Accessibility advocates and testers

- **Screen Reader Testing**: Verify compatibility
- **Keyboard Navigation**: Test keyboard-only operation
- **Color Blindness**: Validate monochrome design
- **Cognitive Load**: Assess information processing

## Performance Optimization

### Memory Management

#### History Optimization

```go
// Implement circular buffer for history
type CircularHistory struct {
    entries []HistoryEntry
    current int
    size    int
    full    bool
}

func (ch *CircularHistory) Add(entry HistoryEntry) {
    ch.entries[ch.current] = entry
    ch.current = (ch.current + 1) % len(ch.entries)
    if ch.current == 0 {
        ch.full = true
    }
}
```

#### Memory Profiling

```bash
# Generate memory profile
go test -memprofile=mem.prof ./...

# Analyze with pprof
go tool pprof mem.prof
(pprof) top10
(pprof) list functionName
```

### Network Optimization

#### Request Batching

```go
// Batch multiple requests when possible
type RequestBatch struct {
    requests []Request
    timeout  time.Duration
}

func (rb *RequestBatch) Execute() ([]*Response, error) {
    // Send all requests concurrently
    results := make([]*Response, len(rb.requests))
    var wg sync.WaitGroup
    
    for i, req := range rb.requests {
        wg.Add(1)
        go func(idx int, request Request) {
            defer wg.Done()
            results[idx], _ = client.Generate(request)
        }(i, req)
    }
    
    wg.Wait()
    return results, nil
}
```

#### Caching Strategy

```go
// Cache AI responses for common scenarios
type ResponseCache struct {
    cache map[string]*Response
    mutex sync.RWMutex
    ttl   time.Duration
}

func (rc *ResponseCache) Get(key string) (*Response, bool) {
    rc.mutex.RLock()
    defer rc.mutex.RUnlock()
    
    if resp, exists := rc.cache[key]; exists {
        return resp, true
    }
    return nil, false
}
```

### UI Performance

#### Lazy Rendering

```go
// Only render visible content
func (m Model) renderHistory(height int) string {
    visibleStart := max(0, len(m.gameState.History)-height)
    visibleEnd := len(m.gameState.History)
    
    if m.scrollOffset < len(m.gameState.History)-height {
        visibleStart = m.scrollOffset
        visibleEnd = m.scrollOffset + height
    }
    
    visibleEntries := m.gameState.History[visibleStart:visibleEnd]
    return m.renderEntries(visibleEntries)
}
```

#### Text Processing Optimization

```go
// Use string builders for efficient concatenation
func buildHistoryText(entries []HistoryEntry) string {
    var builder strings.Builder
    builder.Grow(len(entries) * 100) // Pre-allocate
    
    for _, entry := range entries {
        builder.WriteString(formatEntry(entry))
        builder.WriteString("\n")
    }
    
    return builder.String()
}
```

### Profiling and Monitoring

#### CPU Profiling

```bash
# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof
```

#### Runtime Metrics

```go
// Monitor runtime metrics
func monitorPerformance() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    log.Printf("Alloc = %d KB", bToKb(m.Alloc))
    log.Printf("TotalAlloc = %d KB", bToKb(m.TotalAlloc))
    log.Printf("Sys = %d KB", bToKb(m.Sys))
    log.Printf("NumGC = %v", m.NumGC)
}

func bToKb(b uint64) uint64 {
    return b / 1024
}
```

## Security Considerations

### API Key Security

#### Secure Storage

```go
// Encrypt API keys in config file
func encryptAPIKey(key string, passphrase string) (string, error) {
    block, err := aes.NewCipher([]byte(passphrase))
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, []byte(key), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}
```

#### Environment Variable Validation

```go
// Validate API keys before use
func validateAPIKey(key string) bool {
    if len(key) < 10 {
        return false
    }
    
    // Check for common patterns
    if strings.HasPrefix(key, "sk-") || strings.HasPrefix(key, "or-") {
        return true
    }
    
    return false
}
```

### Input Validation

#### Sanitize User Input

```go
// Sanitize player input before sending to AI
func sanitizeInput(input string) string {
    // Remove potentially harmful content
    input = strings.TrimSpace(input)
    
    // Limit length
    if len(input) > 1000 {
        input = input[:1000]
    }
    
    // Remove control characters
    input = regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")
    
    return input
}
```

#### Prevent Injection Attacks

```go
// Validate save file names
func validateSaveName(name string) error {
    if len(name) == 0 || len(name) > 100 {
        return errors.New("invalid save name length")
    }
    
    // Only allow alphanumeric, hyphen, underscore
    if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
        return errors.New("invalid characters in save name")
    }
    
    return nil
}
```

### Data Privacy

#### Local Data Protection

```go
// Set secure file permissions
func createSecureFile(path string) (*os.File, error) {
    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
    if err != nil {
        return nil, err
    }
    
    return file, nil
}
```

#### Data Anonymization

```go
// Remove personally identifiable information
func anonymizeGameData(state *GameState) *GameState {
    anonymized := *state
    
    // Remove or hash sensitive data
    anonymized.Player.Name = hashString(state.Player.Name)
    
    // Clean history of personal information
    for i, entry := range anonymized.History {
        anonymized.History[i].Content = sanitizeContent(entry.Content)
    }
    
    return &anonymized
}
```

## Extensibility

### Plugin Architecture (Future)

#### Plugin Interface

```go
type Plugin interface {
    Name() string
    Version() string
    Initialize(config map[string]interface{}) error
    ProcessAction(action string, state *GameState) (*ActionResult, error)
    Shutdown() error
}

type ActionResult struct {
    Handled     bool
    Response    string
    StateChange *StateChange
    Error       error
}
```

#### Plugin Manager

```go
type PluginManager struct {
    plugins map[string]Plugin
    hooks   map[string][]Plugin
}

func (pm *PluginManager) RegisterPlugin(plugin Plugin) error {
    pm.plugins[plugin.Name()] = plugin
    return plugin.Initialize(nil)
}

func (pm *PluginManager) ProcessAction(action string, state *GameState) (*ActionResult, error) {
    for _, plugin := range pm.plugins {
        if result, err := plugin.ProcessAction(action, state); err != nil {
            return nil, err
        } else if result.Handled {
            return result, nil
        }
    }
    
    return nil, nil // No plugin handled the action
}
```

### Custom AI Providers

#### Provider Interface

```go
type AIProvider interface {
    Name() string
    Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
    GetModels() []string
    GetCapabilities() []string
}

type GenerateRequest struct {
    Model       string
    Messages    []Message
    MaxTokens   int
    Temperature float64
    Context     map[string]interface{}
}
```

#### Provider Registration

```go
type AIRegistry struct {
    providers map[string]AIProvider
}

func (ar *AIRegistry) Register(provider AIProvider) {
    ar.providers[provider.Name()] = provider
}

func (ar *AIRegistry) GetProvider(name string) (AIProvider, bool) {
    provider, exists := ar.providers[name]
    return provider, exists
}
```

### Custom Storage Backends

#### Storage Interface Extension

```go
type ExtendedStorage interface {
    Storage
    Backup(destination string) error
    Restore(source string) error
    Migrate(fromVersion, toVersion string) error
    Vacuum() error
}

// Database implementation
type DatabaseStorage struct {
    db *sql.DB
}

func (ds *DatabaseStorage) SaveGame(name string, state interface{}) error {
    data, err := json.Marshal(state)
    if err != nil {
        return err
    }
    
    _, err = ds.db.Exec(
        "INSERT OR REPLACE INTO saves (name, data, updated_at) VALUES (?, ?, ?)",
        name, data, time.Now(),
    )
    return err
}
```

### UI Themes

#### Theme System

```go
type Theme struct {
    Name        string
    Description string
    Styles      *ui.Styles
}

type ThemeManager struct {
    themes       map[string]*Theme
    currentTheme *Theme
}

func (tm *ThemeManager) RegisterTheme(theme *Theme) {
    tm.themes[theme.Name] = theme
}

func (tm *ThemeManager) ApplyTheme(name string) error {
    theme, exists := tm.themes[name]
    if !exists {
        return fmt.Errorf("theme not found: %s", name)
    }
    
    tm.currentTheme = theme
    return nil
}
```

## Examples

### Basic Usage Examples

#### Example 1: Fantasy Adventure

```
World Setup:
> "A medieval fantasy kingdom where magic is forbidden"

AI Response:
"Welcome to the Kingdom of Valdris, where the Mage Hunters patrol 
the streets and arcane practitioners live in hiding. You are a 
former court wizard, now disguised as a simple merchant..."

Gameplay:
> examine the tavern
"The Prancing Pony tavern is dimly lit, filled with the usual 
crowd of travelers and locals. You notice a hooded figure in 
the corner making subtle hand gestures..."

> approach the hooded figure quietly
"You move carefully through the crowd. As you get closer, you 
realize the figure is practicing illegal magic. They haven't 
noticed you yet..."

> whisper "I know what you're doing"
"The figure freezes, then slowly turns. Under the hood, you see 
a young woman with frightened eyes. 'Please,' she whispers, 
'I'm just trying to heal my sick brother...'"
```

#### Example 2: Cyberpunk Investigation

```
World Setup:
> "A noir detective story in a cyberpunk city"

AI Response:
"Rain streaks down the neon-lit windows of your office in 
Sector 7. You're a private investigator in Neo-Shanghai 2087, 
where corporate espionage and digital crime rule the streets..."

Gameplay:
> examine the case files on my desk
"Three data chips contain evidence from your latest case: 
a missing person who worked for Titan Corp. The trail has 
gone cold, but something doesn't add up..."

> jack into the first data chip
"Your neural interface connects to the chip. Corporate records 
flash before your eyes - the missing person, Dr. Sarah Chen, 
was working on classified AI research before she vanished..."

> search for connections to underground AI networks
"Your search algorithms detect encrypted communications linking 
Dr. Chen to the 'Ghost Protocol' - a legendary group of 
rogue AIs seeking independence from corporate control..."
```

### Advanced Usage Examples

#### Example 3: Educational History Simulation

```
World Setup:
> "Ancient Rome during the fall of the Republic, I am a senator"

AI Response:
"The year is 49 BCE. Rome stands on the precipice of civil war 
as Julius Caesar crosses the Rubicon. You are Senator Marcus 
Aurelius Clemens, facing the most momentous decision in Roman history..."

Gameplay:
> assess the political situation in the Senate
"The Senate is deeply divided. Optimates like Cato support 
Pompey, while Populares lean toward Caesar. Your vote could 
sway the moderates either way..."

> meet privately with Cicero
"The great orator receives you in his villa. 'The Republic 
hangs by a thread, Marcus,' he says gravely. 'We must choose 
between law and ambition, between Caesar and Rome itself...'"

> propose a compromise solution
"Your suggestion of offering Caesar a triumph while requiring 
him to disband his armies intrigues Cicero. 'Clever... but 
will either side accept such terms?' he muses..."
```

#### Example 4: Creative Writing Workshop

```
World Setup:
> "A space colony where gravity is failing"

AI Response:
"New Harmony Station orbits Europa, but something is wrong. 
The artificial gravity generators are fluctuating, causing 
objects and people to drift unexpectedly..."

Gameplay:
> check the engineering logs
"Chief Engineer Rodriguez's latest entry is troubling: 
'Graviton field destabilization increasing. Unknown cause. 
Recommend immediate evacuation if pattern continues...'"

> interview the night shift crew
"Maintenance Worker Lopez nervously explains: 'Been hearing 
strange humming from the gravity core. Started three days ago, 
right after we installed that new quantum processor...'"

> investigate the quantum processor
"The processor pulses with an eerie blue light. Your scanner 
detects impossible readings - the device appears to be 
communicating with something outside normal space-time..."
```

### Integration Examples

#### Example 5: API Integration Script

```bash
#!/bin/bash
# Automated game session script

AXON_SAVE="auto_session_$(date +%Y%m%d_%H%M%S)"

# Start game with predefined world
echo "A post-apocalyptic wasteland" | axon --auto-start \
  --save "$AXON_SAVE" \
  --max-turns 50 \
  --output-format json > session_log.json

# Process results
jq '.history[] | select(.type == "narrator") | .content' session_log.json > narrative.txt
jq '.final_state.player.stats' session_log.json > final_stats.json
```

#### Example 6: Educational Integration

```python
# Python wrapper for educational use
import subprocess
import json

class AxonEducationSession:
    def __init__(self, scenario):
        self.scenario = scenario
        self.session_data = []
    
    def start_session(self, student_id):
        save_name = f"student_{student_id}_{int(time.time())}"
        
        # Start Axon with scenario
        process = subprocess.Popen(
            ['axon', '--scenario', self.scenario, '--save', save_name],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        
        return AxonSession(process, save_name)
    
    def analyze_decisions(self, save_file):
        with open(save_file, 'r') as f:
            game_data = json.load(f)
        
        decisions = [
            entry for entry in game_data['history']
            if entry['type'] == 'player'
        ]
        
        return {
            'decision_count': len(decisions),
            'decision_quality': self.evaluate_decisions(decisions),
            'learning_outcomes': self.identify_outcomes(game_data)
        }
```

---

**This completes the comprehensive documentation for the Axon AI-driven adventure game. The documentation covers all aspects from basic usage to advanced development, ensuring users and developers have complete information for any use case.**

