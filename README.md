# Go-Microservice-Template

A collection of service templates for Go. Examples include both microservice and serverless architectures.

## Architectural Goals

These architectures are all underpinned by a similar set of goals, outlined below.

### Simplicity

First and foremost, this architcture strives to be as simple as possible, but no simpler. Go is a brutalistic and spartan language. Development in Go favors simple and obvious code, not abstractions, frameworks, or large enterprise style architectures.

The project templates in this repo all follow a simple layered architecture, with a `cmd` folder representing executables, and an `internal` folder holding the application packages. More details on the architecture can be found below.

This architecture strikes a balance between extremely lean and flat architectures that are becoming increasingly popular in the Go community, and architectures that are more familiar to engineers coming from other languages. This balance is struck without compromising on idomatic Go practices.

### Glanceability

### Scaleability

## Requirments

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

## TODO
