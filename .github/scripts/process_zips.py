#!/usr/bin/env python3
"""
Script to process theme and component zip files
"""

import os
import json
import zipfile
import shutil
import glob
from pathlib import Path

# Base paths
REPO_ROOT = Path(__file__).resolve().parent.parent.parent
UPLOADS_DIR = REPO_ROOT / "Uploads"
CATALOG_DIR = REPO_ROOT / "Catalog"  # Changed from "catalog" to "Catalog"
THEMES_DIR = CATALOG_DIR / "Themes"  # Moved into Catalog
COMPONENTS_DIR = CATALOG_DIR / "Components"  # Moved into Catalog

# Component type to directory mapping
COMPONENT_DIRS = {
    ".bg": "Wallpapers",
    ".icon": "Icons",
    ".acc": "Accents",
    ".led": "LEDs",
    ".font": "Fonts",
    ".over": "Overlays"
}

# Flag to keep zip files after processing (set to True to keep them)
KEEP_ZIP_FILES = True

def validate_zip_contents(zip_path):
    """Validate that the zip file contains required files"""
    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
        file_list = zip_ref.namelist()

        # Check if manifest.json exists (at any level)
        manifest_exists = any('manifest.json' in f for f in file_list)

        # Check if preview.png exists (at any level)
        preview_exists = any('preview.png' in f for f in file_list)

        return manifest_exists and preview_exists

def extract_without_nested_dirs(zip_path, dest_path):
    """Extract zip content while handling nested directories"""
    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
        file_list = zip_ref.namelist()

        # Identify if there's a common parent directory in the zip
        # We'll use this to determine if we need to strip a directory level
        common_parent = None
        for item in file_list:
            # Skip __MACOSX entries
            if "__MACOSX" in item:
                continue

            parts = item.split('/')
            if len(parts) > 1 and parts[0] and not common_parent:
                common_parent = parts[0]
            elif len(parts) > 1 and parts[0] and parts[0] != common_parent:
                common_parent = None
                break

        # Extract files, skipping __MACOSX directories and stripping common parent if needed
        for item in file_list:
            # Skip __MACOSX entries
            if "__MACOSX" in item:
                continue

            # Determine the correct extraction path
            if common_parent and item.startswith(common_parent + '/'):
                # Strip the common parent directory
                target_path = os.path.join(dest_path, item[len(common_parent) + 1:])
            else:
                target_path = os.path.join(dest_path, item)

            # Skip directories, we'll create them as needed
            if item.endswith('/'):
                continue

            # Create the directory structure
            os.makedirs(os.path.dirname(target_path), exist_ok=True)

            # Extract the file
            with zip_ref.open(item) as source_file, open(target_path, 'wb') as target_file:
                shutil.copyfileobj(source_file, target_file)

def copy_preview_and_manifest(src_path, component_type=None, is_theme=False):
    """Copy preview.png and manifest.json to their respective directories"""
    if not os.path.exists(src_path):
        return False

    # Get base name of the component or theme
    name = os.path.basename(src_path)

    # Determine destination paths
    if is_theme:
        preview_dest = THEMES_DIR / "previews" / f"{name}.png"
        manifest_dest = THEMES_DIR / "manifests" / f"{name}.json"
    else:
        preview_dest = COMPONENTS_DIR / component_type / "previews" / f"{name}.png"
        manifest_dest = COMPONENTS_DIR / component_type / "manifests" / f"{name}.json"

    # Make sure destination directories exist
    os.makedirs(os.path.dirname(preview_dest), exist_ok=True)
    os.makedirs(os.path.dirname(manifest_dest), exist_ok=True)

    # Copy files
    preview_path = src_path / "preview.png"
    manifest_path = src_path / "manifest.json"

    success = True

    if os.path.exists(preview_path):
        shutil.copy2(preview_path, preview_dest)
    else:
        print(f"Warning: No preview.png found in {src_path}")
        success = False

    if os.path.exists(manifest_path):
        shutil.copy2(manifest_path, manifest_dest)
    else:
        print(f"Warning: No manifest.json found in {src_path}")
        success = False

    return success

def process_theme_zip(zip_path):
    """Process a theme zip file"""
    print(f"Processing theme zip: {zip_path}")

    # Extract theme name from zip filename
    theme_name = os.path.basename(zip_path).replace('.zip', '')

    # Destination path
    dest_path = THEMES_DIR / theme_name

    # Remove destination if it already exists
    if os.path.exists(dest_path):
        shutil.rmtree(dest_path)

    # Create destination directory
    os.makedirs(dest_path, exist_ok=True)

    # Extract zip using our custom function
    extract_without_nested_dirs(zip_path, dest_path)

    print(f"Extracted theme to {dest_path}")

    # Remove any __MACOSX directory that might have been created
    macosx_dir = dest_path / "__MACOSX"
    if os.path.exists(macosx_dir):
        shutil.rmtree(macosx_dir)
        print(f"Removed __MACOSX directory from {dest_path}")

    # Copy preview and manifest to their dedicated directories
    copy_preview_and_manifest(dest_path, is_theme=True)

    # Only remove the processed zip file if KEEP_ZIP_FILES is False
    if not KEEP_ZIP_FILES:
        os.remove(zip_path)
        print(f"Removed processed zip file: {zip_path}")

def process_component_zip(zip_path, component_type):
    """Process a component zip file"""
    print(f"Processing {component_type} component zip: {zip_path}")

    # Extract component name from zip filename
    component_name = os.path.basename(zip_path).replace('.zip', '')

    # Get component directory
    component_dir = COMPONENT_DIRS.get(component_type)
    if not component_dir:
        print(f"Unknown component type: {component_type}")
        return

    # Destination path
    dest_path = COMPONENTS_DIR / component_dir / component_name

    # Remove destination if it already exists
    if os.path.exists(dest_path):
        shutil.rmtree(dest_path)

    # Create destination directory
    os.makedirs(dest_path, exist_ok=True)

    # Extract zip using our custom function
    extract_without_nested_dirs(zip_path, dest_path)

    print(f"Extracted component to {dest_path}")

    # Remove any __MACOSX directory that might have been created
    macosx_dir = dest_path / "__MACOSX"
    if os.path.exists(macosx_dir):
        shutil.rmtree(macosx_dir)
        print(f"Removed __MACOSX directory from {dest_path}")

    # Copy preview and manifest to their dedicated directories
    copy_preview_and_manifest(dest_path, component_dir)

    # Only remove the processed zip file if KEEP_ZIP_FILES is False
    if not KEEP_ZIP_FILES:
        os.remove(zip_path)
        print(f"Removed processed zip file: {zip_path}")

def ensure_directory_structure():
    """Ensure all required directories exist"""
    # Main directories
    os.makedirs(CATALOG_DIR, exist_ok=True)
    os.makedirs(THEMES_DIR, exist_ok=True)
    os.makedirs(COMPONENTS_DIR, exist_ok=True)

    # Preview and manifest directories for themes
    os.makedirs(THEMES_DIR / "previews", exist_ok=True)
    os.makedirs(THEMES_DIR / "manifests", exist_ok=True)

    # Component type directories with their preview and manifest subdirectories
    for comp_dir in COMPONENT_DIRS.values():
        os.makedirs(COMPONENTS_DIR / comp_dir, exist_ok=True)
        os.makedirs(COMPONENTS_DIR / comp_dir / "previews", exist_ok=True)
        os.makedirs(COMPONENTS_DIR / comp_dir / "manifests", exist_ok=True)

def main():
    """Main function to process all zip files"""
    # Ensure the directory structure exists
    ensure_directory_structure()

    # Process theme zips
    theme_zip_pattern = UPLOADS_DIR / "Themes" / "*.theme.zip"
    for zip_path in glob.glob(str(theme_zip_pattern)):
        if validate_zip_contents(zip_path):
            process_theme_zip(zip_path)
        else:
            print(f"Error: {zip_path} is missing required files (manifest.json and/or preview.png)")

    # Process component zips
    for comp_type, comp_dir in COMPONENT_DIRS.items():
        component_zip_pattern = UPLOADS_DIR / "Components" / comp_dir / f"*{comp_type}.zip"
        for zip_path in glob.glob(str(component_zip_pattern)):
            if validate_zip_contents(zip_path):
                process_component_zip(zip_path, comp_type)
            else:
                print(f"Error: {zip_path} is missing required files (manifest.json and/or preview.png)")

if __name__ == "__main__":
    main()