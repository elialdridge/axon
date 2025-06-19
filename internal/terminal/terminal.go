package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"axon/internal/logger"
)

// TerminalInfo contains information about terminal capabilities
type TerminalInfo struct {
	// Basic capabilities
	Width            int
	Height           int
	ColorSupport     bool
	MouseSupport     bool
	AltScreenSupport bool

	// Terminal type information
	TermType    string
	TermProgram string
	IsMinimal   bool
	IsSystemV   bool

	// Specific capabilities
	SupportsUTF8      bool
	SupportsBold      bool
	SupportsUnderline bool
	SupportsReverse   bool
}

// DetectTerminal detects terminal capabilities and returns TerminalInfo
func DetectTerminal() *TerminalInfo {
	info := &TerminalInfo{
		Width:  80, // Default fallback
		Height: 24, // Default fallback
	}

	logger.Debug("Starting terminal detection")

	// Get terminal type
	info.TermType = os.Getenv("TERM")
	info.TermProgram = os.Getenv("TERM_PROGRAM")

	logger.Debug("TERM=%s, TERM_PROGRAM=%s", info.TermType, info.TermProgram)

	// Detect terminal size
	info.Width, info.Height = getTerminalSize()
	logger.Debug("Terminal size: %dx%d", info.Width, info.Height)

	// Detect capabilities based on TERM variable
	info.detectCapabilities()

	// Check for minimal terminal indicators
	info.detectMinimalTerminal()

	// Check for UNIX System V compatibility mode
	info.detectSystemV()

	logger.Info("Terminal detected - Type: %s, Size: %dx%d, Minimal: %v, SystemV: %v",
		info.TermType, info.Width, info.Height, info.IsMinimal, info.IsSystemV)

	return info
}

// getTerminalSize attempts to get terminal size using multiple methods
func getTerminalSize() (width, height int) {
	// Method 1: Try ioctl syscall (POSIX)
	if w, h := getTerminalSizeIoctl(); w > 0 && h > 0 {
		return w, h
	}

	// Method 2: Try environment variables
	if w, h := getTerminalSizeEnv(); w > 0 && h > 0 {
		return w, h
	}

	// Method 3: Try tput command
	if w, h := getTerminalSizeTput(); w > 0 && h > 0 {
		return w, h
	}

	// Method 4: Try stty command
	if w, h := getTerminalSizeStty(); w > 0 && h > 0 {
		return w, h
	}

	// Fallback to conservative defaults
	logger.Debug("Could not detect terminal size, using defaults")
	return 80, 24
}

// getTerminalSizeIoctl gets terminal size using ioctl syscall
func getTerminalSizeIoctl() (width, height int) {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		logger.Debug("ioctl failed: %v", errno)
		return 0, 0
	}

	return int(ws.Col), int(ws.Row)
}

// getTerminalSizeEnv gets terminal size from environment variables
func getTerminalSizeEnv() (width, height int) {
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if c, err := strconv.Atoi(cols); err == nil && c > 0 {
			width = c
		}
	}

	if lines := os.Getenv("LINES"); lines != "" {
		if l, err := strconv.Atoi(lines); err == nil && l > 0 {
			height = l
		}
	}

	if width > 0 && height > 0 {
		logger.Debug("Got terminal size from environment: %dx%d", width, height)
		return width, height
	}

	return 0, 0
}

// getTerminalSizeTput gets terminal size using tput command
func getTerminalSizeTput() (width, height int) {
	// Try tput cols
	if cmd := exec.Command("tput", "cols"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			if c, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil && c > 0 {
				width = c
			}
		}
	}

	// Try tput lines
	if cmd := exec.Command("tput", "lines"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			if l, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil && l > 0 {
				height = l
			}
		}
	}

	if width > 0 && height > 0 {
		logger.Debug("Got terminal size from tput: %dx%d", width, height)
		return width, height
	}

	return 0, 0
}

// getTerminalSizeStty gets terminal size using stty command
func getTerminalSizeStty() (width, height int) {
	cmd := exec.Command("stty", "size")
	if output, err := cmd.Output(); err == nil {
		parts := strings.Fields(strings.TrimSpace(string(output)))
		if len(parts) >= 2 {
			if h, err := strconv.Atoi(parts[0]); err == nil && h > 0 {
				height = h
			}
			if w, err := strconv.Atoi(parts[1]); err == nil && w > 0 {
				width = w
			}
		}
	}

	if width > 0 && height > 0 {
		logger.Debug("Got terminal size from stty: %dx%d", width, height)
		return width, height
	}

	return 0, 0
}

// detectCapabilities detects terminal capabilities based on TERM variable
func (info *TerminalInfo) detectCapabilities() {
	term := strings.ToLower(info.TermType)

	// Default to minimal capabilities
	info.ColorSupport = false
	info.MouseSupport = false
	info.AltScreenSupport = false
	info.SupportsUTF8 = false
	info.SupportsBold = false
	info.SupportsUnderline = false
	info.SupportsReverse = false

	// Check for known terminal types with extended capabilities
	if strings.Contains(term, "xterm") ||
		strings.Contains(term, "screen") ||
		strings.Contains(term, "tmux") ||
		strings.Contains(term, "rxvt") ||
		term == "linux" {

		info.ColorSupport = true
		info.SupportsBold = true
		info.SupportsUnderline = true
		info.SupportsReverse = true

		// UTF-8 support for modern terminals
		if strings.Contains(term, "xterm") ||
			strings.Contains(term, "screen") ||
			strings.Contains(term, "tmux") {
			info.SupportsUTF8 = true
			info.MouseSupport = true
			info.AltScreenSupport = true
		}
	}

	// Check for color support via environment
	if os.Getenv("COLORTERM") != "" {
		info.ColorSupport = true
	}

	logger.Debug("Terminal capabilities - Color: %v, Mouse: %v, AltScreen: %v, UTF8: %v",
		info.ColorSupport, info.MouseSupport, info.AltScreenSupport, info.SupportsUTF8)
}

// detectMinimalTerminal checks if we're running on a minimal terminal
func (info *TerminalInfo) detectMinimalTerminal() {
	term := strings.ToLower(info.TermType)

	// Minimal terminals indicators
	minimalTerminals := []string{
		"dumb",
		"unknown",
		"",
		"vt52",
		"vt100",
		"vt102",
		"ansi",
		"cons25",
	}

	for _, minimal := range minimalTerminals {
		if term == minimal {
			info.IsMinimal = true
			break
		}
	}

	// Additional checks for minimal environments
	if os.Getenv("TERM") == "" ||
		os.Getenv("TERM") == "dumb" ||
		os.Getenv("CI") != "" ||
		os.Getenv("BUILD") != "" {
		info.IsMinimal = true
	}

	// Force minimal mode if size is very small
	if info.Width < 40 || info.Height < 10 {
		info.IsMinimal = true
	}

	logger.Debug("Minimal terminal detection: %v", info.IsMinimal)
}

// detectSystemV checks for UNIX System V compatibility requirements
func (info *TerminalInfo) detectSystemV() {
	// Check for System V indicators
	if os.Getenv("SYSV") != "" {
		info.IsSystemV = true
	}

	// Check for old UNIX systems
	term := strings.ToLower(info.TermType)
	systemVTerminals := []string{
		"vt52",
		"vt100",
		"vt102",
		"vt220",
		"ansi",
		"att",
		"sun",
		"cons25",
	}

	for _, sysv := range systemVTerminals {
		if strings.Contains(term, sysv) {
			info.IsSystemV = true
			break
		}
	}

	logger.Debug("System V compatibility: %v", info.IsSystemV)
}

// GetCompatibleOptions returns Bubble Tea options compatible with terminal
func (info *TerminalInfo) GetCompatibleOptions() []string {
	var options []string

	// Always disable mouse for minimal terminals
	if info.IsMinimal || info.IsSystemV || !info.MouseSupport {
		options = append(options, "no-mouse")
	}

	// Disable alt screen for minimal terminals
	if info.IsMinimal || info.IsSystemV || !info.AltScreenSupport {
		options = append(options, "no-altscreen")
	}

	// Force simple rendering for minimal terminals
	if info.IsMinimal {
		options = append(options, "simple-render")
	}

	logger.Debug("Compatible options: %v", options)
	return options
}

// GetSafeSize returns safe terminal size with bounds checking
func (info *TerminalInfo) GetSafeSize() (width, height int) {
	width = info.Width
	height = info.Height

	// Minimum safe sizes
	minWidth := 40
	minHeight := 10

	// Maximum safe sizes for compatibility
	maxWidth := 200
	maxHeight := 60

	if width < minWidth {
		width = minWidth
	}
	if width > maxWidth {
		width = maxWidth
	}

	if height < minHeight {
		height = minHeight
	}
	if height > maxHeight {
		height = maxHeight
	}

	return width, height
}

// FormatForTerminal formats text appropriately for the terminal type
func (info *TerminalInfo) FormatForTerminal(text string) string {
	// For minimal terminals, strip any formatting
	if info.IsMinimal {
		return stripFormatting(text)
	}

	// For System V, use basic formatting only
	if info.IsSystemV {
		return basicFormatting(text)
	}

	return text
}

// stripFormatting removes all formatting from text
func stripFormatting(text string) string {
	// Remove ANSI escape codes
	result := ""
	inEscape := false

	for _, char := range text {
		if char == '\033' { // ESC character
			inEscape = true
			continue
		}

		if inEscape {
			if char == 'm' || char == 'J' || char == 'H' {
				inEscape = false
			}
			continue
		}

		result += string(char)
	}

	return result
}

// basicFormatting applies only basic System V compatible formatting
func basicFormatting(text string) string {
	// Only allow basic text formatting
	// This is a placeholder - could be extended based on specific System V needs
	return text
}

// PrintCapabilities prints terminal capabilities for debugging
func (info *TerminalInfo) PrintCapabilities() {
	fmt.Printf("Terminal Information:\n")
	fmt.Printf("  Type: %s\n", info.TermType)
	fmt.Printf("  Program: %s\n", info.TermProgram)
	fmt.Printf("  Size: %dx%d\n", info.Width, info.Height)
	fmt.Printf("  Color Support: %v\n", info.ColorSupport)
	fmt.Printf("  Mouse Support: %v\n", info.MouseSupport)
	fmt.Printf("  Alt Screen: %v\n", info.AltScreenSupport)
	fmt.Printf("  UTF-8: %v\n", info.SupportsUTF8)
	fmt.Printf("  Bold: %v\n", info.SupportsBold)
	fmt.Printf("  Underline: %v\n", info.SupportsUnderline)
	fmt.Printf("  Reverse: %v\n", info.SupportsReverse)
	fmt.Printf("  Minimal: %v\n", info.IsMinimal)
	fmt.Printf("  System V: %v\n", info.IsSystemV)
}
