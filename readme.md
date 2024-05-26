# Palindrome messages Service
An application that manages messages and provides details about those
messages, specifically whether or not a message is a palindrome.
# Architecture
### Architecture Overview

The application architecture consists of three main layers: the server layer, the database layer, and the model layer. This design separates concerns, enhancing maintainability and scalability.

#### 1. Server Layer
Handles HTTP requests and responses, using Gorilla Mux for routing. It includes:
- `MessageService`: Manages endpoints and handlers.
- Handlers: Functions for creating, retrieving, updating, deleting, and listing messages.

#### 2. Database Layer
Manages data storage and retrieval. It includes:
- `Repo`: Methods for saving, getting, updating, deleting, and listing messages.

#### 3. Model Layer
Defines the application's core data structures. It includes:
- `Message`: Represents a message with fields like `ID`, `Content`, and `IsPalindrome`.

### Schema

```
+-----------------------+
|     HTTP Client       |
+-----------+-----------+
            |
            v
+-----------+-----------+
|     Server Layer      |
|                       |
| - MessageService      |
| - Handlers (Create,   |
|   Retrieve, Update,   |
|   Delete, List)       |
+-----------+-----------+
            |
            v
+-----------+-----------+
|     Database Layer    |
|                       |
| - Repo (Save, Get,    |
|   Update, Delete,     |
|   List Messages)      |
+-----------+-----------+
            |
            v
+-----------+-------------+
|     Model Layer         |
|                         |
| - Message (ID,          |
|   Content, IsPalindrome)|
+-------------------------+
```

**Flow:**
- **HTTP Client**: Sends requests to the server.
- **Server Layer**: Processes requests and coordinates with the `MessageService`.
- **Database Layer**: `Repo` interacts with the database.
- **Model Layer**: Defines the `Message` structure used throughout the application.

# Project Setup
## Local setup
#### Requirements
* Go >= 1.22.3
1. install dependencies

```
make init
```

2. run the project

```
make run
```

3. build the binary

```
make build
```

## Docker setup

1. build the docker image

```
make docker-build
```
2. run with docker

```
make docker-run
```
3. run with docker compose

```
make docker-compose-run
```

4. stop docker compose

```
make docker-compose-down
```

## Unit test

* run unit test

```
make test
```

* run test race

```
make test-race
```
* test coverage
```
make test-coverage
```

## e2e test

please run the service before performing the e2e test.

```
make test-e2e
```

## Linting

```
make tool-lint
make lint
```

## Mock
* generate mock

```
make tool-moq
make moq
```

## APIs

| Method | Endpoint       | Description                   |
|--------|----------------|-------------------------------|
| POST   | /messages      | Creates a new message         |
| GET    | /messages      | Retrieves all messages        |
| GET    | /messages/{id} | Retrieves a specific message  |
| PUT    | /messages/{id} | Updates a specific message    |
| DELETE | /messages/{id} | Deletes a specific message    |

### Create Message

This API creates a new message.

#### Request

```json
{
  "content": "A man a plan a canal Panama"
}
```

| Field     | Type   | Description                  | Required |
|-----------|--------|------------------------------|----------|
| `content` | string | The content of the message.  | Required |

#### Response

```json
{
  "id": "e9e750a6-fa9f-4942-9917-53f0c79f4546",
  "content": "A man a plan a canal Panama",
  "is_palindrome": true
}
```

It returns the message id, content, and whether the message is a palindrome.

### Retrieve Messages

This API retrieves all messages.

#### Parameters:

None

#### Response

```json
[
  {
    "id": "e9e750a6-fa9f-4942-9917-53f0c79f4546",
    "content": "A man a plan a canal Panama",
    "is_palindrome": true
  },
  {
    "id": "a4b2e4c7-df4f-4e68-9952-57e1b3c4a6a9",
    "content": "Hello world",
    "is_palindrome": false
  }
]
```

It returns a list of messages with their ids, contents, and palindrome statuses.

### Retrieve a Specific Message

This API retrieves a specific message by its ID.

#### Parameters:

| Parameter | Type   | Description                    | Required |
|-----------|--------|--------------------------------|----------|
| `id`      | string | The ID of the message to fetch | Required |

#### Response

```json
{
  "id": "e9e750a6-fa9f-4942-9917-53f0c79f4546",
  "content": "A man a plan a canal Panama",
  "is_palindrome": true
}
```

It returns the message id, content, and whether the message is a palindrome.

### Update Message

This API updates a specific message by its ID.

#### Request

```json
{
  "content": "Updated content"
}
```

| Field     | Type   | Description                  | Required |
|-----------|--------|------------------------------|----------|
| `content` | string | The updated content of the message.  | Required |

#### Parameters:

| Parameter | Type   | Description                    | Required |
|-----------|--------|--------------------------------|----------|
| `id`      | string | The ID of the message to update | Required |

#### Response

```json
{
  "id": "e9e750a6-fa9f-4942-9917-53f0c79f4546",
  "content": "Updated content",
  "is_palindrome": false
}
```

It returns the updated message id, content, and whether the updated content is a palindrome.

### Delete Message

This API deletes a specific message by its ID.

#### Parameters:

| Parameter | Type   | Description                    | Required |
|-----------|--------|--------------------------------|----------|
| `id`      | string | The ID of the message to delete | Required |

#### Response

Status: 204 No Content

If the message does not exist, it returns:

Status: 404 Not Found

```json
{
  "error": "message not found"
}
```

It returns a status indicating the deletion result, or an error message if the message does not exist.