#!/usr/bin/env python3
"""
Script to build the NextUI-Themes catalog from extracted themes and components
"""

import os
import json
from datetime import datetime
from pathlib import Path
import traceback

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

def load_existing_catalog():
    """Load the existing catalog.json file if it exists"""
    catalog_path = CATALOG_DIR / "catalog.json"
    if os.path.exists(catalog_path):
        try:
            with open(catalog_path, 'r') as f:
                return json.load(f)
        except json.JSONDecodeError:
            print(f"Warning: Existing catalog.json is invalid, creating a new one")
        except Exception as e:
            print(f"Error reading existing catalog: {str(e)}")

    # Return default structure if no catalog exists or there was an error
    return {
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

def extract_theme_info(theme_path, existing_info=None):
    """Extract theme information from the theme directory and preview/manifest"""
    theme_name = os.path.basename(theme_path)

    # Start with existing info if available
    theme_info = {}
    if existing_info:
        theme_info = existing_info.copy()
        # If this theme was previously marked as invalid but has been fixed,
        # we'll want to remove the INVALID flag
        if "INVALID" in theme_info:
            del theme_info["INVALID"]

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
        if existing_info:
            theme_info["INVALID"] = "No manifest.json file found"
        return theme_info if existing_info else None

    # Determine which preview to use
    if os.path.exists(theme_preview_path):
        preview_path = theme_preview_path
    elif os.path.exists(previews_dir_path):
        preview_path = previews_dir_path
    else:
        print(f"Warning: No preview found for theme {theme_name}")
        if existing_info:
            theme_info["INVALID"] = "No preview.png file found"
        return theme_info if existing_info else None

    # Read manifest file
    try:
        with open(manifest_path, 'r') as f:
            manifest_data = json.load(f)
    except json.JSONDecodeError as e:
        error_msg = f"Invalid JSON in manifest: {str(e)}"
        print(f"Error: {error_msg}")
        if existing_info:
            theme_info["INVALID"] = error_msg
        return theme_info if existing_info else None
    except Exception as e:
        error_msg = f"Error reading manifest: {str(e)}"
        print(f"Error: {error_msg}")
        if existing_info:
            theme_info["INVALID"] = error_msg
        return theme_info if existing_info else None

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

    # Update theme info
    theme_info.update({
        "preview_path": preview_rel_path,
        "manifest_path": manifest_rel_path,
        "author": author,
        "description": description,
        "URL": url  # Use "URL" property as requested
    })

    return theme_info

def extract_wallpaper_content(component_path):
    """Extract information about wallpaper content types"""
    content_info = {
        "has_system_wallpapers": False,
        "has_list_wallpapers": False,
        "has_collection_wallpapers": False
    }

    # Check for SystemWallpapers directory
    system_wallpapers_path = os.path.join(component_path, "SystemWallpapers")
    if os.path.exists(system_wallpapers_path):
        for file in os.listdir(system_wallpapers_path):
            if file.lower().endswith('.png') and not file.startswith('.'):
                content_info["has_system_wallpapers"] = True
                break

    # Check for ListWallpapers directory
    list_wallpapers_path = os.path.join(component_path, "ListWallpapers")
    if os.path.exists(list_wallpapers_path):
        for file in os.listdir(list_wallpapers_path):
            if file.lower().endswith('.png') and not file.startswith('.'):
                content_info["has_list_wallpapers"] = True
                break

    # Check for CollectionWallpapers directory
    collection_wallpapers_path = os.path.join(component_path, "CollectionWallpapers")
    if os.path.exists(collection_wallpapers_path):
        for file in os.listdir(collection_wallpapers_path):
            if file.lower().endswith('.png') and not file.startswith('.'):
                content_info["has_collection_wallpapers"] = True
                break

    return content_info

def extract_overlay_systems(component_path):
    """Extract supported system tags from an overlay component's directory structure"""
    systems_path = os.path.join(component_path, "Systems")
    supported_systems = []

    # Check if Systems directory exists
    if os.path.exists(systems_path):
        # List all directories in Systems folder, each is a system tag
        try:
            for entry in os.listdir(systems_path):
                entry_path = os.path.join(systems_path, entry)
                if os.path.isdir(entry_path) and not entry.startswith('.'):
                    # Make sure directory actually contains overlay files
                    has_overlays = False
                    for file in os.listdir(entry_path):
                        if file.lower().endswith('.png') and not file.startswith('.'):
                            has_overlays = True
                            break

                    if has_overlays:
                        supported_systems.append(entry)
        except Exception as e:
            print(f"Error scanning overlay systems: {str(e)}")

    # If we found systems, sort them alphabetically
    if supported_systems:
        supported_systems.sort()

    return supported_systems

def extract_component_info(component_path, component_type, existing_info=None):
    """Extract component information from the component directory and preview/manifest"""
    component_name = os.path.basename(component_path)

    # Start with existing info if available
    component_info = {}
    if existing_info:
        component_info = existing_info.copy()
        # If this component was previously marked as invalid but has been fixed,
        # we'll want to remove the INVALID flag
        if "INVALID" in component_info:
            del component_info["INVALID"]

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
        if existing_info:
            component_info["INVALID"] = "No manifest.json file found"
        return component_info if existing_info else None

    # Determine which preview to use
    if os.path.exists(comp_preview_path):
        preview_path = comp_preview_path
    elif os.path.exists(previews_dir_path):
        preview_path = previews_dir_path
    else:
        print(f"Warning: No preview found for component {component_name} of type {component_type}")
        if existing_info:
            component_info["INVALID"] = "No preview.png file found"
        return component_info if existing_info else None

    # Read manifest file
    try:
        with open(manifest_path, 'r') as f:
            manifest_data = json.load(f)
    except json.JSONDecodeError as e:
        error_msg = f"Invalid JSON in manifest: {str(e)}"
        print(f"Error: {error_msg}")
        if existing_info:
            component_info["INVALID"] = error_msg
        return component_info if existing_info else None
    except Exception as e:
        error_msg = f"Error reading manifest: {str(e)}"
        print(f"Error: {error_msg}")
        if existing_info:
            component_info["INVALID"] = error_msg
        return component_info if existing_info else None

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

    # Extract supported systems for overlays
    systems = None
    if component_type == "Overlays":
        systems = extract_overlay_systems(component_path)
        # Also try to get systems from manifest if available
        if not systems and "content" in manifest_data and "systems" in manifest_data["content"]:
            try:
                systems = manifest_data["content"]["systems"]
                if isinstance(systems, list):
                    systems.sort()  # Sort alphabetically
            except Exception as e:
                print(f"Error extracting systems from manifest: {str(e)}")

    # Add this to extract_component_info when processing Wallpapers
    if component_type == "Wallpapers":
        wallpaper_content = extract_wallpaper_content(component_path)
        if wallpaper_content["has_list_wallpapers"]:
            component_info["has_list_wallpapers"] = True

    # Update component info
    component_info.update({
        "preview_path": preview_rel_path,
        "manifest_path": manifest_rel_path,
        "author": author,
        "description": description,
        "URL": url  # Use "URL" property as requested
    })

    # Add systems for overlays
    if systems:
        component_info["systems"] = systems

    return component_info

def main():
    """Main function to build the catalog"""
    # Create necessary directories
    os.makedirs(CATALOG_DIR, exist_ok=True)

    # Load existing catalog
    catalog = load_existing_catalog()

    # Process themes
    if os.path.exists(THEMES_DIR):
        # First, scan filesystem for themes
        processed_themes = set()
        for theme_entry in os.listdir(THEMES_DIR):
            theme_path = THEMES_DIR / theme_entry
            if os.path.isdir(theme_path) and not theme_entry.startswith('.') and theme_entry not in ['previews', 'manifests']:
                # Get existing theme info if available
                existing_info = catalog["themes"].get(theme_entry)

                # Process theme
                theme_info = extract_theme_info(theme_path, existing_info)
                if theme_info:
                    catalog["themes"][theme_entry] = theme_info
                    processed_themes.add(theme_entry)

        # Now handle repository-based themes that might not be on filesystem
        for theme_name, theme_info in list(catalog["themes"].items()):
            if theme_name not in processed_themes and "repository" in theme_info and "commit" in theme_info:
                # Preserve repository-based themes
                print(f"Preserving repository-based theme: {theme_name}")
                processed_themes.add(theme_name)
            elif theme_name not in processed_themes and "INVALID" not in theme_info:
                # Theme was removed from filesystem and isn't repository-based
                print(f"Removing theme that no longer exists: {theme_name}")
                del catalog["themes"][theme_name]

    # Process components
    if os.path.exists(COMPONENTS_DIR):
        for comp_type in COMPONENT_TYPES.keys():
            component_dir = COMPONENTS_DIR / comp_type
            comp_type_lower = COMPONENT_TYPES[comp_type]

            # Track which components we've processed
            processed_components = set()

            if os.path.exists(component_dir):
                for component_entry in os.listdir(component_dir):
                    component_path = component_dir / component_entry
                    if os.path.isdir(component_path) and not component_entry.startswith('.') and component_entry not in ['previews', 'manifests']:
                        # Get existing component info if available
                        existing_info = catalog["components"][comp_type_lower].get(component_entry)

                        # Process component
                        component_info = extract_component_info(component_path, comp_type, existing_info)
                        if component_info:
                            catalog["components"][comp_type_lower][component_entry] = component_info
                            processed_components.add(component_entry)

            # Now handle repository-based components that might not be on filesystem
            for comp_name, comp_info in list(catalog["components"][comp_type_lower].items()):
                if comp_name not in processed_components and "repository" in comp_info and "commit" in comp_info:
                    # Special handling for repository-based overlays
                    if comp_type_lower == "overlays":
                        manifest_path = os.path.join(COMPONENTS_DIR, comp_type, "manifests", f"{comp_name}.json")
                        if os.path.exists(manifest_path):
                            try:
                                with open(manifest_path, 'r') as f:
                                    manifest_data = json.load(f)

                                # Extract systems from manifest content
                                if "content" in manifest_data and "systems" in manifest_data["content"]:
                                    systems = manifest_data["content"]["systems"]
                                    if isinstance(systems, list) and systems:
                                        systems.sort()  # Sort alphabetically
                                        comp_info["systems"] = systems
                                        print(f"Added systems for repository-based overlay: {comp_name}")
                            except Exception as e:
                                print(f"Error processing manifest for {comp_name}: {str(e)}")

                    # Preserve repository-based components
                    print(f"Preserving repository-based component: {comp_name}")
                    processed_components.add(comp_name)
                elif comp_name not in processed_components and "INVALID" not in comp_info:
                    # Component was removed from filesystem and isn't repository-based
                    print(f"Removing component that no longer exists: {comp_name}")
                    del catalog["components"][comp_type_lower][comp_name]

    # Update timestamp
    catalog["last_updated"] = datetime.utcnow().isoformat() + "Z"

    # Write catalog.json
    catalog_path = CATALOG_DIR / "catalog.json"
    with open(catalog_path, 'w') as f:
        json.dump(catalog, f, indent=2)

    print(f"Catalog updated successfully at {datetime.utcnow().isoformat()}Z")

if __name__ == "__main__":
    main()