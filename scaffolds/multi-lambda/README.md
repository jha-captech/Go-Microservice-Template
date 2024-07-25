# Go AWS Lambda Template - Multi Lambda

A multi lambda project scaffold following Go best practices and conventions.

## Instruction

### Setup

Required dependencies:

- Go 1.22+
- Docker
- Docker-Compose
- AWS SAM CLI

### Run

#### SAM Local API

```zsh
make lambda_local_api
```

#### SAM Local - list users event

```zsh
make lambda_local_list_users
```

#### SAM Local - update user event

```zsh
make lambda_local_update_user
```

## Architecture

![system architecture](./diagrams/Go%20Microservice%20Arch-Monolithic%20Lambda.drawio.svg)
