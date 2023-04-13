# GoEvents Recorder API

Event Recorder API is a simple RESTful API service that records events from different clients such as browsers or mobile applications. It provides a single endpoint /record_event to record events and store them into a PostgreSQL database. The API also provides additional information about the client's IP address and User Agent if specified.

## Table of Contents
<!-- vim-markdown-toc GFM -->

* [Requirements](#requirements)
* [Installation](#installation)
* [Configuration](#configuration)
* [Running the API](#running-the-api)
* [Running the tests](#running-the-tests)
* [Integrating with Frontend](#integrating-with-frontend)
  * [Usage](#usage)

<!-- vim-markdown-toc -->

## Requirements
- Go 1.17+
- PostgreSQL 13.0+

## Installation

1. Clone the repository:

```bash
git clone https://github.com/Valpiccola/GoEvents
```

2. Change to the project directory:

```bash
cd GoEvents
```

3. Install the required Go packages:

```bash
go get -d ./...
```

## Configuration

The application expects several environment variables to be set for configuration:

- **DB_USER**: PostgreSQL database user.
- **DB_PASS**: PostgreSQL database password.
- **DB_HOST**: PostgreSQL database host.
- **DB_PORT**: PostgreSQL database port.
- **DB_NAME**: PostgreSQL database name.
- **DB_SCHEMA**: Database schema where the event table is located.
- **ENV**: Environment of the application (i.e., production, staging).
- **IPINFO_TOKEN**: IPInfo API token for IP address details retrieval.

## Running the API
1. Build the API binary:

```bash
go build
```

2. Run the API server:

```bash
./GoEvents 
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

## Integrating with Frontend
To use this API in your frontend application, you can create a registerEvent.js file. This file contains a function called registerEvent that sends an event to the API using the fetch function.

### Usage
1. Create a new file named registerEvent.js in your frontend project:
````javascript
import {
  PUBLIC_API_HTTP_URL,
} from "$env/static/public"

export async function registerEvent(page, event_name, deep, details) {
  fetch(PUBLIC_API_HTTP_URL + "/record_event", {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      Page: page,
      Event_name: event_name,
      Referrer: document.referrer,
      Cookie: document.cookie.match('(^|;)\\s*userId\\s*=\\s*([^;]+)')?.pop() || '',
      Size: innerWidth.toString()+"x"+innerHeight.toString(),
      Language: navigator.language,
      Deep: deep,
      Details: details
    })
  });
}
```

2. Set the PUBLIC_API_HTTP_URL environment variable in your frontend application to the API's base URL (e.g., http://localhost:8080).

3. Import the registerEvent function in your frontend application and call it when you need to record an event:
```javascript
import { registerEvent } from "./registerEvent";

// Example usage
registerEvent("home", "button_click", true, { button_id: "my-button" });
```

This function takes four parameters:

- page: The name of the page where the event occurred (e.g., "home").
- event_name: The name of the event (e.g., "button_click").
- deep: A boolean value indicating if additional information about the client's IP address and User Agent should be retrieved. Set to true to enable this feature.
- details: An object containing any additional details related to the event.

The registerEvent function will send a POST request to the /record_event endpoint of the API with the provided information.

This way, you can easily integrate the Event Recorder API with your frontend application and start recording events in your application.
