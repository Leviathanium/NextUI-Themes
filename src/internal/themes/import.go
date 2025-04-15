// src/internal/themes/import.go
// Implementation of theme import functionality

package themes

// Add at the top of import.go, after existing imports
import (
	"fmt"
	"nextui-themes/internal/icons"  // Add this for icons package
	"nextui-themes/internal/system" // Add this for system paths
	"nextui-themes/internal/ui"
	"os"
	"path/filepath"
	"regexp" // Add this for regex support
	"strings"
)

// ImportWallpapers imports wallpapers from a theme
func ImportWallpapers(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Skip if no wallpapers present
	if !manifest.Content.Wallpapers.Present {
		logger.Printf("No wallpapers found in theme, skipping wallpaper import")
		return nil
	}

	// Get system paths for metadata matching
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logger.Printf("Error getting system paths: %v", err)
		// Continue anyway, we'll use direct path mapping
	}

	// Create maps for better lookups
	systemsByName := make(map[string]system.SystemInfo)
	systemsByTag := make(map[string]system.SystemInfo)
	if systemPaths != nil {
		for _, sys := range systemPaths.Systems {
			systemsByName[sys.Name] = sys
			if sys.Tag != "" {
				systemsByTag[sys.Tag] = sys
			}
		}
	}

	// Regular expression to extract system tag from filename
	reSystemTag := regexp.MustCompile(`\((.*?)\)$`)

	// First, check if we need to handle old-style theme with directories
	// This ensures backward compatibility with older themes
	oldStyleWallpapers := filepath.Join(themePath, "Wallpapers", "Root", "bg.png")
	if _, err := os.Stat(oldStyleWallpapers); err == nil {
		logger.Printf("Detected old-style theme with directory structure, using legacy import")
		return importLegacyWallpapers(themePath, manifest, logger, systemPaths)
	}

	// Process each wallpaper in the manifest
	for _, mapping := range manifest.PathMappings.Wallpapers {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source wallpaper file not found: %s", srcPath)
			continue
		}

		// Special handling based on metadata
		if mapping.Metadata != nil {
			// Handle different system types
			if systemType, ok := mapping.Metadata["SystemType"]; ok {
				switch systemType {
				case "GameSystem":
					// Try to find system by tag
					if tag, ok := mapping.Metadata["SystemTag"]; ok {
						if sys, found := systemsByTag[tag]; found {
							// Update destination path for current system
							dstPath = filepath.Join(sys.MediaPath, "bg.png")
							logger.Printf("Found system by tag %s: %s", tag, sys.Name)
						}
					}

				case "Collection":
					// Handle collection wallpapers
					if collName, ok := mapping.Metadata["CollectionName"]; ok {
						collPath := filepath.Join(systemPaths.Root, "Collections", collName, ".media")
						if err := os.MkdirAll(collPath, 0755); err != nil {
							logger.Printf("Warning: Could not create collection media directory: %v", err)
						} else {
							dstPath = filepath.Join(collPath, "bg.png")
						}
					}
				}
			}
		}

		// If we couldn't determine from metadata, try extracting from filename
		if strings.Contains(srcPath, "SystemWallpapers") && !strings.Contains(dstPath, ".media") {
			// Extract filename
			filename := filepath.Base(srcPath)

			// Check for special cases first
			switch filename {
			case "Root.png":
				// Handle Root wallpaper - needs to go to both locations
				rootPath := filepath.Join(systemPaths.Root, "bg.png")
				rootMediaPath := filepath.Join(systemPaths.Root, ".media", "bg.png")

				// Ensure media directory exists
				if err := os.MkdirAll(filepath.Dir(rootMediaPath), 0755); err != nil {
					logger.Printf("Warning: Could not create root media directory: %v", err)
				}

				// Copy to both locations
				if err := CopyFile(srcPath, rootPath); err != nil {
					logger.Printf("Warning: Could not copy Root wallpaper to root: %v", err)
				} else {
					logger.Printf("Imported Root wallpaper to: %s", rootPath)
				}

				if err := CopyFile(srcPath, rootMediaPath); err != nil {
					logger.Printf("Warning: Could not copy Root wallpaper to media dir: %v", err)
				} else {
					logger.Printf("Imported Root wallpaper to: %s", rootMediaPath)
				}

				// Skip the regular copy below
				continue

			case "Recently Played.png":
				// Handle Recently Played
				dstPath = filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")

			case "Tools.png":
				// Handle Tools
				dstPath = filepath.Join(systemPaths.Tools, ".media", "bg.png")

			case "Collections.png":
				// Handle main Collections directory
				dstPath = filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")

			default:
				// Check if it's a system wallpaper with tag
				matches := reSystemTag.FindStringSubmatch(filename)
				if len(matches) >= 2 {
					tag := matches[1]
					if sys, ok := systemsByTag[tag]; ok {
						dstPath = filepath.Join(sys.MediaPath, "bg.png")
						logger.Printf("Found system by filename tag %s: %s", tag, sys.Name)
					}
				}
			}
		} else if strings.Contains(srcPath, "CollectionWallpapers") {
			// Handle collection wallpapers
			filename := filepath.Base(srcPath)
			collName := strings.TrimSuffix(filename, filepath.Ext(filename))

			// Create collection media path
			collPath := filepath.Join(systemPaths.Root, "Collections", collName, ".media")
			if err := os.MkdirAll(collPath, 0755); err != nil {
				logger.Printf("Warning: Could not create collection media directory: %v", err)
			} else {
				dstPath = filepath.Join(collPath, "bg.png")
			}
		}

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.Printf("Warning: Could not create directory for wallpaper: %v", err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy wallpaper: %v", err)
		} else {
			logger.Printf("Imported wallpaper: %s -> %s", srcPath, dstPath)
		}
	}

	return nil
}

// importLegacyWallpapers handles importing wallpapers from older themes with the directory structure
func importLegacyWallpapers(themePath string, manifest *ThemeManifest, logger *Logger, systemPaths *system.SystemPaths) error {
	logger.Printf("Importing legacy wallpapers with directory structure")

	// Root wallpaper
	rootBgSrc := filepath.Join(themePath, "Wallpapers", "Root", "bg.png")
	if _, err := os.Stat(rootBgSrc); err == nil {
		// Apply to root directory
		rootBgDst := filepath.Join(systemPaths.Root, "bg.png")
		if err := CopyFile(rootBgSrc, rootBgDst); err != nil {
			logger.Printf("Warning: Could not copy Root wallpaper to root: %v", err)
		} else {
			logger.Printf("Imported Root wallpaper to: %s", rootBgDst)
		}

		// Apply to root media directory
		rootMediaDst := filepath.Join(systemPaths.Root, ".media", "bg.png")
		if err := os.MkdirAll(filepath.Dir(rootMediaDst), 0755); err != nil {
			logger.Printf("Warning: Could not create root media directory: %v", err)
		} else if err := CopyFile(rootBgSrc, rootMediaDst); err != nil {
			logger.Printf("Warning: Could not copy Root wallpaper to media dir: %v", err)
		} else {
			logger.Printf("Imported Root wallpaper to: %s", rootMediaDst)
		}
	}

	// Recently Played wallpaper
	rpBgSrc := filepath.Join(themePath, "Wallpapers", "Recently Played", "bg.png")
	if _, err := os.Stat(rpBgSrc); err == nil {
		rpBgDst := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
		if err := os.MkdirAll(filepath.Dir(rpBgDst), 0755); err != nil {
			logger.Printf("Warning: Could not create Recently Played media directory: %v", err)
		} else if err := CopyFile(rpBgSrc, rpBgDst); err != nil {
			logger.Printf("Warning: Could not copy Recently Played wallpaper: %v", err)
		} else {
			logger.Printf("Imported Recently Played wallpaper to: %s", rpBgDst)
		}
	}

	// Tools wallpaper
	toolsBgSrc := filepath.Join(themePath, "Wallpapers", "Tools", "bg.png")
	if _, err := os.Stat(toolsBgSrc); err == nil {
		toolsBgDst := filepath.Join(systemPaths.Tools, ".media", "bg.png")
		if err := os.MkdirAll(filepath.Dir(toolsBgDst), 0755); err != nil {
			logger.Printf("Warning: Could not create Tools media directory: %v", err)
		} else if err := CopyFile(toolsBgSrc, toolsBgDst); err != nil {
			logger.Printf("Warning: Could not copy Tools wallpaper: %v", err)
		} else {
			logger.Printf("Imported Tools wallpaper to: %s", toolsBgDst)
		}
	}

	// Collections wallpaper
	collectionsBgSrc := filepath.Join(themePath, "Wallpapers", "Collections", "bg.png")
	if _, err := os.Stat(collectionsBgSrc); err == nil {
		collectionsBgDst := filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")
		if err := os.MkdirAll(filepath.Dir(collectionsBgDst), 0755); err != nil {
			logger.Printf("Warning: Could not create Collections media directory: %v", err)
		} else if err := CopyFile(collectionsBgSrc, collectionsBgDst); err != nil {
			logger.Printf("Warning: Could not copy Collections wallpaper: %v", err)
		} else {
			logger.Printf("Imported Collections wallpaper to: %s", collectionsBgDst)
		}
	}

	// System wallpapers
	systemsDir := filepath.Join(themePath, "Wallpapers", "Systems")
	if dirEntries, err := os.ReadDir(systemsDir); err == nil {
		// Regular expression to extract system tag from directory name
		reSystemTag := regexp.MustCompile(`^\((.*?)\)$`)

		// Create maps for better lookups
		systemsByTag := make(map[string]system.SystemInfo)
		for _, sys := range systemPaths.Systems {
			if sys.Tag != "" {
				systemsByTag[sys.Tag] = sys
			}
		}

		for _, entry := range dirEntries {
			if entry.IsDir() {
				dirName := entry.Name()

				// Check if directory name is a tag in parentheses
				matches := reSystemTag.FindStringSubmatch(dirName)
				if len(matches) >= 2 {
					tag := matches[1]

					// Try to find system with this tag
					if sys, ok := systemsByTag[tag]; ok {
						// Source wallpaper file
						srcPath := filepath.Join(systemsDir, dirName, "bg.png")
						if _, err := os.Stat(srcPath); err == nil {
							// Destination path
							dstPath := filepath.Join(sys.MediaPath, "bg.png")

							// Ensure media directory exists
							if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
								logger.Printf("Warning: Could not create media directory for system %s: %v", sys.Name, err)
							} else if err := CopyFile(srcPath, dstPath); err != nil {
								logger.Printf("Warning: Could not copy system wallpaper for %s: %v", sys.Name, err)
							} else {
								logger.Printf("Imported system wallpaper for %s to: %s", sys.Name, dstPath)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// ImportIcons imports icons from a theme
func ImportIcons(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Skip if no icons present
	if !manifest.Content.Icons.Present {
		logger.Printf("No icons found in theme, skipping icon import")
		return nil
	}

	// Import each icon using path mappings
	for _, mapping := range manifest.PathMappings.Icons {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source icon file not found: %s", srcPath)
			continue
		}

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.Printf("Warning: Could not create directory for icon: %v", err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy icon: %v", err)
		} else {
			logger.Printf("Imported icon: %s -> %s", srcPath, dstPath)
		}
	}

	return nil
}

// ImportOverlays imports overlays from a theme
func ImportOverlays(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Skip if no overlays present
	if !manifest.Content.Overlays.Present {
		logger.Printf("No overlays found in theme, skipping overlay import")
		return nil
	}

	// Import each overlay using path mappings
	for _, mapping := range manifest.PathMappings.Overlays {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source overlay file not found: %s", srcPath)
			continue
		}

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.Printf("Warning: Could not create directory for overlay: %v", err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy overlay: %v", err)
		} else {
			logger.Printf("Imported overlay: %s -> %s", srcPath, dstPath)
		}
	}

	return nil
}

// ImportFonts imports fonts from a theme
func ImportFonts(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Skip if no fonts present
	if !manifest.Content.Fonts.Present {
		logger.Printf("No fonts found in theme, skipping font import")
		return nil
	}

	// Process each font mapping
	for fontType, mapping := range manifest.PathMappings.Fonts {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source font file not found: %s (%s)", srcPath, fontType)
			continue
		}

		// For backup fonts, only copy if they don't already exist
		if strings.Contains(fontType, "backup") {
			if _, err := os.Stat(dstPath); err == nil {
				logger.Printf("Backup font already exists, skipping: %s", dstPath)
				continue
			}
		}

		// For active fonts, create a backup if one doesn't exist
		if fontType == "og_font" {
			// Backup OG font if needed
			backupPath := filepath.Join(filepath.Dir(dstPath), "font2.backup.ttf")
			if _, err := os.Stat(backupPath); os.IsNotExist(err) {
				if _, err := os.Stat(dstPath); err == nil {
					if err := CopyFile(dstPath, backupPath); err != nil {
						logger.Printf("Warning: Could not create backup of OG font: %v", err)
					} else {
						logger.Printf("Created backup of current OG font: %s", backupPath)
					}
				}
			}
		} else if fontType == "next_font" {
			// Backup Next font if needed
			backupPath := filepath.Join(filepath.Dir(dstPath), "font1.backup.ttf")
			if _, err := os.Stat(backupPath); os.IsNotExist(err) {
				if _, err := os.Stat(dstPath); err == nil {
					if err := CopyFile(dstPath, backupPath); err != nil {
						logger.Printf("Warning: Could not create backup of Next font: %v", err)
					} else {
						logger.Printf("Created backup of current Next font: %s", backupPath)
					}
				}
			}
		}

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.Printf("Warning: Could not create directory for font: %v", err)
			continue
		}

		// Copy the font file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy font file: %v", err)
		} else {
			logger.Printf("Imported font file (%s): %s -> %s", fontType, srcPath, dstPath)
		}
	}

	return nil
}

// ImportSettings imports settings from a theme
func ImportSettings(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Import accent settings if present
	if manifest.Content.Settings.AccentsIncluded {
		if accents, ok := manifest.PathMappings.Settings["accents"]; ok {
			srcPath := filepath.Join(themePath, accents.ThemePath)
			dstPath := accents.SystemPath

			// Check if source file exists
			if _, err := os.Stat(srcPath); os.IsNotExist(err) {
				logger.Printf("Warning: Source accent settings file not found: %s", srcPath)
			} else {
				// Create destination directory
				dstDir := filepath.Dir(dstPath)
				if err := os.MkdirAll(dstDir, 0755); err != nil {
					logger.Printf("Warning: Could not create directory for accent settings: %v", err)
				} else {
					// Copy the file
					if err := CopyFile(srcPath, dstPath); err != nil {
						logger.Printf("Warning: Could not copy accent settings: %v", err)
					} else {
						logger.Printf("Imported accent settings: %s -> %s", srcPath, dstPath)
					}
				}
			}
		}
	}

	// Import LED settings if present
	if manifest.Content.Settings.LEDsIncluded {
		if leds, ok := manifest.PathMappings.Settings["leds"]; ok {
			srcPath := filepath.Join(themePath, leds.ThemePath)
			dstPath := leds.SystemPath

			// Check if source file exists
			if _, err := os.Stat(srcPath); os.IsNotExist(err) {
				logger.Printf("Warning: Source LED settings file not found: %s", srcPath)
			} else {
				// Create destination directory
				dstDir := filepath.Dir(dstPath)
				if err := os.MkdirAll(dstDir, 0755); err != nil {
					logger.Printf("Warning: Could not create directory for LED settings: %v", err)
				} else {
					// Copy the file
					if err := CopyFile(srcPath, dstPath); err != nil {
						logger.Printf("Warning: Could not copy LED settings: %v", err)
					} else {
						logger.Printf("Imported LED settings: %s -> %s", srcPath, dstPath)
					}
				}
			}
		}
	}

	return nil
}

// ImportTheme imports a theme package
func ImportTheme(themeName string) error {
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
		filepath.Join(logsDir, "imports.log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()

	// Create logger
	logger := &Logger{File: logFile}
	logger.Printf("Starting theme import for: %s", themeName)

	// Full path to theme - look in Imports directory
	themePath := filepath.Join(cwd, "Themes", "Imports", themeName)

	// FIRST DELETE ALL EXISTING WALLPAPERS AND ICONS
	// This ensures that the new theme completely replaces previous theme elements
	// rather than just overlaying on top of existing files
	logger.Printf("Deleting all existing wallpapers and icons before applying new theme")

	// Delete all existing backgrounds
	if err := DeleteAllBackgrounds(); err != nil {
		logger.Printf("Warning: Failed to delete existing wallpapers: %v", err)
		// Continue anyway, as we can still try to apply the new theme
	} else {
		logger.Printf("Successfully deleted all existing wallpapers")
	}

	// Delete all existing icons - import from icons package
	if err := icons.DeleteAllIcons(); err != nil {
		logger.Printf("Warning: Failed to delete existing icons: %v", err)
		// Continue anyway, as we can still try to apply the new theme
	} else {
		logger.Printf("Successfully deleted all existing icons")
	}

	// Update the manifest before validation - this scans all the theme files
	// and ensures the manifest reflects the actual contents
	if err := UpdateThemeManifest(themePath); err != nil {
		logger.Printf("Warning: Could not update manifest: %v", err)
		// Continue anyway, as the theme might still be valid
	} else {
		logger.Printf("Successfully updated theme manifest to reflect actual content")
	}

	// Validate theme
	manifest, err := ValidateTheme(themePath, logger)
	if err != nil {
		logger.Printf("Theme validation failed: %v", err)
		return fmt.Errorf("theme validation failed: %w", err)
	}

	// Import theme components
	if err := ImportWallpapers(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing wallpapers: %v", err)
	}

	if err := ImportIcons(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing icons: %v", err)
	}

	if err := ImportOverlays(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing overlays: %v", err)
	}

	if err := ImportFonts(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing fonts: %v", err)
	}

	if err := ImportSettings(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing settings: %v", err)
	}

	logger.Printf("Theme import completed successfully: %s", themeName)

	// Show success message to user
	ui.ShowMessage(fmt.Sprintf("Theme '%s' by %s imported successfully!",
		manifest.ThemeInfo.Name, manifest.ThemeInfo.Author), "5")

	return nil
}
