## Backend - README

<h1 align="center"><b>Nubrik</b></h1>
<h2 align="center"><b>Backend</b></h2>

<p align="center">System for managing medical appointments</p>

- [Backend - README](#backend---readme)
- [Introduction](#introduction)
- [Technologies Used](#technologies-used)
- [Required Dependencies](#required-dependencies)
- [Configuration](#configuration)
- [Installation](#installation)
- [Development \& Kubernetes](#development--kubernetes)
- [APIs](#apis)
- [Contributing](#contributing)
- [License](#license)
- [Improvements](#improvements)

## Introduction

**Nubrik** backend is the RESTful API service for handling authentication, appointment logic, and interaction with the PostgreSQL database.
It is designed as part of a full-stack containerized application, meant to be deployed locally using Kubernetes.

## Technologies Used

- [Node.js](https://nodejs.org/)
- [Express.js](https://expressjs.com/)
- [TypeScript](https://www.typescriptlang.org/)
- [PostgreSQL (client)](https://node-postgres.com/)
- [Docker](https://www.docker.com/)
- [Kubernetes (K8s)](https://kubernetes.io/)

## Required Dependencies

- Node.js >= 18
- npm
- Docker
- Kubernetes + Minikube

## Configuration

Configuration is done through environment variables. For development:

```
PORT=3000
DB_HOST=postgres
DB_PORT=5432
DB_USER=root
DB_PASSWORD=root
DB_NAME=my_database_pg
```

These are configured inside the `ConfigMap` or `.env` file depending on environment. For K8s, config is injected through Kustomize.

## Installation

```bash
npm install
npm run build
npm start
```

## Development & Kubernetes

You can build and deploy this backend using the central Kubernetes project repo:

```bash
make build-backend     # from the Kubernetes deployment directory
make deploy            # deploy entire stack
```

To see logs:
```bash
make logs-backend
```

Access the API at: `http://my_appointment_backend.local/`

## APIs

<!-- TODO -->

## Contributing

1. Fork the repo
2. Create a new branch
3. Commit your changes
4. Push and create a PR

## License

MIT

## Improvements

- Add Swagger/OpenAPI documentation
- Implement role-based access control (RBAC)
- Add unit and integration tests
- Introduce healthcheck route
