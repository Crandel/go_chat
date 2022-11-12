PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS users(
  nick         VARCHAR(50) UNIQUE PRIMARY KEY,
  name         TEXT,
  second_name  TEXT,
  email        TEXT,
  password     TEXT NOT NULL,
  token        TEXT NOT NULL,
  role         TEXT CHECK( role IN ('Admin', 'Member') ) NOT NULL DEFAULT 'Member',
  created      DATETIME DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS rooms (
  name         VARCHAR(50) UNIQUE PRIMARY KEY,
  created      DATETIME DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS user_rooms (
  id           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  user_nick    TEXT NOT NULL UNIQUE,
  room_name    TEXT NOT NULL,
  FOREIGN KEY(user_nick) REFERENCES users(nick),
  FOREIGN KEY(room_name) REFERENCES rooms(name)
);

CREATE TABLE IF NOT EXISTS messages (
  id           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  user_room_id INTEGER NOT NULL,
  payload      TEXT,
  created      DATETIME DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  FOREIGN KEY(user_room_id) REFERENCES user_rooms(id)
);
