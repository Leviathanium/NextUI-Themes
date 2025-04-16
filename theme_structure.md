
# Component Export Directory Structure

Let me explain what happens when you export each component type:

## 1. Wallpapers (.bg)
```
wallpaper_name.bg/
├── manifest.json             # Component manifest with metadata
├── preview.png               # Preview image (copy of Recently Played bg)
├── SystemWallpapers/         # System wallpapers directory
│   ├── Root.png
│   ├── Recently Played.png
│   ├── Tools.png
│   └── Game Boy Advance (GBA).png
└── CollectionWallpapers/     # Collection wallpapers directory
    └── [CollectionName].png
```
**Export location**: `Theme-Manager/Exports/[ExportName].bg/`

## 2. Icons (.icon)
```
icon_name.icon/
├── manifest.json             # Component manifest with metadata
├── SystemIcons/              # System menu icons
│   ├── Collections.png
│   ├── Recently Played.png
│   └── Game Boy Advance (GBA).png
├── ToolIcons/                # Icons for tools
│   └── [ToolName].png
└── CollectionIcons/          # Icons for collections
    └── [CollectionName].png
```
**Export location**: `Theme-Manager/Exports/[ExportName].icon/`

## 3. LEDs (.led)
```
led_name.led/
└── manifest.json             # Component manifest with LED settings embedded
```
Note that LED settings are now embedded directly in the manifest.json file, not in a separate file.

**Export location**: `Theme-Manager/Exports/[ExportName].led/`

## 4. Accents (.acc)
```
accent_name.acc/
└── manifest.json             # Component manifest with accent colors embedded
```
Similarly, accent colors are now embedded directly in the manifest.json file.

**Export location**: `Theme-Manager/Exports/[ExportName].acc/`

## 5. Fonts (.font)
```
font_name.font/
├── manifest.json             # Component manifest
├── OG.ttf                    # OG font replacement
├── OG.backup.ttf             # Backup of original OG font
├── Next.ttf                  # Next font replacement
└── Next.backup.ttf           # Backup of original Next font
```
**Export location**: `Theme-Manager/Exports/[ExportName].font/`

## 6. Full Themes (.theme)
```
theme_name.theme/
├── manifest.json             # Theme manifest with all component information
├── preview.png               # Preview image
├── Wallpapers/               # Contains all wallpapers
├── Icons/                    # Contains all icons
├── Fonts/                    # Contains font files
└── Settings/                 # Contains settings files (legacy support)
```
**Export location**: `Theme-Manager/Exports/[ExportName].theme/`

The big change in this new version is that each component now has a standardized manifest.json file that identifies its type and contains all necessary mapping information. For LEDs and Accents, the settings are now stored directly in the manifest rather than in separate files, making the structure cleaner and more consistent.
