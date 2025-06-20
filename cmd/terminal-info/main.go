package main

import (
	"flag"
	"fmt"
	"os"

	"axon/internal/logger"
	"axon/internal/terminal"
)

func main() {
	var (
		verbose = flag.Bool("v", false, "Enable verbose output")
		json    = flag.Bool("json", false, "Output in JSON format")
	)
	flag.Parse()

	// Initialize minimal logger for this utility
	if *verbose {
		if err := logger.Init(); err != nil {
			fmt.Printf("Warning: Failed to initialize logger: %v\n", err)
		}
		defer logger.Close()
	}

	// Detect terminal capabilities
	termInfo := terminal.DetectTerminal()

	if *json {
		printJSON(termInfo)
	} else {
		printHuman(termInfo)
	}

	// Exit with specific codes for minimal/SystemV detection
	exitCode := 0
	if termInfo.IsMinimal {
		exitCode = 2 // Exit code 2 for minimal terminal
	} else if termInfo.IsSystemV {
		exitCode = 3 // Exit code 3 for System V terminal
	}

	// If logger was initialized, close it before exiting
	if *verbose {
		logger.Close()
	}
	os.Exit(exitCode)
}

func printJSON(info *terminal.TerminalInfo) {
	fmt.Printf(`{
  "width": %d,
  "height": %d,
  "term_type": "%s",
  "term_program": "%s",
  "is_minimal": %t,
  "is_systemv": %t,
  "color_support": %t,
  "mouse_support": %t,
  "alt_screen_support": %t,
  "supports_utf8": %t,
  "supports_bold": %t,
  "supports_underline": %t,
  "supports_reverse": %t
}
`, info.Width, info.Height, info.TermType, info.TermProgram,
		info.IsMinimal, info.IsSystemV, info.ColorSupport, info.MouseSupport,
		info.AltScreenSupport, info.SupportsUTF8, info.SupportsBold,
		info.SupportsUnderline, info.SupportsReverse)
}

func printHuman(info *terminal.TerminalInfo) {
	fmt.Printf("Terminal Detection Results:\n")
	fmt.Printf("==========================\n")
	fmt.Printf("Terminal Type: %s\n", info.TermType)
	if info.TermProgram != "" {
		fmt.Printf("Terminal Program: %s\n", info.TermProgram)
	}
	fmt.Printf("Size: %dx%d\n", info.Width, info.Height)
	fmt.Printf("\n")

	fmt.Printf("Compatibility Mode:\n")
	if info.IsMinimal {
		fmt.Printf("  [✓] Minimal Terminal Detected\n")
	} else {
		fmt.Printf("  [ ] Minimal Terminal\n")
	}
	if info.IsSystemV {
		fmt.Printf("  [✓] UNIX System V Compatible\n")
	} else {
		fmt.Printf("  [ ] UNIX System V Compatible\n")
	}
	fmt.Printf("\n")

	fmt.Printf("Capabilities:\n")
	printCapability("Color Support", info.ColorSupport)
	printCapability("Mouse Support", info.MouseSupport)
	printCapability("Alt Screen", info.AltScreenSupport)
	printCapability("UTF-8 Support", info.SupportsUTF8)
	printCapability("Bold Text", info.SupportsBold)
	printCapability("Underline", info.SupportsUnderline)
	printCapability("Reverse Video", info.SupportsReverse)

	fmt.Printf("\n")
	fmt.Printf("Recommended Settings:\n")
	options := info.GetCompatibleOptions()
	if len(options) == 0 {
		fmt.Printf("  Full features available\n")
	} else {
		for _, opt := range options {
			fmt.Printf("  - %s\n", opt)
		}
	}

	// Safe size recommendation
	safeWidth, safeHeight := info.GetSafeSize()
	if safeWidth != info.Width || safeHeight != info.Height {
		fmt.Printf("  Recommended size: %dx%d\n", safeWidth, safeHeight)
	}
}

func printCapability(name string, supported bool) {
	if supported {
		fmt.Printf("  [✓] %s\n", name)
	} else {
		fmt.Printf("  [ ] %s\n", name)
	}
}
