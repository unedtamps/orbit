# Torrent Indexer API Gateway

This project provides a simple API gateway for searching content (movies, books, TV shows) through Jackett. It uses the `go-jackett` library to interact with a running Jackett instance and exposes a RESTful API.

## Features

*   Search for movies by query.
*   Search for books by query.
*   Search for TV shows by query.
*   Integrated Swagger UI for API documentation.

## Prerequisites

Before running this project, you need:

*   **Go (1.20 or later):** For local development.
*   **Docker and Docker Compose:** For running the application and Jackett in containers.
*   **Jackett Instance:** The API relies on a running Jackett instance. The `docker-compose.yaml` file includes a Jackett service for convenience.

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/your-repo/jacket.git # Replace with your actual repository URL
cd jacket
```

### 2. Configure Environment Variables

Create a `.env` file in the root of the project with the following variables:

```
API_URL=http://jackett:9117 # Or the URL of your Jackett instance
API_KEY=YOUR_JACKETT_API_KEY
HOST_URL=http://localhost:9999 # Or the external URL where your API will be accessible
```

*   **`API_URL`**: The URL of your Jackett instance. If you're using the provided `docker-compose.yaml`, `http://jackett:9117` is the correct value when the API service is running within the same Docker network. For local development where Jackett is running on `localhost`, it would be `http://localhost:9117`.
*   **`API_KEY`**: Your Jackett API key, which can be found on your Jackett dashboard.
*   **`HOST_URL`**: The URL where your API gateway will be externally accessible. This is used by Swagger UI to correctly reference API definitions. For local development, `http://localhost:9999` is suitable.

### 3. Running with Docker Compose (Recommended)

This is the easiest way to get both Jackett and the API Gateway running.

```bash
docker compose up --build
```

This command will:
1.  Build the API gateway Docker image.
2.  Start a Jackett container (accessible at `http://localhost:9117`).
3.  Start the API gateway container (accessible at `http://localhost:9999`).

### 4. Running Locally (Development)

First, ensure you have a Jackett instance running separately (e.g., via Docker or installed directly on your machine).

```bash
go mod tidy
go run main.go
```

The API will be accessible at `http://localhost:9999`.

## API Endpoints

The API provides the following endpoints:

### Search Movies

*   **URL:** `/movies/{query}`
*   **Method:** `GET`
*   **Description:** Searches for movies based on the provided query.
*   **Parameters:**
    *   `query` (path parameter, string, required): The movie title or keyword to search for.
*   **Responses:**
    *   `200 OK`: Returns a `jackett.Result` object containing search results.
    *   `500 Internal Server Error`: If an error occurs during the search.

### Search Books

*   **URL:** `/books/{query}`
*   **Method:** `GET`
*   **Description:** Searches for books based on the provided query.
*   **Parameters:**
    *   `query` (path parameter, string, required): The book title or keyword to search for.
*   **Responses:**
    *   `200 OK`: Returns a `jackett.Result` object containing search results.
    *   `500 Internal Server Error`: If an error occurs during the search.

### Search TV Shows

*   **URL:** `/tv/{query}`
*   **Method:** `GET`
*   **Description:** Searches for TV shows based on the provided query.
*   **Parameters:**
    *   `query` (path parameter, string, required): The TV show title or keyword to search for.
*   **Responses:**
    *   `200 OK`: Returns a `jackett.Result` object containing search results.
    *   `500 Internal Server Error`: If an error occurs during the search.

## API Documentation (Swagger UI)

Once the API is running, you can access the interactive Swagger UI at:

`http://localhost:9999/apidocs/index.html`

This interface allows you to explore the available endpoints, understand their parameters, and even test them directly from your browser.

## Generating API Documentation

To generate automated API documentation, you can use the `swag` command from [github.com/swaggo/swag](https://github.com/swaggo/swag). Then run:

```bash
swag init --parseDependency
```
