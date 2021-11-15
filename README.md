# Chat app in Go (2021)

Based on [GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)
## Structure

### Users
There are two types of users - regular users and admins
Regular users can send messages and join rooms
Admins can create room, moderate messages and users

### Rooms
Room is just a space for users to chat

### Messages and commands
All users can send regular text messages or special commands.
Commands should start with special character `/` and command name, such as `rooms` or `quit`.

1. `User1` connect to `WS Service`.

2. `WS Service` create `User1` goroutine

3. Inside `User1` gouroutine all user messages/commands are read, parsed and send to `command` channel.

4. Another `Run` gouroutine inside `WS Service` is reading `command` channel and depending on command do some actions.

    4.1 Command `Join` will add User to room. If room does not exists it will create new one.

    4.2 Command `Rooms` will send back to user a list of existing rooms.

    4.3 Command `Users` will send back to user a list of users in the same room.

    4.4 Command `Quit` will remove user from room.

    4.5 If the are no commands at all the message will be sent to all users in the room.
