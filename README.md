# PureAPI Core

PureAPI Core is a modular framework for building web APIs in Go. It provides a consistent,
extensible foundation for developing RESTful services by abstracting common concerns
such as error handling, middleware management, database operations, and utility functions.

## Features

- **Custom Error Handling:** Structured API errors that are easy to parse.
- **Middleware Management:** Flexible stacking and ordering of middleware functions.
- **Database Abstraction:** Generic, type-safe CRUD operations and dynamic SQL generation.
- **Transaction Management:** Safe transaction handling with automatic rollback and commit.
- **Utilities:** Context management, event emission, logging, panic recovery, and more.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/)

### Installation

Clone the repository:

```bash
git clone https://github.com/pureapi/pureapi-core
cd pureapi
