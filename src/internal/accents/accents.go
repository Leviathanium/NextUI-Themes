// src/internal/accents/accents.go
// Accent color management for NextUI Theme Manager

package accents

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nextui-themes/internal/logging"
)

// ThemeColor represents a set of colors for a theme
type ThemeColor struct {
	Name   string
	Color1 string // Main UI color
	Color2 string // Primary accent color
	Color3 string // Secondary accent color
	Color4 string // List text color
	Color5 string // Selected list text color
	Color6 string // Hint/info color
}

// Settings file path
const (
	SettingsPath = "/mnt/SDCARD/.userdata/shared/minuisettings.txt"
	AccentsDir   = "Accents" // Directory for external accent theme files
	PresetsDir   = "Presets" // Subdirectory for preset themes
	CustomDir    = "Custom"  // Subdirectory for custom themes
)

// CurrentTheme holds the currently loaded theme settings
var CurrentTheme ThemeColor

// External themes loaded from files
var (
	PresetThemes []ThemeColor
	CustomThemes []ThemeColor
)

// convertHexFormat converts between display format (#RRGGBB) and storage format (0xRRGGBB)
func convertHexFormat(color string, toStorage bool) string {
	if toStorage {
		// Convert from #RRGGBB to 0xRRGGBB
		if strings.HasPrefix(color, "#") {
			return "0x" + color[1:]
		}
		return color // Already in storage format
	} else {
		// Convert from 0xRRGGBB to #RRGGBB
		if strings.HasPrefix(color, "0x") {
			return "#" + color[2:]
		}
		return color // Already in display format
	}
}

// ConvertHexFormat is the exported version of convertHexFormat
func ConvertHexFormat(color string, toStorage bool) string {
	return convertHexFormat(color, toStorage)
}

// InitAccentColors loads the current accent colors from disk
func InitAccentColors() error {
	logging.LogDebug("Initializing accent colors")

	// Set a default theme name
	CurrentTheme.Name = "Current"

	// Try to load settings from disk
	theme, err := GetCurrentColors()
	if err != nil {
		logging.LogDebug("Error loading current colors: %v, using defaults", err)
		// Initialize with default values
		CurrentTheme.Color1 = "#FFFFFF" // White
		CurrentTheme.Color2 = "#9B2257" // Pink
		CurrentTheme.Color3 = "#1E2329" // Dark Blue
		CurrentTheme.Color4 = "#FFFFFF" // White
		CurrentTheme.Color5 = "#000000" // Black
		CurrentTheme.Color6 = "#FFFFFF" // White
	} else {
		// Convert from storage format to display format
		CurrentTheme.Color1 = convertHexFormat(theme.Color1, false)
		CurrentTheme.Color2 = convertHexFormat(theme.Color2, false)
		CurrentTheme.Color3 = convertHexFormat(theme.Color3, false)
		CurrentTheme.Color4 = convertHexFormat(theme.Color4, false)
		CurrentTheme.Color5 = convertHexFormat(theme.Color5, false)
		CurrentTheme.Color6 = convertHexFormat(theme.Color6, false)
	}

	// Load external theme files
	if err := LoadExternalAccentThemes(); err != nil {
		logging.LogDebug("Warning: Could not load external themes: %v", err)
	}

	logging.LogDebug("Current accent colors initialized: %+v", CurrentTheme)
	return nil
}

// LoadExternalAccentThemes loads accent themes from external files
func LoadExternalAccentThemes() error {
	// Clear the current lists of external themes
	PresetThemes = []ThemeColor{}
	CustomThemes = []ThemeColor{}

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Load preset themes
	presetsDir := filepath.Join(cwd, AccentsDir, PresetsDir)
	logging.LogDebug("Loading preset accent themes from: %s", presetsDir)
	if err := loadThemesFromDir(presetsDir, &PresetThemes); err != nil {
		logging.LogDebug("Warning: Could not load preset themes: %v", err)
	}

	// Load custom themes
	customDir := filepath.Join(cwd, AccentsDir, CustomDir)
	logging.LogDebug("Loading custom accent themes from: %s", customDir)
	if err := loadThemesFromDir(customDir, &CustomThemes); err != nil {
		logging.LogDebug("Warning: Could not load custom themes: %v", err)
	}

	logging.LogDebug("Loaded %d preset and %d custom accent themes", len(PresetThemes), len(CustomThemes))
	return nil
}

// loadThemesFromDir loads themes from a specific directory
func loadThemesFromDir(themesDir string, themesList *[]ThemeColor) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		logging.LogDebug("Error creating themes directory: %v", err)
		return fmt.Errorf("error creating themes directory: %w", err)
	}

	// Read the directory
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		logging.LogDebug("Error reading themes directory: %v", err)
		return fmt.Errorf("error reading themes directory: %w", err)
	}

	// Process each file in the directory
	for _, entry := range entries {
		// Skip directories and hidden files
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip files that don't have a .txt extension
		if !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		// Skip placeholder files
		if strings.Contains(entry.Name(), "Place-") && strings.Contains(entry.Name(), "-Here") {
			continue
		}

		// Extract theme name (remove .txt extension)
		themeName := strings.TrimSuffix(entry.Name(), ".txt")

		// Read the theme file
		theme, err := ReadThemeFile(filepath.Join(themesDir, entry.Name()))
		if err != nil {
			logging.LogDebug("Error reading theme file %s: %v", entry.Name(), err)
			continue
		}

		// Set the theme name
		theme.Name = themeName

		// Add the theme to the list
		*themesList = append(*themesList, *theme)
	}

	return nil
}

// ReadThemeFile reads an accent theme from a file
func ReadThemeFile(filepath string) (*ThemeColor, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open theme file: %w", err)
	}
	defer file.Close()

	// Create a new theme
	theme := &ThemeColor{
		Name: "External Theme", // Default name, will be replaced
	}

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Assign values to the theme
		switch key {
		case "color1":
			theme.Color1 = convertHexFormat(value, false)
		case "color2":
			theme.Color2 = convertHexFormat(value, false)
		case "color3":
			theme.Color3 = convertHexFormat(value, false)
		case "color4":
			theme.Color4 = convertHexFormat(value, false)
		case "color5":
			theme.Color5 = convertHexFormat(value, false)
		case "color6":
			theme.Color6 = convertHexFormat(value, false)
		}
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("error reading theme file: %w", scanner.Err())
	}

	return theme, nil
}

// GetCurrentColors reads the current theme colors from the settings file
func GetCurrentColors() (*ThemeColor, error) {
	logging.LogDebug("Reading current accent colors from: %s", SettingsPath)

	// Check if the file exists
	_, err := os.Stat(SettingsPath)
	if os.IsNotExist(err) {
		logging.LogDebug("Settings file does not exist: %s", SettingsPath)
		return nil, fmt.Errorf("settings file not found: %s", SettingsPath)
	}

	// Read the file
	file, err := os.Open(SettingsPath)
	if err != nil {
		logging.LogDebug("Error opening settings file: %v", err)
		return nil, fmt.Errorf("failed to open settings file: %w", err)
	}
	defer file.Close()

	// Parse the file
	colors := &ThemeColor{
		Name: "Current",
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "color1":
			colors.Color1 = value
		case "color2":
			colors.Color2 = value
		case "color3":
			colors.Color3 = value
		case "color4":
			colors.Color4 = value
		case "color5":
			colors.Color5 = value
		case "color6":
			colors.Color6 = value
		}
	}

	if scanner.Err() != nil {
		logging.LogDebug("Error scanning settings file: %v", scanner.Err())
		return nil, fmt.Errorf("error reading settings file: %w", scanner.Err())
	}

	return colors, nil
}

// UpdateCurrentTheme updates the current theme in memory with new values
func UpdateCurrentTheme(themeName string) error {
	logging.LogDebug("Updating current theme to: %s", themeName)

	// First look in preset themes
	for _, theme := range PresetThemes {
		if theme.Name == themeName {
			CurrentTheme.Name = theme.Name
			CurrentTheme.Color1 = theme.Color1
			CurrentTheme.Color2 = theme.Color2
			CurrentTheme.Color3 = theme.Color3
			CurrentTheme.Color4 = theme.Color4
			CurrentTheme.Color5 = theme.Color5
			CurrentTheme.Color6 = theme.Color6

			logging.LogDebug("Theme updated from preset theme: %+v", CurrentTheme)
			return nil
		}
	}

	// Then look in custom themes
	for _, theme := range CustomThemes {
		if theme.Name == themeName {
			CurrentTheme.Name = theme.Name
			CurrentTheme.Color1 = theme.Color1
			CurrentTheme.Color2 = theme.Color2
			CurrentTheme.Color3 = theme.Color3
			CurrentTheme.Color4 = theme.Color4
			CurrentTheme.Color5 = theme.Color5
			CurrentTheme.Color6 = theme.Color6

			logging.LogDebug("Theme updated from custom theme: %+v", CurrentTheme)
			return nil
		}
	}

	logging.LogDebug("Theme not found: %s", themeName)
	return fmt.Errorf("theme not found: %s", themeName)
}

// ApplyThemeColors applies the specified theme colors to the system
func ApplyThemeColors(theme *ThemeColor) error {
	logging.LogDebug("Applying theme colors: %s", theme.Name)

	// Read the current settings file
	_, err := os.Stat(SettingsPath)
	if os.IsNotExist(err) {
		logging.LogDebug("Settings file does not exist, creating: %s", SettingsPath)

		// Create parent directories if needed
		err = os.MkdirAll(filepath.Dir(SettingsPath), 0755)
		if err != nil {
			logging.LogDebug("Error creating parent directories: %v", err)
			return fmt.Errorf("failed to create settings directory: %w", err)
		}

		// Create an empty file
		file, err := os.Create(SettingsPath)
		if err != nil {
			logging.LogDebug("Error creating settings file: %v", err)
			return fmt.Errorf("failed to create settings file: %w", err)
		}
		file.Close()
	}

	// Read existing settings
	settings := make(map[string]string)

	file, err := os.Open(SettingsPath)
	if err != nil {
		logging.LogDebug("Error opening settings file: %v", err)
		return fmt.Errorf("failed to open settings file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		settings[key] = value
	}
	file.Close()

	// Update with new theme colors - convert from display format to storage format
	settings["color1"] = convertHexFormat(theme.Color1, true)
	settings["color2"] = convertHexFormat(theme.Color2, true)
	settings["color3"] = convertHexFormat(theme.Color3, true)
	settings["color4"] = convertHexFormat(theme.Color4, true)
	settings["color5"] = convertHexFormat(theme.Color5, true)
	settings["color6"] = convertHexFormat(theme.Color6, true)

	// Write back to file
	tempFile := SettingsPath + ".tmp"
	outFile, err := os.Create(tempFile)
	if err != nil {
		logging.LogDebug("Error creating temp settings file: %v", err)
		return fmt.Errorf("failed to create temp settings file: %w", err)
	}

	// Write each line back to the file
	for key, value := range settings {
		_, err := fmt.Fprintf(outFile, "%s=%s\n", key, value)
		if err != nil {
			outFile.Close()
			os.Remove(tempFile)
			logging.LogDebug("Error writing to settings file: %v", err)
			return fmt.Errorf("failed to write settings: %w", err)
		}
	}

	outFile.Close()

	// Replace the original file with the new one
	err = os.Rename(tempFile, SettingsPath)
	if err != nil {
		logging.LogDebug("Error replacing settings file: %v", err)
		return fmt.Errorf("failed to update settings file: %w", err)
	}

	logging.LogDebug("Successfully applied theme colors")
	return nil
}

// GetColorPreviewText formats a color value for display
func GetColorPreviewText(colorName string, colorValue string) string {
	return fmt.Sprintf("%s: %s", colorName, colorValue)
}

// ApplyCurrentTheme applies the current in-memory theme to the system
func ApplyCurrentTheme() error {
	logging.LogDebug("Applying current in-memory theme")
	return ApplyThemeColors(&CurrentTheme)
}

// SaveThemeToFile saves a theme to an external file
func SaveThemeToFile(theme *ThemeColor, fileName string, isCustom bool) error {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to themes directory
	var themesDir string
	if isCustom {
		themesDir = filepath.Join(cwd, AccentsDir, CustomDir)
	} else {
		themesDir = filepath.Join(cwd, AccentsDir, PresetsDir)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		logging.LogDebug("Error creating themes directory: %v", err)
		return fmt.Errorf("error creating themes directory: %w", err)
	}

	// Full path to the file
	filePath := filepath.Join(themesDir, fileName)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		logging.LogDebug("Error creating theme file: %v", err)
		return fmt.Errorf("error creating theme file: %w", err)
	}
	defer file.Close()

	// Write theme colors to the file
	_, err = fmt.Fprintf(file, "color1=%s\n", convertHexFormat(theme.Color1, true))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "color2=%s\n", convertHexFormat(theme.Color2, true))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "color3=%s\n", convertHexFormat(theme.Color3, true))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "color4=%s\n", convertHexFormat(theme.Color4, true))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "color5=%s\n", convertHexFormat(theme.Color5, true))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "color6=%s\n", convertHexFormat(theme.Color6, true))
	if err != nil {
		return err
	}

	logging.LogDebug("Successfully saved theme to file: %s", filePath)
	return nil
}

// CreatePlaceholderFiles creates placeholder files in the theme directories
func CreatePlaceholderFiles() error {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create directories
	presetsDir := filepath.Join(cwd, AccentsDir, PresetsDir)
	customDir := filepath.Join(cwd, AccentsDir, CustomDir)

	if err := os.MkdirAll(presetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create presets directory: %w", err)
	}

	if err := os.MkdirAll(customDir, 0755); err != nil {
		return fmt.Errorf("failed to create custom directory: %w", err)
	}

	// Create placeholder file in custom directory if empty
	entries, err := os.ReadDir(customDir)
	if err != nil {
		return fmt.Errorf("failed to read custom directory: %w", err)
	}

	if len(entries) == 0 {
		placeholderPath := filepath.Join(customDir, "Place-Accent-Files-Here.txt")
		file, err := os.Create(placeholderPath)
		if err != nil {
			return fmt.Errorf("failed to create placeholder file: %w", err)
		}

		_, err = file.WriteString("# Place custom accent theme files in this directory\n\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to write placeholder content: %w", err)
		}

		_, err = file.WriteString("# Format should be:\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to write placeholder content: %w", err)
		}

		_, err = file.WriteString("color1=0xRRGGBB\ncolor2=0xRRGGBB\ncolor3=0xRRGGBB\ncolor4=0xRRGGBB\ncolor5=0xRRGGBB\ncolor6=0xRRGGBB\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to write placeholder content: %w", err)
		}

		file.Close()
	}

	return nil
}

// GetCurrentAccentColors retrieves the current accent colors as a map
func GetCurrentAccentColors() (map[string]string, error) {
	// Get current colors as a ThemeColor struct
	colors, err := GetCurrentColors()
	if err != nil {
		return nil, err
	}

	// Convert to a map with storage format (0xRRGGBB)
	colorMap := map[string]string{
		"color1": convertHexFormat(colors.Color1, true),
		"color2": convertHexFormat(colors.Color2, true),
		"color3": convertHexFormat(colors.Color3, true),
		"color4": convertHexFormat(colors.Color4, true),
		"color5": convertHexFormat(colors.Color5, true),
		"color6": convertHexFormat(colors.Color6, true),
	}

	return colorMap, nil
}

// ApplyAccentColors applies accent colors from a map to the system
func ApplyAccentColors(colorMap map[string]string) error {
	// Convert to a ThemeColor struct
	theme := &ThemeColor{
		Name: "ImportedTheme",
	}

	// Set colors, converting from storage format to display format
	if color, ok := colorMap["color1"]; ok {
		theme.Color1 = convertHexFormat(color, false)
	}
	if color, ok := colorMap["color2"]; ok {
		theme.Color2 = convertHexFormat(color, false)
	}
	if color, ok := colorMap["color3"]; ok {
		theme.Color3 = convertHexFormat(color, false)
	}
	if color, ok := colorMap["color4"]; ok {
		theme.Color4 = convertHexFormat(color, false)
	}
	if color, ok := colorMap["color5"]; ok {
		theme.Color5 = convertHexFormat(color, false)
	}
	if color, ok := colorMap["color6"]; ok {
		theme.Color6 = convertHexFormat(color, false)
	}

	// Apply the theme
	return ApplyThemeColors(theme)
}

// ParseAccentColors parses accent colors from a settings file
func ParseAccentColors(filePath string) (map[string]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open accent settings file: %w", err)
	}
	defer file.Close()

	// Read and parse the file
	colorMap := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only process color settings
		if strings.HasPrefix(key, "color") {
			colorMap[key] = value
		}
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("error reading accent settings file: %w", scanner.Err())
	}

	return colorMap, nil
}
