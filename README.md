# Go Microservice Templates

A collection of service templates for Go. Examples include both microservice and serverless architectures.

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

The makeup of `cmd` depends on the deployment style. For a mono-lambda deployment a single `cmd/lambda` executable can be found. For the multi lambda deployment an executable will be defined per function (`cmd/list`, `cmd/update`).

### Application packages

Application packages contain all of our application code and live under `internal`. This prevents outside applications and libraries from importing our application code, keeping our implementation details private.

## Pacakges

| Package      | Description                               | Link                |
| ------------ | ----------------------------------------- | ------------------- |
| `config`     | Configuration definition and loading      | [Link](#config)     |
| `database`   | Database connection with retry logic      | [Link](#database)   |
| `handlers`   | Lambda or net/http Handlers               | [Link](#handlers)   |
| `middleware` | Lambda or net/http middleware             | [Link](#middleware) |
| `models`     | Domain models                             | [Link](#models)     |
| `services`   | Domain services containing business logic | [Link](#services)   |
| `testutil`   | Common test utilities                     | [Link](#testutil)   |

More on architecture [goals](#architecture-goals), [project structure](#project-structure), and [architecture decision](#decisions) can be found below.

## Project structure

The scaffolds in this repo follow a similar layout, with executable placed under `cmd`, and application code placed in packages under `internal`.

A description of common packages can be found [below](#common-packages).

### `cmd`

#### Monolithic Lambda

```
.
└── cmd/
    └── lambda/
        └── main.go               # Monolithic application entrypoint
```

#### Multi Lambda

```
.
└── cmd/
    ├── create_user/
    │   └── main.go               # Create user lambda entrypoint
    ├── read_user/
    │   └── main.go               # Read user lambda entrypoint
    ├── update_user/
    │   └── main.go               # Update user lambda entrypoint
    └── delete_user/
        └── main.go               # Delete user lambda entrypoint
```

### `internal`

The internal layout has very little changes between each scaffold.

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

config contains a struct declaring all application config, as well as a function for loading config from environment variables.

### `database`

database contains standardized logic for connecting to a database and pinging the connection to ensure success.

### `handlers`

handlers contain handler functions. An API will contain handlers that conform to the `net/http Handler` interface, lambda projects will contain handlers that conform to lambda handler signatures.

handlers also contains unexported request and response models, as well as validation and mapping to validate requests before mapping them into a type recognized by the service layer.

Handlers are written as closures, which allows dependencies to be supplied to the handler, without the additional weight of defining a struct and method.

#### Example

##### HTTP

```go
package handlers

type createUserRequest struct {
  Name string `json:"name"`
}

type createUserResponse struct {
  User models.User `json:"user"`
}

func HandleCreateUser(logger *slog.Logger, service services.User) http.Handler {
  return http.HandlerFunc(func(r *http.Request, w http.ResponseWriter) {
    // unmarshal request body into struct, validate, call supplied user service
    // ...
  })
}
```

##### Lambda

```go
package handlers

type createUserRequest struct {
  Name string `json:"name"`
}

type createUserResponse struct {
  User models.User `json:"user"`
}

func HandleCreateUser(logger *slog.Logger, service services.User) http.Handler {
  return func(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
    // unmarshal request body into struct, validate, call supplied user service
    // ...
    return &events.APIGatewayProxyResponse{}, nil
  }
}

```

### `middleware`

middleware contains common middleware functions. Like above, middleware function signatures will differ based on the deployment target of the application.

### `models`

models contains domain models for the application. No other logic should be present within the models package. Struct tags are an appropriate way to enable meta programming on domain models, such as specifying specific DB columns to map values to and from.

### `services`

services contain the core business logic of our application.

### `testutil`

testutil contains common testing utilities for marshaling and unmarshaling data and performing asserts.

### Diagrams

TODO

## Architecture Goals

The following goals underpin many of the decisions for these templates, and help drive architectural decisions. These goals are in no particular order.

### Idiomatic

This architecture follows idomatic Go practices, such as favoring simple and obvious code, embracing small and private interfaces, minimizing abstractions and indirection, developing well named packages, and keeping our both our package surface and API surface clean and free of clutter.

### Simplicity

This architecture strives to be as simple as possible, but no simpler. Go is a brutalist and spartan language. Development in Go favors simple and obvious code, not abstractions, frameworks, or large enterprise style architectures.

The project templates in this repo all follow a simple layered architecture, with a `cmd` folder representing executables, and an `internal` folder holding the application packages. More details on the architecture can be found below.

This architecture strikes a balance between the extremely lean and flat architectures that are becoming more popular in the Go community, and architectures that are more familiar to engineers coming from other languages. This balance is struck without compromising on idomatic Go practices.

### Glanceability

When we say glanceability what we mean is "how easy is to get information at a glance". This architecture strives to be highly glanceable. Packages are well named and flat, file names all follow the same conventions, and code is organized in a way to get the high level details near the top of a file, followed by more specific implementation details.

A few examples of glanceability at work.

- Within the `internal/handlers` package its clear to see which handlers have tests and which do not, simply by the presence or absence of an `_test.go` file.
- Looking in `cmd` in the `lambda` and `mono_lambda` scaffolds we can immediately see how our application is deployed.
  - Under `cmd` in `lambda` we can see the executables representing our individual lambda functions.
  - Under `cmd` in `mono_lambda` we can see that the scaffold contains a single monolithic executable for the entire application.
- Within the `internal/handlers` package we can see a collection of all the operations a service supports.
- All of our application config is located in a single `internal/config` package, which gives us context to the full set of environment variables and configuration needed to run the service.

### Scaleability

For these templates, scaleability refers to the ability to scale the architecture up or down depending on the needs of the project. We've worked to provide sane starting templates, that will conform to the vast majority of projects we'll be developing.

This architecture could also be scaled to extreme levels, such as slimming down to a completely flat architecture (everything in package main), or pushing for higher degrees of structure and package hierarchy.

The provided architecture strikes a balance between those two extremes, by providing enough structure and guardrails to push developers towards the pit of success, without being so opinionated that you become painted into a corner or are stuck writing endless amounts of boilerplate code.

It's important to realize that architecture is about solving organizational challenges more than technical ones, and while we believe these templates will provide a starting point for most projects, you may need modifications depending on your project's circumstances. If ever in doubt, reach out to the maintainers of this repo, or members of the Go COP and have the discussion.

### Maintainability

Developing a solution is only the first step. Software must also be maintained and extended over time. As noted above, the architecture present in these templates favors simplicity, which in turn helps with maintainability.

Because our packages are flat, and we have just a few layers in our architecture its less effort for us to add new functionality. And because we rely less on abstractions and more on concrete implementations, we have less concerns about internal breaking changes.

A simpler solution also means there are less corners for bugs to hide in, helping the overall maintenance burden.

## Decisions

### Flat packages

Application packages declared under `internal` are flat. There are no nested packages.

This keeps code more obvious and clear. Nested packages quickly add additional mental hoops that developers must jump through and begin to make us lose sight of the ultimate goal of development and delivery of functionality.

### `cmd` & `internal`

The choice to use `cmd` to represent application binaries, and `internal` to hold application packages was made in accordance with the Go community's recommendations for module organization. Specifically recommendations for organizing modules that represent deployable artifacts.

[Recommendations on Module Organization](https://go.dev/doc/modules/layout#server-project)

## Scaffold Features

- Utilize idiomatic Go and Industry best practices
  - Limited and judicious use of 3rd party libraries
  - propper error handling
- Include resources for the full lifecycle:
  - Local setup instructions
  - Testing
  - Infrastructure as code
  - CI/CD Pipeline
- Boilerplate code for common needs including but not limited to:
  - Logging
  - Telemetry
  - Database connections
  - Data validation and mapping for request and response models
  - Fully functional handlers
  - Basic service layer
  - Middleware
  - Graceful shutdown
  - Table-driven unit tests
  - E2E Integration tests
