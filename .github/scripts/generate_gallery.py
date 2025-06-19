#!/usr/bin/env python3
"""
Gallery generation script for the new YAML-based theme catalog.
Simplified to only handle themes (no components).
"""

import os
import yaml
import re
import shutil
from string import Template
from datetime import datetime
from pathlib import Path

# Configuration
REPO_URL = "https://github.com/Leviathanium/NextUI-Themes"
RAW_URL = "https://github.com/Leviathanium/NextUI-Themes/raw/main"
CATALOG_PATH = "Catalog/catalog.yml"
OUTPUT_DIR = ".github/index"
TEMPLATES_DIR = ".github/templates"
README_PATH = "README.md"
FEATURED_COUNT = 12  # Number of featured themes to show
COLUMNS_PER_ROW = 3
README_GALLERY_PATTERN = r"<!-- GALLERY_START -->.*?<!-- GALLERY_END -->"

def ensure_dir(path):
    """Ensure a directory exists"""
    os.makedirs(path, exist_ok=True)

def load_catalog():
    """Load the catalog.yml file"""
    try:
        with open(CATALOG_PATH, "r", encoding="utf-8") as f:
            return yaml.safe_load(f)
    except Exception as e:
        print(f"Error loading catalog: {e}")
        return None

def create_templates_dir():
    """Create templates directory and basic templates"""
    ensure_dir(TEMPLATES_DIR)

    # Create basic item template
    item_template = """<td align="center" valign="top" width="{width}%">
<a href="{download_url}">
<img title="{title}" width="480px" src="{preview_url}" /><br/>
<b>{name}</b>
</a><br/>
<sup><i>{author}</i></sup><br>
<sub>
<sup><a title="Last updated: {updated}" href="{history_url}">{updated}</a></sup>
</sub>
</td>"""

    with open(f"{TEMPLATES_DIR}/item_template.md", "w") as f:
        f.write(item_template)

    # Create grid template
    grid_template = """<table align="center"><tr>
{grid_items}
</tr></table>"""

    with open(f"{TEMPLATES_DIR}/grid_template.md", "w") as f:
        f.write(grid_template)

def load_template(template_name):
    """Load a template file from the templates directory"""
    path = os.path.join(TEMPLATES_DIR, template_name)
    try:
        with open(path, "r", encoding="utf-8") as f:
            return f.read()
    except FileNotFoundError:
        print(f"Template {template_name} not found, creating basic template")
        create_templates_dir()
        with open(path, "r", encoding="utf-8") as f:
            return f.read()

def format_date(date_str):
    """Format a date string for display"""
    try:
        date = datetime.fromisoformat(date_str.replace("Z", "+00:00"))
        return date.strftime("%Y-%m-%d")
    except:
        return date_str

def generate_item_card(theme_name, theme_info):
    """Generate a card for a theme"""
    item_template = load_template("item_template.md")

    # Get theme metadata
    name = theme_info.get("description", theme_name.replace(".theme", ""))
    author = theme_info.get("author", "Unknown")
    updated = format_date(theme_info.get("last_updated", ""))

    # URLs - since we don't have zip packages anymore, link to the repository or theme page
    repo_info = theme_info.get("source", {})
    download_url = repo_info.get("repository", f"{REPO_URL}")
    preview_url = f"{RAW_URL}/{theme_info.get('preview_path', '')}"
    history_url = f"{REPO_URL}/commits/main/Themes/{theme_name}"

    # Title for hover
    title_parts = []
    if name:
        title_parts.append(f"Name: {name}")
    if author:
        title_parts.append(f"Author: {author}")
    if updated:
        title_parts.append(f"Last updated: {updated}")
    title_parts.append("(Click to view source)")

    # Format the template
    return item_template.format(
        width=int(100 / COLUMNS_PER_ROW),
        download_url=download_url,
        preview_url=preview_url,
        title="&#013;".join(title_parts),
        name=name,
        author=author,
        updated=updated,
        history_url=history_url
    )

def generate_grid(themes):
    """Generate a grid of theme cards"""
    grid_template = load_template("grid_template.md")

    rows = []
    current_row = []

    theme_items = list(themes.items())

    for i, (theme_name, theme_info) in enumerate(theme_items):
        current_row.append(generate_item_card(theme_name, theme_info))

        if (i + 1) % COLUMNS_PER_ROW == 0 or i == len(theme_items) - 1:
            # Fill the row with empty cells if needed
            while len(current_row) < COLUMNS_PER_ROW:
                current_row.append('<td width="33%"></td>')

            rows.append("\n".join(current_row))
            current_row = []

    grid_items = "</tr><tr>\n".join(rows)
    return grid_template.format(grid_items=grid_items)

def generate_themes_index(catalog):
    """Generate the themes index page"""
    themes = catalog.get("themes", {})

    content = f"# Full Themes\n\n"
    content += f"*{len(themes)} available themes*\n\n"
    content += generate_grid(themes)

    # Create themes directory
    themes_dir = f"{OUTPUT_DIR}/themes"
    ensure_dir(themes_dir)

    with open(f"{themes_dir}/index.md", "w", encoding="utf-8") as f:
        f.write(content)

def generate_main_index(catalog):
    """Generate the main index page"""
    themes = catalog.get("themes", {})

    # Take first N themes for featured section
    featured_themes = dict(list(themes.items())[:FEATURED_COUNT])

    index_content = f"""# NextUI Themes Gallery

Browse and download themes for NextUI devices.

## Browse Themes

- [All Themes (.theme)](/.github/index/themes/index.md) - Complete theme packages

## Featured Themes

{generate_grid(featured_themes)}
"""

    ensure_dir(OUTPUT_DIR)
    with open(f"{OUTPUT_DIR}/index.md", "w", encoding="utf-8") as f:
        f.write(index_content)

def update_readme_gallery(catalog):
    """Update the gallery section in the README"""
    try:
        with open(README_PATH, "r", encoding="utf-8") as f:
            readme_content = f.read()
    except Exception as e:
        print(f"Error reading README: {e}")
        return False

    themes = catalog.get("themes", {})
    featured_themes = dict(list(themes.items())[:FEATURED_COUNT])

    # Generate featured themes grid
    featured_grid = generate_grid(featured_themes)

    # Create the replacement content
    gallery_content = f"""<!-- GALLERY_START -->
## Featured Themes

{featured_grid}

<div align="center"><h3><a href=".github/index/index.md">Browse All Themes</a></h3></div>
<!-- GALLERY_END -->"""

    # Replace the gallery section
    pattern = re.compile(README_GALLERY_PATTERN, re.DOTALL)
    if pattern.search(readme_content):
        updated_readme = pattern.sub(gallery_content, readme_content)

        with open(README_PATH, "w", encoding="utf-8") as f:
            f.write(updated_readme)
        return True
    else:
        print("Gallery section not found in README")
        return False

def main():
    """Main function"""
    print("Generating theme gallery...")

    # Load the catalog
    catalog = load_catalog()
    if not catalog:
        print("Failed to load catalog.yml")
        return

    # Clear output directory
    if os.path.exists(OUTPUT_DIR):
        shutil.rmtree(OUTPUT_DIR)

    # Create templates if they don't exist
    create_templates_dir()

    # Generate gallery pages
    generate_themes_index(catalog)
    generate_main_index(catalog)

    # Update README gallery section
    if update_readme_gallery(catalog):
        print("Updated README gallery section")

    print("Gallery generation complete!")

if __name__ == "__main__":
    main()