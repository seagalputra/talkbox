
# Talkbox

## Database schema

### User

```json
{
  "id": "",
  "first_name": "",
  "last_name": "",
  "username": "",
  "email": "",
  "password": "",
  "created_at": "",
  "updated_at": "",
  "deleted_at": ""
}
```

### Room

```json
{
  "id": "",
  "participants": [
    {
      "id": "",
      "username": "",
      "email": ""
    }
  ], // embeds, should be updated when user update their profile
  "room_type": "private", // private/group
  "created_at": "",
  "updated_at": "",
  "deleted_at": ""
}
```

### Message

```json
{
  "id": "",
  "body": "",
  "attachment": "",
  "user_id": "",
  "room_id": "",
  "created_at": "",
  "updated_at": "",
  "deleted_at": ""
}
```