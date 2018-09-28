CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA public;


CREATE TABLE users
(
  id    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email TEXT NOT NULL UNIQUE,
  hash  TEXT NOT NULL
);

CREATE TABLE notes
(
  id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  owner_id UUID    NOT NULL,
  title    TEXT    NOT NULL,
  visible  BOOLEAN NOT NULL DEFAULT TRUE,
  tags     TEXT [],
  FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE TABLE parts
(
  id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  note_id  UUID NOT NULL,
  position INTEGER,
  title    TEXT,
  kind     TEXT,
  content  TEXT,
  FOREIGN KEY (note_id) REFERENCES notes (id)
);
