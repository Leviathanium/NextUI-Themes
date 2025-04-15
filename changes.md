## Current Streamlined UI Structure

```markdown
NextUI Theme Manager
├── Themes
│   ├── Apply Theme
│   │   ├── Theme Selection
│   │   └── Apply Confirmation
│   ├── Extract Theme Components
│   │   ├── Theme Selection
│   │   ├── Component Selection
│   │   │   ├── All Components
│   │   │   ├── Wallpapers Only
│   │   │   ├── Icons Only
│   │   │   ├── Accents Only
│   │   │   ├── LEDs Only
│   │   │   └── Fonts Only
│   │   └── Extract Confirmation
│   └── Export Current Theme
│       └── Export Confirmation
├── Components
│   ├── Apply Component
│   │   ├── Component Type Selection
│   │   │   ├── Wallpaper Pack
│   │   │   ├── Icon Pack
│   │   │   ├── Accent Pack
│   │   │   ├── LED Pack
│   │   │   └── Font Pack
│   │   ├── Component Selection
│   │   └── Apply Confirmation
│   └── Export Component
│       ├── Component Type Selection
│       │   ├── Wallpaper Pack
│       │   ├── Icon Pack
│       │   ├── Accent Pack
│       │   ├── LED Pack
│       │   └── Font Pack
│       └── Export Confirmation
└── Reset
    └── Reset Options
        ├── Reset Everything
        ├── Reset Wallpapers Only
        ├── Reset Icons Only
        └── Reset Confirmation
```

## New Proposed Structure with a little more clarity

```markdown
NextUI Theme Manager
├── Themes
│   ├── Browse Themes
│   │   ├── Theme Selection
│   │   └── Apply Confirmation
│   ├── Download Themes <-- LOGIC NOT YET IMPLEMENTED
│   │   └── Downloading... 
|   │   │   └── Download complete!
│   ├── Extract Components <-- Removed the word "Theme" for brevity
│   │   ├── Component Selection
│   │   │   ├── All Components
│   │   │   ├── Wallpapers Only
│   │   │   ├── Icons Only
│   │   │   ├── Accents Only
│   │   │   ├── LEDs Only
│   │   │   └── Fonts Only
│   │   └── Extract Confirmation
│   └── Export Theme
│       └── Export Confirmation
├── Components
|   ├── Wallpapers
|   ├── Icons
|   ├── Accents
|   ├── LEDs
|   ├── Fonts
|   └── Component Options
|       ├── Browse <-- Replaced "Apply" with "Browse"
|       ├── Download <-- LOGIC NOT YET IMPLEMENTED
|       |   └── Downloading... 
|       |       └── Download complete!
│       └── Export
│           └── Export Confirmation
└── Reset
    └── Reset Options
        ├── Reset Everything
        ├── Reset Wallpapers Only
        ├── Reset Icons Only
        └── Reset Confirmation
```