# Currency Exchange Go Backend

Welcome to the **Currency Exchange Go Backend** project! This backend service provides real-time currency exchange rates and currency conversion functionalities via RESTful API endpoints. Built with Go, it leverages efficient concurrency and caching mechanisms to deliver fast and reliable data.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
  - [GET /rates](#get-rates)
  - [GET /convert](#get-convert)
- [Testing](#testing)
  - [Running Tests](#running-tests)
  - [Test Coverage](#test-coverage)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **Real-Time Exchange Rates**: Fetches the latest currency exchange rates from a reliable external API.
- **Currency Conversion**: Converts an amount from one currency to another based on real-time rates.
- **Caching Mechanism**: Implements caching to reduce external API calls and improve performance.
- **RESTful API**: Provides easy-to-use API endpoints for integration with front-end applications.
- **Extensive Testing**: Includes unit tests with high code coverage to ensure reliability.

---

## Architecture

The backend is structured into modular packages to separate concerns:

- **`handlers`**: Contains HTTP handlers for API endpoints.
- **`utils`**: Includes utility functions such as HTTP clients for external API calls.
- **`cache`**: Manages caching of exchange rates with expiration handling.
- **`models`**: Defines data models used across the application.

---

## Prerequisites

- **Go**: Version 1.16 or higher
- **Git**: For cloning the repository
- **An External API Key**: (Optional) If you plan to use a specific currency exchange API that requires authentication.

---

## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/paladinknightmaster/currency-exchange-go-backend.git
   cd currency-exchange-go-backend
   ```

2. **Install Dependencies**

   Ensure you have Go modules enabled and download the required packages:

   ```bash
   go mod download
   ```

---

## Configuration

### Environment Variables

The application can be configured using environment variables. You can create a `.env` file in the project root to set these variables.

**Example `.env` file:**

```env
API_KEY=your_api_key_here
```

- **`API_KEY`**: Your API key for the external currency exchange service.

---

## Running the Application

1. **Build the Application**

   ```bash
   go build -o currency-exchange-backend
   ```

2. **Run the Application**

   ```bash
   ./currency-exchange-backend
   ```

   The server will start and listen on port `8080` by default.

---

## API Endpoints

### GET /rates

Fetches the latest currency exchange rates.

- **URL**: `/rates`
- **Method**: `GET`
- **Success Response**:
  - **Code**: `200 OK`
  - **Content**:

    ```json
    {
      "base_code": "USD",
      "rates": {
        "EUR": 0.85,
        "GBP": 0.75,
        ...
      }
    }
    ```

- **Error Response**:
  - **Code**: `500 Internal Server Error`
  - **Content**:

    ```json
    {
      "error": "Error fetching rates"
    }
    ```

### GET /convert

Converts an amount from one currency to another.

- **URL**: `/convert`
- **Method**: `GET`
- **Query Parameters**:
  - **`from`**: The currency code to convert from (e.g., `USD`).
  - **`to`**: The currency code to convert to (e.g., `EUR`).
  - **`amount`**: The amount to convert (e.g., `100`).
- **Success Response**:
  - **Code**: `200 OK`
  - **Content**:

    ```json
    {
      "from": "USD",
      "to": "EUR",
      "amount": 100,
      "converted_amount": 85
    }
    ```

- **Error Responses**:
  - **Code**: `400 Bad Request` (missing parameters)
  - **Content**:

    ```json
    {
      "error": "Missing required query parameters"
    }
    ```

  - **Code**: `500 Internal Server Error` (conversion error)
  - **Content**:

    ```json
    {
      "error": "Error converting currency"
    }
    ```

---

## Testing

### Running Tests

The project includes unit tests to ensure functionality and reliability.

**Run all tests:**

```bash
go test ./... -v
```

### Test Coverage

Generate a test coverage report to identify untested parts of the codebase.

**Generate coverage report:**

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in your browser to view the coverage report.

---

## Project Structure

```
currency-exchange-go-backend/
├── cache/
│   ├── cache.go
│   └── cache_test.go
├── handlers/
│   ├── convert.go
│   ├── rates.go
│   └── handlers_test.go
├── models/
│   └── models.go
├── utils/
│   ├── http_client.go
│   ├── http_client_test.go
│   └── mock_http_client.go
├── main.go
├── main_test.go
├── go.mod
├── go.sum
└── README.md
```

- **`cache/`**: Caching logic and cache management.
- **`handlers/`**: HTTP handlers for API endpoints.
- **`models/`**: Data models and structures.
- **`utils/`**: Utility functions and HTTP client logic.
- **`main.go`**: Entry point of the application.
- **`main_test.go`**: Tests for the main package.

---

## Additional Notes

### Testing Strategy

- **Isolation of Tests**: Tests are designed to be independent and not interfere with each other.
- **Mocking External Dependencies**: External API calls are mocked to ensure tests do not rely on external services.
- **Avoiding Side Effects**: Tests do not modify or delete real files (e.g., `.env` file) to prevent unintended side effects.

### Environment Loading

- **Customizable Environment File**: The application allows specifying a custom environment file for flexibility.
- **Error Handling**: If the environment file is not found, the application logs a message but continues running, allowing for default configurations.

---