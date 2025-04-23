# NextUI-Themes

Repository for NextUI themes and component packages.

<!-- GALLERY_START -->
## Featured Themes

<!-- This will be replaced by the gallery generator -->

[Browse All Themes](.github/index/index.md)
<!-- GALLERY_END -->

## How to Install

### Themes

1. Download the theme package (click on a theme image above)
2. Extract the theme to your TrimUI Brick device in the `Tools/tg5040/Theme-Manager.pak/Themes/` directory
3. Open the Theme Manager app on your device
4. Select the theme and apply it

### Components

1. Download the component package
2. Extract to your TrimUI Brick device in the appropriate component directory:
   - `Tools/tg5040/Theme-Manager.pak/Components/[Type]/`
3. Open the Theme Manager app on your device
4. Select "Components" and the appropriate component type
5. Choose the component and apply it

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
