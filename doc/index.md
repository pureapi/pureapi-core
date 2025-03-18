# PureAPI Core Documentation

This is the documentation for PureAPI Core. It is designed to help you get started quickly, understand the architecture and features of the system.

## Overview

PureAPI Core is a modular framework for building robust, scalable RESTful APIs in Golang. It abstracts common tasks such as error handling, middleware management, database operations, and more to help you focus on your business logic.

## Getting Started

This section will help you set up your development environment, install PureAPI Core, and run your first example.

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

# Design Principles

PureAPI Core is built with a focus on simplicity, modularity, and extendability. These design principles guide the architecture and implementation of the system, ensuring that it remains easy to use, maintain, and extend over time.

## Core Principles

These are the core design principles:

- **Simplicity and Modularity:**  Separation of concerns makes the codebase easier to understand, test, and extend.

- **Extendability:** Flexible extension and customization.

- **Minimal Third-Party Dependencies:** Minimize reliance on 3rd party libraries. Currently such dependencies are used primarily for testing and examples..

## Testing Strategy

The first rule is that we will have tests... eventually. Apart from that, the following testing principles guide development:

- **Public Contracts First:** Testing and development focus on the public contracts of the code. Changes must ensure that external behavior remains consistent and predictable.

- **Security-Critical Testing:** More critical sections should have additional tests to fight against potential vulnerabilities.

- **Private Function Testing:** Direct testing of private functions is secondary and only done when absolutely necessary, and reserved for more stable and critical code.

- **Incremental Test Coverage:** Strive for "high enough" test coverage. Higher coverage all the way to 100% can be an ambitious end goal for the most stable code only.

- **Integration Testing:** Where it makes sense, integration tests are employed to ensure that different components work together seamlessly.

# Architectural Overview

## Endpoint Package

The **Endpoint Package** is for defining and managing your API's endpoints. It provides a structured way to register endpoints, apply middleware, and encapsulate common endpoint logic for a consistent and extensible API design.

## Key Components

### Endpoints

Endpoints represent the basic building blocks of your API. Each endpoint is defined by:
- **URL:** The route at which the endpoint is accessible.
- **HTTP Method:** The type of request (e.g., GET, POST).
- **Handler:** A function that processes the request and generates a response.
- **Middlewares:** Functions that wrap around the handler to perform shared tasks like logging, authentication, or input validation.

*Example:*  
For an API endpoint that creates a new user, you would register it with a URL such as `/users`, use the `POST` method, and assign a handler that validates input and creates the user. Additional middlewares can be applied to enforce security (e.g., verifying API tokens) and logging.

### Middlewares

Middlewares are functions that wrap an HTTP handler to extend or modify its behavior. They can:
- Execute pre-processing before the main handler (e.g., authentication checks).
- Perform post-processing after the main handler (e.g., logging or response formatting).
- Handle errors and enforce policies.

*Example:*  
A middleware can intercept a request to check if a valid API token is present. If the token is missing or invalid, the middleware can reject the request before it reaches the endpoint’s handler.

### Middleware Wrappers

Middleware wrappers encapsulate individual middlewares along with an identifier and optional metadata. They allow you to:
- Assign a unique ID to each middleware for reference.
- Attach additional configuration or data to the middleware.
- Reorder or selectively enable/disable middlewares based on their identifiers.

*Example:*  
Suppose you have two logging middlewares with different purposes. By wrapping them with unique IDs and metadata, you can easily manage their order or replace one without changing the other.

### Stacks

A stack is a collection of middleware wrappers that forms a custom chain of middlewares. Instead of attaching each middleware individually, you can:
- Define a middleware stack that bundles multiple wrappers.
- Apply the entire stack to an endpoint.
- Create different stacks for various API needs (e.g., public vs. authenticated endpoints).

*Example:*  
For a public API endpoint, you might create a stack that includes rate limiting, input validation, and logging. For a private endpoint, the stack could additionally include authentication and authorization wrappers.

### Endpoint Definitions

Endpoint definitions offer a declarative approach to specify endpoints. Similar to endpoints they combine they use the endpoint URL, HTTP method, handler, but use a stack instead of middlewares. This abstraction simplifies:
- Consistent endpoint registration with reusable middleware stacks.
- Separation of concerns between the endpoint’s logic and its specific middleware configuration.

*Example:*  
Define an endpoint for fetching orders with the URL `/orders`, method `GET`, and assign a middleware stack that includes caching and logging. This definition simplifies the process of applying the same set of middlewares across similar endpoints.

### Endpoint Handler

The endpoint handler implements a generic, reusable logic flow for processing requests. Its responsibilities include:
1. **Input Processing:** Using an input handler to validate and parse request data.
2. **Business Logic Execution:** Calling a handler logic function to perform the core operation.
3. **Output Processing:** Using an output handler to format and send the response.
4. **Error Handling:** Mapping and handling errors via an error handler.

*Example:*  
A generic handler for a "create resource" endpoint might first validate input, then call a service function to create the resource, and finally format the response. Developers can implement their own input and output handlers to customize behavior while reusing the common flow provided by the generic handler.

# Database Package

The **Database Package** is designed to simplify interactions with SQL databases by providing a consistent, abstracted interface for connecting, querying, managing transactions, and handling errors.

## Key Components

### Connection Management

The package uses a `ConnectConfig` structure to hold all necessary configuration details such as:
- **Driver and Credentials:** Information about the database driver, user, password, host, and port.
- **Connection Parameters:** Including the database name, connection limits, and optional DSN format.
- **Runtime Settings:** Configuration of connection lifetimes and pool sizes.

The `Connect` function utilizes this configuration to:
- Establish a connection to the database.
- Configure the connection with appropriate settings.
- Validate the connection by pinging the database.

*Example:*  
You can connect to an in-memory SQLite database for testing or switch to a production-grade MySQL/PostgreSQL database by simply adjusting the `ConnectConfig` parameters.

### Common Database Operations

A suite of functions is provided to perform common database operations in a simplified manner:
- **Executing Queries:** Functions like `Exec` and `ExecRaw` run queries without returning rows.
- **Querying Data:** Functions such as `Query`, `QueryRaw`, and `QuerySingleValue` help in retrieving data.
- **Result Handling:** Helper functions like `RowToEntity` and `RowsToEntities` convert raw SQL results into Go data structures.

*Example:*  
To retrieve the number of users in the database, you might use `QuerySingleValue` to execute a count query, automatically handling preparation, execution, and scanning of the result.

### SQL Abstraction

The system abstracts the standard SQL types from the `database/sql` package by defining a set of interfaces in the `types` package. The `sqlDB` implementation wraps the native `*sql.DB` to:
- Conform to these interfaces.
- Provide a uniform error handling and connection management layer.
- Enable easier testing by allowing mock implementations.

*Example:*  
By abstracting the SQL layer, you can swap the actual database connection with a mock during testing.

### Transaction Management

The package simplifies the handling of business logic within transactions through the `Transaction` function, which:
- Wraps your transactional business logic (`TxFn`) into a safe execution context.
- Automatically commits the transaction on success or rolls it back if an error occurs.
- Recovers from panics to prevent the database from entering an inconsistent state.

*Example:*  
When you need to perform multiple interdependent operations—such as creating an order and updating stock levels—you can wrap them in a transaction to ensure atomicity, where either all operations succeed or none do.

### Error Checking

An optional **Error Checker** interface allows you to implement custom logic to translate database-specific errors into meaningful application errors. This approach:
- Ensures consistency in error handling across your application.
- Provides more informative error messages that can be used in client responses.

*Example:*  
If an insert operation violates a unique constraint, an error checker can catch the raw SQL error and translate it into a custom error message, such as "Username already exists," which is more meaningful to the end user.

# Server Package

The **Server Package** is responsible for managing the HTTP server layer of your API. It provides a robust and configurable framework for:

- **Setting Up the HTTP Server:** Configuring a server with sensible defaults for timeouts and header sizes.
- **Request Routing and Endpoint Multiplexing:** Mapping incoming requests to the correct endpoint handlers based on URL and HTTP method.
- **Graceful Shutdown:** Listening for OS signals to ensure the server shuts down cleanly, allowing active requests to finish.
- **Panic Recovery:** Catching unexpected panics during request handling to prevent server crashes.
- **Event Emission and Logging:** Emitting lifecycle events and logging key actions, such as server start, endpoint registration, and errors.

## Key Components

### HTTP Server Setup

The `DefaultHTTPServer` function creates and returns a default HTTP server instance configured with default values, such as:
- **Read/Write Timeouts:** Set to 10 seconds to protect against slow clients.
- **Idle Timeout:** Set to 60 seconds to manage persistent connections.
- **Maximum Header Size:** Limited to 64KB to prevent excessive resource usage.

*Example:*  
In production, use `DefaultHTTPServer` to initialize a server that listens on a designated port with secure default settings, ensuring reliable performance under load.

### Request Routing

The server package implements a custom HTTP handler (`Handler`) that:
- Registers endpoints by URL and HTTP method with the provided handlers and middlewares.
- Registers a default "not found" handler when no endpoint matches the request.

*Example:*  
When you have multiple endpoints (e.g., `/users`, `/orders`), the handler maps requests to the appropriate handler based on both the URL and the method (GET, POST, etc.).

### Graceful Shutdown

Graceful shutdown is handled by:
- Listening for OS interrupt signals (e.g., SIGTERM, Interrupt).
- Initiating a shutdown process that allows existing requests to complete within a specified timeout.
- Emitting events to log the shutdown process and any errors encountered.

*Example:*  
During a deployment or restart, the graceful shutdown mechanism ensures that the server stops accepting new requests and cleanly terminates existing ones.

### Panic Recovery

To ensure stability, the server package wraps request handlers with a panic recovery mechanism that:
- Recovers from panics during request processing.
- Logs the panic event along with a stack trace.
- Returns a 500 Internal Server Error response to the client.

*Example:*  
If an unforeseen error occurs in an endpoint handler, the panic recovery mechanism catches the error, logs detailed diagnostics, and prevents the entire server from crashing.

### Event Emission and Logging

The server package utilizes an event emitter and logger to record significant server events, such as:
- Server startup and shutdown.
- Endpoint registration.
- Errors during startup, shutdown, or request processing.
- Panic occurrences and recovery actions.

*Example:*  
Integrate these events with your monitoring system to track server health, troubleshoot issues, and gain insights into server events.

# Getting Help

If you encounter issues or have suggestions, please refer to the Contributing Guidelines or open an issue or discussion on our GitHub repository.

Please also [check the examples](#running-an-example) section for help running examples.
