// src/internal/app/app.go
// Application initialization and setup

package app

import (
	"os"

	"nextui-themes/internal/accents"
	"nextui-themes/internal/icons"
	"nextui-themes/internal/leds"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/themes"
)

// Initialize sets up the application
func Initialize() error {
	// Initialize app state
	state.CurrentScreen = MainMenu

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Set up environment variables for the TrimUI brick
	logging.LogDebug("Setting environment variables")

	_ = os.Setenv("DEVICE", "brick")
	_ = os.Setenv("PLATFORM", "tg5040")

	// Add current directory to PATH instead of replacing it
	existingPath := os.Getenv("PATH")
	newPath := cwd + ":" + existingPath
	_ = os.Setenv("PATH", newPath)
	logging.LogDebug("Updated PATH: %s", newPath)

	_ = os.Setenv("LD_LIBRARY_PATH", "/mnt/SDCARD/.system/tg5040/lib:/usr/trimui/lib")

	// Create directory structure
	logging.LogDebug("Creating theme directories")

	// Create theme directory structure
	if err := themes.EnsureThemeDirectoryStructure(); err != nil {
		logging.LogDebug("Warning: Could not create theme directories: %v", err)
	}

	// Create placeholder files
	if err := themes.CreatePlaceholderFiles(); err != nil {
		logging.LogDebug("Warning: Could not create placeholder files: %v", err)
	}

	// Initialize new component system
	if err := themes.InitComponentSystem(); err != nil {
		logging.LogDebug("Warning: Could not initialize component system: %v", err)
	}

	// Create Icons directory and placeholder
	if err := icons.CreatePlaceholderFile(); err != nil {
		logging.LogDebug("Error creating icons placeholder: %v", err)
	}

	// Initialize accent colors
	if err := accents.InitAccentColors(); err != nil {
		logging.LogDebug("Error initializing accent colors: %v", err)
	}

	// Initialize LED settings
	if err := leds.InitLEDSettings(); err != nil {
		logging.LogDebug("Error initializing LED settings: %v", err)
	}

	// Reset component state to defaults
	ResetComponentState()

	logging.LogDebug("Initialization complete")
	return nil
}
