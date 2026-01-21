# Repertoire Web

- [Repertoire Web](#repertoire-web)
  - [Abstract](#abstract)
  - [Prerequisites](#prerequisites)
    - [NodeJs](#nodejs)
  - [Get Started](#get-started)
    - [Environment Variables](#environment-variables)
    - [Restore dependencies](#restore-dependencies)
    - [Run](#run)
  - [Build](#build)

## Abstract

This is the Web application, which is written completely in Typescript using **Vite** and **React**.
This platform project entirely uses the [UI](../repertoire.ui/README.md) package.
So we recommend that, after you finish with this setup, to get back on the UI and run the project from there.

## Prerequisites

Before you can get started, there are some things you need to have installed on your system.

### NodeJs

First, the application is written in Typescript (which compiles to JavaScript), so you need Node.js.
If you don't have it already installed, please download it from here:
`https://nodejs.org/en/download/prebuilt-installer`.

Note that the recommended version is: `v24.13.0` (**Nvm** <Node Version Manager> can be used to achieve that).

## Get Started

To get started on contributing to this project, you can do so by following the next steps.

### Environment Variables

Duplicate the `.env.dev` file and rename it to `.env` (the UI uses some variables too).

### Restore dependencies

To restore the dependencies, type the following command in the terminal:

```sh
npm ci
```

### Run

To run the application in development, type the following command:

```sh
npm run dev
```

## Build

To build the application, type the following:

```bash
npm run build
```
