# Websocket for Consultation (Telemedicine) Feature

The documentation is put as a markdown file since Postman does not support documentation for
Websocket request.

## Creating a Room or a Consultation Session

To create a room or a consultation session use the REST API endpoint `POST /v1/chats`.

The request body must contain

```json
{
  "doctor_id": 1,
  "user_id": 2
}
```

As of right now, there is no validation yet whether the `doctor_id` and/or the `user_id` exist
in the database. **So, make sure you give the right values**.

The response body will be something like this

```json
{
  "id": 1,
  "user_id": 2,
  "doctor_id": 1,
  "consultation_session_status_id": 1,
  "created_at": "2024-01-21T13:53:33.629983+07:00",
  "updated_at": "2024-01-21T13:53:33.629983+07:00",
  "messages": []
}
```

## Getting All Rooms or Consultation Sessions

To get all rooms or consultation sessions you are currently having and the ones that have already been ended, use the
REST API endpoint `GET /v1/chats`.

By default, it will give you a result of consultation sessions sorted by date, descending. If for some reason, you want
to see the consultation session
from the oldest ones, you can do so by providing the query parameter `?sort_by=date&sort=asc`

Here are the statuses available for a room:

1. `1`: this represents a chat is in an `Ongoing` state. Users can still send messages in this room and join to the
   websocket.
2. `2`: this represents a chat is in an `Ended` state. Users can no longer send any message in this room. However, any
   alert message is still possible to be shown and updated in users' view. Alerts can be emitted when the doctor gives
   or update a prescription, or when the doctor gives a sick leave form for the user.

Complete query parameters you can put

1. `status`: the value must be oneof `1` or `2`
    1. `1`: represents `Ongoing` chat
    2. `2`: represents `Ended` chat
2. `sort_by`: the value must be `date`. Not date as in `2024-01-01`, but literally `date`.  
   However, even if you don't set this, the endpoint will give you a result of sorted by date, descending, by default.
3. `sort`: the value must be oneof `asc` or `desc`. By default it will be set to `desc`. If no `sort_by` given then any
   value given for `sort` is ignored.

The response should contain

```json
{
  "data": {
    "total_items": 1,
    "total_pages": 1,
    "current_page_total_items": 1,
    "current_page": 1,
    "items": [
      {
        "id": 1,
        "user_id": 6,
        "doctor_id": 5,
        "consultation_session_status_id": 1,
        "created_at": "2024-01-21T13:53:33.629983+07:00",
        "updated_at": "2024-01-21T13:53:33.629983+07:00",
        "consultation_session_status": {
          "name": "Ongoing"
        },
        "user": {
          "user_id": 6,
          "name": "lumban boy",
          "profile_photo": ""
        },
        "doctor": {
          "user_id": 5,
          "name": "dokter wasik",
          "profile_photo": ""
        },
        "messages": [
          {
            "is_typing": false,
            "message_type": 1,
            "message": "This is the last message from a chat room or a consultation session",
            "attachment": "",
            "created_at": "2024-01-21T13:54:19.519447+07:00",
            "sender_id": 5,
            "session_id": 1
          }
        ]
      }
    ]
  }
}
```

## Get Room by Id

To get a room by its Id you can use the REST API endpoint `GET /v1/chats/:id`

The response should contain the details of the room including the user's profile and the doctor's profile. It should
also include all the messages that happened in the room. The question is whether to retrieve the messages by pagination
or give all of them as a whole.

This is how the response looks like

```json
{
  "data": {
    "id": 1,
    "user_id": 6,
    "doctor_id": 5,
    "consultation_session_status_id": 1,
    "created_at": "2024-01-21T13:53:33.629983+07:00",
    "updated_at": "2024-01-21T13:53:33.629983+07:00",
    "consultation_session_status": {
      "name": "Ongoing"
    },
    "user": {
      "user_id": 6,
      "name": "lumban boy",
      "profile_photo": ""
    },
    "doctor": {
      "user_id": 5,
      "name": "dokter wasik",
      "profile_photo": ""
    },
    "messages": [
      {
        "is_typing": false,
        "message_type": 1,
        "message": "This is the last message from a chat room or a consultation session",
        "attachment": "",
        "created_at": "2024-01-21T13:54:19.519447+07:00",
        "sender_id": 5,
        "session_id": 1
      },
      {
        "is_typing": false,
        "message_type": 1,
        "message": "This is the new last message",
        "attachment": "",
        "created_at": "2024-01-21T14:27:20.454924+07:00",
        "sender_id": 6,
        "session_id": 1
      }
    ]
  }
}
```

## Joining a Room

To be able to join a room, the room id to be joined must be first using the Get Room by Id endpoint.

If the status is `2` or `Ended`, then the doctor or the user trying to do this action should not be allowed.
You will get `400 Bad Request` in the response header.

If the status is `1` or `Ongoing`, then the doctor or the user is allowed to join and send messages.

The endpoint to join the room is the following where it is using the `websocket` or `ws` protocol
`ws://{{ADDRESS}}/v1/chats/:id/join?token={{YOUR_TOKEN}}`

## Sending Message in The Room
This is done after successfully joining a room.

Messages you send in the room must be in the format of `JSON`.

Currently, here are the `key-value` pairs recognized

```json
{
  "is_typing": true,
  "message_type": 1,
  "message": "Any message in string",
  "attachment": "data:image/jpg;base64,base64+encoded+image+or+pdf"
}
```

1. `is_typing`: `boolean`
2. `message_type`: one of `1` or `2` in `integer` type
    1. `1`: represents that the message will be a regular message sent by user's typing
    2. `2`: represents that the message will be an alert message sent by user's action (specifically doctor's actions
       such as creating prescription).
3. `message`: `string`
4. `attachment`: `string` containing `base64` encoded `string` of an `image` or `pdf`. Use the
   [`dataurl` standard](https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/Data_URLs).
5. `sender_id`: `integer` use the current user sending the message here. It could be the doctor's or the user's
   `user_id`.
6. `session_id`: `integer` use the `id` of the object you get when getting the room by id before joining the room.


## Ending a Room

To end a consultation session or room, you can use the endpoint `PUT /v1/chats/:id`.
**This can only be done one way, from `Ongoing` or `1` to `Ended` or `2`**.

If a room's status is already at `Ended`, then it will return an error saying `chat already ended`.
