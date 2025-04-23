#!/usr/bin/env python3
"""
Script to build the NextUI-Themes catalog from extracted themes and components
"""

import os
import json
from datetime import datetime
from pathlib import Path

# Base paths
REPO_ROOT = Path(__file__).resolve().parent.parent.parent
CATALOG_DIR = REPO_ROOT / "Catalog"  # Changed from "catalog" to "Catalog"
THEMES_DIR = CATALOG_DIR / "Themes"  # Moved into Catalog
COMPONENTS_DIR = CATALOG_DIR / "Components"  # Moved into Catalog
UPLOADS_DIR = REPO_ROOT / "Uploads"  # Path to uploaded zips

# GitHub raw content base URL
GITHUB_RAW_URL = "https://github.com/Leviathanium/NextUI-Themes/raw/main/Uploads"

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
    """Extract theme information from the theme directory and preview/manifest"""
    theme_name = os.path.basename(theme_path)

    # Check for manifest in the theme directory
    theme_manifest_path = os.path.join(theme_path, "manifest.json")
    theme_preview_path = os.path.join(theme_path, "preview.png")

    # Check for manifest in the manifests directory
    manifests_dir_path = os.path.join(THEMES_DIR, "manifests", f"{theme_name}.json")

    # Check for preview in the previews directory
    previews_dir_path = os.path.join(THEMES_DIR, "previews", f"{theme_name}.png")

    # Determine which manifest to use
    if os.path.exists(theme_manifest_path):
        manifest_path = theme_manifest_path
    elif os.path.exists(manifests_dir_path):
        manifest_path = manifests_dir_path
    else:
        print(f"Warning: No manifest found for theme {theme_name}")
        return None

    # Determine which preview to use
    if os.path.exists(theme_preview_path):
        preview_path = theme_preview_path
    elif os.path.exists(previews_dir_path):
        preview_path = previews_dir_path
    else:
        print(f"Warning: No preview found for theme {theme_name}")
        return None

    # Read manifest file
    try:
        with open(manifest_path, 'r') as f:
            manifest_data = json.load(f)
    except json.JSONDecodeError:
        print(f"Error: Invalid JSON in {manifest_path}")
        return None
    except Exception as e:
        print(f"Error reading manifest {manifest_path}: {str(e)}")
        return None

    # Extract author and description
    author = manifest_data.get("theme_info", {}).get("author", "Unknown")
    description = manifest_data.get("theme_info", {}).get("name", theme_name)

    # Create relative paths for catalog
    preview_rel_path = f"Catalog/Themes/previews/{theme_name}.png"
    manifest_rel_path = f"Catalog/Themes/manifests/{theme_name}.json"

    # Find corresponding ZIP file in Uploads directory
    zip_filename = f"{theme_name}.zip"
    zip_path = os.path.join(UPLOADS_DIR, "Themes", zip_filename)

    # Generate URL for the ZIP file
    url = f"{GITHUB_RAW_URL}/Themes/{zip_filename}"

    # Create theme info
    theme_info = {
        "preview_path": preview_rel_path,
        "manifest_path": manifest_rel_path,
        "author": author,
        "description": description,
        "URL": url  # Use "URL" property as requested
    }

    # If this theme was in the old catalog and had repository info, preserve it
    old_catalog_path = os.path.join(CATALOG_DIR, "catalog.json")
    if os.path.exists(old_catalog_path):
        try:
            with open(old_catalog_path, 'r') as f:
                old_catalog = json.load(f)

            if theme_name in old_catalog.get("themes", {}):
                old_theme_info = old_catalog["themes"][theme_name]
                if "repository" in old_theme_info:
                    theme_info["repository"] = old_theme_info["repository"]
                if "commit" in old_theme_info:
                    theme_info["commit"] = old_theme_info["commit"]
        except:
            pass  # If anything goes wrong, just continue without the repository info

    return theme_info

def extract_component_info(component_path, component_type):
    """Extract component information from the component directory and preview/manifest"""
    component_name = os.path.basename(component_path)

    # Check for manifest in the component directory
    comp_manifest_path = os.path.join(component_path, "manifest.json")
    comp_preview_path = os.path.join(component_path, "preview.png")

    # Check for manifest in the manifests directory
    manifests_dir_path = os.path.join(COMPONENTS_DIR, component_type, "manifests", f"{component_name}.json")

    # Check for preview in the previews directory
    previews_dir_path = os.path.join(COMPONENTS_DIR, component_type, "previews", f"{component_name}.png")

    # Determine which manifest to use
    if os.path.exists(comp_manifest_path):
        manifest_path = comp_manifest_path
    elif os.path.exists(manifests_dir_path):
        manifest_path = manifests_dir_path
    else:
        print(f"Warning: No manifest found for component {component_name} of type {component_type}")
        return None

    # Determine which preview to use
    if os.path.exists(comp_preview_path):
        preview_path = comp_preview_path
    elif os.path.exists(previews_dir_path):
        preview_path = previews_dir_path
    else:
        print(f"Warning: No preview found for component {component_name} of type {component_type}")
        return None

    # Read manifest file
    try:
        with open(manifest_path, 'r') as f:
            manifest_data = json.load(f)
    except json.JSONDecodeError:
        print(f"Error: Invalid JSON in {manifest_path}")
        return None
    except Exception as e:
        print(f"Error reading manifest {manifest_path}: {str(e)}")
        return None

    # Extract author and description
    author = manifest_data.get("component_info", {}).get("author", "Unknown")
    description = manifest_data.get("component_info", {}).get("name", component_name)

    # Create relative paths for catalog
    preview_rel_path = f"Catalog/Components/{component_type}/previews/{component_name}.png"
    manifest_rel_path = f"Catalog/Components/{component_type}/manifests/{component_name}.json"

    # Find corresponding ZIP file in Uploads directory
    zip_filename = f"{component_name}.zip"
    zip_path = os.path.join(UPLOADS_DIR, "Components", component_type, zip_filename)

    # Generate URL for the ZIP file
    url = f"{GITHUB_RAW_URL}/Components/{component_type}/{zip_filename}"

    # Return component info
    component_info = {
        "preview_path": preview_rel_path,
        "manifest_path": manifest_rel_path,
        "author": author,
        "description": description,
        "URL": url  # Use "URL" property as requested
    }

    # If this component was in the old catalog and had repository info, preserve it
    old_catalog_path = os.path.join(CATALOG_DIR, "catalog.json")
    if os.path.exists(old_catalog_path):
        try:
            with open(old_catalog_path, 'r') as f:
                old_catalog = json.load(f)

            comp_type_lower = COMPONENT_TYPES[component_type]
            if comp_name in old_catalog.get("components", {}).get(comp_type_lower, {}):
                old_comp_info = old_catalog["components"][comp_type_lower][comp_name]
                if "repository" in old_comp_info:
                    component_info["repository"] = old_comp_info["repository"]
                if "commit" in old_comp_info:
                    component_info["commit"] = old_comp_info["commit"]
        except:
            pass  # If anything goes wrong, just continue without the repository info

    return component_info

def main():
    """Main function to build the catalog"""
    # Create necessary directories
    os.makedirs(CATALOG_DIR, exist_ok=True)

    # Process themes
    if os.path.exists(THEMES_DIR):
        for theme_entry in os.listdir(THEMES_DIR):
            theme_path = THEMES_DIR / theme_entry
            if os.path.isdir(theme_path) and not theme_entry.startswith('.') and theme_entry not in ['previews', 'manifests']:
                theme_info = extract_theme_info(theme_path)
                if theme_info:
                    catalog["themes"][theme_entry] = theme_info

    # Process components
    if os.path.exists(COMPONENTS_DIR):
        for comp_type in COMPONENT_TYPES.keys():
            component_dir = COMPONENTS_DIR / comp_type
            comp_type_lower = COMPONENT_TYPES[comp_type]

            if os.path.exists(component_dir):
                for component_entry in os.listdir(component_dir):
                    component_path = component_dir / component_entry
                    if os.path.isdir(component_path) and not component_entry.startswith('.') and component_entry not in ['previews', 'manifests']:
                        component_info = extract_component_info(component_path, comp_type)
                        if component_info:
                            catalog["components"][comp_type_lower][component_entry] = component_info

    # Write catalog.json
    catalog_path = CATALOG_DIR / "catalog.json"
    with open(catalog_path, 'w') as f:
        json.dump(catalog, f, indent=2)

    print(f"Catalog updated successfully at {datetime.utcnow().isoformat()}Z")

if __name__ == "__main__":
    main()
