// src/internal/themes/component_import.go
// Implementation of component-specific import functionality

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nextui-themes/internal/accents"
	"nextui-themes/internal/fonts"
	"nextui-themes/internal/system"
	"nextui-themes/internal/ui"
)

// ComponentType represents different component types that can be imported
type ComponentType int

const (
	ComponentTypeFullTheme ComponentType = iota + 1
	ComponentTypeAccent
	ComponentTypeLED
	ComponentTypeWallpaper
	ComponentTypeIcon
	ComponentTypeFont
)

// ImportComponent imports a specific component type
func ImportComponent(componentType ComponentType, itemName string, selectedComponents map[ComponentType]bool) error {
	// Create logging directory if it doesn't exist
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	logsDir := filepath.Join(cwd, "Logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("error creating logs directory: %w", err)
	}

	// Create log file
	logFile, err := os.OpenFile(
		filepath.Join(logsDir, "component_imports.log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()

	// Create logger
	logger := &Logger{logFile}
	logger.Printf("Starting component import: Type=%d, Item=%s", componentType, itemName)

	// Handle different component types
	switch componentType {
	case ComponentTypeFullTheme:
		// For full themes, check if we're importing all components or specific ones
		if len(selectedComponents) == 0 {
			// Import all components (legacy mode)
			return ImportTheme(itemName)
		}

		// Import specific components from the theme
		return ImportThemeComponents(itemName, selectedComponents, logger)

	case ComponentTypeAccent:
		return ImportAccentPack(itemName, logger)

	case ComponentTypeLED:
		return ImportLEDPack(itemName, logger)

	case ComponentTypeWallpaper:
		return ImportWallpaperPack(itemName, logger)

	case ComponentTypeIcon:
		return ImportIconPack(itemName, logger)

	case ComponentTypeFont:
		return ImportFontPack(itemName, logger)

	default:
		return fmt.Errorf("unsupported component type: %d", componentType)
	}
}

// ImportThemeComponents imports specific components from a theme
func ImportThemeComponents(themeName string, selectedComponents map[ComponentType]bool, logger *Logger) error {
	logger.Printf("Importing specific components from theme: %s", themeName)

	// Full path to theme - look in Themes/Imports directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	themePath := filepath.Join(cwd, "Themes", "Imports", themeName)

	// Validate theme
	manifest, err := ValidateTheme(themePath, logger)
	if err != nil {
		logger.Printf("Theme validation failed: %v", err)
		return fmt.Errorf("theme validation failed: %w", err)
	}

	// Import selected components
	if selectedComponents[ComponentTypeWallpaper] {
		logger.Printf("Importing wallpapers from theme")
		if err := ImportWallpapers(themePath, manifest, logger); err != nil {
			logger.Printf("Error importing wallpapers: %v", err)
			// Continue with other components
		}
	}

	if selectedComponents[ComponentTypeIcon] {
		logger.Printf("Importing icons from theme")
		if err := ImportIcons(themePath, manifest, logger); err != nil {
			logger.Printf("Error importing icons: %v", err)
			// Continue with other components
		}
	}

	if selectedComponents[ComponentTypeFont] {
		logger.Printf("Importing fonts from theme")
		if err := ImportFonts(themePath, manifest, logger); err != nil {
			logger.Printf("Error importing fonts: %v", err)
			// Continue with other components
		}
	}

	if selectedComponents[ComponentTypeAccent] {
		logger.Printf("Importing accent settings from theme")
		// Extract just accent settings from the theme
		if manifest.Content.Settings.AccentsIncluded {
			if err := ImportSettings(themePath, manifest, logger); err != nil {
				logger.Printf("Error importing settings: %v", err)
				// Continue with other components
			}
		} else {
			logger.Printf("No accent settings found in theme")
		}
	}

	if selectedComponents[ComponentTypeLED] {
		logger.Printf("Importing LED settings from theme")
		// Extract just LED settings from the theme
		if manifest.Content.Settings.LEDsIncluded {
			if err := ImportSettings(themePath, manifest, logger); err != nil {
				logger.Printf("Error importing settings: %v", err)
				// Continue with other components
			}
		} else {
			logger.Printf("No LED settings found in theme")
		}
	}

	logger.Printf("Component import completed")
	return nil
}

// ImportAccentPack imports an accent pack (.acc)
func ImportAccentPack(packName string, logger *Logger) error {
	logger.Printf("Importing accent pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to accent pack
	packPath := filepath.Join(cwd, "Accents", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("accent pack not found: %s", packPath)
	}

	// Find settings file in the pack - should be minuisettings.txt
	settingsPath := filepath.Join(packPath, "minuisettings.txt")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return fmt.Errorf("settings file not found in accent pack: %s", settingsPath)
	}

	// Copy settings file to system location
	systemSettingsPath := accents.SettingsPath

	// Create system settings directory if it doesn't exist
	systemSettingsDir := filepath.Dir(systemSettingsPath)
	if err := os.MkdirAll(systemSettingsDir, 0755); err != nil {
		return fmt.Errorf("failed to create system settings directory: %w", err)
	}

	// Copy settings file
	if err := CopyFile(settingsPath, systemSettingsPath); err != nil {
		return fmt.Errorf("failed to copy accent settings: %w", err)
	}

	logger.Printf("Successfully imported accent pack: %s", packName)
	ui.ShowMessage(fmt.Sprintf("Accent pack '%s' imported successfully", packName), "3")

	return nil
}

// ImportLEDPack imports a LED pack (.led)
func ImportLEDPack(packName string, logger *Logger) error {
	logger.Printf("Importing LED pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to LED pack
	packPath := filepath.Join(cwd, "LEDs", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("LED pack not found: %s", packPath)
	}

	// Find settings file in the pack - should be ledsettings_brick.txt
	settingsPath := filepath.Join(packPath, "ledsettings_brick.txt")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return fmt.Errorf("settings file not found in LED pack: %s", settingsPath)
	}

	// Copy settings file to system location
	systemSettingsPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"

	// Create system settings directory if it doesn't exist
	systemSettingsDir := filepath.Dir(systemSettingsPath)
	if err := os.MkdirAll(systemSettingsDir, 0755); err != nil {
		return fmt.Errorf("failed to create system settings directory: %w", err)
	}

	// Copy settings file
	if err := CopyFile(settingsPath, systemSettingsPath); err != nil {
		return fmt.Errorf("failed to copy LED settings: %w", err)
	}

	logger.Printf("Successfully imported LED pack: %s", packName)
	ui.ShowMessage(fmt.Sprintf("LED pack '%s' imported successfully", packName), "3")

	return nil
}

// ImportWallpaperPack imports a wallpaper pack (.bg)
func ImportWallpaperPack(packName string, logger *Logger) error {
	logger.Printf("Importing wallpaper pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to wallpaper pack
	packPath := filepath.Join(cwd, "Wallpapers", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("wallpaper pack not found: %s", packPath)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Create a minimal manifest for wallpaper import
	manifest := &ThemeManifest{}
	manifest.Content.Wallpapers.Present = true
	manifest.PathMappings.Wallpapers = []PathMapping{}

	// First check SystemWallpapers directory
	sysWallDir := filepath.Join(packPath, "SystemWallpapers")
	if entries, err := os.ReadDir(sysWallDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				// Add to manifest
				baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

				// Try to match to system path based on name
				var dstPath string
				var metadata map[string]string

				if baseName == "Root" {
					dstPath = filepath.Join(systemPaths.Root, "bg.png")
					metadata = map[string]string{
						"SystemName":    "Root",
						"WallpaperType": "Main",
					}
				} else if baseName == "Root-Media" {
					dstPath = filepath.Join(systemPaths.Root, ".media", "bg.png")
					metadata = map[string]string{
						"SystemName":    "Root",
						"WallpaperType": "Media",
					}
				} else if baseName == "Recently Played" {
					dstPath = filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
					metadata = map[string]string{
						"SystemName":    "Recently Played",
						"WallpaperType": "Media",
					}
				} else if baseName == "Tools" {
					dstPath = filepath.Join(systemPaths.Tools, ".media", "bg.png")
					metadata = map[string]string{
						"SystemName":    "Tools",
						"WallpaperType": "Media",
					}
				} else if baseName == "Collections" {
					dstPath = filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")
					metadata = map[string]string{
						"SystemName":    "Collections",
						"WallpaperType": "Media",
					}
				} else {
					// Try to match by system tag
					for _, sys := range systemPaths.Systems {
						if strings.Contains(baseName, fmt.Sprintf("(%s)", sys.Tag)) {
							dstPath = filepath.Join(sys.MediaPath, "bg.png")
							metadata = map[string]string{
								"SystemName":    sys.Name,
								"SystemTag":     sys.Tag,
								"WallpaperType": "System",
							}
							break
						}
					}
				}

				if dstPath != "" {
					manifest.PathMappings.Wallpapers = append(manifest.PathMappings.Wallpapers, PathMapping{
						ThemePath:  filepath.Join("SystemWallpapers", entry.Name()),
						SystemPath: dstPath,
						Metadata:   metadata,
					})
				}
			}
		}
	}

	// Also check CollectionWallpapers directory
	collWallDir := filepath.Join(packPath, "CollectionWallpapers")
	if entries, err := os.ReadDir(collWallDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				// Add to manifest
				baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
				dstPath := filepath.Join(systemPaths.Root, "Collections", baseName, ".media", "bg.png")

				manifest.PathMappings.Wallpapers = append(manifest.PathMappings.Wallpapers, PathMapping{
					ThemePath:  filepath.Join("CollectionWallpapers", entry.Name()),
					SystemPath: dstPath,
					Metadata: map[string]string{
						"CollectionName": baseName,
						"WallpaperType":  "Collection",
					},
				})
			}
		}
	}

	// Now import the wallpapers
	if err := ImportWallpapers(packPath, manifest, logger); err != nil {
		return fmt.Errorf("error importing wallpapers: %w", err)
	}

	logger.Printf("Successfully imported wallpaper pack: %s", packName)
	ui.ShowMessage(fmt.Sprintf("Wallpaper pack '%s' imported successfully", packName), "3")

	return nil
}

// ImportIconPack imports an icon pack (.icon)
func ImportIconPack(packName string, logger *Logger) error {
	logger.Printf("Importing icon pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to icon pack
	packPath := filepath.Join(cwd, "Icons", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("icon pack not found: %s", packPath)
	}

	// Create a minimal manifest for icon import
	manifest := &ThemeManifest{}
	manifest.Content.Icons.Present = true
	manifest.PathMappings.Icons = []PathMapping{}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Check SystemIcons directory
	sysIconsDir := filepath.Join(packPath, "SystemIcons")
	if entries, err := os.ReadDir(sysIconsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				// Add to manifest
				baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

				// Try to match to system path based on name
				var dstPath string
				var metadata map[string]string

				if baseName == "Collections" {
					dstPath = filepath.Join(systemPaths.Root, ".media", "Collections.png")
					metadata = map[string]string{
						"SystemName": "Collections",
						"IconType":   "Special",
					}
				} else if baseName == "Recently Played" {
					dstPath = filepath.Join(systemPaths.Root, ".media", "Recently Played.png")
					metadata = map[string]string{
						"SystemName": "Recently Played",
						"IconType":   "Special",
					}
				} else if baseName == "Tools" {
					toolsBaseDir := filepath.Dir(systemPaths.Tools)
					dstPath = filepath.Join(toolsBaseDir, ".media", "tg5040.png")
					metadata = map[string]string{
						"SystemName": "Tools",
						"IconType":   "Special",
					}
				} else {
					// Try to match by system tag
					for _, sys := range systemPaths.Systems {
						if strings.Contains(baseName, fmt.Sprintf("(%s)", sys.Tag)) {
							dstPath = filepath.Join(systemPaths.Roms, ".media", fmt.Sprintf("%s.png", sys.Name))
							metadata = map[string]string{
								"SystemName": sys.Name,
								"SystemTag":  sys.Tag,
								"IconType":   "System",
							}
							break
						}
					}
				}

				if dstPath != "" {
					manifest.PathMappings.Icons = append(manifest.PathMappings.Icons, PathMapping{
						ThemePath:  filepath.Join("SystemIcons", entry.Name()),
						SystemPath: dstPath,
						Metadata:   metadata,
					})
				}
			}
		}
	}

	// Check ToolIcons directory
	toolIconsDir := filepath.Join(packPath, "ToolIcons")
	if entries, err := os.ReadDir(toolIconsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				// Tool icons go in the Tools/.media directory
				dstPath := filepath.Join(systemPaths.Tools, ".media", entry.Name())
				baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

				manifest.PathMappings.Icons = append(manifest.PathMappings.Icons, PathMapping{
					ThemePath:  filepath.Join("ToolIcons", entry.Name()),
					SystemPath: dstPath,
					Metadata: map[string]string{
						"ToolName": baseName,
						"IconType": "Tool",
					},
				})
			}
		}
	}

	// Check CollectionIcons directory
	collIconsDir := filepath.Join(packPath, "CollectionIcons")
	if entries, err := os.ReadDir(collIconsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				// Collection icons go in the Collections/.media directory
				dstPath := filepath.Join(systemPaths.Root, "Collections", ".media", entry.Name())
				baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

				manifest.PathMappings.Icons = append(manifest.PathMappings.Icons, PathMapping{
					ThemePath:  filepath.Join("CollectionIcons", entry.Name()),
					SystemPath: dstPath,
					Metadata: map[string]string{
						"CollectionName": baseName,
						"IconType":       "Collection",
					},
				})
			}
		}
	}

	// Now import the icons
	if err := ImportIcons(packPath, manifest, logger); err != nil {
		return fmt.Errorf("error importing icons: %w", err)
	}

	logger.Printf("Successfully imported icon pack: %s", packName)
	ui.ShowMessage(fmt.Sprintf("Icon pack '%s' imported successfully", packName), "3")

	return nil
}

// ImportFontPack imports a font pack (.font)
func ImportFontPack(packName string, logger *Logger) error {
	logger.Printf("Importing font pack: %s", packName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to font pack
	packPath := filepath.Join(cwd, "Fonts", "Imports", packName)

	// Ensure directory exists
	_, err = os.Stat(packPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("font pack not found: %s", packPath)
	}

	// Create a minimal manifest for font import
	manifest := &ThemeManifest{}
	manifest.Content.Fonts.Present = true
	manifest.PathMappings.Fonts = make(map[string]PathMapping)

	// Check for OG font
	ogFontPath := filepath.Join(packPath, "OG.ttf")
	if _, err := os.Stat(ogFontPath); err == nil {
		manifest.PathMappings.Fonts["og_font"] = PathMapping{
			ThemePath:  "OG.ttf",
			SystemPath: fonts.OGFontPath,
		}
	}

	// Check for OG font backup
	ogBackupPath := filepath.Join(packPath, "OG.backup.ttf")
	if _, err := os.Stat(ogBackupPath); err == nil {
		manifest.PathMappings.Fonts["og_backup"] = PathMapping{
			ThemePath:  "OG.backup.ttf",
			SystemPath: filepath.Join(filepath.Dir(fonts.OGFontPath), fonts.OGFontBackupName),
		}
		manifest.Content.Fonts.OGReplaced = true
	}

	// Check for Next font
	nextFontPath := filepath.Join(packPath, "Next.ttf")
	if _, err := os.Stat(nextFontPath); err == nil {
		manifest.PathMappings.Fonts["next_font"] = PathMapping{
			ThemePath:  "Next.ttf",
			SystemPath: fonts.NextFontPath,
		}
	}

	// Check for Next font backup
	nextBackupPath := filepath.Join(packPath, "Next.backup.ttf")
	if _, err := os.Stat(nextBackupPath); err == nil {
		manifest.PathMappings.Fonts["next_backup"] = PathMapping{
			ThemePath:  "Next.backup.ttf",
			SystemPath: filepath.Join(filepath.Dir(fonts.NextFontPath), fonts.NextFontBackupName),
		}
		manifest.Content.Fonts.NextReplaced = true
	}

	// Now import the fonts
	if err := ImportFonts(packPath, manifest, logger); err != nil {
		return fmt.Errorf("error importing fonts: %w", err)
	}

	logger.Printf("Successfully imported font pack: %s", packName)
	ui.ShowMessage(fmt.Sprintf("Font pack '%s' imported successfully", packName), "3")

	return nil
}

// Connect the UI's performImport function to our backend implementation
func PerformImport(componentType ComponentType, itemName string, selectedComponents map[ComponentType]bool) error {
	return ImportComponent(componentType, itemName, selectedComponents)
}
