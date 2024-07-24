# Go-Microservice-Template

A collection of service templates for Go. Examples include both microservice and serverless architectures.

## Architecture

### Layout

TODO

### Diagrams

TODO

## Architectural Goals

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

## Architecture Decisions

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
