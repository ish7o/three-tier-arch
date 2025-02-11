# User Management API ✨

welcome to my smol user management API written in Go~

## what does it do? uwu

this API manages user data through these endpoints:
- `GET /users` - shows all users
- `GET /users/<id>` - peeks at one specific user
- `POST /users` - adds a new user to our database
- `PATCH /users/<id>` - updates user info
- `DELETE /users/<id>` - removes a user (pls don't, they're all precious)

## user model ⭐️

when creating a new user:
```json
{
    "firstName": "string",
    "lastName": "string",
    "birthYear": 42069,
    "group": "user/premium/admin"
}
```

when getting them back:
```json
{
    "id": 1337,
    "firstName": "string",
    "lastName": "string",
    "age": 420,
    "group": "user/premium/admin"
}
```

## running this bby~

1. clone this repo
2. `go run main.go` and watch the magic happen ✨

made with <3 in Go because python is too mainstream
