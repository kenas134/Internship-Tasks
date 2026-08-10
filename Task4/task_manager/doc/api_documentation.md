# Task Management REST API

## Overview

The Task Management REST API is a RESTful web service built with **Go** and the **Gin** web framework.

The API provides CRUD operations for managing tasks:

* Create a task
* Retrieve all tasks
* Retrieve a task by ID
* Update a task
* Delete a task

The application currently uses an **in-memory data store**, meaning that task data is stored only while the application is running. Any changes will be lost when the application is restarted.

---

## Base URL

```text
http://localhost:8080
```

---

## API Endpoints

| Method   | Endpoint     | Description              |
| -------- | ------------ | ------------------------ |
| `GET`    | `/tasks`     | Retrieve all tasks       |
| `GET`    | `/tasks/:id` | Retrieve a specific task |
| `POST`   | `/tasks`     | Create a new task        |
| `PUT`    | `/tasks/:id` | Update an existing task  |
| `DELETE` | `/tasks/:id` | Delete a task            |

---

# Task Object

A task has the following properties:

| Field         | Type   | Description                        |
| ------------- | ------ | ---------------------------------- |
| `id`          | string | Unique identifier of the task      |
| `title`       | string | Name of the task                   |
| `description` | string | Detailed description of the task   |
| `due_date`    | string | Date and time when the task is due |
| `status`      | string | Current status of the task         |

### Example Task

```json
{
    "id": "1",
    "title": "Clean the room",
    "description": "Organize books, clothes, and clean the floor",
    "due_date": "2026-08-10T08:00:00Z",
    "status": "Pending"
}
```

The `due_date` uses the standard ISO 8601/RFC 3339 date-time format.

---

# 1. Get All Tasks

Retrieves all tasks currently stored in the in-memory database.

## Request

**Method:** `GET`

**Endpoint:** `/tasks`

**Full URL:**

```text
http://localhost:8080/tasks
```

### Request Body

No request body is required.

## Successful Response

**Status Code:** `200 OK`

### Response Body

```json
{
    "tasks": [
        {
            "id": "1",
            "title": "Clean the room",
            "description": "Organize books, clothes, and clean the floor",
            "due_date": "2026-08-10T08:00:00Z",
            "status": "Pending"
        },
        {
            "id": "2",
            "title": "Wash dishes",
            "description": "Clean all dishes",
            "due_date": "2026-08-10T19:00:00Z",
            "status": "Completed"
        }
    ]
}
```

---

# 2. Get Task By ID

Retrieves a specific task using its unique task ID.

## Request

**Method:** `GET`

**Endpoint:** `/tasks/:id`

**Example URL:**

```text
http://localhost:8080/tasks/1
```

## Path Parameter

| Parameter | Type   | Description           |
| --------- | ------ | --------------------- |
| `id`      | string | Unique ID of the task |

### Example

```text
GET http://localhost:8080/tasks/1
```

## Request Body

No request body is required.

## Successful Response

**Status Code:** `200 OK`

### Response Body

```json
{
    "id": "1",
    "title": "Clean the room",
    "description": "Organize books, clothes, and clean the floor",
    "due_date": "2026-08-10T08:00:00Z",
    "status": "Pending"
}
```

## Task Not Found

If the specified task does not exist:

**Status Code:** `404 Not Found`

### Response Body

```json
{
    "error": "task not found"
}
```

---

# 3. Create Task

Creates a new task and stores it in the in-memory database.

## Request

**Method:** `POST`

**Endpoint:** `/tasks`

**Full URL:**

```text
http://localhost:8080/tasks
```

## Request Headers

```text
Content-Type: application/json
```

## Request Body

The request body must contain the task information in JSON format.

### Example Request

```json
{
    "id": "6",
    "title": "Take out the trash",
    "description": "Take all trash bags outside",
    "due_date": "2026-08-18T08:00:00Z",
    "status": "Pending"
}
```

## Request Fields

| Field         | Type   | Description                        |
| ------------- | ------ | ---------------------------------- |
| `id`          | string | Unique ID of the task              |
| `title`       | string | Name of the task                   |
| `description` | string | Description of the task            |
| `due_date`    | string | Date and time when the task is due |
| `status`      | string | Current status of the task         |

## Successful Response

**Status Code:** `201 Created`

### Response Body

```json
{
    "message": "Task created"
}
```

## Invalid Request

If the request contains invalid JSON:

**Status Code:** `400 Bad Request`

### Response Body

The response contains an error message describing the invalid request.

Example:

```json
{
    "error": "invalid request body"
}
```

---

# 4. Update Task

Updates an existing task using its unique task ID.

## Request

**Method:** `PUT`

**Endpoint:** `/tasks/:id`

**Example URL:**

```text
http://localhost:8080/tasks/1
```

## Path Parameter

| Parameter | Type   | Description                     |
| --------- | ------ | ------------------------------- |
| `id`      | string | Unique ID of the task to update |

## Request Headers

```text
Content-Type: application/json
```

## Request Body

The updated task information should be provided as JSON.

### Example Request

```json
{
    "title": "Clean the entire room",
    "description": "Clean the floor and organize everything",
    "due_date": "2026-08-11T08:00:00Z",
    "status": "Completed"
}
```

## Successful Response

**Status Code:** `200 OK`

### Response Body

```json
{
    "message": "Task updated"
}
```

## Task Not Found

If the specified task does not exist:

**Status Code:** `404 Not Found`

### Response Body

```json
{
    "error": "task not found"
}
```

## Invalid Request

If the request contains invalid JSON:

**Status Code:** `400 Bad Request`

---

# 5. Delete Task

Deletes an existing task using its unique task ID.

## Request

**Method:** `DELETE`

**Endpoint:** `/tasks/:id`

**Example URL:**

```text
http://localhost:8080/tasks/1
```

## Path Parameter

| Parameter | Type   | Description                     |
| --------- | ------ | ------------------------------- |
| `id`      | string | Unique ID of the task to delete |

## Request Body

No request body is required.

## Successful Response

**Status Code:** `200 OK`

### Response Body

```json
{
    "message": "Task removed"
}
```

## Task Not Found

If the specified task does not exist:

**Status Code:** `404 Not Found`

### Response Body

```json
{
    "message": "Task not found"
}
```

---

# HTTP Status Codes

The API uses standard HTTP status codes to indicate the result of each request.

| Status Code       | Meaning                                   |
| ----------------- | ----------------------------------------- |
| `200 OK`          | Request completed successfully            |
| `201 Created`     | A new task was successfully created       |
| `400 Bad Request` | The request contains invalid data or JSON |
| `404 Not Found`   | The requested task does not exist         |

---

# Error Handling

Errors are returned as JSON responses.

### Example: Task Not Found

```json
{
    "error": "task not found"
}
```

### Example: Invalid Request

```json
{
    "error": "invalid request body"
}
```

The exact error message may vary depending on the error returned by the application.

---

# Postman Testing

The API can be tested using **Postman**.

Create a Postman collection named:

```text
Task Management API
```

Add the following requests:

```text
Task Management API
│
├── GET All Tasks
├── GET Task By ID
├── POST Create Task
├── PUT Update Task
└── DELETE Task
```

## Recommended Test Sequence

### 1. Retrieve all tasks

```text
GET http://localhost:8080/tasks
```

Expected status:

```text
200 OK
```

### 2. Retrieve a task

```text
GET http://localhost:8080/tasks/1
```

Expected status:

```text
200 OK
```

### 3. Create a task

```text
POST http://localhost:8080/tasks
```

Send a JSON request body containing the new task.

Expected status:

```text
201 Created
```

### 4. Update the task

```text
PUT http://localhost:8080/tasks/6
```

Send the updated task information as JSON.

Expected status:

```text
200 OK
```

### 5. Delete the task

```text
DELETE http://localhost:8080/tasks/6
```

Expected status:

```text
200 OK
```

### 6. Verify deletion

```text
GET http://localhost:8080/tasks/6
```

Expected status:

```text
404 Not Found
```

---

# Running the API

## Prerequisites

Make sure Go is installed on your system.

## Start the Application

From the project directory, run:

```bash
go run main.go
```

The server will start at:

```text
http://localhost:8080
```

The API can then be tested using Postman or another HTTP client.

---

# Data Storage

This project uses an in-memory data store instead of a persistent database.

Tasks are stored in a Go slice while the application is running.

For example:

```go
var tasks = []models.Task{
    {
        ID:     "1",
        Title:  "Clean the room",
        Status: "Pending",
    },
    {
        ID:     "2",
        Title:  "Wash dishes",
        Status: "Completed",
    },
}
```

Because the data is stored in memory, all changes made through the API are lost when the application is stopped or restarted.

---

# Project Structure

```text
task_manager/
│
├── main.go
│
├── controllers/
│   └── task_controller.go
│
├── models/
│   └── task.go
│
├── data/
│   └── task_service.go
│
├── router/
│   └── router.go
│
├── docs/
│   └── api_documentation.md
│
└── go.mod
```

## Component Responsibilities

| Component      | Responsibility                         |
| -------------- | -------------------------------------- |
| `main.go`      | Entry point of the application         |
| `router/`      | Defines API routes                     |
| `controllers/` | Handles HTTP requests and responses    |
| `models/`      | Defines the Task structure             |
| `data/`        | Manages task data and business logic   |
| `docs/`        | Contains API documentation             |
| `go.mod`       | Defines the Go module and dependencies |

---

# Architecture

The application follows a layered architecture:

```text
Client / Postman
       │
       ▼
    Router
       │
       ▼
  Controller
       │
       ▼
  Data Service
       │
       ▼
 In-Memory Data
```

The router receives the request and directs it to the appropriate controller.

The controller processes the HTTP request and calls the appropriate data service function.

The data service performs the required operation on the in-memory task data.

The controller then sends the appropriate JSON response and HTTP status code back to the client.

---

