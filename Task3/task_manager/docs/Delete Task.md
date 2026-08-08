# Delete Task

## Description

Deletes an existing task from the application's in-memory database using its unique task ID.

The task ID is provided as a path parameter.

## Request

**Method:** DELETE

**Endpoint:** `/tasks/:id`

**Full URL:** `http://localhost:8080/tasks/1`

## Path Parameter

**Parameter:** `id`

The unique ID of the task to delete.

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
    "message": "Task removed"
}
```

## Task Not Found

If the specified task ID does not exist, the API returns:

**Status Code:** `404 Not Found`

### Response Body

```json
{
    "message": "Task not found"
}
```
