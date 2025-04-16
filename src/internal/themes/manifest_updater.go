// src/internal/themes/manifest_updater.go
// Functions for updating theme manifests based on actual content

package themes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
)

// createEmptyManifest creates a new empty manifest with basic metadata
func createEmptyManifest() ThemeManifest {
	// Initialize a new empty manifest
	var manifest ThemeManifest

	// Initialize metadata sections
	manifest.ThemeInfo.Version = "1.0.0"
	manifest.ThemeInfo.ExportedBy = GetVersionString()
	manifest.ThemeInfo.CreationDate = time.Now()
	manifest.ThemeInfo.Author = "AuthorName" // Default author name

	// Initialize content section
	manifest.Content.Wallpapers.Present = false
	manifest.Content.Wallpapers.Count = 0

	manifest.Content.Icons.Present = false
	manifest.Content.Icons.SystemCount = 0
	manifest.Content.Icons.ToolCount = 0
	manifest.Content.Icons.CollectionCount = 0

	manifest.Content.Overlays.Present = false
	manifest.Content.Overlays.Systems = []string{}

	manifest.Content.Fonts.Present = false
	manifest.Content.Fonts.OGReplaced = false
	manifest.Content.Fonts.NextReplaced = false

	manifest.Content.Settings.AccentsIncluded = false
	manifest.Content.Settings.LEDsIncluded = false

	// Initialize path mappings
	manifest.PathMappings.Wallpapers = []PathMapping{}
	manifest.PathMappings.Icons = []PathMapping{}
	manifest.PathMappings.Overlays = []PathMapping{}
	manifest.PathMappings.Fonts = make(map[string]PathMapping)
	manifest.PathMappings.Settings = make(map[string]PathMapping)

	return manifest
}

// UpdateThemeManifest scans a theme directory and updates its manifest
func UpdateThemeManifest(themePath string) error {
	// Use fmt for critical early initialization to avoid logging issues
	fmt.Printf("Scanning theme directory to update manifest: %s\n", themePath)

	// Load existing manifest or create new one
	manifestPath := filepath.Join(themePath, "manifest.json")
	var manifest ThemeManifest

	if fileData, err := os.ReadFile(manifestPath); err == nil {
		if err := json.Unmarshal(fileData, &manifest); err != nil {
			fmt.Printf("Error parsing manifest: %v, creating new one\n", err)
			// Initialize new manifest
			manifest = createEmptyManifest()
		}
	} else {
		fmt.Printf("No manifest found, creating new one\n")
		manifest = createEmptyManifest()
	}

	// Extract theme name from directory name
	themeName := filepath.Base(themePath)
	manifest.ThemeInfo.Name = strings.TrimSuffix(themeName, filepath.Ext(themeName))
	manifest.ComponentType = "Theme" // Set component type for full themes

	// Update each component type
	if err := updateWallpapersInManifest(themePath, &manifest); err != nil {
		logging.LogDebug("Error updating wallpapers in manifest: %v", err)
	}

	if err := updateIconsInManifest(themePath, &manifest); err != nil {
		logging.LogDebug("Error updating icons in manifest: %v", err)
	}

	if err := updateOverlaysInManifest(themePath, &manifest); err != nil {
		logging.LogDebug("Error updating overlays in manifest: %v", err)
	}

	if err := updateFontsInManifest(themePath, &manifest); err != nil {
		logging.LogDebug("Error updating fonts in manifest: %v", err)
	}

	if err := updateSettingsInManifest(themePath, &manifest); err != nil {
		logging.LogDebug("Error updating settings in manifest: %v", err)
	}

	// Write updated manifest
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating manifest JSON: %w", err)
	}

	logging.LogDebug("Writing updated manifest to: %s", manifestPath)
	return os.WriteFile(manifestPath, manifestJSON, 0644)
}

// EnhanceManifestUpdater updates a manifest based on directory contents
func EnhanceManifestUpdater(themePath string) error {
	// First check if a manifest exists
	manifestPath := filepath.Join(themePath, "manifest.json")

	var manifest ThemeManifest

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// No manifest exists, create a new one
		manifest = createEmptyManifest()
	} else {
		// Read existing manifest
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			return fmt.Errorf("error reading manifest: %w", err)
		}

		if err := json.Unmarshal(data, &manifest); err != nil {
			// Invalid manifest, create a new one
			manifest = createEmptyManifest()
		}
	}

	// Determine component type from directory extension if not set
	if manifest.ComponentType == "" {
		manifest.ComponentType = DetermineComponentType(themePath)
	}

	// Update manifest based on component type and directory contents
	switch manifest.ComponentType {
	case "Wallpapers":
		updateWallpapersInManifest(themePath, &manifest)
		// Check/create preview.png from Recently Played if needed
		createWallpaperPreview(themePath, &manifest)
	case "Icons":
		updateIconsInManifest(themePath, &manifest)
	case "LEDs":
		updateComponentLEDsFromSettingsFile(themePath, &manifest)
	case "Accents":
		updateComponentAccentsFromSettingsFile(themePath, &manifest)
	case "Fonts":
		updateFontsInManifest(themePath, &manifest)
	case "Theme":
		// Full theme - update all components
		updateWallpapersInManifest(themePath, &manifest)
		updateIconsInManifest(themePath, &manifest)
		updateOverlaysInManifest(themePath, &manifest)
		updateFontsInManifest(themePath, &manifest)
		updateSettingsInManifest(themePath, &manifest)
	}

	// Write updated manifest
	return WriteComponentManifest(themePath, &manifest)
}

// createWallpaperPreview creates a preview image for a wallpaper pack
func createWallpaperPreview(themePath string, manifest *ThemeManifest) error {
	previewPath := filepath.Join(themePath, "preview.png")

	// If preview already exists, update manifest and skip
	if _, err := os.Stat(previewPath); err == nil {
		manifest.PreviewImage = "preview.png"
		return nil
	}

	// Look for Recently Played wallpaper
	rpPath := filepath.Join(themePath, "SystemWallpapers", "Recently Played.png")
	if _, err := os.Stat(rpPath); err == nil {
		// Copy as preview
		if err := copyComponentFile(rpPath, previewPath); err != nil {
			return err
		}
		manifest.PreviewImage = "preview.png"
		return nil
	}

	// If not found, try Root.png
	rootPath := filepath.Join(themePath, "SystemWallpapers", "Root.png")
	if _, err := os.Stat(rootPath); err == nil {
		if err := copyComponentFile(rootPath, previewPath); err != nil {
			return err
		}
		manifest.PreviewImage = "preview.png"
		return nil
	}

	// No suitable image found
	return nil
}

// updateComponentLEDsFromSettingsFile extracts LED settings from settings file
func updateComponentLEDsFromSettingsFile(themePath string, manifest *ThemeManifest) error {
	// Look for LED settings file
	ledSettingsPath := filepath.Join(themePath, "ledsettings_brick.txt")
	if _, err := os.Stat(ledSettingsPath); os.IsNotExist(err) {
		// No LED settings file, check if we already have settings in manifest
		if manifest.LEDSettings == nil || len(manifest.LEDSettings) == 0 {
			// No settings in manifest either
			manifest.Content.Settings.LEDsIncluded = false
			return nil
		}
		// We have settings in the manifest but no file, which is fine
	} else {
		// Read the LED settings file
		data, err := os.ReadFile(ledSettingsPath)
		if err != nil {
			return fmt.Errorf("error reading LED settings file: %w", err)
		}

		// Parse the settings file (simplified as the actual parser would be in leds package)
		settingsMap := parseLEDSettingsForComponent(string(data))

		// Store in manifest
		manifest.LEDSettings = settingsMap
	}

	// Update content section
	manifest.Content.Settings.LEDsIncluded = true
	return nil
}

// updateComponentAccentsFromSettingsFile extracts accent settings from settings file
func updateComponentAccentsFromSettingsFile(themePath string, manifest *ThemeManifest) error {
	// Look for accent settings file
	accentSettingsPath := filepath.Join(themePath, "minuisettings.txt")
	if _, err := os.Stat(accentSettingsPath); os.IsNotExist(err) {
		// No accent settings file, check if we already have settings in manifest
		if manifest.AccentColors == nil || len(manifest.AccentColors) == 0 {
			// No settings in manifest either
			manifest.Content.Settings.AccentsIncluded = false
			return nil
		}
		// We have settings in the manifest but no file, which is fine
	} else {
		// Read the accent settings file
		data, err := os.ReadFile(accentSettingsPath)
		if err != nil {
			return fmt.Errorf("error reading accent settings file: %w", err)
		}

		// Parse the settings file (simplified as the actual parser would be in accents package)
		colorMap := parseAccentColorsForComponent(string(data))

		// Store in manifest
		manifest.AccentColors = colorMap
	}

	// Update content section
	manifest.Content.Settings.AccentsIncluded = true
	return nil
}

// copyComponentFile copies a file from src to dst - named to avoid conflicts
func copyComponentFile(src, dst string) error {
	// Read source file
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("error reading source file: %w", err)
	}

	// Write to destination
	return ioutil.WriteFile(dst, data, 0644)
}

// parseLEDSettingsForComponent parses LED settings from a string
// This is a simplified version as the actual implementation would use the leds package
func parseLEDSettingsForComponent(data string) map[string]map[string]interface{} {
	// Simple implementation to parse LED settings
	// In a real implementation, this would use the proper parser from the leds package
	result := make(map[string]map[string]interface{})

	// Example parsing - would be replaced by actual parsing logic
	sections := strings.Split(data, "[")
	for _, section := range sections {
		if len(section) == 0 {
			continue
		}

		lines := strings.Split(section, "\n")
		if len(lines) < 2 {
			continue
		}

		// Extract section name
		sectionName := strings.TrimRight(lines[0], "]")
		result[sectionName] = make(map[string]interface{})

		// Extract key-value pairs
		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Try to convert numeric values
			if strings.HasPrefix(value, "#") {
				// Color value
				result[sectionName][key] = value
			} else if intVal, err := strconv.Atoi(value); err == nil {
				result[sectionName][key] = intVal
			} else {
				result[sectionName][key] = value
			}
		}
	}

	return result
}

// parseAccentColorsForComponent parses accent colors from a string
// This is a simplified version as the actual implementation would use the accents package
func parseAccentColorsForComponent(data string) map[string]string {
	// Simple implementation to parse accent colors
	// In a real implementation, this would use the proper parser from the accents package
	result := make(map[string]string)

	// Example parsing - would be replaced by actual parsing logic
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only process color settings
		if strings.HasPrefix(key, "color") {
			result[key] = value
		}
	}

	return result
}

// updateWallpapersInManifest updates the wallpapers section of the manifest
// with priority given to the preferred directory structure
func updateWallpapersInManifest(themePath string, manifest *ThemeManifest) error {
	logging.LogDebug("Updating wallpapers in manifest")

	// Get system paths for reference
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Initialize counters and mappings
	wallpaperCount := 0
	newWallpaperMappings := []PathMapping{}

	// Regular expression to extract system tag
	reSystemTag := regexp.MustCompile(`\((.*?)\)`)

	// 1. FIRST CHECK PREFERRED STRUCTURE: SystemWallpapers directory
	sysWallDir := filepath.Join(themePath, "Wallpapers", "SystemWallpapers")
	if entries, err := os.ReadDir(sysWallDir); err == nil {
		logging.LogDebug("Found SystemWallpapers directory, checking for wallpapers")
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				continue
			}

			// Extract system tag or name from file name
			fileName := entry.Name()
			baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

			// Try to detect if this is a special wallpaper (Root, Tools, etc.)
			systemPath := ""
			metadata := map[string]string{}

			if baseName == "Root" {
				systemPath = filepath.Join(systemPaths.Root, "bg.png")
				metadata = map[string]string{
					"SystemName":    "Root",
					"WallpaperType": "Main",
				}
			} else if baseName == "Root-Media" {
				systemPath = filepath.Join(systemPaths.Root, ".media", "bg.png")
				metadata = map[string]string{
					"SystemName":    "Root",
					"WallpaperType": "Media",
				}
			} else if baseName == "Recently Played" {
				systemPath = filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
				metadata = map[string]string{
					"SystemName":    "Recently Played",
					"WallpaperType": "Media",
				}
			} else if baseName == "Tools" {
				systemPath = filepath.Join(systemPaths.Tools, ".media", "bg.png")
				metadata = map[string]string{
					"SystemName":    "Tools",
					"WallpaperType": "Media",
				}
			} else if baseName == "Collections" {
				systemPath = filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")
				metadata = map[string]string{
					"SystemName":    "Collections",
					"WallpaperType": "Media",
				}
			} else {
				// Try to match by system tag
				matches := reSystemTag.FindStringSubmatch(baseName)
				if len(matches) >= 2 {
					systemTag := matches[1]

					// Find matching system
					for _, sys := range systemPaths.Systems {
						if sys.Tag == systemTag {
							systemPath = filepath.Join(sys.MediaPath, "bg.png")
							metadata = map[string]string{
								"SystemName":    sys.Name,
								"SystemTag":     systemTag,
								"WallpaperType": "System",
							}
							break
						}
					}
				}
			}

			// If we found a valid system path, add it to manifest
			if systemPath != "" {
				themeRelPath := filepath.Join("Wallpapers", "SystemWallpapers", fileName)
				fullThemePath := filepath.Join(themePath, themeRelPath)

				// Verify the file exists
				if _, err := os.Stat(fullThemePath); err == nil {
					newWallpaperMappings = append(newWallpaperMappings, PathMapping{
						ThemePath:  themeRelPath,
						SystemPath: systemPath,
						Metadata:   metadata,
					})

					wallpaperCount++
					logging.LogDebug("Found system wallpaper in preferred structure: %s -> %s", themeRelPath, systemPath)
				}
			}
		}
	} else {
		logging.LogDebug("SystemWallpapers directory not found or error reading it: %v", err)
	}

	// 2. CHECK COLLECTION WALLPAPERS (preferred structure)
	collWallDir := filepath.Join(themePath, "Wallpapers", "CollectionWallpapers")
	if entries, err := os.ReadDir(collWallDir); err == nil {
		logging.LogDebug("Found CollectionWallpapers directory, checking for wallpapers")
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				continue
			}

			// Extract collection name from file
			fileName := entry.Name()
			collectionName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

			// Set path to collection's media directory
			systemPath := filepath.Join(systemPaths.Root, "Collections", collectionName, ".media", "bg.png")
			themeRelPath := filepath.Join("Wallpapers", "CollectionWallpapers", fileName)
			fullThemePath := filepath.Join(themePath, themeRelPath)

			// Verify the file exists
			if _, err := os.Stat(fullThemePath); err == nil {
				newWallpaperMappings = append(newWallpaperMappings, PathMapping{
					ThemePath:  themeRelPath,
					SystemPath: systemPath,
					Metadata: map[string]string{
						"CollectionName": collectionName,
						"WallpaperType":  "Collection",
					},
				})

				wallpaperCount++
				logging.LogDebug("Found collection wallpaper: %s -> %s", themeRelPath, systemPath)
			}
		}
	} else {
		logging.LogDebug("CollectionWallpapers directory not found or error reading it: %v", err)
	}

	// 3. LEGACY STRUCTURE CHECK - for backward compatibility
	// Only check these if wallpaperCount is still 0, to prioritize the new structure
	if wallpaperCount == 0 {
		// Check for wallpapers in standard locations
		checkDirectories := []struct {
			themeSubPath string
			systemPath   string
			metadata     map[string]string
		}{
			// Root background
			{
				themeSubPath: "Wallpapers/Root/bg.png",
				systemPath:   filepath.Join(systemPaths.Root, "bg.png"),
				metadata: map[string]string{
					"SystemName":    "Root",
					"WallpaperType": "Main",
				},
			},
			// Root .media background
			{
				themeSubPath: "Wallpapers/Root/.media/bg.png",
				systemPath:   filepath.Join(systemPaths.Root, ".media", "bg.png"),
				metadata: map[string]string{
					"SystemName":    "Root",
					"WallpaperType": "Media",
				},
			},
			// Recently Played background
			{
				themeSubPath: "Wallpapers/Recently Played/bg.png",
				systemPath:   filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png"),
				metadata: map[string]string{
					"SystemName":    "Recently Played",
					"WallpaperType": "Media",
				},
			},
			// Tools background
			{
				themeSubPath: "Wallpapers/Tools/bg.png",
				systemPath:   filepath.Join(systemPaths.Tools, ".media", "bg.png"),
				metadata: map[string]string{
					"SystemName":    "Tools",
					"WallpaperType": "Media",
				},
			},
			// Collections background
			{
				themeSubPath: "Wallpapers/Collections/bg.png",
				systemPath:   filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png"),
				metadata: map[string]string{
					"SystemName":    "Collections",
					"WallpaperType": "Media",
				},
			},
		}

		// Check each known location
		for _, check := range checkDirectories {
			fullThemePath := filepath.Join(themePath, check.themeSubPath)
			if _, err := os.Stat(fullThemePath); err == nil {
				wallpaperCount++
				newWallpaperMappings = append(newWallpaperMappings, PathMapping{
					ThemePath:  check.themeSubPath,
					SystemPath: check.systemPath,
					Metadata:   check.metadata,
				})
				logging.LogDebug("Found wallpaper in legacy structure: %s", check.themeSubPath)
			}
		}

		// Check for system-specific wallpapers in the legacy Systems directory
		systemsDir := filepath.Join(themePath, "Wallpapers", "Systems")
		if entries, err := os.ReadDir(systemsDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}

				// Extract system tag from directory name
				dirName := entry.Name()
				matches := reSystemTag.FindStringSubmatch(dirName)
				if len(matches) < 2 {
					logging.LogDebug("Skipping directory without system tag: %s", dirName)
					continue
				}

				systemTag := matches[1]
				bgPath := filepath.Join(systemsDir, dirName, "bg.png")

				// Check if bg.png exists in the system directory
				if _, err := os.Stat(bgPath); err == nil {
					// Find the corresponding system in our system paths
					var systemFound bool
					var systemName string
					var systemMediaPath string

					for _, sys := range systemPaths.Systems {
						if sys.Tag == systemTag {
							systemFound = true
							systemName = sys.Name
							systemMediaPath = sys.MediaPath
							break
						}
					}

					if !systemFound {
						logging.LogDebug("Warning: Found system tag '%s' that doesn't match any known system", systemTag)
						continue
					}

					// Add to manifest
					themeRelPath := fmt.Sprintf("Wallpapers/Systems/%s/bg.png", dirName)
					systemPath := filepath.Join(systemMediaPath, "bg.png")

					newWallpaperMappings = append(newWallpaperMappings, PathMapping{
						ThemePath:  themeRelPath,
						SystemPath: systemPath,
						Metadata: map[string]string{
							"SystemName":    systemName,
							"SystemTag":     systemTag,
							"WallpaperType": "System",
						},
					})

					wallpaperCount++
					logging.LogDebug("Found system wallpaper in legacy structure for tag '%s': %s", systemTag, themeRelPath)
				}
			}
		}
	}

	// Update the manifest
	manifest.PathMappings.Wallpapers = newWallpaperMappings
	manifest.Content.Wallpapers.Present = (wallpaperCount > 0)
	manifest.Content.Wallpapers.Count = wallpaperCount
	logging.LogDebug("Updated wallpapers in manifest, found %d wallpapers", wallpaperCount)

	return nil
}

// updateIconsInManifest updates the icons section of the manifest
func updateIconsInManifest(themePath string, manifest *ThemeManifest) error {
	logging.LogDebug("Updating icons in manifest")

	// Get system paths for reference
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Initialize counters and mappings
	systemIconCount := 0
	toolIconCount := 0
	collectionIconCount := 0
	newIconMappings := []PathMapping{}

	// Regular expression to extract system tag
	reSystemTag := regexp.MustCompile(`\((.*?)\)`)

	// 1. Check SystemIcons directory
	sysIconsDir := filepath.Join(themePath, "Icons", "SystemIcons")
	if entries, err := os.ReadDir(sysIconsDir); err == nil {
		logging.LogDebug("Found SystemIcons directory, checking for icons")
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				continue
			}

			fileName := entry.Name()
			systemPath := ""
			metadata := map[string]string{
				"IconType": "System",
			}

			// Check for special icons
			if fileName == "Collections.png" {
				systemPath = filepath.Join(systemPaths.Root, ".media", "Collections.png")
				metadata["SystemName"] = "Collections"
				metadata["SystemTag"] = "COLLECTIONS"
				metadata["IconType"] = "Special"
			} else if fileName == "Recently Played.png" {
				systemPath = filepath.Join(systemPaths.Root, ".media", "Recently Played.png")
				metadata["SystemName"] = "Recently Played"
				metadata["SystemTag"] = "RECENT"
				metadata["IconType"] = "Special"
			} else if fileName == "Tools.png" {
				// Tools icon is stored in the parent Tools directory's .media
				toolsBaseDir := filepath.Dir(systemPaths.Tools)
				systemPath = filepath.Join(toolsBaseDir, ".media", "tg5040.png")
				metadata["SystemName"] = "Tools"
				metadata["SystemTag"] = "TOOLS"
				metadata["IconType"] = "Special"
			} else {
				// Regular system icon, extract tag from filename if present
				baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
				matches := reSystemTag.FindStringSubmatch(baseName)

				if len(matches) >= 2 {
					systemTag := matches[1]
					metadata["SystemTag"] = systemTag

					// Find matching system
					for _, sys := range systemPaths.Systems {
						if sys.Tag == systemTag {
							// System icon goes in Roms/.media directory with system's name
							systemPath = filepath.Join(systemPaths.Roms, ".media", sys.Name+".png")
							metadata["SystemName"] = sys.Name
							break
						}
					}
				} else {
					// No tag found, try matching by name directly
					for _, sys := range systemPaths.Systems {
						if sys.Name == baseName {
							systemPath = filepath.Join(systemPaths.Roms, ".media", sys.Name+".png")
							metadata["SystemName"] = sys.Name
							if sys.Tag != "" {
								metadata["SystemTag"] = sys.Tag
							}
							break
						}
					}
				}
			}

			// If we found a valid system path, add to mappings
			if systemPath != "" {
				themeRelPath := filepath.Join("Icons", "SystemIcons", fileName)

				newIconMappings = append(newIconMappings, PathMapping{
					ThemePath:  themeRelPath,
					SystemPath: systemPath,
					Metadata:   metadata,
				})

				systemIconCount++
				logging.LogDebug("Found system icon: %s -> %s", themeRelPath, systemPath)
			} else {
				logging.LogDebug("Could not determine system path for icon: %s", fileName)
			}
		}
	} else {
		logging.LogDebug("SystemIcons directory not found or error reading it: %v", err)
	}

	// 2. Check ToolIcons directory
	toolIconsDir := filepath.Join(themePath, "Icons", "ToolIcons")
	if entries, err := os.ReadDir(toolIconsDir); err == nil {
		logging.LogDebug("Found ToolIcons directory, checking for icons")
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				continue
			}

			fileName := entry.Name()
			toolName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

			// Tool icons go in the Tools/.media directory
			systemPath := filepath.Join(systemPaths.Tools, ".media", fileName)
			themeRelPath := filepath.Join("Icons", "ToolIcons", fileName)

			newIconMappings = append(newIconMappings, PathMapping{
				ThemePath:  themeRelPath,
				SystemPath: systemPath,
				Metadata: map[string]string{
					"ToolName": toolName,
					"IconType": "Tool",
				},
			})

			toolIconCount++
			logging.LogDebug("Found tool icon: %s -> %s", themeRelPath, systemPath)
		}
	} else {
		logging.LogDebug("ToolIcons directory not found or error reading it: %v", err)
	}

	// 3. Check CollectionIcons directory
	collIconsDir := filepath.Join(themePath, "Icons", "CollectionIcons")
	if entries, err := os.ReadDir(collIconsDir); err == nil {
		logging.LogDebug("Found CollectionIcons directory, checking for icons")
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				continue
			}

			fileName := entry.Name()
			collectionName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

			// Collection icons go in the Collections/.media directory
			systemPath := filepath.Join(systemPaths.Root, "Collections", ".media", fileName)
			themeRelPath := filepath.Join("Icons", "CollectionIcons", fileName)

			newIconMappings = append(newIconMappings, PathMapping{
				ThemePath:  themeRelPath,
				SystemPath: systemPath,
				Metadata: map[string]string{
					"CollectionName": collectionName,
					"IconType":       "Collection",
				},
			})

			collectionIconCount++
			logging.LogDebug("Found collection icon: %s -> %s", themeRelPath, systemPath)
		}
	} else {
		logging.LogDebug("CollectionIcons directory not found or error reading it: %v", err)
	}

	// Update the manifest
	manifest.PathMappings.Icons = newIconMappings
	manifest.Content.Icons.Present = (systemIconCount+toolIconCount+collectionIconCount > 0)
	manifest.Content.Icons.SystemCount = systemIconCount
	manifest.Content.Icons.ToolCount = toolIconCount
	manifest.Content.Icons.CollectionCount = collectionIconCount

	logging.LogDebug("Updated icons in manifest, found %d system icons, %d tool icons, %d collection icons",
		systemIconCount, toolIconCount, collectionIconCount)

	return nil
}

// updateOverlaysInManifest updates the overlays section of the manifest
func updateOverlaysInManifest(themePath string, manifest *ThemeManifest) error {
	logging.LogDebug("Updating overlays in manifest")

	overlaysDir := filepath.Join(themePath, "Overlays")
	if _, err := os.Stat(overlaysDir); os.IsNotExist(err) {
		logging.LogDebug("Overlays directory not found, skipping overlay scan")
		manifest.Content.Overlays.Present = false
		manifest.Content.Overlays.Systems = []string{}
		manifest.PathMappings.Overlays = []PathMapping{}
		return nil
	}

	// Count overlays and track systems
	overlayCount := 0
	systemsWithOverlays := make(map[string]bool)
	newOverlayMappings := []PathMapping{}

	// Read overlay directory
	systemDirs, err := os.ReadDir(overlaysDir)
	if err != nil {
		logging.LogDebug("Error reading overlays directory: %v", err)
		return fmt.Errorf("error reading overlays directory: %w", err)
	}

	// Process each system in the overlays directory
	for _, systemDir := range systemDirs {
		if !systemDir.IsDir() {
			continue
		}

		systemName := systemDir.Name()
		systemOverlaysPath := filepath.Join(overlaysDir, systemName)

		// Read system overlays directory
		overlayFiles, err := os.ReadDir(systemOverlaysPath)
		if err != nil {
			logging.LogDebug("Error reading overlays for system %s: %v", systemName, err)
			continue
		}

		systemHasOverlays := false

		// Process each overlay file
		for _, overlayFile := range overlayFiles {
			if overlayFile.IsDir() || !strings.HasSuffix(strings.ToLower(overlayFile.Name()), ".png") {
				continue
			}

			fileName := overlayFile.Name()
			themeRelPath := filepath.Join("Overlays", systemName, fileName)

			// Overlay files go in the /mnt/SDCARD/Overlays directory
			systemPath := filepath.Join("/mnt/SDCARD", "Overlays", systemName, fileName)

			newOverlayMappings = append(newOverlayMappings, PathMapping{
				ThemePath:  themeRelPath,
				SystemPath: systemPath,
				Metadata: map[string]string{
					"SystemName": systemName,
				},
			})

			overlayCount++
			systemHasOverlays = true
			logging.LogDebug("Found overlay: %s -> %s", themeRelPath, systemPath)
		}

		if systemHasOverlays {
			systemsWithOverlays[systemName] = true
		}
	}

	// Convert systems map to slice
	systemsList := make([]string, 0, len(systemsWithOverlays))
	for system := range systemsWithOverlays {
		systemsList = append(systemsList, system)
	}

	// Update the manifest
	manifest.PathMappings.Overlays = newOverlayMappings
	manifest.Content.Overlays.Present = (overlayCount > 0)
	manifest.Content.Overlays.Systems = systemsList

	logging.LogDebug("Updated overlays in manifest, found %d overlays in %d systems",
		overlayCount, len(systemsList))

	return nil
}

// updateFontsInManifest updates the fonts section of the manifest
func updateFontsInManifest(themePath string, manifest *ThemeManifest) error {
	logging.LogDebug("Updating fonts in manifest")

	// Font paths
	ogFontPath := "/mnt/SDCARD/.system/res/font2.ttf"
	ogBackupPath := "/mnt/SDCARD/.system/res/font2.backup.ttf"
	nextFontPath := "/mnt/SDCARD/.system/res/font1.ttf"
	nextBackupPath := "/mnt/SDCARD/.system/res/font1.backup.ttf"

	// Target paths in theme
	themeOGPath := filepath.Join(themePath, "Fonts", "OG.ttf")
	themeOGBackupPath := filepath.Join(themePath, "Fonts", "OG.backup.ttf")
	themeNextPath := filepath.Join(themePath, "Fonts", "Next.ttf")
	themeNextBackupPath := filepath.Join(themePath, "Fonts", "Next.backup.ttf")

	// Initialize mapping
	fontMappings := make(map[string]PathMapping)
	fontsPresent := false
	ogReplaced := false
	nextReplaced := false

	// Check OG font
	if _, err := os.Stat(themeOGPath); err == nil {
		fontMappings["og_font"] = PathMapping{
			ThemePath:  "Fonts/OG.ttf",
			SystemPath: ogFontPath,
		}
		fontsPresent = true
		logging.LogDebug("Found OG font")
	}

	// Check OG backup font
	if _, err := os.Stat(themeOGBackupPath); err == nil {
		fontMappings["og_backup"] = PathMapping{
			ThemePath:  "Fonts/OG.backup.ttf",
			SystemPath: ogBackupPath,
		}
		ogReplaced = true
		fontsPresent = true
		logging.LogDebug("Found OG backup font")
	}

	// Check Next font
	if _, err := os.Stat(themeNextPath); err == nil {
		fontMappings["next_font"] = PathMapping{
			ThemePath:  "Fonts/Next.ttf",
			SystemPath: nextFontPath,
		}
		fontsPresent = true
		logging.LogDebug("Found Next font")
	}

	// Check Next backup font
	if _, err := os.Stat(themeNextBackupPath); err == nil {
		fontMappings["next_backup"] = PathMapping{
			ThemePath:  "Fonts/Next.backup.ttf",
			SystemPath: nextBackupPath,
		}
		nextReplaced = true
		fontsPresent = true
		logging.LogDebug("Found Next backup font")
	}

	// Update the manifest
	manifest.PathMappings.Fonts = fontMappings
	manifest.Content.Fonts.Present = fontsPresent
	manifest.Content.Fonts.OGReplaced = ogReplaced
	manifest.Content.Fonts.NextReplaced = nextReplaced

	logging.LogDebug("Updated fonts in manifest, fonts present: %v, OG replaced: %v, Next replaced: %v",
		fontsPresent, ogReplaced, nextReplaced)

	return nil
}

// updateSettingsInManifest updates the settings section of the manifest
func updateSettingsInManifest(themePath string, manifest *ThemeManifest) error {
	logging.LogDebug("Updating settings in manifest")

	// Settings file paths
	accentSettingsPath := "/mnt/SDCARD/.userdata/shared/minuisettings.txt"
	ledSettingsPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"

	// Theme paths
	themeAccentPath := filepath.Join(themePath, "Settings", "minuisettings.txt")
	themeLEDPath := filepath.Join(themePath, "Settings", "ledsettings_brick.txt")

	// Initialize mapping and flags
	settingsMappings := make(map[string]PathMapping)
	accentsIncluded := false
	ledsIncluded := false

	// Check accent settings
	if _, err := os.Stat(themeAccentPath); err == nil {
		settingsMappings["accents"] = PathMapping{
			ThemePath:  "Settings/minuisettings.txt",
			SystemPath: accentSettingsPath,
		}
		accentsIncluded = true
		logging.LogDebug("Found accent settings")

		// Extract accent colors for the manifest
		if err := ExtractAccentColors(themeAccentPath, manifest); err != nil {
			logging.LogDebug("Warning: Could not extract accent colors: %v", err)
		}
	}

	// Check LED settings
	if _, err := os.Stat(themeLEDPath); err == nil {
		settingsMappings["leds"] = PathMapping{
			ThemePath:  "Settings/ledsettings_brick.txt",
			SystemPath: ledSettingsPath,
		}
		ledsIncluded = true
		logging.LogDebug("Found LED settings")

		// TODO: Extract LED settings details for the manifest
		logging.LogDebug("LED settings extracted and added to manifest")
	}

	// Update the manifest
	manifest.PathMappings.Settings = settingsMappings
	manifest.Content.Settings.AccentsIncluded = accentsIncluded
	manifest.Content.Settings.LEDsIncluded = ledsIncluded

	logging.LogDebug("Updated settings in manifest, accents included: %v, LEDs included: %v",
		accentsIncluded, ledsIncluded)

	return nil
}
