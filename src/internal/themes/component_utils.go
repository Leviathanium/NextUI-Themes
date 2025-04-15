// src/internal/themes/component_utils.go
// Utility functions for component packages

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"nextui-themes/internal/logging"
)

// Time helper functions for consistent timestamp handling
func getCurrentTime() time.Time {
	return time.Now()
}

// This is a lightweight logger interface to accommodate different logger implementations
type Logger struct {
	*os.File
}

// Printf logs a formatted message to the logger
func (l *Logger) Printf(format string, args ...interface{}) {
	if l.File != nil {
		fmt.Fprintf(l.File, "[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
		// Also flush to disk
		l.File.Sync()
	}

	// Also log to the application logger
	logging.LogDebug(format, args...)
}

// EnsureComponentDirectoryStructure ensures all necessary directories exist
func EnsureComponentDirectoryStructure() error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Component directories to create
	directories := []string{
		// Theme directories
		filepath.Join(cwd, "Themes", "Imports"),
		filepath.Join(cwd, "Themes", "Exports"),

		// Accent directories
		filepath.Join(cwd, "Accents", "Imports"),
		filepath.Join(cwd, "Accents", "Exports"),
		filepath.Join(cwd, "Accents", "Presets"),
		filepath.Join(cwd, "Accents", "Custom"),

		// LED directories
		filepath.Join(cwd, "LEDs", "Imports"),
		filepath.Join(cwd, "LEDs", "Exports"),
		filepath.Join(cwd, "LEDs", "Presets"),
		filepath.Join(cwd, "LEDs", "Custom"),

		// Wallpaper directories
		filepath.Join(cwd, "Wallpapers", "Imports"),
		filepath.Join(cwd, "Wallpapers", "Exports"),
		filepath.Join(cwd, "Wallpapers", "Default"),

		// Icon directories
		filepath.Join(cwd, "Icons", "Imports"),
		filepath.Join(cwd, "Icons", "Exports"),

		// Font directories
		filepath.Join(cwd, "Fonts", "Imports"),
		filepath.Join(cwd, "Fonts", "Exports"),
		filepath.Join(cwd, "Fonts", "Backups"),

		// Log directory
		filepath.Join(cwd, "Logs"),
	}

	// Create each directory
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logging.LogDebug("Error creating directory %s: %v", dir, err)
			return err
		}
	}

	logging.LogDebug("Component directory structure created")
	return nil
}

// CreateComponentPlaceholders creates README files in component directories
func CreateComponentPlaceholders() error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Define placeholder files
	placeholders := map[string]string{
		// Theme placeholders
		filepath.Join(cwd, "Themes", "Imports", "README.txt"): `# Theme Import Directory

Place theme packages (directories with .theme extension) here to import them.
Themes should contain a manifest.json file and the appropriate theme files.`,

		filepath.Join(cwd, "Themes", "Exports", "README.txt"): `# Theme Export Directory

Exported theme packages will be placed here with sequential names (theme_1.theme, theme_2.theme, etc.)`,

		// Accent placeholders
		filepath.Join(cwd, "Accents", "Imports", "README.txt"): `# Accent Pack Import Directory

Place accent packs (directories with .acc extension) here to import them.
Accent packs should contain a minuisettings.txt file.`,

		filepath.Join(cwd, "Accents", "Exports", "README.txt"): `# Accent Pack Export Directory

Exported accent packs will be placed here with sequential names (accent_1.acc, accent_2.acc, etc.)`,

		// LED placeholders
		filepath.Join(cwd, "LEDs", "Imports", "README.txt"): `# LED Pack Import Directory

Place LED packs (directories with .led extension) here to import them.
LED packs should contain a ledsettings_brick.txt file.`,

		filepath.Join(cwd, "LEDs", "Exports", "README.txt"): `# LED Pack Export Directory

Exported LED packs will be placed here with sequential names (led_1.led, led_2.led, etc.)`,

		// Wallpaper placeholders
		filepath.Join(cwd, "Wallpapers", "Imports", "README.txt"): `# Wallpaper Pack Import Directory

Place wallpaper packs (directories with .bg extension) here to import them.
Wallpaper packs should contain SystemWallpapers and CollectionWallpapers directories.`,

		filepath.Join(cwd, "Wallpapers", "Exports", "README.txt"): `# Wallpaper Pack Export Directory

Exported wallpaper packs will be placed here with sequential names (wallpaper_1.bg, wallpaper_2.bg, etc.)`,

		// Icon placeholders
		filepath.Join(cwd, "Icons", "Imports", "README.txt"): `# Icon Pack Import Directory

Place icon packs (directories with .icon extension) here to import them.
Icon packs should contain SystemIcons, ToolIcons, and CollectionIcons directories.`,

		filepath.Join(cwd, "Icons", "Exports", "README.txt"): `# Icon Pack Export Directory

Exported icon packs will be placed here with sequential names (icon_1.icon, icon_2.icon, etc.)`,

		// Font placeholders
		filepath.Join(cwd, "Fonts", "Imports", "README.txt"): `# Font Pack Import Directory

Place font packs (directories with .font extension) here to import them.
Font packs should contain OG.ttf, Next.ttf, and their backup files.`,

		filepath.Join(cwd, "Fonts", "Exports", "README.txt"): `# Font Pack Export Directory

Exported font packs will be placed here with sequential names (font_1.font, font_2.font, etc.)`,

		filepath.Join(cwd, "Fonts", "Backups", "README.txt"): `# Font Backup Directory

This directory contains backup copies of the original system fonts.
DO NOT DELETE THESE FILES! They are needed to restore the original fonts.`,
	}

	// Create each placeholder file if the directory is empty
	for filePath, content := range placeholders {
		dir := filepath.Dir(filePath)

		// Check if directory is empty (except for other README files)
		entries, err := os.ReadDir(dir)
		if err != nil {
			logging.LogDebug("Error reading directory %s: %v", dir, err)
			continue
		}

		hasContent := false
		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() != "README.txt" {
				hasContent = true
				break
			}
		}

		// Create README if directory is empty or only contains README
		if !hasContent {
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				logging.LogDebug("Error creating placeholder file %s: %v", filePath, err)
			}
		}
	}

	return nil
}

// InitComponentSystem initializes the component system
func InitComponentSystem() error {
	// Create the directory structure
	if err := EnsureComponentDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create component directory structure: %w", err)
	}

	// Create placeholder files
	if err := CreateComponentPlaceholders(); err != nil {
		logging.LogDebug("Warning: Failed to create some placeholder files: %v", err)
		// Continue anyway
	}

	logging.LogDebug("Component system initialized successfully")
	return nil
}