#!/usr/bin/env python3
"""
Script to build the NextUI-Themes catalog from extracted themes and components
"""

import os
import json
import shutil
from datetime import datetime
from pathlib import Path

# Base paths
REPO_ROOT = Path(__file__).resolve().parent.parent.parent
THEMES_DIR = REPO_ROOT / "Themes"
COMPONENTS_DIR = REPO_ROOT / "Components"
CATALOG_DIR = REPO_ROOT / "Catalog"
PREVIEWS_DIR = CATALOG_DIR / "previews"
MANIFESTS_DIR = CATALOG_DIR / "manifests"

# Component types mapping
COMPONENT_TYPES = {
    "Accents": "accents",
    "LEDs": "leds",
    "Icons": "icons",
    "Fonts": "fonts",
    "Wallpapers": "wallpapers",
    "Overlays": "overlays"
}

# Initialize catalog structure
catalog = {
    "last_updated": datetime.utcnow().isoformat() + "Z",
    "themes": {},
    "components": {
        "accents": {},
        "leds": {},
        "icons": {},
        "fonts": {},
        "wallpapers": {},
        "overlays": {}
    }
}

def extract_theme_info(theme_path):
    """Extract theme information and copy preview/manifest"""
    theme_name = os.path.basename(theme_path)
    
    # Check if this is a valid theme directory (has manifest.json and preview.png)
    manifest_path = os.path.join(theme_path, "manifest.json")
    preview_path = os.path.join(theme_path, "preview.png")
    
    if not os.path.exists(manifest_path) or not os.path.exists(preview_path):
        print(f"Warning: {theme_path} is missing manifest.json or preview.png")
        return None
    
    # Read manifest file
    try:
        with open(manifest_path, 'r') as f:
            manifest_data = json.load(f)
    except json.JSONDecodeError:
        print(f"Error: Invalid JSON in {manifest_path}")
        return None
    
    # Extract author and description
    author = manifest_data.get("theme_info", {}).get("author", "Unknown")
    description = manifest_data.get("theme_info", {}).get("name", theme_name)
    
    # Create destination paths
    preview_dest = PREVIEWS_DIR / "themes" / theme_name
    manifest_dest = MANIFESTS_DIR / "themes" / theme_name
    
    # Create directories
    os.makedirs(preview_dest, exist_ok=True)
    os.makedirs(manifest_dest, exist_ok=True)
    
    # Copy files
    shutil.copy2(preview_path, preview_dest / "preview.png")
    shutil.copy2(manifest_path, manifest_dest / "manifest.json")
    
    # Create relative paths for catalog
    preview_rel_path = f"Catalog/previews/themes/{theme_name}/preview.png"
    manifest_rel_path = f"Catalog/manifests/themes/{theme_name}/manifest.json"
    
    # Return theme info
    theme_info = {
        "preview_path": preview_rel_path,
        "manifest_path": manifest_rel_path,
        "author": author,
        "description": description
    }
    
    return theme_info

def extract_component_info(component_path, component_type):
    """Extract component information and copy preview/manifest"""
    component_name = os.path.basename(component_path)
    
    # Check if this is a valid component directory (has manifest.json and preview.png)
    manifest_path = os.path.join(component_path, "manifest.json")
    preview_path = os.path.join(component_path, "preview.png")
    
    if not os.path.exists(manifest_path) or not os.path.exists(preview_path):
        print(f"Warning: {component_path} is missing manifest.json or preview.png")
        return None
    
    # Read manifest file
    try:
        with open(manifest_path, 'r') as f:
            manifest_data = json.load(f)
    except json.JSONDecodeError:
        print(f"Error: Invalid JSON in {manifest_path}")
        return None
    
    # Extract author and description
    author = manifest_data.get("component_info", {}).get("author", "Unknown")
    description = manifest_data.get("component_info", {}).get("name", component_name)
    
    # Create destination paths
    preview_dest = PREVIEWS_DIR / "components" / component_type / component_name
    manifest_dest = MANIFESTS_DIR / "components" / component_type / component_name
    
    # Create directories
    os.makedirs(preview_dest, exist_ok=True)
    os.makedirs(manifest_dest, exist_ok=True)
    
    # Copy files
    shutil.copy2(preview_path, preview_dest / "preview.png")
    shutil.copy2(manifest_path, manifest_dest / "manifest.json")
    
    # Create relative paths for catalog
    preview_rel_path = f"Catalog/previews/components/{component_type}/{component_name}/preview.png"
    manifest_rel_path = f"Catalog/manifests/components/{component_type}/{component_name}/manifest.json"
    
    # Return component info
    component_info = {
        "preview_path": preview_rel_path,
        "manifest_path": manifest_rel_path,
        "author": author,
        "description": description
    }
    
    return component_info

def main():
    """Main function to build the catalog"""
    # Process themes
    for theme_name in os.listdir(THEMES_DIR):
        theme_path = THEMES_DIR / theme_name
        if os.path.isdir(theme_path) and not theme_name.startswith('.'):
            theme_info = extract_theme_info(theme_path)
            if theme_info:
                catalog["themes"][theme_name] = theme_info
    
    # Process components
    for comp_dir, comp_type in COMPONENT_TYPES.items():
        component_dir = COMPONENTS_DIR / comp_dir
        if os.path.exists(component_dir):
            for component_name in os.listdir(component_dir):
                component_path = component_dir / component_name
                if os.path.isdir(component_path) and not component_name.startswith('.'):
                    component_info = extract_component_info(component_path, comp_type)
                    if component_info:
                        catalog["components"][comp_type][component_name] = component_info
    
    # Write catalog.json
    with open(CATALOG_DIR / "catalog.json", 'w') as f:
        json.dump(catalog, f, indent=2)
    
    print(f"Catalog updated successfully at {datetime.utcnow().isoformat()}Z")

if __name__ == "__main__":
    main()
