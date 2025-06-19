package terminal

import (
	"os"
	"testing"
)

func TestDetectTerminal(t *testing.T) {
	// Save original environment
	originalTerm := os.Getenv("TERM")
	originalTermProgram := os.Getenv("TERM_PROGRAM")
	originalCI := os.Getenv("CI")
	defer func() {
		os.Setenv("TERM", originalTerm)
		os.Setenv("TERM_PROGRAM", originalTermProgram)
		os.Setenv("CI", originalCI)
	}()

	tests := []struct {
		name           string
		term           string
		termProgram    string
		expectedMinimal bool
		expectedSystemV bool
	}{
		{
			name:           "xterm modern terminal",
			term:           "xterm-256color",
			termProgram:    "",
			expectedMinimal: false,
			expectedSystemV: false,
		},
		{
			name:           "dumb terminal",
			term:           "dumb",
			termProgram:    "",
			expectedMinimal: true,
			expectedSystemV: false,
		},
		{
			name:           "vt100 terminal",
			term:           "vt100",
			termProgram:    "",
			expectedMinimal: true,
			expectedSystemV: true,
		},
		{
			name:           "empty terminal",
			term:           "",
			termProgram:    "",
			expectedMinimal: true,
			expectedSystemV: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TERM", tt.term)
			os.Setenv("TERM_PROGRAM", tt.termProgram)
			// Clear CI variable to avoid interference with test expectations
			os.Unsetenv("CI")
			os.Unsetenv("BUILD")
			
			info := DetectTerminal()
			
			if info.IsMinimal != tt.expectedMinimal {
				t.Errorf("Expected IsMinimal %v, got %v", tt.expectedMinimal, info.IsMinimal)
			}
			
			if info.IsSystemV != tt.expectedSystemV {
				t.Errorf("Expected IsSystemV %v, got %v", tt.expectedSystemV, info.IsSystemV)
			}
		})
	}
}

func TestGetSafeSize(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		height         int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "normal size",
			width:          80,
			height:         24,
			expectedWidth:  80,
			expectedHeight: 24,
		},
		{
			name:           "too small",
			width:          20,
			height:         5,
			expectedWidth:  40,
			expectedHeight: 10,
		},
		{
			name:           "too large",
			width:          300,
			height:         100,
			expectedWidth:  200,
			expectedHeight: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &TerminalInfo{
				Width:  tt.width,
				Height: tt.height,
			}
			
			w, h := info.GetSafeSize()
			
			if w != tt.expectedWidth {
				t.Errorf("Expected width %d, got %d", tt.expectedWidth, w)
			}
			
			if h != tt.expectedHeight {
				t.Errorf("Expected height %d, got %d", tt.expectedHeight, h)
			}
		})
	}
}

func TestStripFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "text with ANSI color codes",
			input:    "\033[31mRed text\033[0m",
			expected: "Red text",
		},
		{
			name:     "complex ANSI codes",
			input:    "\033[1;31mBold Red\033[0m normal \033[4munderline\033[0m",
			expected: "Bold Red normal underline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripFormatting(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFormatForTerminal(t *testing.T) {
	tests := []struct {
		name      string
		minimal   bool
		systemV   bool
		input     string
		expected  string
	}{
		{
			name:     "normal terminal",
			minimal:  false,
			systemV:  false,
			input:    "\033[31mRed text\033[0m",
			expected: "\033[31mRed text\033[0m",
		},
		{
			name:     "minimal terminal",
			minimal:  true,
			systemV:  false,
			input:    "\033[31mRed text\033[0m",
			expected: "Red text",
		},
		{
			name:     "system V terminal",
			minimal:  false,
			systemV:  true,
			input:    "\033[31mRed text\033[0m",
			expected: "\033[31mRed text\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &TerminalInfo{
				IsMinimal: tt.minimal,
				IsSystemV: tt.systemV,
			}
			
			result := info.FormatForTerminal(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetCompatibleOptions(t *testing.T) {
	tests := []struct {
		name        string
		minimal     bool
		systemV     bool
		mouseSupport bool
		altScreen   bool
		expectedOpts []string
	}{
		{
			name:         "modern terminal",
			minimal:      false,
			systemV:      false,
			mouseSupport: true,
			altScreen:    true,
			expectedOpts: []string{},
		},
		{
			name:         "minimal terminal",
			minimal:      true,
			systemV:      false,
			mouseSupport: false,
			altScreen:    false,
			expectedOpts: []string{"no-mouse", "no-altscreen", "simple-render"},
		},
		{
			name:         "system V terminal",
			minimal:      false,
			systemV:      true,
			mouseSupport: false,
			altScreen:    false,
			expectedOpts: []string{"no-mouse", "no-altscreen"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &TerminalInfo{
				IsMinimal:        tt.minimal,
				IsSystemV:        tt.systemV,
				MouseSupport:     tt.mouseSupport,
				AltScreenSupport: tt.altScreen,
			}
			
			options := info.GetCompatibleOptions()
			
			if len(options) != len(tt.expectedOpts) {
				t.Errorf("Expected %d options, got %d", len(tt.expectedOpts), len(options))
			}
			
			for i, expected := range tt.expectedOpts {
				if i >= len(options) || options[i] != expected {
					t.Errorf("Expected option %q at index %d, got %q", expected, i, options[i])
				}
			}
		})
	}
}

func TestDetectCapabilities(t *testing.T) {
	// Save original environment
	originalColorTerm := os.Getenv("COLORTERM")
	defer func() {
		os.Setenv("COLORTERM", originalColorTerm)
	}()

	tests := []struct {
		name         string
		term         string
		colorTerm    string
		expectedColor bool
		expectedMouse bool
		expectedUTF8  bool
	}{
		{
			name:         "xterm with color",
			term:         "xterm-256color",
			colorTerm:    "",
			expectedColor: true,
			expectedMouse: true,
			expectedUTF8:  true,
		},
		{
			name:         "dumb terminal",
			term:         "dumb",
			colorTerm:    "",
			expectedColor: false,
			expectedMouse: false,
			expectedUTF8:  false,
		},
		{
			name:         "linux console",
			term:         "linux",
			colorTerm:    "",
			expectedColor: true,
			expectedMouse: false,
			expectedUTF8:  false,
		},
		{
			name:         "forced color support",
			term:         "dumb",
			colorTerm:    "truecolor",
			expectedColor: true,
			expectedMouse: false,
			expectedUTF8:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("COLORTERM", tt.colorTerm)
			
			info := &TerminalInfo{TermType: tt.term}
			info.detectCapabilities()
			
			if info.ColorSupport != tt.expectedColor {
				t.Errorf("Expected ColorSupport %v, got %v", tt.expectedColor, info.ColorSupport)
			}
			
			if info.MouseSupport != tt.expectedMouse {
				t.Errorf("Expected MouseSupport %v, got %v", tt.expectedMouse, info.MouseSupport)
			}
			
			if info.SupportsUTF8 != tt.expectedUTF8 {
				t.Errorf("Expected SupportsUTF8 %v, got %v", tt.expectedUTF8, info.SupportsUTF8)
			}
		})
	}
}

