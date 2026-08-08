# Update Task

## Description

Updates an existing task in the application's in-memory database using its unique task ID.

The task ID is provided as a path parameter, and the new task information is provided in the request body as JSON.

## Request

**Method:** PUT

**Endpoint:** `/tasks/:id`

**Full URL:** `http://localhost:8080/tasks/1`

## Path Parameter

**Parameter:** `id`

The unique ID of the task to update.

Example:

```text
/tasks/1
```

## Request Body

A JSON object containing the fields to update is required.

**Content-Type:** `application/json`

### Example Request Body

```json
{
    "title": "Clean the entire room",
    "description": "Clean the floor and organize everything",
    "due_date": "2026-08-11T08:00:00Z",
    "status": "Completed"
}
```

## Request Fields

| Field         | Type   | Description                     |
| ------------- | ------ | ------------------------------- |
| `title`       | string | Updated name of the task        |
| `description` | string | Updated description of the task |
| `due_date`    | string | Updated due date and time       |
| `status`      | string | Updated status of the task      |

## Successful Response

**Status Code:** `200 OK`

### Response Body

```json
{
    "message": "Task updated"
}
```

## Task Not Found

If the specified task ID does not exist, the API returns:

**Status Code:** `404 Not Found`

### Response Body

```json
{
    "error": "task not found"
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
