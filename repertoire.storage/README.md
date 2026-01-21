# Repertoire Storage

## Abstract

This is the Storage Service of the **Repertoire** application,
which is written completely in Go using the Gin Gonic Framework.

## Prerequisites

Before you can get started, there are some things you need to have installed on your system.

### Go

First, the application is written in Golang, so you need the Go SDK.
If you don't have it already installed, please download it from here: `https://go.dev/dl/`.

## Get Started

To get started on contributing to this project, you can do so by following the next steps.

### Environment Variables

Duplicate the `.env.dev` file and rename it to `.env`.

### Restore dependencies

To restore the dependencies, type the following command in the terminal:

```sh
go mod download
```

### Run

To run the application from the CLI, type the following command:

```sh
go run main.go
```

## Build Executable

If you want to build an executable, you can do so by typing the following command in the terminal:

```sh
go build
```

Usually, it will detect the system you are on and build an accustomed executable.

## Testing

### Unit Testing

The unit tests normally reside close to the unit under testing.

To run the tests you can run:

```sh
go test
```

## Live Reload

To use live reload in the application, install *Air* on your system, by using the following command:

```sh
go install github.com/air-verse/air@latest
```

Next, use the following command to build and run the application and start the live reloading:

```sh
air
```

The above command will look in the directory for a `.air.toml` file for configuration,
but that's already included in the application.