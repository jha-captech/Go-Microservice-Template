# Go Microservice Templates

A collection of service templates for Go. Examples include both microservice and serverless architectures.

## Table of contents

- [Features](#features)
- [Templates](#templates)
- [Architecture](#architecture)
  - [Executables](#executables)
  - [Application Packages](#application-packages)
- [Project Structure](#project-structure)
- [Packages](#packages)
- [Architecture Goals](#architecture-goals)
- [Architecture Decisions](#architecture-decisions)
  - [Flat Packages](#flat-packages)
  - [`cmd` & `internal`](#cmd--internal)
- [Architecture Diagrams](#architecture-diagrams)
- [References](#references)

## Features

Each template in this repo is a fully deployable application and has the following features.

- Utilizes idiomatic Go and Industry best practices
  - Limited and judicious use of 3rd party libraries
  - Proper error handling
- Includes resources for the full lifecycle:
  - Local setup instructions
  - Testing
  - Infrastructure as code
  - CI/CD Pipeline
- Boilerplate code for common needs including but not limited to:
  - Logging
  - Database connections
  - Data validation and mapping for request and response models
  - Fully functional handlers
  - Basic service layer
  - Middleware
  - Graceful shutdown
  - Table-driven unit tests
  - E2E Integration tests

## Templates

| Template       | Description                                                            | Link                                                                             |
| -------------- | ---------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `multi-lambda` | A microservice comprised of multiple lambda functions. _**Preferred**_ | [Link](./scaffolds/multi-lambda/README.md#go-aws-lambda-template---multi-lambda) |
| `mono-lambda`  | A microservice comprised of a single monolithic lambda function.       | [Link](./scaffolds/mono-lambda/README.md#go-aws-lambda-template---mono-lambda)   |

## Architecture

The templates in this repo use a simple layered architecture. Executables are located under the `cmd` directory, and application code is located in packages under `internal`.

### Executables

Each executable under `cmd` is structured as a folder with the name of the executable, containing a single `main.go` file. This file has no logic outside of initializing dependencies and starting the runtime (Lambda, or an HTTP server).

```
.
└── cmd/
    └── lambda/
        └── main.go
```

### Application Packages

Application packages contain all of our application code and live under `internal`. This prevents outside applications and libraries from importing our application code, keeping our implementation details private.

| Package      | Description                               | Link                |
| ------------ | ----------------------------------------- | ------------------- |
| `config`     | Configuration definition and loading      | [Link](#config)     |
| `database`   | Database connection with retry logic      | [Link](#database)   |
| `handlers`   | Lambda or net/http handlers               | [Link](#handlers)   |
| `middleware` | Lambda or net/http middleware             | [Link](#middleware) |
| `models`     | Domain models                             | [Link](#models)     |
| `services`   | Domain services containing business logic | [Link](#services)   |
| `testutil`   | Common test utilities                     | [Link](#testutil)   |

## Project Structure

### `cmd`

#### Monolithic Lambda

```
.
└── cmd/
    └── api/
        └── main.go               # Monolithic application entrypoint
```

#### Multi Lambda

```
.
└── cmd/
    ├── create_user/
    │   └── main.go               # Create user lambda entrypoint
    ├── update_user/
    │   └── main.go               # Update user lambda entrypoint
    └── ...
```

### `internal`

The internal layout is mostly static between each scaffold. Internal contains the packages that implementing our application logic.

```
.
└── internal/
    ├── config/
    │   └── config.go             # Application config definition and loading
    ├── database/
    │   └── database.go           # Database connection logic
    ├── handlers/
    │   ├── api.go                # Monolithic lambda handler
    │   ├── create_user.go        # Individual lambda handler, includes request and response structs and validation
    │   ├── create_user_test.go   # Lambda handler test
    │   └── ...
    ├── middleware/
    │   ├── middleware.go         # Common type definitions and helpers for middleware
    │   └── recovery.go           # Sample lambda middleware for panic recovery
    ├── models/
    │   └── user.go               # Plain go structs representing domain models
    ├── services/
    │   └── user.go               # Domain service, center of all business logic
    └── testutil/
        └── ...                   # Common utilities for tests
```

## Packages

Each scaffold has common packages it implements, the details of which can be found below.

### `config`

config contains all application config as well as a function for loading config from environment variables.

### `database`

database contains standardized logic for connecting to a database and pinging the connection to ensure success.

### `handlers`

handlers contains handler functions. The style of handlers will depend on the service type. A Lambda service will have [Lambda style handlers](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-signatures), an HTTP service will have [net/http style handlers](https://pkg.go.dev/net/http#Handler).

We recommend writing handlers as closures, which allows dependencies to be passed without the additional boilerplate of defining a struct and method. This works especially well as handlers typically do not have many dependencies, and writing a closure is a very succinct way to supply dependencies.

#### HTTP Handler Example

```go
package handlers

type createUserRequest struct {
  Name string `json:"name"`
}

type createUserResponse struct {
  User models.User `json:"user"`
}

type userCreator interface {
  CreateUser(name string) (models.User, error)
}

func HandleCreateUser(logger *slog.Logger, service userCreator) http.Handler {
  return http.HandlerFunc(func(r *http.Request, w http.ResponseWriter) {
    // unmarshal, validate, call service method ...
  })
}
```

#### Lambda Handler Example

```go
package handlers

type createUserRequest struct {
  Name string `json:"name"`
}

type createUserResponse struct {
  User models.User `json:"user"`
}

type userCreator interface {
  CreateUser(name string) (models.User, error)
}

func HandleCreateUser(logger *slog.Logger, service userCreator) http.Handler {
  return func(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
    // unmarshal, validate, call service method ...
    return &events.APIGatewayProxyResponse{}, nil
  }
}

```

### `middleware`

middleware contains common middleware functions. Middleware style will differ between Lambda and HTTP services.

### `models`

models contains domain models for the application.

We recommend writing models as plain structs to keep them light and flexible. Struct tags can used to attach metadata for operations such as marshaling data to a database row.

### `services`

services contains our application services, where the core business logic of our application is defined.

We recommend defining services as structs. Defining services as structs allows us to share common dependencies such as database connections across methods.

We strongly advise against exporting interfaces for your services from this package. Interfaces in Go are more often used at point of consumption. This follows the "accept interfaces return structs" idiom for Go.

### `testutil`

testutil contains common testing utilities for marshaling and unmarshaling data and performing asserts.

## Architecture Goals

The following goals underpin many of the decisions for these templates, and help drive architectural decisions. These goals are in no particular order.

### Idiomatic

This architecture follows idomatic Go practices, such as favoring simple and obvious code, embracing small and private interfaces, minimizing abstractions and indirection, developing well named packages, developing code for tesability, and focusing on functionality first.

### Simple

This architecture strives to be as simple as possible, but no simpler. Go is a brutalist and spartan language. Development in Go favors simple and obvious code, not abstractions, frameworks, or large enterprise style architectures.

This architecture strikes a balance between extremely lean and flat architectures that are becoming more popular in the Go community, and architectures that are more familiar to engineers coming from other languages. This balance is struck without compromising on idomatic Go practices.

### Glanceable

Glanceability means "how easy is to get information at a glance". This architecture strives to be highly glanceable. Packages are well named and flat, file names all follow the same conventions, and code is organized in a way to get the high level details near the top of a file, followed by more specific implementation details.

A few examples of glanceability at work.

- Within `internal/handlers` its clear to see which handlers have tests and which do not, simply by the presence or absence of a `_test.go` file.
- The `cmd` folder makes it obvious how many executables a service contains.
  - `cmd` in `multi-lambda` contains multiple executables, one per lambda functions.
  - `cmd` in `mono-lambda` contains a single executable.
- Within `internal/handlers` we can see all the operations a service supports based on file names alone.
- Within `internal/config` we can find all the configuration our service needs.

### Scalable

Being scalable refers to the ability to scale the architecture up or down depending on the needs of the project. We've worked to provide a strong architecture that will conform to the vast majority of projects we'll be developing.

This architecture can be scaled to extreme levels, either down to a flat architecture for a simple project, or up to meet the needs of a larger and more complex product.

We recommend using an architecture that works best for the environment you're in, which often means starting with a middle ground approach like this until there is enough data to push in another direction.

### Maintainable

Developing a solution is only the first step. Software must also be maintained and extended over time. As noted above, the architecture present in these templates favors simplicity, which in turn helps with maintainability.

Because our packages are flat, and we have just a few layers in our architecture its less effort for us to add new functionality. And because we rely less on abstractions and more on concrete implementations, we have less concerns about internal breaking changes.

A simpler solution also means there are less corners for bugs to hide in, helping the overall maintenance burden.

## Architecture Decisions

### Flat Packages

Application packages declared under `internal` are flat. There are no nested packages.

This keeps code more obvious and clear. Nested packages quickly add additional mental hoops that developers must jump through and begin to make us lose sight of the ultimate goal of development and delivery of functionality.

### `cmd` & `internal`

The choice to use `cmd` to represent application binaries, and `internal` to hold application packages was made in accordance with the Go community's recommendations for module organization. Specifically recommendations for organizing modules that represent deployable artifacts.

[Recommendations on Module Organization](https://go.dev/doc/modules/layout#server-project)

## Architecture Diagrams

### Mono Lambda

![mono lambda system architecture](./scaffolds/mono-lambda/diagrams/Go%20Microservice%20Arch-Monolithic%20Lambda.drawio.svg)

### Multi Lambda

![multi lambda system architecture](./scaffolds/multi-lambda/diagrams/Go%20Microservice%20Arch-Multi%20Lambda.drawio.svg)

## References

- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
  - Great general guidance for writing Go applications.
- [Organizing a Go Module](https://go.dev/doc/modules/layout)
  - Great guidance for various ways to organize modules.
- [How I write HTTP services in Go after 13 years](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/)
  - Great tips and tricks for writing HTTP services. Matt Ryer falls in the lean and minimal Go camp, and some of his patterns may be too advanced for larger and less mature teams.

The provided projects use an architecture that strikes a balance between heavy enterprise architecture approaches and extremely lean approaches. By providing enough structure and guardrails to push developers towards the pit of success, without being so opinionated that you become painted into a corner and get stuck writing endless boilerplate code.

It's important to realize that architecture is about addressing organizational challenges as much as technical challenges. We believe these templates will provide a comprehensive starting point for most projects but you may need modifications depending on your project's circumstances. If ever in doubt, reach out to the maintainers of this repo, or members of the Go COP and have the discussion.
