# Repertoire

![Typescript](https://img.shields.io/badge/TypeScript-3178C6?style=flat&logo=typescript&logoColor=white)
![Electron](https://img.shields.io/badge/Electron-47848F?style=flat&logo=electron&logoColor=white)
![React](https://img.shields.io/badge/React-0081A3?style=flat&logo=react&logoColor=white)
![Redux](https://img.shields.io/badge/Redux-764ABC?style=flat&logo=redux&logoColor=white)
![Mantine](https://img.shields.io/badge/Mantine-339AF0?style=flat&logo=mantine&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Go Gin](https://img.shields.io/badge/Gin-008ECF?style=flat&logo=gin&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-38B6FF?style=flat&logo=go&logoColor=white)
![Goose](https://img.shields.io/badge/Goose-00ADD8?style=flat&logo=duckduckgo&logoColor=white)
![Postgres](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white)
![Meilisearch](https://img.shields.io/badge/Meilisearch-FF5CAA?style=flat&logo=meilisearch&logoColor=white)
![Centrifugo](https://img.shields.io/badge/Centrifugo-AF5365?style=flat)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)
![Oracle](https://img.shields.io/badge/Oracle-FF7900?logo=icloud&style=flat&logoColor=white)

![Version](https://img.shields.io/badge/version-0.12.0-1098ad)
 
* [Repertoire](#repertoire)
  * [Abstract](#abstract)
  * [Get Started](#get-started)
    * [Server](#server)
    * [Storage](#storage)
    * [UI](#ui)
    * [Platform](#platform)
      * [Desktop](#desktop)
      * [Web](#web)
  * [Git](#git)
    * [Git Issues](#git-issues)
    * [Git Branches](#git-branches)
    * [Git Commits](#git-commits)
  * [Deployment](#deployment)
    * [Dockerized Backend](#dockerized-backend)
    * [Desktop](#desktop-1)
  * [License](#license)

## Abstract

**Repertoire** is a Music Application written in Typescript and Go.
It is an application intended for beginner musicians to organize their '_repertoire_' (songs, albums, etc.).
Its features include adding songs, organized in albums and artists.
It is also possible to add those songs on playlists.
Songs can be measured how rehearsed they have been, whether they have been recorded or not and other features.

## Get Started

To get started on the application, you need:
- the main server
- the authentication server
- the storage server
- the UI 
- and one of the platforms up and running

IF you have Go and Node installed, you can take the shortest route. 
Fire the following shell script: `startup-web.sh` or `startup-desktop.sh`.

Otherwise, a full setup on each project has to be made 
(except for Desktop or Web, depending on where you prefer to host the app).

Additionally, if you prefer to run everything in containers and not install GO or Node, 
you can run the `run-container.sh` script to run all the docker containers and apply the database migrations.

### Auth

This is the Server that takes care of the authentication of the application with components and platforms.
<br>
For a detailed setup on the auth of the application, check out the [documentation](repertoire.auth/README.md).

### Server

This is the HTTP Server of the application, and it does all the business logic and calls to the database.
<br>
For a detailed setup on the server of the application, check out the [documentation](repertoire.server/README.md).

### Storage

This is the Server that takes care of the storage of the application. It takes care of images, audios and other type of files.
<br>
For a detailed setup on the storage of the application, check out the [documentation](repertoire.storage/README.md).

### UI

This is the UI package of the application, and it renders all the screens, 
makes all the calls and takes care of the style of the application using **React**, **Mantine**, and more.
<br>
For a detailed setup on the UI of the application, check out the [documentation](repertoire.ui/README.md).

### Platform

You are not obliged to set up both of the below platforms, you can only choose one of them, depending on preferences.

#### Desktop

This is the desktop project that takes care of rendering the **UI** in a desktop environment using **Electron**, 
**React**, and **Vite**.
<br>
For a detailed setup on the server of the application, check out the [documentation](repertoire.server/README.md).

#### Web

This is the web project that takes care of rendering the **UI** in a web environment using **React** and **Vite**.
<br>
For a detailed setup on the server of the application, check out the [documentation](repertoire.server/README.md).

## Git

The application is stored on the well-known cloud-based service and version control **GitHub**.

### Git Issues

Usually, whenever the code has to be modified, an issue will be created on **GitHub**.
The title of the issue will be the developer's choice
(for preference, it should resemble the git commits to a higher level overview without, of course, the issue tag),
however, the labels attached to it should be significant.
The title should also include the **project** that it's being worked on and the layer, or screen/component.

The expected labels to use are:

- **project** that is being updated:
  - **auth** if the auth is the topic of the issue
  - **desktop** if the desktop is the topic of the issue
  - **server** if the server is the topic of the issue
  - **storage** if the storage is the topic of the issue
  - **ui** if the ui is the topic of the issue
  - **web** if the web is the topic of the issue
- **type** of the issue:
  - **build** if the issue is supposed to update the dependencies
  - **bug** if the issue is supposed to solve a problem
  - **code-quality** when the code is refactored for better readability (interchangeable with **refactor**)
  - **devOps** when the changes are directed to operations outside the code, i.e., Docker, GitHub Actions
  - **documentation** if the source code is not the target, but its documentation
  - **feature** for new features that are added to the application
  - **performance** when the code is being optimized
  - **refactor** when the code is just being cleaned up (interchangeable with **code-quality**)
  - **styling** when only the style is being affected (usually used for the front end components)
  - **test** if a unit or integration test for the source code is being written

### Git Branches

Based on the above label, the branch will have similar prefixes:

- **auth** for an _auth_ issue
- **desktop** for a _desktop_ issue
- **server** for a _server_ issue
- **storage** for a _storage_ issue
- **ui** for an _ui_ issue
- **web** for a _web_ issue
- **build** for a _build_ issue
- **bug** for a _bug_ issue
- **code** for a _code-quality_ issue
- **devOps** for a _devOps_ issue
- **doc** for a _documentation_ issue
- **feature** for a _feature_ issue
- **perform** for a _performance_ issue
- **refactor** for a _refactor_ issue
- **style** for a _styling_ issue
- **test** for a _test_ issue

For example, if your ticket's name is "\[Server\] Data - New Repository for Menus"
and it is a feature issue, then the branch name should look something like this: `feature/139-server-data-...`.

Also, the main development branch is **develop**, and **master** is used for releases.

### Git Commits

The commits on branches should follow a structure (similar to the issues) 
(*#TAG PROJECT - LAYER/SCREEN: COMPONENT - CHANGES*).
First, there should be a hashtag followed by the number of the issue (e.g., #3)
followed by the project's name that is the target of the commit (e.g., web, server)
followed by a dash and the changes done on it (or the component that was modifier, dash, then the changes)
(by preference, try to include a meaningful verb like add or update).
For example, `#3 Server - Data - Add New Repository for Menu` 
or `#3 Server - Data: User Repository - Add New method to return by id`.

Note: if the target of the commit has nothing to do with any layer or screen, 
no worries, just try to be as concise as possible.
However, if multiple layers or projects are affected, try submitting more smaller commits.

## Deployment

### Dockerized Backend

The deployment is fairly simple. 
As soon as a pull request from *develop* is pushed onto *master* a CI/CD pipeline will be triggered.
This Github Action will build and push the server docker image to **Docker Hub**.

Afterwards, the cloud server detects the latest docker images and pulls them.

The _storage_ or the authentication services are not udnergoing changes so often, 
so there is not benefit for automating the deploying process.

### Desktop

Another CI/CD pipeline is in charge of deploying a new windows installer of the application.
It will be included in a new Release on **GitHub**.

## License

The project is licensed under the [MIT](https://opensource.org/license/mit) license.

---

**Repertoire**© 2024-2025