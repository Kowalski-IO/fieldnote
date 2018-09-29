CREATE TABLE users
(
  id       BYTES PRIMARY KEY DEFAULT uuid_v4(),
  username TEXT NOT NULL UNIQUE,
  hash     TEXT NOT NULL
);

CREATE TABLE notes
(
  id       BYTES PRIMARY KEY DEFAULT uuid_v4(),
  owner_id BYTES   NOT NULL,
  title    TEXT    NOT NULL,
  visible  BOOLEAN NOT NULL  DEFAULT TRUE,
  tags     TEXT [],
  FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE TABLE parts
(
  id       BYTES PRIMARY KEY DEFAULT uuid_v4(),
  note_id  BYTES NOT NULL,
  position INTEGER,
  title    TEXT,
  kind     TEXT,
  content  TEXT,
  FOREIGN KEY (note_id) REFERENCES notes (id)
);
