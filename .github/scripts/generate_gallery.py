#!/usr/bin/env python3
"""
NextUI Theme Gallery Generator

This script generates a beautiful README gallery for NextUI themes based on the catalog.json file.
It creates visual galleries for both themes and component packages.
"""

import os
import json
import math
import datetime
from string import Template
from pathlib import Path

# Constants
REPO_URL = "https://github.com/Leviathanium/NextUI-Themes"
RAW_URL = "https://github.com/Leviathanium/NextUI-Themes/raw/main"
CATALOG_FILE = "Catalog/catalog.json"
OUTPUT_DIR = "generated"
THEMES_PER_PAGE = 12
COMPONENTS_PER_PAGE = 12
COLUMNS = 3
MAX_RECENTS = 6  # Number of recent themes to show on the front page

# Icons and templates
PREVIEW_ICON = "https://user-images.githubusercontent.com/44569252/194037184-ae453506-2536-4c6f-8a19-4a6c1de6ce32.png"
AUTHOR_ICON = "https://user-images.githubusercontent.com/44569252/194037581-698a5004-8b75-4da6-a63d-b41d541ebde2.png"
README_ICON = "https://user-images.githubusercontent.com/44569252/215358455-b6a1348b-8161-40d6-9cc1-cc31720377c4.png"
HAS_ICONPACK_ICON = "https://user-images.githubusercontent.com/44569252/215106002-fbcf1815-8080-447c-94c2-61f161efb503.png"
PREV_PAGE_ICON = "https://github.com/OnionUI/Themes/assets/44569252/a0717376-2b5b-4534-9eba-4d2d3961f06b"
NEXT_PAGE_ICON = "https://github.com/OnionUI/Themes/assets/44569252/52e7a7af-9561-4345-b146-63105925d326"

# Header definition
HEADER_LINKS = {
    "index": ["Start", "../../README.md"],
    "themes": ["Themes", "../themes/index.md"],
    "components": ["Components", "../components/index.md"],
    "contributing": ["Contributing", "../../CONTRIBUTING.md"]
}

PAGE_TITLES = {
    "themes": "Browse Available Themes",
    "components_wallpapers": "Browse Wallpaper Components",
    "components_icons": "Browse Icon Components",
    "components_accents": "Browse Accent Components",
    "components_leds": "Browse LED Components",
    "components_fonts": "Browse Font Components",
    "components_overlays": "Browse Overlay Components"
}

# Templates
README_HEADER = """
# NextUI Themes Repository

A collection of themes and component packages for NextUI devices. This gallery allows you to browse and download themes directly.

"""

ITEM_TEMPLATE = Template("""
<td align="center" valign="top" width="${COLUMN_WIDTH}">
${COLUMN_SPANNER}<br/>
<a href="${DOWNLOAD_URL}">
<img title="${TITLE}" width="480px" src="${PREVIEW_URL}" /><br/>
<b>${NAME}</b>
</a><br/>
<sup><i>${AUTHOR}</i></sup><br>
<sub>
<sup><a title="Last updated: ${UPDATED}" href="${REPO_URL}">${UPDATED}</a></sup> ${AUTHOR_BTN}${PREVIEW_BTN}${README_BTN}${HAS_ICONPACK}
</sub>
</td>
""")

GRID_TEMPLATE = Template("""
<table align="center"><tr>
${GRID_ITEMS}
</tr></table>
""")

PAGINATION_TEMPLATE = Template("""
---

<table align="center"><tr>
${PREV_PAGE}
<td align="center" valign="middle">

${PAGE_LINKS}

</td>
${NEXT_PAGE}
</tr></table>
""")

def main():
    # Create output directory
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    os.makedirs(os.path.join(OUTPUT_DIR, "themes"), exist_ok=True)
    os.makedirs(os.path.join(OUTPUT_DIR, "components"), exist_ok=True)

    # Load catalog
    with open(CATALOG_FILE, 'r', encoding='utf-8') as f:
        catalog = json.load(f)

    # Extract themes and components
    themes = catalog.get("themes", {})
    components = catalog.get("components", {})

    # Get recently updated themes
    recent_themes = sort_items_by_date(themes, max_items=MAX_RECENTS)
    
    # Create the main README
    generate_main_readme(themes, components, recent_themes)
    
    # Generate theme pages
    generate_theme_pages(themes)
    
    # Generate component pages by type
    for comp_type, comps in components.items():
        generate_component_pages(comp_type, comps)

    print("Gallery generation complete!")

def sort_items_by_date(items, max_items=None):
    """Sort items by their last updated date (uses current date as placeholder)"""
    # In a real implementation, you would extract dates from commit history
    # For now, we'll just use the items as is
    sorted_items = list(items.items())
    # In a real implementation, sort by date: sorted_items.sort(key=lambda x: x[1].get("last_updated", ""), reverse=True)
    
    if max_items:
        return sorted_items[:max_items]
    return sorted_items

def generate_main_readme(themes, components, recent_themes):
    """Generate the main README with theme counts and recent themes"""
    with open(os.path.join(OUTPUT_DIR, "README.md"), 'w', encoding='utf-8') as f:
        # Write header
        f.write(generate_header("index"))
        f.write(README_HEADER)
        
        # Write theme and component counts
        f.write("## Browse Themes and Components\n\n")
        f.write(f"### [Browse All Themes ({len(themes)})](generated/themes/index.md)\n\n")
        
        # Component counts by type
        component_types = {
            "wallpapers": "Wallpaper",
            "icons": "Icon",
            "accents": "Accent",
            "leds": "LED",
            "fonts": "Font",
            "overlays": "Overlay"
        }
        
        for comp_type, name in component_types.items():
            count = len(components.get(comp_type, {}))
            if count > 0:
                f.write(f"### [Browse {name} Components ({count})](generated/components/{comp_type}/index.md)\n\n")
        
        # Recently updated themes
        f.write("## Recently Updated Themes\n\n")
        f.write(generate_theme_grid(recent_themes))
        
        # How to install section
        f.write("""
## How to Install

### Themes

1. Download the theme package (click on a theme image above)
2. Extract the theme to your TrimUI Brick device in the `Tools/tg5040/Theme-Manager.pak/Themes/` directory
3. Open the Theme Manager app on your device
4. Select the theme and apply it

### Components

1. Download the component package
2. Extract to your TrimUI Brick device in the appropriate component directory:
   - `Tools/tg5040/Theme-Manager.pak/Components/[Type]/`
3. Open the Theme Manager app on your device
4. Select "Components" and the appropriate component type
5. Choose the component and apply it
""")

def generate_theme_pages(themes):
    """Generate pages for theme browsing"""
    themes_list = list(themes.items())
    total_pages = math.ceil(len(themes_list) / THEMES_PER_PAGE)
    
    os.makedirs(os.path.join(OUTPUT_DIR, "themes"), exist_ok=True)
    
    for page in range(total_pages):
        start = page * THEMES_PER_PAGE
        end = start + THEMES_PER_PAGE
        page_themes = themes_list[start:end]
        
        filename = "index.md" if page == 0 else f"page-{page+1}.md"
        filepath = os.path.join(OUTPUT_DIR, "themes", filename)
        
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(generate_header("themes"))
            f.write(f"\n## {PAGE_TITLES['themes']}\n\n")
            f.write(f"*Page {page+1} of {total_pages} — {len(themes)} themes available*\n\n")
            f.write(generate_theme_grid(page_themes))
            
            if total_pages > 1:
                f.write(generate_pagination(page, total_pages, "themes"))

def generate_component_pages(comp_type, components):
    """Generate pages for component browsing by type"""
    components_list = list(components.items())
    total_pages = math.ceil(len(components_list) / COMPONENTS_PER_PAGE)
    
    component_dir = os.path.join(OUTPUT_DIR, "components", comp_type)
    os.makedirs(component_dir, exist_ok=True)
    
    for page in range(total_pages):
        start = page * COMPONENTS_PER_PAGE
        end = start + COMPONENTS_PER_PAGE
        page_components = components_list[start:end]
        
        filename = "index.md" if page == 0 else f"page-{page+1}.md"
        filepath = os.path.join(component_dir, filename)
        
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(generate_header("components"))
            f.write(f"\n## {PAGE_TITLES.get(f'components_{comp_type}', f'Browse {comp_type.title()} Components')}\n\n")
            f.write(f"*Page {page+1} of {total_pages} — {len(components)} components available*\n\n")
            f.write(generate_component_grid(comp_type, page_components))
            
            if total_pages > 1:
                f.write(generate_pagination(page, total_pages, f"components/{comp_type}"))

def generate_header(current_page):
    """Generate the header with navigation links"""
    header = "# NextUI Themes\n\n"
    links = []
    
    for key, (text, link) in HEADER_LINKS.items():
        if key == current_page:
            links.append(f"**{text}**")
        else:
            links.append(f"[{text}]({link})")
    
    header += " • ".join(links) + "\n\n"
    return header

def generate_theme_grid(themes):
    """Generate a grid of theme items"""
    items = []
    row = []
    
    for i, (theme_name, theme_info) in enumerate(themes):
        row.append(generate_theme_item(theme_name, theme_info, i))
        
        if len(row) == COLUMNS or i == len(themes) - 1:
            items.append("".join(row))
            row = []
            
            if i < len(themes) - 1:  # Don't add </tr><tr> after last row
                items.append("</tr><tr>")
    
    grid_items = "\n".join(items)
    return GRID_TEMPLATE.substitute(GRID_ITEMS=grid_items)

def generate_component_grid(comp_type, components):
    """Generate a grid of component items"""
    items = []
    row = []
    
    for i, (comp_name, comp_info) in enumerate(components):
        row.append(generate_component_item(comp_type, comp_name, comp_info, i))
        
        if len(row) == COLUMNS or i == len(components) - 1:
            items.append("".join(row))
            row = []
            
            if i < len(components) - 1:  # Don't add </tr><tr> after last row
                items.append("</tr><tr>")
    
    grid_items = "\n".join(items)
    return GRID_TEMPLATE.substitute(GRID_ITEMS=grid_items)

def generate_theme_item(theme_name, theme_info, index):
    """Generate a single theme item for the grid"""
    preview_url = f"{RAW_URL}/{theme_info.get('preview_path', '')}"
    download_url = f"{RAW_URL}/Uploads/Themes/{theme_name}.zip"
    repo_url = f"{REPO_URL}/commits/main/Catalog/Themes/{theme_name}"
    
    # Extract info
    author = theme_info.get("author", "Unknown")
    description = theme_info.get("description", "")
    
    # You would get this from git history in a real implementation
    last_updated = "2025-04-23"  # Placeholder
    
    # Build title tooltip
    title = f"Name: {theme_name}\n"
    title += f"Author: {author}\n"
    title += f"Last updated: {last_updated}\n"
    if description:
        title += f"Description: {description}\n"
    title += "(Click to download)"
    
    # Calculate column spacing
    column_spanner = "&nbsp;" * 50 if index < COLUMNS else ""
    column_width = f"{100 / COLUMNS:.2f}%"
    
    # Author button
    author_btn = f'&nbsp;&nbsp;<a href="{REPO_URL}/search?l=ZIP&q=filename%3A%22{author}%22"><img src="{AUTHOR_ICON}" width="16" title="Search themes by this author"></a>' if author else ""
    
    # Preview button
    preview_btn = f'&nbsp;&nbsp;<a href="{preview_url}"><img title="View full-size preview" src="{PREVIEW_ICON}" width="16"></a>'
    
    # Readme button (would need to determine if readme exists)
    readme_btn = ""  # f'&nbsp;&nbsp;<a href="{readme_path}"><img src="{README_ICON}" height="16" title="README"></a>'
    
    # Icon pack indicator (would need to determine if theme has icons)
    has_iconpack = ""  # f' &nbsp;<a href="{icon_url}"><img src="{HAS_ICONPACK_ICON}" height="16" title="This theme contains an icon pack"></a>'
    
    return ITEM_TEMPLATE.substitute(
        NAME=theme_name,
        AUTHOR=author,
        TITLE=title,
        UPDATED=last_updated,
        PREVIEW_URL=preview_url,
        DOWNLOAD_URL=download_url,
        REPO_URL=repo_url,
        COLUMN_SPANNER=column_spanner,
        COLUMN_WIDTH=column_width,
        AUTHOR_BTN=author_btn,
        PREVIEW_BTN=preview_btn,
        README_BTN=readme_btn,
        HAS_ICONPACK=has_iconpack
    )

def generate_component_item(comp_type, comp_name, comp_info, index):
    """Generate a single component item for the grid"""
    preview_url = f"{RAW_URL}/{comp_info.get('preview_path', '')}"
    download_url = f"{RAW_URL}/Uploads/Components/{comp_type}/{comp_name}.zip"
    repo_url = f"{REPO_URL}/commits/main/Catalog/Components/{comp_type}/{comp_name}"
    
    # Extract info
    author = comp_info.get("author", "Unknown")
    description = comp_info.get("description", "")
    
    # You would get this from git history in a real implementation
    last_updated = "2025-04-23"  # Placeholder
    
    # Build title tooltip
    title = f"Name: {comp_name}\n"
    title += f"Author: {author}\n"
    title += f"Type: {comp_type.title()}\n"
    title += f"Last updated: {last_updated}\n"
    if description:
        title += f"Description: {description}\n"
    title += "(Click to download)"
    
    # Calculate column spacing
    column_spanner = "&nbsp;" * 50 if index < COLUMNS else ""
    column_width = f"{100 / COLUMNS:.2f}%"
    
    # Author button
    author_btn = f'&nbsp;&nbsp;<a href="{REPO_URL}/search?l=ZIP&q=filename%3A%22{author}%22"><img src="{AUTHOR_ICON}" width="16" title="Search components by this author"></a>' if author else ""
    
    # Preview button
    preview_btn = f'&nbsp;&nbsp;<a href="{preview_url}"><img title="View full-size preview" src="{PREVIEW_ICON}" width="16"></a>'
    
    # Readme button (would need to determine if readme exists)
    readme_btn = ""  # f'&nbsp;&nbsp;<a href="{readme_path}"><img src="{README_ICON}" height="16" title="README"></a>'
    
    return ITEM_TEMPLATE.substitute(
        NAME=comp_name,
        AUTHOR=author,
        TITLE=title,
        UPDATED=last_updated,
        PREVIEW_URL=preview_url,
        DOWNLOAD_URL=download_url,
        REPO_URL=repo_url,
        COLUMN_SPANNER=column_spanner,
        COLUMN_WIDTH=column_width,
        AUTHOR_BTN=author_btn,
        PREVIEW_BTN=preview_btn,
        README_BTN=readme_btn,
        HAS_ICONPACK=""
    )

def generate_pagination(current_page, total_pages, directory):
    """Generate pagination links"""
    # Previous page button
    prev_page = ""
    if current_page > 0:
        prev_link = "index.md" if current_page == 1 else f"page-{current_page}.md"
        prev_page = f"""<td align="right">\n\n[![Previous page]({PREV_PAGE_ICON})]({prev_link})\n\n</td>"""
    
    # Next page button
    next_page = ""
    if current_page < total_pages - 1:
        next_link = f"page-{current_page+2}.md"
        next_page = f"""<td>\n\n[![Next page]({NEXT_PAGE_ICON})]({next_link})\n\n</td>"""
    
    # Page links
    page_links = []
    for page in range(total_pages):
        page_num = page + 1
        if page == current_page:
            page_links.append(f"&nbsp;**{page_num}**&nbsp;")
        else:
            link = "index.md" if page == 0 else f"page-{page_num}.md"
            page_links.append(f"[&nbsp;{page_num}&nbsp;]({link})")
    
    page_links_str = " ".join(page_links)
    
    return PAGINATION_TEMPLATE.substitute(
        PREV_PAGE=prev_page,
        PAGE_LINKS=page_links_str,
        NEXT_PAGE=next_page
    )

if __name__ == "__main__":
    main()
