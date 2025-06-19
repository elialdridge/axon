# Terminal Compatibility Guide

Axon now includes comprehensive terminal detection and compatibility features to ensure the game works properly across a wide range of terminal emulators and operating systems, including minimal terminals and UNIX System V systems.

## Supported Terminal Types

### Modern Terminals (Full Features)
- **xterm and variants** (xterm-256color, xterm-88color, etc.)
- **GNOME Terminal**, **Konsole**, **iTerm2**, **Windows Terminal**
- **tmux** and **screen** sessions
- **Warp Terminal** (current environment)

Features available:
- Full color support
- Mouse interaction
- Alt screen buffer
- UTF-8 text rendering
- Advanced text formatting (bold, underline, reverse)

### Minimal Terminals (Basic Features)
- **dumb** terminal
- **Console redirects** (CI/CD environments)
- **Build systems** and **automated scripts**
- **Very small terminals** (less than 40x10)

Features provided:
- Plain text output only
- No mouse support
- No alt screen buffer
- ASCII-only rendering
- Essential functionality preserved

### UNIX System V Compatible
- **vt52**, **vt100**, **vt102**, **vt220**
- **ATT terminals**
- **Sun workstation** terminals
- **cons25** (FreeBSD console)

Features provided:
- Basic ANSI escape sequences only
- Compatible with legacy UNIX systems
- Reduced feature set for maximum compatibility
- Reliable text rendering

## Auto-Detection Features

### Terminal Capabilities Detection
Axon automatically detects:

1. **Terminal Size** - Uses multiple methods:
   - `ioctl` system calls (POSIX)
   - Environment variables (`COLUMNS`, `LINES`)
   - `tput` command
   - `stty` command
   - Safe fallback to 80x24

2. **Terminal Type** - Analyzes:
   - `$TERM` environment variable
   - `$TERM_PROGRAM` for specific programs
   - `$COLORTERM` for color support
   - Special environment indicators (`$CI`, `$BUILD`)

3. **Feature Support** - Determines:
   - Color rendering capabilities
   - Mouse event support
   - Alt screen buffer support
   - UTF-8 text encoding
   - Advanced formatting (bold, underline, reverse)

### Automatic Adjustments

When auto-detection is enabled (default), Axon automatically:

- **Disables mouse support** for terminals that don't support it
- **Disables alt screen** for minimal terminals
- **Strips ANSI formatting** for dumb terminals
- **Adjusts text wrapping** based on actual terminal width
- **Uses safe rendering** for System V terminals

## Manual Configuration

You can override auto-detection in the configuration file (`~/.axon/config.json`):

```json
{
  "terminal": {
    "width": 80,
    "height": 24,
    "color_enabled": false,
    "force_minimal": false,
    "force_systemv": false,
    "auto_detect": true,
    "mouse_enabled": true,
    "alt_screen_enabled": true
  }
}
```

### Configuration Options

- `auto_detect`: Enable/disable automatic terminal detection
- `force_minimal`: Force minimal mode regardless of detection
- `force_systemv`: Force System V compatibility mode
- `width`/`height`: Override detected terminal size
- `color_enabled`: Enable/disable color output
- `mouse_enabled`: Enable/disable mouse support
- `alt_screen_enabled`: Enable/disable alt screen buffer

## Testing Terminal Compatibility

Use the included terminal detection utility:

```bash
# Basic detection
./terminal-info

# JSON output for scripting
./terminal-info -json

# Verbose logging
./terminal-info -v
```

The utility returns different exit codes:
- `0`: Modern terminal with full features
- `2`: Minimal terminal detected
- `3`: System V compatible terminal detected

## Examples

### Testing Different Terminal Types

```bash
# Test with minimal terminal
TERM=dumb axon

# Test with System V terminal  
TERM=vt100 axon

# Test with small terminal size
COLUMNS=40 LINES=15 axon

# Force minimal mode
echo '{"terminal":{"force_minimal":true}}' > ~/.axon/config.json
```

### Integration with Scripts

```bash
#!/bin/bash
# Check terminal compatibility before launching
./terminal-info > /dev/null
case $? in
    0) echo "Full features available" ;;
    2) echo "Running in minimal mode" ;;
    3) echo "System V compatibility mode" ;;
esac
```

## Troubleshooting

### Terminal Not Detected Correctly

1. Check your `$TERM` environment variable:
   ```bash
   echo $TERM
   ```

2. Verify terminal size detection:
   ```bash
   ./terminal-info
   ```

3. Force specific mode if needed:
   ```json
   {"terminal": {"force_minimal": true, "auto_detect": false}}
   ```

### Performance Issues

For very slow terminals or systems:

1. Disable mouse support:
   ```json
   {"terminal": {"mouse_enabled": false}}
   ```

2. Use minimal mode:
   ```json
   {"terminal": {"force_minimal": true}}
   ```

3. Reduce terminal size:
   ```json
   {"terminal": {"width": 80, "height": 24}}
   ```

### Character Encoding Issues

For terminals with limited character support:

1. Force ASCII mode by setting `TERM=dumb`
2. Or use System V mode: `TERM=vt100`
3. Check locale settings: `locale`

## Implementation Details

### Detection Algorithm

1. **Size Detection**: Multiple fallback methods ensure size detection works even on constrained systems
2. **Capability Detection**: Conservative approach - features are disabled unless positively detected
3. **Compatibility Modes**: Minimal and System V modes provide maximum compatibility
4. **Safe Defaults**: Conservative defaults ensure the game works even if detection fails

### Performance Considerations

- Detection runs once at startup
- Results are cached for the session
- Minimal overhead for compatibility checks
- Graceful degradation when features unavailable

## Future Enhancements

Planned improvements:
- Runtime terminal capability changes detection
- Additional terminal types support
- Performance optimizations for slow terminals
- Extended System V features support
- Better integration with terminal multiplexers

## Contributing

When adding new terminal support:

1. Add detection logic to `internal/terminal/terminal.go`
2. Add corresponding tests in `internal/terminal/terminal_test.go`
3. Update this documentation
4. Test with actual target terminal if possible

