# Go URL Shortener - Codebase Guidelines

This document outlines the guidelines and best practices for developing and maintaining the Go URL Shortener service in a production environment.

## 1. Code Style and Conventions

*   **Go Idioms:** Adhere strictly to Go's idiomatic style. Use `gofmt` to format code, and regularly run `golint` and `staticcheck` to catch common errors and style issues.
*   **Readability:** Prioritize clear, concise, and self-documenting code. Avoid cleverness over clarity.
*   **Comments:** Avoid comments on self-explanatory code. Add comments for:
    *   Complex algorithms or business logic.
    *   Public APIs, functions, and structs (using GoDoc format).
    *   Explanations of *why* a particular approach was taken, especially for non-obvious decisions or workarounds.
*   **Consistency:** Maintain consistent naming, formatting, and structural patterns across the entire codebase.

## 2. Testing (Test-Driven Development - TDD)

*   **TDD Approach:** Implement features using a Test-Driven Development (TDD) workflow:
    1.  Write a failing test.
    2.  Write the minimum code to make the test pass.
    3.  Refactor the code.
*   **Comprehensive Testing:**
    *   **Unit Tests:** Thoroughly test individual functions and methods in isolation, covering happy paths, edge cases, and error conditions. Aim for high code coverage.
    *   **Integration Tests:** Verify the interaction between different components (e.g., database, external services, API handlers).
    *   **End-to-End Tests:** Simulate user flows to ensure the entire system works as expected from a user's perspective.
*   **Maintainable Tests:** Tests should be fast, reliable, independent, and easy to understand. Use mocking or stubbing for external dependencies.

## 3. Logging

*   **Structured Logging:** Use structured logging (e.g., JSON format) to make logs easily parsable and searchable by log aggregation systems.
*   **Log Levels:** Utilize appropriate log levels:
    *   `DEBUG`: Detailed information for debugging.
    *   `INFO`: General operational messages, indicating progress or state changes.
    *   `WARN`: Potentially harmful situations that are not errors yet (e.g., retries, degraded performance).
    *   `ERROR`: Error events that require attention.
    *   `FATAL`: Critical errors that cause the application to terminate.
*   **Contextual Logging:** Include relevant contextual information in logs (e.g., request IDs, user IDs, trace IDs, function names) to aid in debugging and tracing issues across services.
*   **Avoid Sensitive Data:** Never log sensitive information (passwords, API keys, PII).
*   **Centralized Logging:** Design for integration with centralized log management systems (e.g., ELK Stack, Splunk, DataDog).

## 4. Committing Practices

*   **Small, Atomic Commits:** Each commit should represent a single, logical change. Avoid large, multi-purpose commits.
*   **Descriptive Commit Messages:** Write clear, concise, and descriptive commit messages that explain *what* changed and *why*. Follow conventional commit guidelines (e.g., `feat:`, `fix:`, `chore:`).
*   **Version Control:** Utilize Git effectively for branching, merging, and pull requests.

## 5. Naming Conventions

*   **Go Naming Conventions:** Follow Go's official naming conventions (e.g., PascalCase for exported identifiers, camelCase for unexported identifiers).
*   **Clarity and Brevity:** Choose names that are clear, unambiguous, and reasonably brief. Avoid abbreviations unless widely understood.
*   **Consistency:** Be consistent with naming within domains and across the codebase.

## 6. Production-Grade Package Structure

*   **Standard Go Project Layout:** Adhere to the standard Go project layout (e.g., `cmd`, `internal`, `pkg`, `api`, `web`, `build`, `deployments`, `scripts`).
*   **Separation of Concerns:** Organize code into logical packages based on functionality and responsibility. Each package should have a clear purpose and minimal external dependencies.
*   **Internal Package:** Use `internal/` for application-specific private code that other projects should not import.
*   **API Package:** Define API contracts, request/response models, and handler interfaces in a dedicated `api/` package.

## 7. Error Handling

*   **Idiomatic Go Errors:** Use Go's built-in `error` interface. Return errors explicitly as the last return value.
*   **Error Wrapping:** Use `fmt.Errorf` with `%w` to wrap errors, providing a traceable chain of error context. This allows for programmatic inspection of error types.
*   **Distinguish Errors:** Differentiate between transient errors (e.g., network issues, temporary service unavailability) and permanent errors (e.g., invalid input, unrecoverable database errors).
*   **Centralized Error Handling:** Implement a consistent strategy for handling and reporting errors at different layers of the application.
*   **Avoid Panics:** Panics should generally be reserved for unrecoverable program states (e.g., nil dereferences, out-of-memory). Use errors for recoverable situations.

## 8. Security Best Practices

*   **Input Validation:** Validate all input at the application's boundaries (e.g., API endpoints, form submissions) to prevent injection attacks (SQL, XSS), buffer overflows, and other vulnerabilities.
*   **Authentication & Authorization:** Implement robust authentication and authorization mechanisms. Use secure token management (e.g., JWT with proper signing/validation) and enforce least privilege.
*   **Secure Configuration:** Store sensitive configuration (database credentials, API keys) securely, ideally using environment variables, secrets management services (e.g., HashiCorp Vault, AWS Secrets Manager), and never hardcode them or commit them to source control.
*   **Dependency Management:** Regularly scan and update third-party dependencies to patch known vulnerabilities.
*   **HTTPS:** Enforce HTTPS for all communication.
*   **Rate Limiting:** Implement rate limiting to protect against abuse and DDoS attacks.

## 9. Performance Considerations

*   **Benchmarking:** Write benchmarks (`go test -bench`) for critical code paths to identify performance bottlenecks.
*   **Profiling:** Use Go's `pprof` tool to profile CPU, memory, and goroutine usage to optimize resource consumption.
*   **Concurrency:** Use goroutines and channels effectively, but avoid over-optimization prematurely. Profile to identify areas where concurrency helps.
*   **Database Optimization:** Optimize database queries, use appropriate indexing, and consider connection pooling.
*   **Caching:** Implement caching strategies (e.g., Redis) for frequently accessed, immutable data to reduce database load and improve response times.

## 10. Documentation

*   **API Documentation:** Provide clear and up-to-date API documentation (e.g., using OpenAPI/Swagger) for external consumers and internal frontends.
*   **README.md:** Keep the project's `README.md` updated with setup instructions, deployment guides, and key architectural decisions.
*   **Architectural Decision Records (ADRs):** Document significant architectural decisions and their justifications.
*   **GoDoc:** Leverage Go's built-in documentation tool for package, type, function, and method explanations.
