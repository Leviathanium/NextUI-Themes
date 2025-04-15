// src/internal/themes/export.go
// Implementation of theme export functionality

package themes

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	// Removed unused import: "nextui-themes/internal/logging"
	"nextui-themes/internal/accents"
	"nextui-themes/internal/fonts"
	"nextui-themes/internal/system"
	"nextui-themes/internal/ui"
)

// CreateThemeExportDirectory creates a new theme directory with sequential naming
func CreateThemeExportDirectory() (string, error) {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to Themes/Exports directory
	exportsDir := filepath.Join(cwd, "Themes", "Exports")

	// Ensure directory exists
	if err := os.MkdirAll(exportsDir, 0755); err != nil {
		return "", fmt.Errorf("error creating exports directory: %w", err)
	}

	// Generate sequential theme name
	themeNumber := 1
	var themePath string

	for {
		themeName := fmt.Sprintf("theme_%d.theme", themeNumber)
		themePath = filepath.Join(exportsDir, themeName)

		if _, err := os.Stat(themePath); os.IsNotExist(err) {
			// Theme directory doesn't exist, we can use this name
			break
		}

		themeNumber++
	}

	// Create the theme directory and subdirectories using the preferred structure
	// that mirrors the icon directory structure
	subDirs := []string{
		// Primary wallpaper structure
		"Wallpapers/SystemWallpapers",
		"Wallpapers/CollectionWallpapers",

		// Legacy structure (for backward compatibility)
		"Wallpapers/Root",
		"Wallpapers/Collections",
		"Wallpapers/Recently Played",
		"Wallpapers/Tools",
		"Wallpapers/Systems",

		// Icons directory structure
		"Icons/SystemIcons",
		"Icons/ToolIcons",
		"Icons/CollectionIcons",

		// Other directories
		"Overlays",
		"Fonts",
		"Settings",
	}

	for _, dir := range subDirs {
		path := filepath.Join(themePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", fmt.Errorf("error creating theme subdirectory %s: %w", dir, err)
		}
	}

	return themePath, nil
}

// ExportWallpapers exports wallpapers to the theme directory using the preferred structure
func ExportWallpapers(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logger.Printf("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Initialize wallpapers section in manifest
	manifest.Content.Wallpapers.Present = false
	manifest.Content.Wallpapers.Count = 0

	// Directory for system wallpapers in the preferred structure
	sysWallDir := filepath.Join(themePath, "Wallpapers", "SystemWallpapers")
	collWallDir := filepath.Join(themePath, "Wallpapers", "CollectionWallpapers")

	// Ensure directories exist
	if err := os.MkdirAll(sysWallDir, 0755); err != nil {
		logger.Printf("Error creating SystemWallpapers directory: %v", err)
		return fmt.Errorf("error creating SystemWallpapers directory: %w", err)
	}

	if err := os.MkdirAll(collWallDir, 0755); err != nil {
		logger.Printf("Error creating CollectionWallpapers directory: %v", err)
		return fmt.Errorf("error creating CollectionWallpapers directory: %w", err)
	}

	// Root wallpaper (bg.png at root level)
	rootBg := filepath.Join(systemPaths.Root, "bg.png")
	if _, err := os.Stat(rootBg); err == nil {
		// Using the new preferred structure: SystemWallpapers/Root.png
		targetPath := filepath.Join(sysWallDir, "Root.png")
		if err := CopyFile(rootBg, targetPath); err != nil {
			logger.Printf("Warning: Could not copy root bg.png: %v", err)
		} else {
			// Add to manifest path mappings
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Root.png",
					SystemPath: rootBg,
					Metadata: map[string]string{
						"SystemName":    "Root",
						"WallpaperType": "Main",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.Printf("Exported root wallpaper: %s", rootBg)
		}
	}

	// Root .media wallpaper
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", "bg.png")
	if _, err := os.Stat(rootMediaBg); err == nil {
		// Using the new preferred structure: SystemWallpapers/Root-Media.png
		targetPath := filepath.Join(sysWallDir, "Root-Media.png")
		if err := CopyFile(rootMediaBg, targetPath); err != nil {
			logger.Printf("Warning: Could not copy root .media/bg.png: %v", err)
		} else {
			// Add to manifest path mappings
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Root-Media.png",
					SystemPath: rootMediaBg,
					Metadata: map[string]string{
						"SystemName":    "Root",
						"WallpaperType": "Media",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.Printf("Exported root .media wallpaper: %s", rootMediaBg)
		}
	}

	// Recently Played wallpaper
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
	if _, err := os.Stat(rpBg); err == nil {
		// Using the new preferred structure: SystemWallpapers/Recently Played.png
		targetPath := filepath.Join(sysWallDir, "Recently Played.png")
		if err := CopyFile(rpBg, targetPath); err != nil {
			logger.Printf("Warning: Could not copy Recently Played bg.png: %v", err)
		} else {
			// Add to manifest path mappings
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Recently Played.png",
					SystemPath: rpBg,
					Metadata: map[string]string{
						"SystemName":    "Recently Played",
						"WallpaperType": "Media",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.Printf("Exported Recently Played wallpaper: %s", rpBg)
		}
	}

	// Tools wallpaper
	toolsBg := filepath.Join(systemPaths.Tools, ".media", "bg.png")
	if _, err := os.Stat(toolsBg); err == nil {
		// Using the new preferred structure: SystemWallpapers/Tools.png
		targetPath := filepath.Join(sysWallDir, "Tools.png")
		if err := CopyFile(toolsBg, targetPath); err != nil {
			logger.Printf("Warning: Could not copy Tools bg.png: %v", err)
		} else {
			// Add to manifest path mappings
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Tools.png",
					SystemPath: toolsBg,
					Metadata: map[string]string{
						"SystemName":    "Tools",
						"WallpaperType": "Media",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.Printf("Exported Tools wallpaper: %s", toolsBg)
		}
	}

	// Collections wallpaper
	collectionsPath := filepath.Join(systemPaths.Root, "Collections")
	collectionsBg := filepath.Join(collectionsPath, ".media", "bg.png")
	if _, err := os.Stat(collectionsBg); err == nil {
		// Using the new preferred structure: SystemWallpapers/Collections.png
		targetPath := filepath.Join(sysWallDir, "Collections.png")
		if err := CopyFile(collectionsBg, targetPath); err != nil {
			logger.Printf("Warning: Could not copy Collections bg.png: %v", err)
		} else {
			// Add to manifest path mappings
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Collections.png",
					SystemPath: collectionsBg,
					Metadata: map[string]string{
						"SystemName":    "Collections",
						"WallpaperType": "Media",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.Printf("Exported Collections wallpaper: %s", collectionsBg)
		}
	}

	// Check for collection-specific wallpapers
	if _, err := os.Stat(collectionsPath); err == nil {
		// Read the Collections directory to find individual collections
		collections, err := os.ReadDir(collectionsPath)
		if err == nil {
			for _, collection := range collections {
				if collection.IsDir() && !strings.HasPrefix(collection.Name(), ".") {
					collectionName := collection.Name()
					collectionBg := filepath.Join(collectionsPath, collectionName, ".media", "bg.png")

					// Check if this collection has a background
					if _, err := os.Stat(collectionBg); err == nil {
						// Using the new preferred structure: CollectionWallpapers/[CollectionName].png
						targetPath := filepath.Join(collWallDir, collectionName+".png")
						if err := CopyFile(collectionBg, targetPath); err != nil {
							logger.Printf("Warning: Could not copy collection %s bg.png: %v", collectionName, err)
						} else {
							// Add to manifest path mappings
							manifest.PathMappings.Wallpapers = append(
								manifest.PathMappings.Wallpapers,
								PathMapping{
									ThemePath:  "Wallpapers/CollectionWallpapers/" + collectionName + ".png",
									SystemPath: collectionBg,
									Metadata: map[string]string{
										"CollectionName": collectionName,
										"WallpaperType":  "Collection",
									},
								},
							)
							manifest.Content.Wallpapers.Present = true
							manifest.Content.Wallpapers.Count++
							logger.Printf("Exported collection wallpaper for %s: %s", collectionName, collectionBg)
						}
					}
				}
			}
		}
	}

	// System wallpapers
	for _, system := range systemPaths.Systems {
		systemBg := filepath.Join(system.MediaPath, "bg.png")
		if _, err := os.Stat(systemBg); err == nil {
			// Skip systems without a tag
			if system.Tag == "" {
				logger.Printf("Skipping system with no tag: %s", system.Name)
				continue
			}

			// Using the new preferred structure: SystemWallpapers/[SystemName (TAG)].png
			fileName := fmt.Sprintf("%s (%s).png", system.Name, system.Tag)
			targetPath := filepath.Join(sysWallDir, fileName)

			if err := CopyFile(systemBg, targetPath); err != nil {
				logger.Printf("Warning: Could not copy system %s bg.png: %v", system.Name, err)
			} else {
				// Add to manifest path mappings with extended metadata
				manifest.PathMappings.Wallpapers = append(
					manifest.PathMappings.Wallpapers,
					PathMapping{
						ThemePath:  "Wallpapers/SystemWallpapers/" + fileName,
						SystemPath: systemBg,
						// Add system metadata to improve matching
						Metadata: map[string]string{
							"SystemName":    system.Name,
							"SystemTag":     system.Tag,
							"WallpaperType": "System",
						},
					},
				)
				manifest.Content.Wallpapers.Present = true
				manifest.Content.Wallpapers.Count++
				logger.Printf("Exported system wallpaper for %s (%s): %s", system.Name, system.Tag, systemBg)
			}
		}
	}

	return nil
}

// ExportIcons exports icons to the theme directory
func ExportIcons(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logger.Printf("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Initialize icons section in manifest
	manifest.Content.Icons.Present = false
	manifest.Content.Icons.SystemCount = 0
	manifest.Content.Icons.ToolCount = 0
	manifest.Content.Icons.CollectionCount = 0

	// Root media directory for special icons
	rootMediaPath := filepath.Join(systemPaths.Root, ".media")

	// Collections icon
	collectionsIcon := filepath.Join(rootMediaPath, "Collections.png")
	if _, err := os.Stat(collectionsIcon); err == nil {
		targetPath := filepath.Join(themePath, "Icons", "SystemIcons", "Collections.png")
		if err := CopyFile(collectionsIcon, targetPath); err != nil {
			logger.Printf("Warning: Could not copy Collections icon: %v", err)
		} else {
			// Add to manifest path mappings
			manifest.PathMappings.Icons = append(
				manifest.PathMappings.Icons,
				PathMapping{
					ThemePath:  "Icons/SystemIcons/Collections.png",
					SystemPath: collectionsIcon,
					Metadata: map[string]string{
						"SystemName": "Collections",
						"SystemTag":  "COLLECTIONS", // Special tag for Collections
						"IconType":   "Special",
					},
				},
			)
			manifest.Content.Icons.Present = true
			manifest.Content.Icons.SystemCount++
			logger.Printf("Exported Collections icon: %s", collectionsIcon)
		}
	}

	// Recently Played icon
	rpIcon := filepath.Join(rootMediaPath, "Recently Played.png")
	if _, err := os.Stat(rpIcon); err == nil {
		targetPath := filepath.Join(themePath, "Icons", "SystemIcons", "Recently Played.png")
		if err := CopyFile(rpIcon, targetPath); err != nil {
			logger.Printf("Warning: Could not copy Recently Played icon: %v", err)
		} else {
			// Add to manifest path mappings
			manifest.PathMappings.Icons = append(
				manifest.PathMappings.Icons,
				PathMapping{
					ThemePath:  "Icons/SystemIcons/Recently Played.png",
					SystemPath: rpIcon,
					Metadata: map[string]string{
						"SystemName": "Recently Played",
						"SystemTag":  "RECENT", // Special tag for Recently Played
						"IconType":   "Special",
					},
				},
			)
			manifest.Content.Icons.Present = true
			manifest.Content.Icons.SystemCount++
			logger.Printf("Exported Recently Played icon: %s", rpIcon)
		}
	}

	// Tools icon
	toolsBaseDir := filepath.Dir(systemPaths.Tools)
	toolsMediaPath := filepath.Join(toolsBaseDir, ".media")
	toolsIcon := filepath.Join(toolsMediaPath, "tg5040.png")
	if _, err := os.Stat(toolsIcon); err == nil {
		targetPath := filepath.Join(themePath, "Icons", "SystemIcons", "Tools.png")
		if err := CopyFile(toolsIcon, targetPath); err != nil {
			logger.Printf("Warning: Could not copy Tools icon: %v", err)
		} else {
			// Add to manifest path mappings
			manifest.PathMappings.Icons = append(
				manifest.PathMappings.Icons,
				PathMapping{
					ThemePath:  "Icons/SystemIcons/Tools.png",
					SystemPath: toolsIcon,
					Metadata: map[string]string{
						"SystemName": "Tools",
						"SystemTag":  "TOOLS", // Special tag for Tools
						"IconType":   "Special",
					},
				},
			)
			manifest.Content.Icons.Present = true
			manifest.Content.Icons.SystemCount++
			logger.Printf("Exported Tools icon: %s", toolsIcon)
		}
	}

	// Regular expression to extract system tag
	re := regexp.MustCompile(`\((.*?)\)`)

	// System icons
	romsMediaPath := filepath.Join(systemPaths.Roms, ".media")
	if entries, err := os.ReadDir(romsMediaPath); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				// Skip bg.png
				if entry.Name() == "bg.png" {
					continue
				}

				iconPath := filepath.Join(romsMediaPath, entry.Name())
				targetPath := filepath.Join(themePath, "Icons", "SystemIcons", entry.Name())

				// Extract system name and tag
				systemName := entry.Name()
				systemName = strings.TrimSuffix(systemName, ".png") // Remove the .png extension

				// Extract the tag using regex
				systemTag := ""
				matches := re.FindStringSubmatch(systemName)
				if len(matches) >= 2 {
					systemTag = matches[1]
				}

				if err := CopyFile(iconPath, targetPath); err != nil {
					logger.Printf("Warning: Could not copy system icon %s: %v", entry.Name(), err)
				} else {
					// Add to manifest path mappings with metadata
					manifest.PathMappings.Icons = append(
						manifest.PathMappings.Icons,
						PathMapping{
							ThemePath:  fmt.Sprintf("Icons/SystemIcons/%s", entry.Name()),
							SystemPath: iconPath,
							Metadata: map[string]string{
								"SystemName": systemName,
								"SystemTag":  systemTag,
								"IconType":   "System",
							},
						},
					)
					manifest.Content.Icons.Present = true
					manifest.Content.Icons.SystemCount++
					logger.Printf("Exported system icon: %s", iconPath)
				}
			}
		}
	}

	// Tool icons
	toolsDir := filepath.Join(systemPaths.Tools, ".media")
	if entries, err := os.ReadDir(toolsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				// Skip bg.png and tg5040.png
				if entry.Name() == "bg.png" || entry.Name() == "tg5040.png" {
					continue
				}

				iconPath := filepath.Join(toolsDir, entry.Name())
				targetPath := filepath.Join(themePath, "Icons", "ToolIcons", entry.Name())

				// Get tool name without extension
				toolName := strings.TrimSuffix(entry.Name(), ".png")

				if err := CopyFile(iconPath, targetPath); err != nil {
					logger.Printf("Warning: Could not copy tool icon %s: %v", entry.Name(), err)
				} else {
					// Add to manifest path mappings with metadata
					manifest.PathMappings.Icons = append(
						manifest.PathMappings.Icons,
						PathMapping{
							ThemePath:  fmt.Sprintf("Icons/ToolIcons/%s", entry.Name()),
							SystemPath: iconPath,
							Metadata: map[string]string{
								"ToolName": toolName,
								"IconType": "Tool",
							},
						},
					)
					manifest.Content.Icons.Present = true
					manifest.Content.Icons.ToolCount++
					logger.Printf("Exported tool icon: %s", iconPath)
				}
			}
		}
	}

	// Collection icons
	collectionsMediaPath := filepath.Join(systemPaths.Root, "Collections", ".media")
	if entries, err := os.ReadDir(collectionsMediaPath); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				// Skip bg.png
				if entry.Name() == "bg.png" {
					continue
				}

				iconPath := filepath.Join(collectionsMediaPath, entry.Name())
				targetPath := filepath.Join(themePath, "Icons", "CollectionIcons", entry.Name())

				// Get collection name without extension
				collectionName := strings.TrimSuffix(entry.Name(), ".png")

				if err := CopyFile(iconPath, targetPath); err != nil {
					logger.Printf("Warning: Could not copy collection icon %s: %v", entry.Name(), err)
				} else {
					// Add to manifest path mappings with metadata
					manifest.PathMappings.Icons = append(
						manifest.PathMappings.Icons,
						PathMapping{
							ThemePath:  fmt.Sprintf("Icons/CollectionIcons/%s", entry.Name()),
							SystemPath: iconPath,
							Metadata: map[string]string{
								"CollectionName": collectionName,
								"IconType":       "Collection",
							},
						},
					)
					manifest.Content.Icons.Present = true
					manifest.Content.Icons.CollectionCount++
					logger.Printf("Exported collection icon: %s", iconPath)
				}
			}
		}
	}

	return nil
}

// ExportOverlays exports overlays to the theme directory
func ExportOverlays(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Initialize overlays section in manifest
	manifest.Content.Overlays.Present = false
	manifest.Content.Overlays.Systems = []string{}

	// Path to Overlays directory
	overlaysPath := filepath.Join("/mnt/SDCARD", "Overlays")
	if _, err := os.Stat(overlaysPath); os.IsNotExist(err) {
		logger.Printf("Overlays directory not found, skipping overlay export")
		return nil
	}

	// Read Overlays directory
	entries, err := os.ReadDir(overlaysPath)
	if err != nil {
		logger.Printf("Error reading Overlays directory: %v", err)
		return nil // Skip but don't fail
	}

	// Process each system in the Overlays directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		systemName := entry.Name()
		systemOverlaysPath := filepath.Join(overlaysPath, systemName)
		targetDir := filepath.Join(themePath, "Overlays", systemName)

		// Create system overlays directory in theme
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			logger.Printf("Warning: Could not create overlays directory for system %s: %v", systemName, err)
			continue
		}

		// Read system overlays directory
		overlayFiles, err := os.ReadDir(systemOverlaysPath)
		if err != nil {
			logger.Printf("Warning: Could not read overlays for system %s: %v", systemName, err)
			continue
		}

		systemHasOverlays := false

		// Copy each overlay file
		for _, overlayFile := range overlayFiles {
			if overlayFile.IsDir() || !strings.HasSuffix(strings.ToLower(overlayFile.Name()), ".png") {
				continue
			}

			srcPath := filepath.Join(systemOverlaysPath, overlayFile.Name())
			targetPath := filepath.Join(targetDir, overlayFile.Name())

			if err := CopyFile(srcPath, targetPath); err != nil {
				logger.Printf("Warning: Could not copy overlay %s for system %s: %v",
					overlayFile.Name(), systemName, err)
			} else {
				// Add to manifest path mappings
				manifest.PathMappings.Overlays = append(
					manifest.PathMappings.Overlays,
					PathMapping{
						ThemePath:  fmt.Sprintf("Overlays/%s/%s", systemName, overlayFile.Name()),
						SystemPath: srcPath,
					},
				)
				systemHasOverlays = true
				logger.Printf("Exported overlay: %s", srcPath)
			}
		}

		if systemHasOverlays {
			manifest.Content.Overlays.Present = true
			manifest.Content.Overlays.Systems = append(manifest.Content.Overlays.Systems, systemName)
		}
	}

	return nil
}

// ExportFonts exports fonts to the theme directory
func ExportFonts(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Initialize fonts section in manifest
	manifest.Content.Fonts.Present = false
	manifest.Content.Fonts.OGReplaced = false
	manifest.Content.Fonts.NextReplaced = false

	// Initialize fonts mapping in manifest
	manifest.PathMappings.Fonts = make(map[string]PathMapping)

	// Copy and check OG font
	ogFontPath := fonts.OGFontPath
	ogBackupPath := filepath.Join(filepath.Dir(fonts.OGFontPath), fonts.OGFontBackupName)

	targetOGPath := filepath.Join(themePath, "Fonts", "OG.ttf")
	targetOGBackupPath := filepath.Join(themePath, "Fonts", "OG.backup.ttf")

	// Copy OG font
	if _, err := os.Stat(ogFontPath); err == nil {
		if err := CopyFile(ogFontPath, targetOGPath); err != nil {
			logger.Printf("Warning: Could not copy OG font: %v", err)
		} else {
			manifest.Content.Fonts.Present = true
			manifest.PathMappings.Fonts["og_font"] = PathMapping{
				ThemePath:  "Fonts/OG.ttf",
				SystemPath: ogFontPath,
			}
			logger.Printf("Exported OG font: %s", ogFontPath)
		}
	}

	// Check if OG font has been replaced (backup exists)
	if _, err := os.Stat(ogBackupPath); err == nil {
		manifest.Content.Fonts.OGReplaced = true

		// Copy OG backup font
		if err := CopyFile(ogBackupPath, targetOGBackupPath); err != nil {
			logger.Printf("Warning: Could not copy OG font backup: %v", err)
		} else {
			manifest.PathMappings.Fonts["og_backup"] = PathMapping{
				ThemePath:  "Fonts/OG.backup.ttf",
				SystemPath: ogBackupPath,
			}
			logger.Printf("Exported OG font backup: %s", ogBackupPath)
		}
	} else {
		// No backup found, copy our default backup
		cwd, err := os.Getwd()
		if err != nil {
			logger.Printf("Error getting current directory: %v", err)
		} else {
			defaultOGBackup := filepath.Join(cwd, "Fonts", "Backups", "font2.backup.ttf")
			if _, err := os.Stat(defaultOGBackup); err == nil {
				if err := CopyFile(defaultOGBackup, targetOGBackupPath); err != nil {
					logger.Printf("Warning: Could not copy default OG font backup: %v", err)
				} else {
					manifest.PathMappings.Fonts["og_backup"] = PathMapping{
						ThemePath:  "Fonts/OG.backup.ttf",
						SystemPath: ogBackupPath,
					}
					logger.Printf("Exported default OG font backup")
				}
			} else {
				logger.Printf("Warning: No default OG font backup found: %s", defaultOGBackup)
			}
		}
	}

	// Copy and check Next font
	nextFontPath := fonts.NextFontPath
	nextBackupPath := filepath.Join(filepath.Dir(fonts.NextFontPath), fonts.NextFontBackupName)

	targetNextPath := filepath.Join(themePath, "Fonts", "Next.ttf")
	targetNextBackupPath := filepath.Join(themePath, "Fonts", "Next.backup.ttf")

	// Copy Next font
	if _, err := os.Stat(nextFontPath); err == nil {
		if err := CopyFile(nextFontPath, targetNextPath); err != nil {
			logger.Printf("Warning: Could not copy Next font: %v", err)
		} else {
			manifest.Content.Fonts.Present = true
			manifest.PathMappings.Fonts["next_font"] = PathMapping{
				ThemePath:  "Fonts/Next.ttf",
				SystemPath: nextFontPath,
			}
			logger.Printf("Exported Next font: %s", nextFontPath)
		}
	}

	// Check if Next font has been replaced (backup exists)
	if _, err := os.Stat(nextBackupPath); err == nil {
		manifest.Content.Fonts.NextReplaced = true

		// Copy Next backup font
		if err := CopyFile(nextBackupPath, targetNextBackupPath); err != nil {
			logger.Printf("Warning: Could not copy Next font backup: %v", err)
		} else {
			manifest.PathMappings.Fonts["next_backup"] = PathMapping{
				ThemePath:  "Fonts/Next.backup.ttf",
				SystemPath: nextBackupPath,
			}
			logger.Printf("Exported Next font backup: %s", nextBackupPath)
		}
	} else {
		// No backup found, copy our default backup
		cwd, err := os.Getwd()
		if err != nil {
			logger.Printf("Error getting current directory: %v", err)
		} else {
			defaultNextBackup := filepath.Join(cwd, "Fonts", "Backups", "font1.backup.ttf")
			if _, err := os.Stat(defaultNextBackup); err == nil {
				if err := CopyFile(defaultNextBackup, targetNextBackupPath); err != nil {
					logger.Printf("Warning: Could not copy default Next font backup: %v", err)
				} else {
					manifest.PathMappings.Fonts["next_backup"] = PathMapping{
						ThemePath:  "Fonts/Next.backup.ttf",
						SystemPath: nextBackupPath,
					}
					logger.Printf("Exported default Next font backup")
				}
			} else {
				logger.Printf("Warning: No default Next font backup found: %s", defaultNextBackup)
			}
		}
	}

	return nil
}

// ExportSettings exports settings to the theme directory
func ExportSettings(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Initialize settings section in manifest
	manifest.Content.Settings.AccentsIncluded = false
	manifest.Content.Settings.LEDsIncluded = false

	// Initialize settings mapping in manifest
	manifest.PathMappings.Settings = make(map[string]PathMapping)

	// Export accent settings
	accentSettingsPath := accents.SettingsPath
	if _, err := os.Stat(accentSettingsPath); err == nil {
		targetPath := filepath.Join(themePath, "Settings", "minuisettings.txt")
		if err := CopyFile(accentSettingsPath, targetPath); err != nil {
			logger.Printf("Warning: Could not copy accent settings: %v", err)
		} else {
			manifest.Content.Settings.AccentsIncluded = true
			manifest.PathMappings.Settings["accents"] = PathMapping{
				ThemePath:  "Settings/minuisettings.txt",
				SystemPath: accentSettingsPath,
			}
			logger.Printf("Exported accent settings: %s", accentSettingsPath)

			// Extract accent colors for the manifest
			if err := ExtractAccentColors(accentSettingsPath, manifest); err != nil {
				logger.Printf("Warning: Could not extract accent colors: %v", err)
			}
		}
	}

	// Export LED settings
	ledSettingsPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"
	if _, err := os.Stat(ledSettingsPath); err == nil {
		targetPath := filepath.Join(themePath, "Settings", "ledsettings_brick.txt")
		if err := CopyFile(ledSettingsPath, targetPath); err != nil {
			logger.Printf("Warning: Could not copy LED settings: %v", err)
		} else {
			manifest.Content.Settings.LEDsIncluded = true
			manifest.PathMappings.Settings["leds"] = PathMapping{
				ThemePath:  "Settings/ledsettings_brick.txt",
				SystemPath: ledSettingsPath,
			}
			logger.Printf("Exported LED settings: %s", ledSettingsPath)

			// Extract LED settings for the manifest
			if err := ExtractLEDSettings(ledSettingsPath, manifest); err != nil {
				logger.Printf("Warning: Could not extract LED settings: %v", err)
			} else {
				logger.Printf("LED settings extracted and added to manifest")
			}
		}
	}

	return nil
}

// ExtractAccentColors extracts accent colors from the settings file
func ExtractAccentColors(settingsPath string, manifest *ThemeManifest) error {
	// Initialize accent colors map
	manifest.AccentColors = make(map[string]string)

	// Read the settings file
	file, err := os.Open(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to open settings file: %w", err)
	}
	defer file.Close()

	// Read each line and extract color settings
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "color") && len(line) > 6 {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				colorKey := parts[0]
				colorValue := parts[1]

				// Convert to display format (#RRGGBB)
				if strings.HasPrefix(colorValue, "0x") {
					colorValue = "#" + colorValue[2:]
				}

				manifest.AccentColors[colorKey] = colorValue
			}
		}
	}

	return scanner.Err()
}

// ExtractLEDSettings extracts LED settings from the LED settings file
func ExtractLEDSettings(settingsPath string, manifest *ThemeManifest) error {
	// Initialize LED settings map
	manifest.LEDSettings = make(map[string]map[string]string)

	// Read the settings file
	file, err := os.Open(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to open LED settings file: %w", err)
	}
	defer file.Close()

	// Read each line and extract settings
	scanner := bufio.NewScanner(file)
	var currentSection string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check if this is a section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Extract section name (LED name)
			currentSection = line[1 : len(line)-1]
			// Initialize the map for this section
			manifest.LEDSettings[currentSection] = make(map[string]string)
			continue
		}

		// If we have a current section and this is a key-value pair
		if currentSection != "" && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// Convert color values to display format (#RRGGBB) if needed
				if (key == "color1" || key == "color2") && strings.HasPrefix(value, "0x") {
					value = "#" + value[2:]
				}

				// Add to the current section's map
				manifest.LEDSettings[currentSection][key] = value
			}
		}
	}

	return scanner.Err()
}

// GeneratePreview generates a preview image for the theme
func GeneratePreview(themePath string, logger *Logger) error {
	// For now, we'll copy the root background as the preview
	// In a more advanced version, we could generate a composite image

	rootBgPath := filepath.Join("/mnt/SDCARD", "bg.png")
	previewPath := filepath.Join(themePath, "preview.png")

	if _, err := os.Stat(rootBgPath); err == nil {
		if err := CopyFile(rootBgPath, previewPath); err != nil {
			logger.Printf("Warning: Could not create preview image: %v", err)
			return err
		}
		logger.Printf("Created preview image from root background")
	} else {
		logger.Printf("Warning: Root background not found, no preview image created")
	}

	return nil
}

// ExportTheme exports the current theme settings
func ExportTheme() error {
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
		filepath.Join(logsDir, "exports.log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()

	// Create logger
	logger := &Logger{File: logFile}
	logger.Printf("Starting theme export")

	// Create theme directory
	themePath, err := CreateThemeExportDirectory()
	if err != nil {
		logger.Printf("Error creating theme directory: %v", err)
		return fmt.Errorf("error creating theme directory: %w", err)
	}

	logger.Printf("Created theme directory: %s", themePath)

	// Initialize manifest
	manifest := &ThemeManifest{}

	// Export theme components
	if err := ExportWallpapers(themePath, manifest, logger); err != nil {
		logger.Printf("Error exporting wallpapers: %v", err)
	}

	if err := ExportIcons(themePath, manifest, logger); err != nil {
		logger.Printf("Error exporting icons: %v", err)
	}

	if err := ExportOverlays(themePath, manifest, logger); err != nil {
		logger.Printf("Error exporting overlays: %v", err)
	}

	if err := ExportFonts(themePath, manifest, logger); err != nil {
		logger.Printf("Error exporting fonts: %v", err)
	}

	if err := ExportSettings(themePath, manifest, logger); err != nil {
		logger.Printf("Error exporting settings: %v", err)
	}

	// Generate preview image
	if err := GeneratePreview(themePath, logger); err != nil {
		logger.Printf("Error generating preview: %v", err)
	}

	// Write manifest
	if err := WriteManifest(themePath, manifest, logger); err != nil {
		logger.Printf("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	logger.Printf("Theme export completed successfully: %s", themePath)

	// Show success message to user
	themeName := filepath.Base(themePath)
	ui.ShowMessage(fmt.Sprintf("Theme exported successfully: %s", themeName), "3")

	return nil
}
