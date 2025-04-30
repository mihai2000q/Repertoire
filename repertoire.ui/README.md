# Repertoire UI

* [Repertoire UI](#repertoire-ui)
  * [Abstract](#abstract)
  * [Prerequisites](#prerequisites)
    * [NodeJs](#nodejs)
  * [Get Started](#get-started)
    * [Restore dependencies](#restore-dependencies)
    * [Typecheck](#typecheck)
  * [Platform Configuration](#platform-configuration)
  * [Folder Structure](#folder-structure)
  * [Testing](#testing)
    * [Unit Testing](#unit-testing)

## Abstract

This is the UI of the **Repertoire** application, which is written completely in TypeScript.
Notable libraries used to write this project are: **React**, **Mantine**, and **Redux**.
The idea behind this package is to have a centralized piece of code that can be reused for different platforms,
i.e., desktop or web.

## Prerequisites

Before you can get started, there are some things you need to have installed on your system.

### NodeJs

First, the application is written in TypeScript (which compiles to JavaScript), so you need Node.js.
If you don't have it already installed, please download it from here:
`https://nodejs.org/en/download/prebuilt-installer`.

Note that the recommended version is: `v20.18.0` (**Nvm** <Node Version Manager> can be used to achieve that).

## Get Started

To get started on contributing to this project, you can do so by following the next steps.

### Restore dependencies

To restore the dependencies, type the following command in the terminal:

```sh
npm ci
```

### Typecheck

As this is a package, it cannot be run standalone, so all you can do now is to run the typecheck for the application:

```sh
npm run typecheck
```

## Platform Configuration

As previously mentioned, this package cannot be run as it requires a platform, so, if you use _WebStorm_,
you can create a configuration for a platform of your choice (desktop or web).

Before you can do that, you should set up any platform.
See
[desktop documentation](../repertoire.desktop/README.md)
or
[web documentation](../repertoire.web/README.md).

Now that you have a platform up and running, you can run it while editing the UI project as well by doing the following:

1. Go to **Edit Configuration**
2. Click on **Add New Configuration**
3. Choose **npm**
4. Then, on the Configuration Tab, under the **package.json** search for the directory of the platform you want
5. Then choose the command **dev** for development

Congratulations! Now you can run the desktop application from the UI package.

## Folder Structure

As this is a small-sized project, the screens, or the **views** are all in one directory.

Inside the **components**, if they stand on the root, they are commonly used for all the views.
If the component is in another package, then the package name designates the view that those components are used for.

The **data** package is used for constant variables for UI information that doesn't change based on the server
(e.g., sidebar links).

Inside the **hooks** there are the custom hooks of the application.

Inside the **router** there are the components that are doing the router security redirection
(e.g., if the user is not authenticated, then redirect to sign in page).

Inside the **state** package, the setup of the **Redux** apis and store reside.
There are also some separate files for each type (i.e., song) if there are too many endpoints (about 3-4 is the limit).

Inside the **theme** package, the theme of the application is defined.

Inside the **types** package, there are the types of the application, such as models (i.e., Song), or http requests.

Inside the **validation** package, there are forms for each view and their specific validation schema
(designed using **Zod**).

## Testing

### Unit Testing

The unit tests are done in isolation on the components, so, naturally, those reside in the **UI** project.
The choice for this small-sized project was to write unit tests as close as possible to the target component.

To be noted, that in the root **src** folder you can find the `test-utils` file that provide render wrappers.

The framework for unit testing is **Vitest** (also, testing library).
