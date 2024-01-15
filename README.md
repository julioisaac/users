# users

![ci status](https://github.com/julioisaac/users/actions/workflows/ci.yaml/badge.svg)

## Description

Project to handle with Users CSV file upload, Users API and integrations stuff

## Requirements

- [Docker](https://docker.com)
- [Go](https://go.dev)
  - [golangci-lint](https://golangci-lint.run/) (optional for linting)

## Instructions

### Environment

### Development

- **Install dependencies**: `make install`
- **Run locally**: `make run`
- **Run with Docker**: `make docker/up`
- **Lint**: `make lint`
- **Tests**: `make test`
- **Show available make commands**: `make help`

## How to Deploy

## Built With

- Golang 1.21
- Docker
- Make

## Application Checklist

- [ ] [Swagger](http://localhost:9000/)
- [ ] [APM]()
- [ ] [CI](https://github.com/julioisaac/users/actions)
- [ ] [CD](https://github.com/julioisaac/users/actions)
- [ ] [Sonar]()
- [ ] [Grafana]()
- [ ] [Logs]()

## Architecture


### Ingestion sequence flow

```mermaid
sequenceDiagram
  Worker->>+Broker: Ingest users csv file
  Broker-->>+Consumer: Consume users data
  Consumer-->>+Cache: Check cache

  alt Not found
    Consumer-->>Database: Save in database
    Consumer->>+Cache: Save in cache
  end
```

### Client request sequence flow

```mermaid
sequenceDiagram

  Client->>+API: GET user data
  API-->>+Cache: Check cache
  alt Not found
    API-->>Database: Check in database
    Database-->>API: Response from database
    alt Found
      API-->>Cache: Save in cache
    end
    API-->>Client: User data response from database
  else Found
    Cache-->>API: Response from cache
    API-->>Client: User data response from cache
  end

```

