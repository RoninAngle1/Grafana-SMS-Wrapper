# Grafana-SMS-Wrapper - With Authentication
This is a Webhook for Grafana SMS Wrapper
# Overview

This Go code creates a server that listens for incoming Grafana webhook alerts and processes them by sending formatted messages to a configured API. Below is a breakdown of the code.

## Imports and Setup

The program uses various Go packages:
- **`bytes`**: To handle byte slices, particularly for request bodies.
- **`encoding/json`**: For JSON encoding/decoding.
- **`fmt`**: For formatted I/O.
- **`io/ioutil`**: To read HTTP responses.
- **`log`**: For logging errors and other messages.
- **`net/http`**: For HTTP server and client functionalities.
- **`os`**: To interact with the operating system, particularly for reading the config file.

## Structures

### Config

Holds the configuration details, including:
- **Server settings**:
  - `IP`: Server's IP address.
  - `Port`: Server's port.
  - `Path`: Server's endpoint path for receiving webhook requests.
- **PhoneNumbers**: List of phone numbers to which alerts will be sent.
- **URL**: The API URL to which alerts will be forwarded.
- **UserName**: Username for API authentication.
- **Password**: Password for API authentication.

### GrafanaAlert

Defines the structure of the incoming Grafana alert payload, which includes:
- **Alerts**: A list of alerts, each containing:
  - `Title`: Title of the alert.
  - `Annotations`: Contains a `Description` field with additional details about the alert.

## Main Function

### Workflow:
1. **Reads Configuration**: The configuration file (`config.json`) is loaded using the `readConfig` function.
2. **Sets up HTTP Server**: The server listens for HTTP POST requests at the path defined in the config (`config.Server.Path`).
3. **Handles Incoming Requests**: 
   - On receiving a POST request, it decodes the incoming JSON payload (Grafana alert data).
   - For each alert, it formats a message by combining the `title` and `description`.
   - The message is sent to each phone number in the configuration via the `sendMessage` function.
   - If successful, the server responds with an HTTP 200 status and a success message. Otherwise, it returns an error.

## Supporting Functions

### `readConfig`

Reads the `config.json` file and returns a `Config` structure. If there is an issue reading or decoding the file, an error is returned.

### `sendMessage`

Sends the formatted alert message to the specified API URL for each phone number. The function:
- Constructs a JSON body containing the message and phone number.
- Marshals the body into JSON.
- Sends an HTTP POST request to the API with custom headers (`UserName` and `Password` for authentication).
- Logs the request and the response for debugging.
- If the API responds with a non-OK status, an error is returned.

## Key Workflow

1. **Server Start**: The server listens at the IP and port defined in the configuration.
2. **Handling Webhook**: When a Grafana webhook is triggered, the server decodes the payload and sends alerts to each phone number via an API call.
3. **API Communication**: Each alert is sent to an external API (as defined in the `config.json`) using the username and password for authentication.
4. **Error Handling**: The program handles various errors:
   - Invalid JSON in the incoming request.
   - Issues reading or parsing the configuration file.
   - Errors when sending messages (e.g., network issues, API errors).
   - Non-OK HTTP responses from the API are logged and returned as errors.

## Example Workflow

1. A Grafana alert is triggered and sent to the server's webhook endpoint (e.g., `http://<server_ip>:<port>/alert`).
2. The server decodes the alert's JSON payload and extracts the title and description of the alert.
3. The server sends the alert message (formatted with the title and description) to each phone number defined in the configuration by making an HTTP request to the destination API.
4. The server responds to Grafana with a success message if all messages are sent successfully or an error if something fails.

## Overall Purpose

This program acts as a bridge that receives Grafana alerts and forwards them as notifications (via API calls) to a set of phone numbers defined in the configuration file.
