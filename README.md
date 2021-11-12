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
