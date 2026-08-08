# Task Management REST API

## Overview

The Task Management REST API is a RESTful web service developed using **Go** and the **Gin framework**.

The API allows users to create, retrieve, update, and delete tasks.

The application currently uses an **in-memory database** to store tasks. This means that task data is stored only while the application is running. All changes are lost when the application is restarted.

## Base URL

```text
http://localhost:8080
```

## API Endpoints

| Method | Endpoint     | Description               |
| ------ | ------------ | ------------------------- |
| GET    | `/tasks`     | Retrieves all tasks       |
| GET    | `/tasks/:id` | Retrieves a specific task |
| POST   | `/tasks`     | Creates a new task        |
| PUT    | `/tasks/:id` | Updates an existing task  |
| DELETE | `/tasks/:id` | Deletes an existing task  |

## Task Object

A task contains the following fields:

| Field         | Type   | Description                        |
| ------------- | ------ | ---------------------------------- |
| `id`          | string | Unique identifier of the task      |
| `title`       | string | Name of the task                   |
| `description` | string | Description of the task            |
| `due_date`    | string | Date and time when the task is due |
| `status`      | string | Current status of the task         |

## Example Task

```json
{
    "id": "1",
    "title": "Clean the room",
    "description": "Organize books, clothes, and clean the floor",
    "due_date": "2026-08-10T08:00:00Z",
    "status": "Pending"
}
```

## Content Type

Requests containing a JSON body must use:

```text
Content-Type: application/json
```

This applies to:

* `POST /tasks`
* `PUT /tasks/:id`

GET and DELETE requests do not require a request body.

## HTTP Status Codes

The API uses standard HTTP status codes.

| Status Code       | Meaning                                           |
| ----------------- | ------------------------------------------------- |
| `200 OK`          | Request completed successfully                    |
| `201 Created`     | A new task was successfully created               |
| `400 Bad Request` | The request contains invalid data or invalid JSON |
| `404 Not Found`   | The requested task does not exist                 |

## Error Response

When an error occurs, the API returns a JSON response containing an error message.

Example:

```json
{
    "error": "task not found"
}
```

## Example Workflow

A typical interaction with the API can be performed in the following order.

### 1. Retrieve All Tasks

```http
GET http://localhost:8080/tasks
```

### 2. Retrieve a Specific Task

```http
GET http://localhost:8080/tasks/1
```

### 3. Create a New Task

```http
POST http://localhost:8080/tasks
```

Request body:

```json
{
    "id": "6",
    "title": "Take out the trash",
    "description": "Take all trash bags outside",
    "due_date": "2026-08-18T08:00:00Z",
    "status": "Pending"
}
```

### 4. Update the Task

```http
PUT http://localhost:8080/tasks/6
```

### 5. Delete the Task

```http
DELETE http://localhost:8080/tasks/6
```

## Data Storage

The API uses an **in-memory slice** as its data store.

Tasks are stored while the Go application is running.

The API does not currently use a persistent database such as **MySQL**, **PostgreSQL**, or **MongoDB**.

When the application is stopped or restarted, changes made through the API are lost and the initial task data is restored.

## API Architecture

The application is organized into separate layers:

```text
Client / Postman
       |
       v
    Router
       |
       v
  Controller
       |
       v
 Data Service
       |
       v
In-Memory Task Data
       |
       v
   Task Model
```

### Router

Defines the API routes and determines which controller handles each incoming request.

### Controller

Handles HTTP requests, reads path parameters and JSON request bodies, calls the data service, and sends HTTP responses.

### Data Service

Contains the logic for retrieving, creating, updating, and deleting tasks from the in-memory data store.

### Model

Defines the structure of a task and its fields.

## Testing

The API can be tested using **Postman**.

The Postman collection contains requests for all supported CRUD operations:

* Get all tasks
* Get a task by ID
* Create a task
* Update a task
* Delete a task

Both successful requests and error scenarios should be tested.
