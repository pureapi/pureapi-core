# PureAPI Core

PureAPI Core is a modular framework for building robust, efficient, and
maintainable RESTful APIs in Golang. It abstracts common API concerns such as
error handling, middleware management, and database operations, allowing
developers to focus on delivering business logic.

## Introduction

PureAPI Core is the foundational repository of the PureAPI framework.
It provides a consistent, extensible base for developing REST APIs by offering:
- Pre-built modules for handling common tasks.
- A focus on type safety and clarity.
- Robust testing for both public interfaces and critical functions.

Whether you're building microservices, enterprise-grade web
applications, or prototyping new ideas, PureAPI Core gives you the tools to
get started quickly and scale with confidence.

## Features

- **Modular Architecture:**  
  Easily extend and integrate custom middleware, plugins, and extensions.
- **Custom Error Handling:**  
  Structured error messages that simplify debugging and client error
  processing.
- **Middleware Management:**  
  Flexible stacking and ordering to handle cross-cutting concerns.
- **Database Abstraction:**  
  Generic, type-safe CRUD operations, dynamic SQL generation, and safe
  transaction management.

## Use Cases

PureAPI Core serves as a solid foundation for a variety of REST API projects,
including but not limited to:

- **Microservices:** Spin up lightweight services with a consistent development experience.
- **Enterprise Web Services:** Build scalable and maintainable services.
- **Rapid Prototyping:** Quickly prototype ideas with minimal setup and clear, tested modules.
- **API Gateways:** Aggregate and manage multiple services through a unified API layer.

## Quick Start

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.24.0 or higher recommended)

### Installation

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/pureapi/pureapi-core.git
cd pureapi-core
```

### Running an Example

PureAPI Core comes with several examples. For instance, to run the basic HTTP server example:

```bash
cd doc/examples/serverbasic/
go run example.go
```

### Documentation

For more detailed information, please refer to the documentation provided in the [doc](./doc/index.md) directory.
