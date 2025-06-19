#!/usr/bin/env python3
"""
Script to process the push.yml file with theme submissions.
Updated for YAML format and direct cloning to Themes/ directory.
"""

import os
import yaml
import shutil
import tempfile
import zipfile
import subprocess
import re
import json
from pathlib import Path
from datetime import datetime
import urllib.request
import ssl

# Base paths
REPO_ROOT = Path(__file__).resolve().parent.parent.parent
UPLOAD_DIR = REPO_ROOT / "Upload"
THEMES_DIR = REPO_ROOT / "Themes"
CATALOG_DIR = REPO_ROOT / "Catalog"
PUSH_YML_PATH = UPLOAD_DIR / "push.yml"
CATALOG_YML_PATH = CATALOG_DIR / "catalog.yml"

def load_push_yml():
    """Load the push.yml file"""
    if not PUSH_YML_PATH.exists():
        print(f"Error: {PUSH_YML_PATH} does not exist")
        return None

    try:
        with open(PUSH_YML_PATH, 'r') as f:
            push_data = yaml.safe_load(f)
        return push_data
    except yaml.YAMLError as e:
        print(f"Error parsing push.yml: {e}")
        return None
    except Exception as e:
        print(f"Error loading push.yml: {e}")
        return None

def load_catalog_yml():
    """Load the catalog.yml file"""
    if not CATALOG_YML_PATH.exists():
        # Create initial catalog if it doesn't exist
        initial_catalog = {
            "info": {
                "last_updated": datetime.utcnow().isoformat() + "Z",
                "repository_url": "https://github.com/Leviathanium/NextUI-Themes.git",
                "branch": "main"
            },
            "themes": {}
        }
        with open(CATALOG_YML_PATH, 'w') as f:
            yaml.dump(initial_catalog, f, default_flow_style=False, sort_keys=False)
        return initial_catalog

    try:
        with open(CATALOG_YML_PATH, 'r') as f:
            return yaml.safe_load(f)
    except Exception as e:
        print(f"Error loading catalog.yml: {e}")
        return None

def save_catalog_yml(catalog):
    """Save the catalog.yml file"""
    try:
        # Update last_updated timestamp
        catalog["info"]["last_updated"] = datetime.utcnow().isoformat() + "Z"

        with open(CATALOG_YML_PATH, 'w') as f:
            yaml.dump(catalog, f, default_flow_style=False, sort_keys=False, indent=2)
        print(f"Updated catalog.yml")
        return True
    except Exception as e:
        print(f"Error saving catalog.yml: {e}")
        return False

def validate_submission(submission):
    """Validate a submission entry"""
    required_fields = ["type", "name", "author", "submission_method"]

    # Check for required fields
    for field in required_fields:
        if field not in submission:
            print(f"Error: Missing required field '{field}' in submission")
            return False

    # Validate type (only themes supported)
    if submission["type"] != "theme":
        print(f"Error: Invalid type '{submission['type']}'. Only 'theme' is supported")
        return False

    # Validate name
    if not submission["name"].endswith(".theme"):
        print(f"Error: Name '{submission['name']}' must end with '.theme'")
        return False

    # Validate submission_method
    valid_methods = ["repository", "zip", "url"]
    if submission["submission_method"] not in valid_methods:
        print(f"Error: Invalid submission_method '{submission['submission_method']}'. Must be one of: {', '.join(valid_methods)}")
        return False

    # Validate method-specific fields
    if submission["submission_method"] == "repository":
        if "url" not in submission or not submission["url"]:
            print(f"Error: 'url' is required for repository submissions")
            return False
        if "commit" not in submission or not submission["commit"]:
            print(f"Error: 'commit' is required for repository submissions")
            return False
        if "branch" not in submission:
            submission["branch"] = "main"  # Default to main if not specified

    elif submission["submission_method"] == "zip":
        zip_path = UPLOAD_DIR / f"{submission['name']}.zip"
        if not zip_path.exists():
            print(f"Error: Zip file '{zip_path}' not found for zip submission")
            return False

    elif submission["submission_method"] == "url":
        if "url" not in submission or not submission["url"]:
            print(f"Error: 'url' is required for url submissions")
            return False
        if not submission["url"].startswith(("http://", "https://")):
            print(f"Error: Invalid URL format '{submission['url']}'")
            return False

    return True

def clean_existing_theme(theme_name):
    """Clean up existing theme from Themes/ directory and catalog"""
    print(f"Cleaning up existing theme: {theme_name}")

    # Remove from Themes/ directory
    theme_path = THEMES_DIR / theme_name
    if theme_path.exists():
        shutil.rmtree(theme_path)
        print(f"Removed existing theme directory: {theme_path}")

    # Remove preview and manifest files
    preview_path = CATALOG_DIR / "previews" / f"{theme_name}.png"
    manifest_path = CATALOG_DIR / "manifests" / f"{theme_name}.json"

    if preview_path.exists():
        os.remove(preview_path)
        print(f"Removed existing preview: {preview_path}")

    if manifest_path.exists():
        os.remove(manifest_path)
        print(f"Removed existing manifest: {manifest_path}")

def clone_repository(repo_url, commit_hash, branch, target_dir):
    """Clone a repository at a specific commit and branch"""
    print(f"Cloning {repo_url} (branch: {branch}, commit: {commit_hash})")

    try:
        # Clone repository
        subprocess.run(['git', 'clone', '--no-checkout', repo_url, target_dir], check=True)

        # Change to repository directory
        cwd = os.getcwd()
        os.chdir(target_dir)

        try:
            # Fetch branch if specified and different from main
            if branch and branch != "main":
                subprocess.run(['git', 'fetch', 'origin', branch], check=True)
                subprocess.run(['git', 'checkout', branch], check=True)
            else:
                subprocess.run(['git', 'checkout', 'main'], check=True)

            # Checkout specific commit
            subprocess.run(['git', 'checkout', commit_hash], check=True)

            # Remove .git directory to save space
            shutil.rmtree('.git', ignore_errors=True)
        finally:
            # Return to original directory
            os.chdir(cwd)

        print(f"Repository cloned successfully to {target_dir}")
        return True
    except Exception as e:
        print(f"Error cloning repository: {e}")
        return False

def validate_theme_structure(theme_dir):
    """Validate that the theme directory contains required files"""
    required_files = ["manifest.yml", "preview.png"]

    for required_file in required_files:
        file_path = theme_dir / required_file
        if not file_path.exists():
            print(f"Error: Required file '{required_file}' not found in theme directory")
            return False

    return True

def copy_metadata_files(theme_dir, theme_name):
    """Copy preview.png and manifest files to Catalog directories"""
    try:
        # Source files
        preview_src = theme_dir / "preview.png"
        manifest_src = theme_dir / "manifest.yml"

        # Destination files
        preview_dest = CATALOG_DIR / "previews" / f"{theme_name}.png"
        manifest_dest = CATALOG_DIR / "manifests" / f"{theme_name}.json"

        # Create directories if they don't exist
        os.makedirs(os.path.dirname(preview_dest), exist_ok=True)
        os.makedirs(os.path.dirname(manifest_dest), exist_ok=True)

        # Copy preview
        shutil.copy2(preview_src, preview_dest)
        print(f"Copied preview: {preview_src} -> {preview_dest}")

        # Convert YAML manifest to JSON and copy
        with open(manifest_src, 'r') as f:
            manifest_data = yaml.safe_load(f)

        with open(manifest_dest, 'w') as f:
            json.dump(manifest_data, f, indent=2)
        print(f"Converted and copied manifest: {manifest_src} -> {manifest_dest}")

        return True
    except Exception as e:
        print(f"Error copying metadata files: {e}")
        return False

def extract_zip_to_theme_dir(zip_path, theme_name):
    """Extract a zip file to the Themes directory"""
    theme_dir = THEMES_DIR / theme_name

    try:
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            # Create theme directory
            os.makedirs(theme_dir, exist_ok=True)

            # Extract files, handling nested directories
            file_list = zip_ref.namelist()

            # Check for common root directory
            common_root = None
            for item in file_list:
                if "__MACOSX" in item:
                    continue
                parts = item.split('/')
                if len(parts) > 1 and parts[0] and not common_root:
                    common_root = parts[0]
                elif len(parts) > 1 and parts[0] and parts[0] != common_root:
                    common_root = None
                    break

            # Extract files
            for item in file_list:
                if "__MACOSX" in item:
                    continue

                # Determine target path
                if common_root and item.startswith(common_root + '/'):
                    relative_path = item[len(common_root) + 1:]
                    if not relative_path:  # Skip root directory entry
                        continue
                    target_path = theme_dir / relative_path
                else:
                    target_path = theme_dir / item

                # Skip directories
                if item.endswith('/'):
                    continue

                # Create directory structure
                os.makedirs(os.path.dirname(target_path), exist_ok=True)

                # Extract file
                with zip_ref.open(item) as source_file, open(target_path, 'wb') as target_file:
                    shutil.copyfileobj(source_file, target_file)

        print(f"Extracted zip to: {theme_dir}")
        return theme_dir
    except Exception as e:
        print(f"Error extracting zip file: {e}")
        return None

def download_file(url, local_path):
    """Download a file from a URL"""
    try:
        # Create directory structure
        os.makedirs(os.path.dirname(local_path), exist_ok=True)

        # Create SSL context for some hosting services
        context = ssl._create_unverified_context()

        # Download file
        urllib.request.urlretrieve(url, local_path)
        print(f"Downloaded: {url} -> {local_path}")
        return True
    except Exception as e:
        print(f"Error downloading file: {e}")
        return False

def update_catalog(submission, theme_name):
    """Update catalog.yml with the new theme entry"""
    catalog = load_catalog_yml()
    if not catalog:
        return False

    # Create theme entry
    theme_entry = {
        "author": submission["author"],
        "description": theme_name.replace(".theme", ""),  # Use theme name as description
        "last_updated": datetime.utcnow().isoformat() + "Z",
        "preview_path": f"Catalog/previews/{theme_name}.png",
        "manifest_path": f"Catalog/manifests/{theme_name}.json",
        "theme_path": f"Themes/{theme_name}"
    }

    # Add source information if it's a repository submission
    if submission["submission_method"] == "repository":
        theme_entry["source"] = {
            "repository": submission["url"],
            "commit": submission["commit"],
            "branch": submission.get("branch", "main")
        }

    # Add theme to catalog (at the beginning to maintain insertion order)
    themes = catalog.get("themes", {})
    new_themes = {theme_name: theme_entry}

    # Add all other themes
    for name, info in themes.items():
        if name != theme_name:  # Skip if already exists
            new_themes[name] = info

    catalog["themes"] = new_themes

    return save_catalog_yml(catalog)

def process_repository_submission(submission):
    """Process a repository submission"""
    theme_name = submission["name"]
    print(f"Processing repository submission: {theme_name}")

    # Clean existing theme
    clean_existing_theme(theme_name)

    # Clone to temporary directory
    with tempfile.TemporaryDirectory() as temp_dir:
        if not clone_repository(submission["url"], submission["commit"],
                                submission.get("branch", "main"), temp_dir):
            return False

        temp_path = Path(temp_dir)

        # Validate theme structure
        if not validate_theme_structure(temp_path):
            return False

        # Copy to Themes directory
        theme_dir = THEMES_DIR / theme_name
        os.makedirs(THEMES_DIR, exist_ok=True)
        shutil.copytree(temp_path, theme_dir)
        print(f"Copied theme to: {theme_dir}")

        # Copy metadata files
        if not copy_metadata_files(theme_dir, theme_name):
            return False

        # Update catalog
        if not update_catalog(submission, theme_name):
            return False

        print(f"Successfully processed repository submission: {theme_name}")
        return True

def process_zip_submission(submission):
    """Process a zip submission"""
    theme_name = submission["name"]
    print(f"Processing zip submission: {theme_name}")

    # Clean existing theme
    clean_existing_theme(theme_name)

    # Source zip file
    source_zip = UPLOAD_DIR / f"{theme_name}.zip"

    # Extract to Themes directory
    theme_dir = extract_zip_to_theme_dir(source_zip, theme_name)
    if not theme_dir:
        return False

    # Validate theme structure
    if not validate_theme_structure(theme_dir):
        return False

    # Copy metadata files
    if not copy_metadata_files(theme_dir, theme_name):
        return False

    # Update catalog
    if not update_catalog(submission, theme_name):
        return False

    # Remove original zip file
    try:
        os.remove(source_zip)
        print(f"Removed original zip file: {source_zip}")
    except Exception as e:
        print(f"Warning: Could not remove original zip file: {e}")

    print(f"Successfully processed zip submission: {theme_name}")
    return True

def process_url_submission(submission):
    """Process a URL submission"""
    theme_name = submission["name"]
    print(f"Processing URL submission: {theme_name}")

    # Clean existing theme
    clean_existing_theme(theme_name)

    # Download to temporary file
    with tempfile.NamedTemporaryFile(delete=False, suffix='.zip') as temp_file:
        temp_path = temp_file.name

        try:
            if not download_file(submission["url"], temp_path):
                return False

            # Extract to Themes directory
            theme_dir = extract_zip_to_theme_dir(temp_path, theme_name)
            if not theme_dir:
                return False

            # Validate theme structure
            if not validate_theme_structure(theme_dir):
                return False

            # Copy metadata files
            if not copy_metadata_files(theme_dir, theme_name):
                return False

            # Update catalog
            if not update_catalog(submission, theme_name):
                return False

            print(f"Successfully processed URL submission: {theme_name}")
            return True
        finally:
            # Clean up temp file
            try:
                os.unlink(temp_path)
            except:
                pass

def reset_push_yml():
    """Reset push.yml after successful processing"""
    empty_push = {
        "submission": []
    }

    try:
        with open(PUSH_YML_PATH, 'w') as f:
            yaml.dump(empty_push, f, default_flow_style=False)
        print("Reset push.yml to empty state")
        return True
    except Exception as e:
        print(f"Error resetting push.yml: {e}")
        return False

def main():
    """Main function to process push.yml"""
    print("Starting push.yml processing")

    # Create necessary directories
    os.makedirs(UPLOAD_DIR, exist_ok=True)
    os.makedirs(THEMES_DIR, exist_ok=True)
    os.makedirs(CATALOG_DIR / "previews", exist_ok=True)
    os.makedirs(CATALOG_DIR / "manifests", exist_ok=True)

    # Load push.yml
    push_data = load_push_yml()
    if not push_data:
        print("Failed to load push.yml, aborting")
        return False

    # Check if push.yml has submissions
    if "submission" not in push_data or not push_data["submission"]:
        print("No submissions found in push.yml")
        return True

    # Process each submission
    success = True
    for submission in push_data["submission"]:
        # Validate submission
        if not validate_submission(submission):
            print(f"Validation failed for submission: {submission.get('name', 'unknown')}")
            success = False
            continue

        # Process submission based on method
        if submission["submission_method"] == "repository":
            if not process_repository_submission(submission):
                print(f"Failed to process repository submission: {submission['name']}")
                success = False
        elif submission["submission_method"] == "zip":
            if not process_zip_submission(submission):
                print(f"Failed to process zip submission: {submission['name']}")
                success = False
        elif submission["submission_method"] == "url":
            if not process_url_submission(submission):
                print(f"Failed to process URL submission: {submission['name']}")
                success = False

    # Reset push.yml if all submissions were processed successfully
    if success:
        if not reset_push_yml():
            print("Warning: Failed to reset push.yml")
            success = False

    print("Push.yml processing complete" if success else "Push.yml processing completed with errors")
    return success

if __name__ == "__main__":
    success = main()
    if not success:
        exit(1)