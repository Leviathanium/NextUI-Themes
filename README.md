# NextUI-Themes

Repository for NextUI themes and component packages.

## Gallery

### Themes

| Name | Author | Preview |
|------|--------|---------|
| Consolized.theme | Gamnrd | [Preview](Catalog/Themes/previews/Consolized.theme.png) |
| Pearly Gate.theme | Shin | [Preview](Catalog/Themes/previews/Pearly%20Gate.theme.png) |
| Default.theme | NextUI | [Preview](Catalog/Themes/previews/Default.theme.png) |

### Wallpapers

| Name | Author | Preview |
|------|--------|---------|
| Retro-Programmer.bg | Leviathan | [Preview](Catalog/Components/Wallpapers/previews/Retro-Programmer.bg.png) |
| Sunset-Forest.bg | Leviathan | [Preview](Catalog/Components/Wallpapers/previews/Sunset-Forest.bg.png) |
| Retro-Mario-Chill.bg | Leviathan | [Preview](Catalog/Components/Wallpapers/previews/Retro-Mario-Chill.bg.png) |
| Arcade-Dark.bg | Leviathan | [Preview](Catalog/Components/Wallpapers/previews/Arcade-Dark.bg.png) |
| Firewatch.bg | Leviathan | [Preview](Catalog/Components/Wallpapers/previews/Firewatch.bg.png) |
| Cozy.bg | Leviathan | [Preview](Catalog/Components/Wallpapers/previews/Cozy.bg.png) |
| Blackstreets.bg | Leviathan | [Preview](Catalog/Components/Wallpapers/previews/Blackstreets.bg.png) |
| Default.bg | Leviathan | [Preview](Catalog/Components/Wallpapers/previews/Default.bg.png) |

### Overlays

| Name | Author | Preview |
|------|--------|---------|
| Vanilla.over | Shin | [Preview](Catalog/Components/Overlays/previews/Vanilla.over.png) |

### Icons

*No icon components found in the catalog.*

### Accents

*No accent components found in the catalog.*

## About This Repository

This repository hosts themes and component packages for NextUI devices. It uses an automated workflow to process uploaded zip files.

## Directory Structure

- `Themes/` - Extracted theme packages
- `Components/` - Extracted component packages by type
  - `Accents/` - Accent color packages (.acc)
  - `LEDs/` - LED settings packages (.led)
  - `Icons/` - Icon packages (.icon)
  - `Fonts/` - Font packages (.font)
  - `Wallpapers/` - Wallpaper packages (.bg)
  - `Overlays/` - Overlay packages (.over)
- `Uploads/` - Temporary location for zip files (automatically processed)
- `catalog/` - Automatically generated catalog of available themes and components

## How to Submit Your Theme

### Option 1: Submit a Zip File

1. Prepare your theme or component package following the guidelines
2. Create a zip file of your theme folder with the appropriate extension:
   - Full theme: `your-theme-name.theme.zip`
   - Wallpaper component: `your-component-name.bg.zip`
   - Icon component: `your-component-name.icon.zip`
   - Accent component: `your-component-name.acc.zip`
   - LED component: `your-component-name.led.zip`
   - Font component: `your-component-name.font.zip`
   - Overlay component: `your-component-name.over.zip`
3. Submit your zip file through:
   - Discord submission in the #theme-submissions channel
   - Open an issue in this repository with a link to your zip file

### Option 2: Submit as a Git Repository (Advanced)

For developers familiar with Git:

1. Create a Git repository with your theme or component
2. Open an issue in this repository with your repository URL
3. A maintainer will add it as a submodule

## Theme Structure Guidelines

All themes and components must include:

1. A valid `manifest.json` file
2. A `preview.png` image for the catalog

### Full Theme Package (.theme)

A full theme package should contain:
- `manifest.json` - Theme metadata
- `preview.png` - Theme preview image
- `Wallpapers/` - Directory for wallpaper images
- `Icons/` - Directory for icon images
- `Fonts/` - Directory for font files
- `Overlays/` - Directory for overlays

### Component Packages

A component package should contain only the relevant files for that component type:
- `manifest.json` - Component metadata
- `preview.png` - Component preview image
- Component-specific directories and files

## File Naming Guidelines

- System files must include the system tag in parentheses
  - Example: `Game Boy Advance (GBA).png`
- Special system files use literal names
  - `Root.png`, `Collections.png`, etc.
- Collection files use the collection name
  - Example: `Favorites.png`

## Need Help?

If you need assistance, open an issue in this repository or join our Discord community.
