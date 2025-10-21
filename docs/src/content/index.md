---
title: "Home"
publish: true
template: index
date: 2025-10-21
author: author
tags:
    - tag-on-index-page
---

# Garlic: Fast and Simple Static Site Generator in Go

_By Shreyas Kaundinya_


---

## Table of Contents

- [0. Contact Me](#0-contact-me)
- [1. Intro & Installation](#1-intro--installation)
- [2. Running the project](#2-running-the-project)
- [3. Folder Structure](#3-folder-structure)
- [4. Architecture](#4-architecture)
- [5. Deployment](#5-deployment)
- [6. TODO & Limitations](#6-todo--limitations)

---

## 0. Contact Me 

If you have any questions, feel free to drop an DM on [X (Formerly Twitter)](https://x.com/shreyassk08)

Checkout my website [here](https://www.kaundz.com) built using Garlic SSG.

---

[Back to top](#table-of-contents)

## 1. Intro & Installation

### 1.1 Features

- Classic Static Site Generator with some opinionated features
- Converts **markdown content into HTML** which can be served as a static website.
- **HTML components** using JSX like syntax (with limited support as of now)
- **Templates** for content
- **Tags** using frontmatter
- **Hot reloading** support for development (reloads when content is changed)

### 1.2 Installation / Build Instructions

> NOTE: Since I currently don't have a release, you'll have to build the project yourself.

Step 1: Clone the [repository](https://github.com/shreyaskaundinya/garlic)

Step 2: Run the following commands

```bash
go mod tidy

# windows
go build -o garlic.exe main.go

# linux
go build -o garlic main.go

# macos
go build -o garlic main.go
```

---

[Back to top](#table-of-contents)

## 2. Running the project

### 2.1 Command Line Arguments

- `--src-folder`: The source folder of the project
- `--dest-folder`: The destination folder of the project
- `--serve`: Whether to serve the project [serves the project at `http://localhost:8084`]. This also enables hot reloading support for when the content is changed.
- `--seed-files`: Whether to seed the project [adds the default files to your source folder]

### 2.2 Examples

```powershell
.\garlic.exe --src-folder "S:\src" --dest-folder "S:\dest" --serve --seed-files
```

```bash
./garlic --src-folder ./src --dest-folder ./dest --serve --seed-files
```

---

[Back to top](#table-of-contents)

## 3. Folder Structure

The folder structure is fixed, the base folders that are required in the source folder are :

- [`src/content`](#32-content-folder--srccontent)
- [`src/templates`](#33-templates-folder--srctemplates)
- [`src/assets`](#34-assets-folder--srcassets)
- [`src/components`](#35-components-folder--srccomponents)

We will dive deeper into what each folder signifies in the next section.

> Don't worry about the boilerplate code, it will be seeded for you if you use the `--seed-files` flag.

### 3.1 Tree

```text
src/
├── content/
   ├── index.md
├── templates/
   ├── index.html
   ├── _tags.html
   ├── _individual_tag.html
├── components/
   ├── Footerbar.html
   ├── Navbar.html
   ├── Tags.html
├── assets/
   ├── styles/
       ├── global.css
└── dest/
    ├── assets/
       ├── styles/
           ├── global.css
    ├── index.html
    ├── tags/
       ├── index.html
       ├── tag-on-index-page/
           ├── index.html
```

### 3.2 Content Folder : `src/content`  

The content folder is where you put your markdown files. The routing of the website is based on the file structure inside this folder.

So if you want to create a page at `/about`, you can create a file at `src/content/about/index.md`.

The rendered HTML will be in `dest/about/index.html`.

Each markdown file will contain something called as **frontmatter**. This is a way to add metadata to the markdown file. It is written between `---` lines. Follows the YAML syntax.

#### Content Example

```markdown
---
title: "Home"
publish: true
template: index
date: 2025-10-21
author: author
tags:
    - tag-on-index-page
---

<!-- rest of the content of the page -->

# Home
```

- `title`: The title of the page, this can be used in the template
- `publish`: Whether to publish the page, if false, the page will not be rendered
- `template`: The template to use for the page, this is the template that will be used to render the page
- `date`: The date of the page. 
- `author`: The author of the page.
- `tags`: The tags of the page. **Atleast one tag is required per page**.

> NOTE: currently there is no support for `date` and `author` in the templates.

### 3.3 Templates Folder : `src/templates`

Each markdown file in the `src/content` folder will be rendered into an HTML file in the `dest` folder. 

Templates contain HTML that will wrap around the markdown content. 

> NOTE : Currently only supports flat file structure in this folder.

Content of the template is injected into the `{{ $content }}` placeholder.

```html
<main>{{ $content }}</main>
```

For the above [example](#content-example), the rendered HTML will be:


```html
<main>
  <h1>Home</h1>
</main>
```

Title of the template from the frontmatter is injected into the `{{ $title }}` placeholder.

```html
<title>{{ $title }}</title>
```

For the above [example](#content-example), the rendered HTML will be:

```html
<title>Home</title>
```

Some special templates required for internal purposes are:

> NOTE: working actively to make these templates more flexible and powerful. A default template will be provided during `seeding`

- `_tags.html` : This is the template where all the tags will be listed with their count.
  - `{{ $content }}` will be replaced with `<ul>...</ul>` tag with each tag as a `<li>`. Each `<li>` tag will have a `<a>` tag with the tag name and the count of number of pages that have that tag. The `<a>` tag will have the href to the individual tag page where all the pages with that tag will be listed.
- `_individual_tag.html` : This is the template where a single tag and all the pages with that tag will be listed.
  - `{{ $content }}` will be replaced with the list of pages with that tag.

### 3.4 Assets Folder : `src/assets`

Contains all the static assets of the website. This folder will be directly copy pasted into your destination folder. The files can be accessed using the `/assets/` prefix.

For example, if you want to access the `global.css` file, you can do it by using the `/assets/styles/global.css` path.

```html
<link rel="stylesheet" href="/assets/styles/global.css" />
```

### 3.5 Components Folder : `src/components`

You can add smaller HTML components here. These will be injected into the templates.

> NOTE : Currently only supports flat file structure in this folder.

> WARN: The components cannot have the same name as a native html element such as `div`, `span`, `a`, `img`, etc.


Let's say we create a component called `Navbar.html` 


To use it in a template, the syntax is as follows:

```html
<Navbar></Navbar>
```

Limitations (which could later be supported):
- Self closing tags are not supported.
- No support for conditional rendering.
- No support for loops
- No support for props for components

---

[Back to top](#table-of-contents)

## 4. Architecture

> This section is a work in progress.

Garlic is a static site generator that is built in Go. It is designed to be fast and simple to use.

Markdown files written in the `src/content` folder are parsed and rendered into HTML, this is then injected into the templates. 

Components are small HTML files that can be injected into the templates based on the file name, similar to JSX.

`src/components/Navbar.html`

```html
<nav>
  <a href="/">Home</a>
</nav>
```

`src/templates/index.html`

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{{ $title }}</title>
    <link rel="stylesheet" href="/assets/styles/global.css" />
  </head>
  <style></style>
  <body>
    <Navbar></Navbar>
    <main>{{ $content }}</main>
  </body>
</html>
```

`src/content/index.md`

```markdown
---
title: "Home"
publish: true
template: index
author: author
tags:
    - tags
---

# Home

This is the home page.
```

The rendered HTML will be:

`dest/index.html`


```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Home</title>
    <link rel="stylesheet" href="/assets/styles/global.css" />
  </head>
  <body>
    <nav>
      <a href="/">Home</a>
    </nav>

    <main>
      <h1>Home</h1>
      <p>This is the home page.</p>
    </main>
  </body>
</html>
```

---

[Back to top](#table-of-contents)

## 5. Deployment

Once you follow these steps and run Garlic on your content, you have freshly baked static website 🧄🥖 ready to be deployed.

These are some ways you can deploy your website:

- simply copying the `dest` folder to your web server.
- pushing to a Github repository and deploying using Github Pages
- pushing to a git repository and deploying using Vercel, Netlify, etc which provide you with a free domain and a way to do continuous deployment.

---

[Back to top](#table-of-contents)

## 6. TODO & Limitations

- [ ] Not require atleast one tag per page
- [ ] Add concurrency support to rendering
- [ ] Deleting unused files from destination folder
- [ ] Being able to use author from frontmatter in templates
- [ ] Being able to iterate over posts and tags in templates
- [ ] Being able to use date from frontmatter in templates
- [ ] RSS feed generation
- [ ] Components
  - [ ] No support for self closing tags
  - [ ] No support for conditional rendering
  - [ ] No support for loops
  - [ ] No support for props for components

---