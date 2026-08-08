# Get Task By ID

## Description

Retrieves a specific task from the application's in-memory database using its unique task ID.

The task ID is provided as a path parameter.

## Request

**Method:** GET

**Endpoint:** `/tasks/:id`

**Full URL:** `http://localhost:8080/tasks/1`

## Path Parameter

**Parameter:** `id`

The unique ID of the task to retrieve.

Example:

```text
/tasks/1
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

If the specified task ID does not exist, the API returns:

**Status Code:** `404 Not Found`

### Response Body

```json
{
    "error": "task not found"
}
```
