# Create Task

## Description

Creates a new task and stores it in the application's in-memory database.

The task information must be provided in the request body as JSON.

## Request

**Method:** POST

**Endpoint:** `/tasks`

**Full URL:** `http://localhost:8080/tasks`

## Request Body

A JSON object containing the information for the new task is required.

**Content-Type:** `application/json`

### Example Request Body

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

If the request body contains invalid JSON, the API returns:

**Status Code:** `400 Bad Request`

### Response Body

```json
{
    "error": "invalid request body"
}
```

The actual error message may vary depending on the invalid JSON that was provided.
