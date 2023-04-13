# GoEvents Recorder API

Event Recorder API is a lightweight, efficient, and scalable RESTful service designed to capture user events from web and mobile applications. It provides a straightforward and secure way to log events and store them in a PostgreSQL database for further analysis. The API also supports enriching event data with client IP address and User Agent information when requested.

## Table of Contents
<!-- vim-markdown-toc GFM -->

* [Features](#features)
* [Use Cases](#use-cases)
* [Getting Started](#getting-started)
  * [Requirements](#requirements)
  * [Installation](#installation)
  * [Configuration](#configuration)
  * [Preparing the database](#preparing-the-database)
* [Running the API](#running-the-api)
  * [Running the tests](#running-the-tests)
* [Integrating with Frontend](#integrating-with-frontend)
  * [Usage](#usage)
* [Deploy it on DigitalOcean](#deploy-it-on-digitalocean)

<!-- vim-markdown-toc -->

## Features
- Simple and easy-to-use RESTful API
- Secure event recording with CORS support
- IP address and User Agent enrichment (optional)
- Extensive test coverage
- Seamless frontend integration using a provided JavaScript function
- Scalable and adaptable design


## Use Cases
- Web analytics and user behavior tracking
- A/B testing and feature rollout monitoring
- Performance and error reporting
- Custom event logging for tailored insights

## Getting Started

### Requirements
- Go 1.17+
- PostgreSQL 13.0+

### Installation

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

### Configuration

The application expects several environment variables to be set for configuration:

- **DB_USER**: PostgreSQL database user.
- **DB_PASS**: PostgreSQL database password.
- **DB_HOST**: PostgreSQL database host.
- **DB_PORT**: PostgreSQL database port.
- **DB_NAME**: PostgreSQL database name.
- **DB_SCHEMA**: Database schema where the event table is located.
- **ALLOWED_ORIGINS**: Endpoints of your frontend framework
- **ENV**: Environment of the application (i.e., production, staging).
- **IPINFO_TOKEN**: IPInfo API token for IP address details retrieval.

### Preparing the database
Before you can start recording events and analyzing the captured data, it is essential to create a proper query to extract the relevant information from the stored events. 

```sql
CREATE TABLE IF NOT EXISTS event (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    details JSONB NOT NULL
);
```

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

### Running the tests

To run the tests, execute the following command:
```
go test -v ./...
```

The tests will automatically set the DB_SCHEMA environment variable to test and create a test database connection. Make sure to have a proper test database and schema configured before running the tests.

## Integrating with Frontend
To use this API in your frontend application, you can create a registerEvent.js file. This file contains a function called registerEvent that sends an event to the API using the fetch function.

### Usage
1. Create a new file named registerEvent.js in your frontend project:
```javascript
const PUBLIC_API_HTTP_URL = '';

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

## Deploy it on DigitalOcean

To deploy this API endpoint on DigitalOcean, you'll need to create a folder named ".do" in your project's root directory and add an appropriate app.yaml file within it.

```yaml
name: YOUR_APP_NAME
services:
- name: YOUR_SERVICE_NAME
  github:
    repo: YOUR_REPO
    branch: YOUR_BRANCH
```

