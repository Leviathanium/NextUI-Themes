#!/usr/bin/env python3

import os
import json
import math
import re
import shutil
from string import Template
from datetime import datetime
from urllib.parse import quote

# Configuration
REPO_URL = "https://github.com/Leviathanium/NextUI-Themes"
RAW_URL = "https://github.com/Leviathanium/NextUI-Themes/raw/main"
CATALOG_PATH = "Catalog/catalog.json"
OUTPUT_DIR = ".github/index"
TEMPLATES_DIR = ".github/templates"
README_PATH = "README.md"
MAX_FEATURED = 3
MAX_RECENT = 6
ITEMS_PER_PAGE = 12
COLUMNS_PER_ROW = 3
README_GALLERY_PATTERN = r"<!-- GALLERY_START -->.*?<!-- GALLERY_END -->"

# Icons and metadata
PREVIEW_ICON = '<img title="View full-size preview" src="https://user-images.githubusercontent.com/44569252/194037184-ae453506-2536-4c6f-8a19-4a6c1de6ce32.png" width="16">'
AUTHOR_ICON = '<img src="https://user-images.githubusercontent.com/44569252/194037581-698a5004-8b75-4da6-a63d-b41d541ebde2.png" width="16" title="Search themes by this author">'
README_ICON = '<img src="https://user-images.githubusercontent.com/44569252/215358455-b6a1348b-8161-40d6-9cc1-cc31720377c4.png" height="16" title="README">'
ICONPACK_ICON = '<img src="https://user-images.githubusercontent.com/44569252/215106002-fbcf1815-8080-447c-94c2-61f161efb503.png" height="16" title="This theme contains an icon pack">'
PREV_PAGE_ICON = "https://github.com/OnionUI/Themes/assets/44569252/38f4d157-d58f-49fc-b006-2d131e9e3ecf"
NEXT_PAGE_ICON = "https://github.com/OnionUI/Themes/assets/44569252/a0717376-2b5b-4534-9eba-4d2d3961f06b"

# Component types
COMPONENT_TYPES = {
    "themes": {
        "title": "Full Themes",
        "extension": ".theme"
    },
    "wallpapers": {
        "title": "Wallpapers",
        "extension": ".bg"
    },
    "icons": {
        "title": "Icons",
        "extension": ".icon"
    },
    "accents": {
        "title": "Accent Colors",
        "extension": ".acc"
    },
    "fonts": {
        "title": "Fonts",
        "extension": ".font"
    },
    "overlays": {
        "title": "Overlays",
        "extension": ".over"
    }
}

def load_template(template_name):
    """Load a template file from the templates directory"""
    path = os.path.join(TEMPLATES_DIR, template_name)
    with open(path, "r", encoding="utf-8") as f:
        return Template(f.read())

def load_catalog():
    """Load the catalog.json file"""
    try:
        with open(CATALOG_PATH, "r", encoding="utf-8") as f:
            return json.load(f)
    except Exception as e:
        print(f"Error loading catalog: {e}")
        return None

def ensure_dir(path):
    """Ensure a directory exists"""
    os.makedirs(path, exist_ok=True)

def url_encode(path):
    """URL encode a path for use in URLs"""
    return quote(path, safe='')

def paginate_items(items, items_per_page=ITEMS_PER_PAGE):
    """Split items into pages"""
    return [items[i:i + items_per_page] for i in range(0, len(items), items_per_page)]

def format_date(date_str):
    """Format a date string for display"""
    try:
        date = datetime.fromisoformat(date_str.replace("Z", "+00:00"))
        return date.strftime("%Y-%m-%d")
    except:
        return date_str

def generate_page_filename(page, total_pages):
    """Generate a filename for a page"""
    if page == 0:
        return "index.md"
    else:
        return f"page-{page+1}.md"

def generate_pagination(current_page, total_pages, base_path):
    """Generate pagination links"""
    if total_pages <= 1:
        return ""
    
    pagination_template = load_template("pagination_template.md")
    
    # Previous page link
    prev_page = ""
    if current_page > 0:
        prev_file = generate_page_filename(current_page - 1, total_pages)
        prev_page = f'<td align="right">\n\n[![Previous page]({PREV_PAGE_ICON})]({prev_file})\n\n</td>'
    
    # Next page link
    next_page = ""
    if current_page < total_pages - 1:
        next_file = generate_page_filename(current_page + 1, total_pages)
        next_page = f'<td>\n\n[![Next page]({NEXT_PAGE_ICON})]({next_file})\n\n</td>'
    
    # Page number links
    page_links = []
    for i in range(total_pages):
        if i == current_page:
            page_links.append(f"**{i+1}**")
        else:
            page_file = generate_page_filename(i, total_pages)
            page_links.append(f"[{i+1}]({page_file})")
    
    return pagination_template.substitute({
        "PREV_PAGE": prev_page,
        "PAGE_LINKS": " ".join(page_links),
        "NEXT_PAGE": next_page
    })

def generate_item_card(item, type_key, width=None):
    """Generate a card for a theme or component"""
    item_template = load_template("item_template.md")
    
    # Get item metadata
    name = item.get("description", os.path.basename(item.get("manifest_path", "")).replace(".json", ""))
    author = item.get("author", "Unknown")
    updated = format_date(item.get("last_updated", ""))
    
    # URLs
    preview_url = f"{RAW_URL}/{item.get('preview_path', '')}"
    download_url = item.get("URL", "")
    history_url = f"{REPO_URL}/commits/main/Catalog/{type_key}/{name}"
    
    # Special features
    has_readme = False  # We'd need to check if README exists
    has_icons = type_key == "themes" and name in item.get("icons", [])
    
    # Generate metadata icons
    metadata_icons = [
        f"&nbsp;&nbsp;<a href=\"{preview_url}\">{PREVIEW_ICON}</a>",
        f"&nbsp;&nbsp;<a href=\"https://github.com/search?q=filename%3A%22{url_encode(author)}%22+repo%3ALeviathanium%2FNextUI-Themes\">{AUTHOR_ICON}</a>" if author != "Unknown" else ""
    ]
    
    if has_readme:
        metadata_icons.append(f"&nbsp;&nbsp;<a href=\"#\">{README_ICON}</a>")
    
    if has_icons:
        metadata_icons.append(f" &nbsp;<a href=\"#\">{ICONPACK_ICON}</a>")
    
    # Title for hover
    title_parts = []
    if name:
        title_parts.append(f"Name: {name}")
    if author:
        title_parts.append(f"Author: {author}")
    if updated:
        title_parts.append(f"Last updated: {updated}")
    title_parts.append("(Click to download)")
    
    # Substitute values in template
    return item_template.substitute({
        "COLUMN_WIDTH": width or int(100 / COLUMNS_PER_ROW),
        "COLUMN_SPANNER": "<br/>" if width else "",
        "DOWNLOAD_URL": download_url,
        "PREVIEW_URL": preview_url,
        "TITLE": "&#013;".join(title_parts),
        "NAME": name,
        "AUTHOR": author,
        "UPDATED": updated,
        "HISTORY_URL": history_url,
        "METADATA_ICONS": "".join(metadata_icons)
    })

def generate_grid(items, type_key):
    """Generate a grid of items"""
    grid_template = load_template("grid_template.md")
    
    rows = []
    current_row = []
    
    for i, item in enumerate(items):
        current_row.append(generate_item_card(item, type_key))
        
        if (i + 1) % COLUMNS_PER_ROW == 0 or i == len(items) - 1:
            # Fill the row with empty cells if needed
            while len(current_row) < COLUMNS_PER_ROW:
                current_row.append("")
            
            rows.append("\n".join(current_row))
            current_row = []
    
    grid_items = "</tr><tr>\n".join(rows)
    return grid_template.substitute({"GRID_ITEMS": grid_items})

def generate_category_page(items, type_key, current_page, total_pages, total_items):
    """Generate a category page"""
    category_template = load_template("category_template.md")
    
    grid = generate_grid(items, type_key)
    pagination = generate_pagination(current_page, total_pages, "")
    
    return category_template.substitute({
        "CATEGORY_TITLE": COMPONENT_TYPES[type_key]["title"],
        "CURRENT_PAGE": current_page + 1,
        "TOTAL_PAGES": total_pages,
        "TOTAL_ITEMS": total_items,
        "GRID_ITEMS": grid,
        "PAGINATION": pagination
    })

def generate_category_pages(catalog, type_key):
    """Generate pages for a category"""
    type_info = COMPONENT_TYPES[type_key]
    category_dir = f"{OUTPUT_DIR}/{type_key}"
    ensure_dir(category_dir)
    
    # Get items for this type
    if type_key == "themes":
        items = catalog.get("themes", {}).values()
    else:
        items = catalog.get("components", {}).get(type_key, {}).values()
    
    items = sorted(items, key=lambda x: x.get("last_updated", ""), reverse=True)
    total_items = len(items)
    
    if total_items == 0:
        print(f"No items found for {type_info['title']}")
        return []
    
    # Split items into pages
    pages = paginate_items(items)
    total_pages = len(pages)
    
    print(f"Generating {total_pages} pages for {type_info['title']}")
    
    # Generate each page
    for i, page_items in enumerate(pages):
        page_content = generate_category_page(
            page_items, type_key, i, total_pages, total_items
        )
        
        page_filename = generate_page_filename(i, total_pages)
        page_path = f"{category_dir}/{page_filename}"
        
        with open(page_path, "w", encoding="utf-8") as f:
            f.write(page_content)
    
    return items

def generate_main_index(catalog, featured_themes, recently_added, recently_updated):
    """Generate the main index page"""
    index_template = load_template("index_template.md")
    
    # Generate featured themes grid
    featured_grid = generate_grid(featured_themes[:MAX_FEATURED], "themes")
    
    # Generate recently added grid
    recent_grid = generate_grid(recently_added[:MAX_RECENT], "themes")
    
    # Generate recently updated grid
    updated_grid = generate_grid(recently_updated[:MAX_RECENT], "themes")
    
    index_content = index_template.substitute({
        "FEATURED_THEMES": featured_grid,
        "RECENTLY_ADDED": recent_grid,
        "RECENTLY_UPDATED": updated_grid
    })
    
    ensure_dir(OUTPUT_DIR)
    with open(f"{OUTPUT_DIR}/index.md", "w", encoding="utf-8") as f:
        f.write(index_content)

def update_readme_gallery(featured_themes):
    """Update the gallery section in the README"""
    # Read the current README
    try:
        with open(README_PATH, "r", encoding="utf-8") as f:
            readme_content = f.read()
    except Exception as e:
        print(f"Error reading README: {e}")
        return False
    
    # Generate featured themes grid
    featured_grid = generate_grid(featured_themes[:MAX_FEATURED], "themes")
    
    # Create the replacement content
    gallery_content = f"<!-- GALLERY_START -->\n## Featured Themes\n\n{featured_grid}\n\n[Browse All Themes](.github/index/index.md)\n<!-- GALLERY_END -->"
    
    # Replace the gallery section
    pattern = re.compile(README_GALLERY_PATTERN, re.DOTALL)
    if pattern.search(readme_content):
        updated_readme = pattern.sub(gallery_content, readme_content)
        
        # Write the updated README
        with open(README_PATH, "w", encoding="utf-8") as f:
            f.write(updated_readme)
        return True
    else:
        print("Gallery section not found in README")
        return False

def select_featured_themes(themes, count=MAX_FEATURED):
    """Select featured themes based on recency and quality"""
    # For now, just use the most recently updated themes
    return sorted(themes, key=lambda x: x.get("last_updated", ""), reverse=True)[:count]

def main():
    """Main function"""
    print("Generating theme gallery...")
    
    # Load the catalog
    catalog = load_catalog()
    if not catalog:
        print("Failed to load catalog.json")
        return
    
    # Clear output directory
    if os.path.exists(OUTPUT_DIR):
        shutil.rmtree(OUTPUT_DIR)
    
    # Generate pages for each category
    all_themes = []
    recently_updated = []
    
    for type_key in COMPONENT_TYPES.keys():
        items = generate_category_pages(catalog, type_key)
        
        # For themes, also build lists for the main page
        if type_key == "themes" and items:
            all_themes = items
            recently_updated = sorted(items, key=lambda x: x.get("last_updated", ""), reverse=True)
    
    # Select featured and recent items
    featured_themes = select_featured_themes(all_themes)
    recently_added = sorted(all_themes, key=lambda x: x.get("creation_date", ""), reverse=True)
    
    # Generate main index
    generate_main_index(catalog, featured_themes, recently_added, recently_updated)
    
    # Update README gallery section
    if update_readme_gallery(featured_themes):
        print("Updated README gallery section")
    
    print("Gallery generation complete!")

if __name__ == "__main__":
    main()
