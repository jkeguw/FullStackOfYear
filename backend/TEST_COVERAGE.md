# Backend Test Coverage

This document provides information about the test coverage for the backend services.

## Service Tests

| Service          | Coverage | Notes                                |
|------------------|----------|--------------------------------------|
| `auth`           | N/A      | Login, token management, authentication tests |
| `device`         | ~85-90%  | Device management, user device configuration tests |
| `email`          | ~69%     | Email sending tests |
| `jwt`            | ~90%     | JWT token generation and validation tests |
| `measurement`    | ~62%     | Measurement data management tests |
| `review`         | ~80-85%  | Review creation, management, and moderation tests |
| `sensitivity`    | 90.6%    | Sensitivity calculation tests for various calibration methods |
| `token`          | ~90%     | Token management, revocation tests |
| `user`           | N/A      | User management tests |

## Test Strategy

The testing strategy focuses on:

1. **Unit Tests**: Testing individual functions and methods for correctness
2. **Service Layer Tests**: Testing service interfaces with mocked dependencies
3. **Integration Tests**: Testing interactions between services (where applicable)

## Mocking Approach

- MongoDB repositories are mocked using interfaces and mock implementations
- External services are mocked to avoid network calls
- Testify/mock is used to create and verify expectations

## Areas for Improvement

- Add more tests for user service
- Improve test coverage for email service
- Add integration tests with embedded MongoDB for full-stack testing
- Add more edge case tests for error handling

## How to Run Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./services/sensitivity

# Run tests with coverage
go test -cover ./...

# Run a specific test
go test -run TestServiceName_MethodName ./services/package
```

## Sample Test Sessions

```bash
$ go test -cover ./services/sensitivity
ok      FullStackOfYear/backend/services/sensitivity     0.006s  coverage: 90.6% of statements
```