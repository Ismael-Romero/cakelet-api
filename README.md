# Cakelet Api

## Getting Started

### Running the Application

To run the application locally, execute the following command from the root directory of the project:

```bash
go run .
```

## Tech Stack & Dependencies

This project is built using the following core language version and dependencies:

### Environment
* **Go Version:** `1.26.5`

### Core Dependencies
| Package | Version | Purpose / Description |
| :--- | :---: | :--- |
| `github.com/golang-jwt/jwt/v5` | `v5.3.1` | JSON Web Token implementation for Go |
| `go.uber.org/fx` | `v1.24.0` | Application framework and dependency injection |
| `go.uber.org/zap` | `v1.26.0` | Blazing fast, structured, logged logging framework |
| `go.uber.org/dig` | `v1.19.0` | Underlying graph-based dependency injection toolkit |
| `go.uber.org/multierr` | `v1.10.0` | Combining multiple Go errors into a single error |
| `golang.org/x/crypto` | `v0.54.0` | Supplementary Go cryptography libraries |
| `golang.org/x/sys` | `v0.47.0` | Low-level operating system primitive interactions |


# Development Standards & Architecture

This project is built following industry standards and modern architectural patterns:

*   **Code Style & Guidelines:** The source code adheres to the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) to ensure high readability, consistency, and maintainability across the codebase.
*   **Dependency Injection:** The application's lifecycle and dependency management are fully orchestrated using [Uber Fx](https://uber-go.github.io/fx/), promoting modularity, clean architecture, and decoupled components.
