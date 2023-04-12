# GoEvents Recorder API

Event Recorder API is a simple RESTful API service that records events from different clients such as browsers or mobile applications. It provides a single endpoint /record_event to record events and store them into a PostgreSQL database. The API also provides additional information about the client's IP address and User Agent if specified.

## Table of Contents
Requirements
Installation
Configuration
Running the API
Running the tests

## Requirements
- Go 1.17+
- PostgreSQL 13.0+

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/event-recorder-api.git
```

2. Change to the project directory:

```bash
cd event-recorder-api
```

3. Install the required Go packages:

```bash
go get -d ./...
```

## Configuration

The application expects several environment variables to be set for configuration:

DB_USER: PostgreSQL database user.
DB_PASS: PostgreSQL database password.
DB_HOST: PostgreSQL database host.
DB_PORT: PostgreSQL database port.
DB_NAME: PostgreSQL database name.
DB_SCHEMA: Database schema where the event table is located.
ENV: Environment of the application (i.e., production, staging).
IPINFO_TOKEN: IPInfo API token for IP address details retrieval.

## Running the API
1. Build the API binary:

```bash
go build -o event-recorder-api
```

2. Run the API server:

```bash
./event-recorder-api
```

The API server will start listening on port 8080. You can test the /record_event endpoint with a simple curl command:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"Event_name": "test_event"}' http://localhost:8080/record_event
```

## Running the tests

To run the tests, execute the following command:
```
go test -v ./...
```

The tests will automatically set the DB_SCHEMA environment variable to test and create a test database connection. Make sure to have a proper test database and schema configured before running the tests.
