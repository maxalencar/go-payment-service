# Payment Service

Microservice that integrates with multiple payment gateways. It manages deposit and withdrawal operations, support multiple data interchange formats and handle asynchronous callbacks. It currently support two payment gateways (gatewayA and gatewayB), but it can be easily extended to handle many more as this service provides a PaymentGateway interface which is protocol-agnostic.

All endpoints support currently JSON and XML formats.

### Design Decisions

- [Go](https://golang.org/) 
    - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.
- [Circuit Breaker Pattern](github.com/sony/gobreaker)
    - A fault-tolerance mechanism designed to prevent repeated attempts to access an external service when it is experiencing failures or downtime.
    - It monitors the success/failure of requests to an external service.
    - Improves system stability by stopping wasteful calls and freeing resources.
- [Exponential Backoff Strategy](github.com/cenkalti/backoff/v4)
    - A retry strategy where failed requests are retried with progressively increasing delay between attempts.
    - Each retry increases the delay exponentially (e.g., 2 seconds, 4 seconds, 8 seconds, etc.), up to a maximum limit.
    - Improves resilience by allowing temporary failures (e.g., network blips) to self-recover without immediate user impact.
- [Table-driven tests using subtests](https://blog.golang.org/subtests) 
    - TDT were used as the approach to reduce the amount of repetitive code compared to repeating the same code for each test and makes it straightforward to add more test cases.

[Testify](https://github.com/stretchr/testify)
  
- [assert](https://github.com/stretchr/testify#assert-package) The assert package provides some helpful methods that allow you to write better test code in Go.
- [suite](https://github.com/stretchr/testify#suite-package) The suite package provides functionality that you might be used to from more common object-oriented languages. With it, you can build a testing suite as a struct, build setup/teardown methods and testing methods on your struct, and run them with 'go test' as per normal.
- [require](https://github.com/stretchr/testify#require-package) The require package provides same global functions as the assert package, but instead of returning a boolean result they terminate current test.

## Getting Set Up

Before running the application, you will need to ensure that you have a few requirements installed.

These instructions cover how to setup your development environment.

1. [Go](https://golang.org/)
2. [Make](https://www.gnu.org/software/make/) (optional for windows)

### Software Package Manager (Alternative if you don't have the requirements installed and would like to use a SPM)

#### MacOS

1. Install [homebrew](https://brew.sh/)
2. Install Go `brew install go`

#### Windows

1. Install [chocolatey](https://chocolatey.org/)
2. Install Go `choco install golang`
3. Install Make `choco install make`

## Project Structure

Following [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

### `/api`

OpenAPI/Swagger specs.

### `/cmd`

Main application for this project.

### `/internal`

Internal application logic.

### `/pkg`

Content that can be shared with external applications.

### `/test`

Additional e2e test and payment gateway emulator.

## Running the server

### How to Run:

1. Open a terminal.
2. Go to project's root path
3. Run the following command: 
    1. Running using defaults: `go run .\cmd\app\main.go`

### Alternatively using Makefile

1. Open the terminal
2. Go to project's root path
3. Run the following command: `make run`

### Alternatively using docker-compose

1. Open the terminal
2. Go to project's root path
3. Run the following command: `docker compose up --build`. (update the port if you are not running in default port 8080)

## Running the tests

### Running all tests

Follow the same steps `1` and `2` from the previous section and run the following command:

``` shell
    go test ./... -v
```

### Alternatively using Makefile

``` shell
    make test
```

### Endpoints

#### POST /deposit

It manages a deposit

Example:

Request:

    curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"amount": {"amount": 10, "currency": "EUR"},"cardDetails": {"number": "4111111111111111", "name": "Test", "expiryMonth": 10, "expiryYear": 2030, "cvv": "123"}, "gatewayDetails": {"id": "gatewayA", "callbackUrl": "http://localhost:8080/callback"}}' \
    http://localhost:8080/deposit

Response:
    
    {"transactionId":"35eda736-f26c-476a-a1c9-2f3ee0870d26","status":"pending","processedAt":"0001-01-01T00:00:00Z"}

#### POST /withdrawal

It manages a withdrawal.

Example:

Request:

    curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"amount": {"amount": 10, "currency": "EUR"},"cardDetails": {"number": "4111111111111111", "name": "Test", "expiryMonth": 10, "expiryYear": 2030, "cvv": "123"}, "gatewayDetails": {"id": "gatewayA", "callbackUrl": "http://localhost:8080/callback"}}' \
    http://localhost:8080/withdrawal

Response:
    
    {"transactionId":"a66c584c-b36d-4e2e-a55d-477a5b961e9f","status":"pending","processedAt":"0001-01-01T00:00:00Z"}

#### GET /transactions/{id}

Get a transaction by ID.

Example:

Request:

    curl -X GET http://localhost:8080/b78946ba-80ad-432b-9a13-38598c680095

Response:
    
    {"id":"60526b13-3260-4b28-aaa6-edeefa68eb6f","amount":{"amount":10,"currency":"EUR"},"cardDetails":{"name":"Test","number":"4111111111111111","type":"","expiryMonth":10,"expiryYear":2030,"cvv":"123"},"gatewayDetails":{"id":"gatewayA","name":"","callbackUrl":"http://localhost:8080/callback"},"type":"deposit","status":"succeeded","externalId":"da0b91e4-331b-43e1-ad53-4d046105c210","createdAt":"2024-09-30T15:28:40.364145671Z","updatedAt":"2024-09-30T15:28:40.365270855Z"}

## Future Improvements

- Add account feature to manage customers / balances
- Implement tokenization service for card holder PCI information
- Add mask sensitive data.
- Improve validation bad request response.
- Improve error response when processing payment.
- Discover card type based on the card number.
- Add support to more data formats.
- Add config layer.
- Add more test cases.

[@maxalencar](https://github.com/maxalencar)