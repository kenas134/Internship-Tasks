# Get All Tasks

## Description

Retrieves all tasks currently stored in the application's
in-memory database.

This endpoint does not require any request body or parameters.

## Request

**Method:** GET

**Endpoint:** `/tasks`

**Full URL:** `http://localhost:8080/tasks`

## Request Body

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