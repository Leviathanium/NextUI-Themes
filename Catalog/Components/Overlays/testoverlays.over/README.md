# NextUI Overlay Component Template

This repository contains a template for creating overlay components for NextUI devices.

## What are Overlays?

Overlays are system-specific transparent PNG images that can be displayed over games to show decorative elements.

## Directory Structure

- `manifest.json` - Component metadata and path mappings
- `preview.png` - Component preview image (replace with your preview)
- `Systems/` - Contains system-specific overlay directories
  - `MGBA/` - Overlays for Game Boy Advance
  - `SFC/` - Overlays for Super Nintendo/Super Famicom
  - `MD/` - Overlays for Sega Genesis/Mega Drive
  - (Add more system directories as needed)

## Getting Started

1. Replace the placeholder `preview.png` with a preview of your overlay component
2. Update the `manifest.json` with your component information (name, author, etc.)
3. Add your overlay PNG files to the appropriate system directories

## Important Notes

- System tags should match those used by NextUI (e.g., MGBA, SFC, MD, etc.)
- Overlay PNG files should be transparent where needed
- Make sure the resolution matches the target device screen (1024x768 for TrimUI Brick)

## For more information, see the full component documentation
