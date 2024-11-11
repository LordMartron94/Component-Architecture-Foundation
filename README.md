# Component-Based Architecture Foundation

## Table of Contents
- [Introduction](#introduction)
- [Architecture Overview](#architecture-overview)
  - [Example Sequences](#example-sequences)
- [JSON Scheme](#json-scheme)
  - [Registration](#registration)
  - [Keep Alive](#keep-alive)
  - [Example Request](#example-request)
- [Changelog](#changelog)

## Introduction
The Component-Based Architecture Foundation (CBAf) framework aims to provide a foundation for building scalable,
modular, and maintainable systems.
It focuses on creating a robust and efficient communication infrastructure between components,
ensuring that they can interact with each other without relying on external dependencies.

This ensures also that components can be built in any language or OS, so long as it has a compatible network stack.
Changing from a local Component-Based Architecture to a distributed system involves minimal changes to the existing codebase.
In fact, it only requires changing host addresses and ports.
One idea using this architecture framework is
to have a game which is very calculation-heavy but spread it across various computers or servers
and only have the main client handle the UI component.
This way, even the most computationally expensive games can be played on a potato with minimal performance impact.

Using this system, developers can focus on building reusable and maintainable components,
with a focus on separation of concerns and loose coupling.
To prevent wasting time and energy, developers can reuse virtually anything they build,
and each component can be built further independently.

## Architecture Overview
The CBAf framework consists of several components:
1. **Components**: These are the building blocks of the system.
   They can be written in any language or OS,
   and they communicate with each other using a network stack (IPv4 TCP as of now).
2. **Requests**: These are messages sent between components. Routed using the server (this project/component).
3. **Responses**: These are messages sent in response to requests, either by the server directly or forwarded from the target component back to the client component.
4. **Server**: This is the central component that manages the communication between components. The server is responsible for routing requests to the appropriate components and handling responses.

### Example Sequences
1. **Client sends a 'Login' request to the server**
2. **Server forwards the request to the appropriate component (e.g., Authentication component)**
3. **Authentication component receives the 'Login' request and validates the credentials**
4. **Authentication component sends an 'LoginResponse' to the client**
5. **Client receives the 'LoginResponse' and handles the result**

____________

1. **UI component sends a 'LoadGame' request to the server**
2. **Server forwards the request to the appropriate component (e.g., Database Management Component)**
3. **Database Management Component retrieves the game data from the database**
4. **Database Management Component sends a 'LoadGameResponse' to the UI component**
5. **UI component receives the 'LoadGameResponse' and displays the game data**

## JSON Scheme
Each message must have the following fields present or they will be invalid:
- `requester_id` (string, UUIDv4): A unique ID for the component.
- `time_sent` (string, RFC3339Nano): The timestamp when the request was sent.
- `payload` (message payload): The data to be sent with the request.
- `request_id` (string, UUIDv4): A unique identifier for the request. (GUIDv4 recommended)


### Registration

It is crucial that every component registers itself first. This can be done in the following way:
```json
{
  "requester_id": "ea1973db-31e7-4fe4-bd57-e217f246f6a1",
  "payload": {
    "action": "register",
    "args": [
      {
        "type": "string",
        "value": "Common"
      },
      {
        "type": "string",
        "value": "0.1.14"
      }
    ]
  }
}
```
_Example from my Python Common library for logging_

Note that the actual capabilities for the registration payload are generated programmatically because it needs to be formatted as a proper json string:
```python
with open(Path(__file__).parent.joinpath('component_signature.json'), 'r') as f:
    json_string = json.loads(f.read())
    json_string = json.dumps(json_string)
    message["payload"]["args"].append({"type": "list", "value": json_string})
```

```json
[
  {
    "name": "log_debug",
    "signature": {
      "num_of_args": 4,
      "type_of_args": [
        "string",
        "bool",
        "string",
        "string"
      ]
    }
  },
  {
    "name": "log_info",
    "signature": {
      "num_of_args": 4,
      "type_of_args": [
        "string",
        "bool",
        "string",
        "string"
      ]
    }
  },
  {
    "name": "log_warn",
    "signature": {
      "num_of_args": 4,
      "type_of_args": [
        "string",
        "bool",
        "string",
        "string"
      ]
    }
  },
  {
    "name": "log_error",
    "signature": {
      "num_of_args": 4,
      "type_of_args": [
        "string",
        "bool",
        "string",
        "string"
      ]
    }
  },
  {
    "name": "log_critical",
    "signature": {
      "num_of_args": 4,
      "type_of_args": [
        "string",
        "bool",
        "string",
        "string"
      ]
    }
  },
  {
    "name": "shutdown",
    "signature": {
      "num_of_args": 0,
      "type_of_args": []
    }
  }
]
```
_component_registration.json_

As you can see, each capability has a name (the action/method) and a signature.
The signature consists of a number of arguments and a list of their types.
This way, when a request is sent,
the server automatically finds methods that fit the request (action name must be equal and the signature too).
If a suitable target component isn't found, the server sends an error response back to the client.

### Keep Alive

Components should send a 'KeepAlive' request to the server every 5 seconds to inform the server that the connection is still alive.
This is necessary to prevent the server from timing out the connection after a period of inactivity.
Here is an example of a 'KeepAlive' request:
```json
{
  "requester_id": "5b4b25d1-9538-40ce-8147-4d712013543e",
  "time_sent": "2024-11-11T13:35:32.4440646+01:00",
  "payload": {
    "action": "keep_alive",
    "args": null
  },
  "unique_id": "93fbce0d-d801-428c-8426-0c2a6c1bb1a6"
}
```

### Example Request

An example for a `log_debug` request sent from one of my Windows UI components:
```json
{
  "requester": {
    "title": "Windows Frontend",
    "version": "0.0.0",
    "capabilities": []
  },
  "target": {
    "title": "Unknown",
    "version": "0.0.0",
    "capabilities": []
  },
  "time_sent": "2024-11-08T08:17:01.3762364+01:00",
  "payload": {
    "action": "log_debug",
    "args": [
      {
        "type": "string",
        "value": "Calculated square position: (0, 568,75)"
      },
      {
        "type": "bool",
        "value": "False"
      },
      {
        "type": "string",
        "value": ""
      },
      {
        "type": "string",
        "value": "UI.Windows.Debugging"
      }
    ]
  }
}
```
_Note, this was before I added the UUIDs and changed the requester format._

## Changelog
- **10 November 2024:** Version 0.1.0 released
  - Component Architecture Foundation framework introduced
  - Basic server/client communication over IPv4 TCP connection using JSON protocol
  - Proper shutdown sequence: server receives shutdown signal from any component, forwards to all other registered components, waits for them to send the 'unregister' signal, and then shuts itself down.
