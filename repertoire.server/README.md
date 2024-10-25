# Repertoire Server

* [Repertoire Server](#repertoire-server)
  * [Abstract](#abstract)
  * [Prerequisites](#prerequisites)
    * [Go](#go)
    * [Docker](#docker)
  * [Get Started](#get-started)
    * [Docker Containers](#docker-containers)
    * [Restore dependencies](#restore-dependencies)
    * [Build](#build)
    * [Run](#run)
    * [Quick Test](#quick-test)
  * [Architecture](#architecture)
  * [Api](#api)
    * [Authentication](#authentication)
  * [Database](#database)
    * [Migration](#migration)
  * [Testing](#testing)
    * [Unit Testing](#unit-testing)
  * [Live Reload](#live-reload)

## Abstract

This is the HTTP Server (or backend) of the **Repertoire** application, 
which is written completely in Go using the Gin Gonic Framework.

## Prerequisites

Before you can get started, there are some things you need to have installed on your system.

### Go

First, the application is written in Golan, so you need the Go SDK.
If you don't have it already installed, please download it from here: `https://go.dev/dl/`.

### Docker

Developing becomes easier while using isolated containers provided by Docker, so,
if you do not have it installed already, please do so from here:
`https://www.docker.com/products/docker-desktop/` (you may have to restart afterward).

## Get Started

To get started on contributing to this project, you can do so by following the next steps.

### Docker Containers

If you decide to run the application in Docker instead, a `docker-compose.yml` file has been provided.
This file uses variables from the `.env` file, which is locally defined for each user.
Therefore, you can copy and paste the `.env.dev` file and rename it to `.env`.

Once you have the environment files setup, type in the following command:

```sh
docker compose up -d
```

This command will open up two containers: one for the database and the other for the application,
which is running by default in development mode.

If you decide to use a local database instance instead of Docker, then you are free to do so, 
but that's not included in the documentation.

If you decide that you only want the database from the docker compose file, you can run the following command:

```sh
docker compose up -d database
```

If you decide to run the backend application locally, follow the next steps.

### Restore dependencies

To restore the dependencies, type the following command in the terminal:

```sh
go mod download
```

### Build

To build the application, type the following command in the terminal:

```sh
go build
```

### Run

To tun the application from the CLI, type the following command:

```sh
go run main.go
```

### Quick Test

Now to be sure that everything works accordingly, try sending an HTTP Request (e.g., via Postman) to
`http://localhost:1123/api`.

```js
PUT {{host}}/auth/sign-in
```

And the body of the request shall be:

```json
{
    "email": "Some@Example.com",
    "password": "Password123"
}
```

It will return an Invalid Credentials Error, however,
now you know that you have a working connection to the API and to the database.

## Architecture

The architecture implemented in this application is following the principles of Clean Coding and Clean Architecture.
A few principles of the Domain-Driven-Design are also applied (i.e., each model has its own service or handler).

The four layers of this type of architecture feature the:
- **Domain** Layer—the place where the domain models are defined (defined in the _models_ package)
- **Infrastructure** Layer—the data access layer for the database (defined in the _data_ package)
- **Application** Layer—the layer that makes abstract calls to the infrastructure and makes business logic decisions 
(defined in the _domain_ package)
- **Presentation** Layer—the endpoint where the application becomes exposed (defined in the _api_ package)

To put it in simple terms, the workflow of the application would be the following:
- the http request comes through an _Api_ handler, where it is also validated
- then the Api sends it to a _Domain_ Service
- then it is being sent to a specific Use Case (still on the _domain_)
- which make calls to a _Data_ Repository or Service
- then the result is being passed back to the api and then to the user (the error middleware takes care of errors)

## Api

### Authentication

The developer shall provide a token for endpoints that do not allow anonymous requests.
One way to add one inside Postman is to go to the *Authorization* tab and under the type *JWT Bearer* add the following:
- **Secret Key**, which can be found in `.env`
- **Payload** which should include the following:
  - **sub**, the user id of an existing user (coincide with db data)
  - **jti**, the id of the token (any uuid)
  - **exp**, the expiration date in ms
  - **iss**, issuer of the application, which can be found in `.env`
  - **aud**, audience of the application, which can be found in `.env`

## Database

### Migration

When in development mode, the application will run **Auto Migrate** from GORM on startup.

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