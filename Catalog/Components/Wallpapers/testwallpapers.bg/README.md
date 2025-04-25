# NextUI Wallpaper Component Template

This repository contains a template for creating wallpaper components for NextUI devices.

## What are Wallpaper Components?

Wallpaper components provide background images for various screens throughout the NextUI interface, including the main menu, system screens, and collections.

## Directory Structure

- `manifest.json` - Component metadata and wallpaper path mappings
- `preview.png` - Component preview image (replace with your preview)
- `SystemWallpapers/` - Wallpapers for main menu and systems
  - `Root.png` - Main menu background
  - `Root-Media.png` - Alternative main menu background
  - `Recently Played.png` - Recently played list background
  - `Tools.png` - Tools menu background
  - `Collections.png` - Collections menu background
  - `Game Boy Advance (GBA).png` - Main menu system background (include tag)
- `ListWallpapers/` - System-list specific backgrounds
  - `Game Boy Advance (GBA)-list.png` - System-list specific background (include tag and `-list`)
- `CollectionWallpapers/` - Wallpapers for user collections
  - `Favorites.png` - Background for the Favorites collection

## Getting Started

1. Replace the placeholder `preview.png` with a preview that demonstrates your wallpaper style
2. Update the `manifest.json` with your component information (name, author, etc.)
3. Replace the placeholder wallpaper files with your custom wallpapers
4. Add additional wallpapers as needed, following the naming conventions
5. Update the `content` and `path_mappings` sections in manifest.json for any new wallpapers
6. When complete, commit your changes and update the `NextUI-Themes` catalog

## Wallpaper Requirements

- PNG format
- Resolution should match your device's screen (1024x768 for TrimUI Brick)
- All wallpapers should follow a consistent style or theme
- System wallpapers must include the system tag in parentheses (e.g., Game Boy Advance (GBA).png)
- Collection wallpapers must match the collection folder name exactly (e.g., Favorites.png)

## Important Wallpapers

Make sure to include these critical wallpapers:
- `Root.png` - Main menu background
- `Recently Played.png` - Recently played list background
- `Tools.png` - Tools menu background
- `Collections.png` - Collections menu background

## Tips for Good Wallpapers

- Avoid overly busy or cluttered backgrounds that make text difficult to read
- Ensure sufficient contrast with UI elements
- Consider how your wallpapers will look with different accent colors
- Creating matching sets for all systems creates a cohesive experience

## For more information, see the [full component documentation.](https://github.com/Leviathanium/NextUI-Theme-Manager/blob/main/documents/COMPONENT_BUILDING.md)
