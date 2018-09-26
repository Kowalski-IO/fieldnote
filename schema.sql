CREATE TABLE users
(
  id INTEGER PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  hash TEXT NOT NULL
);


CREATE TABLE notes
(
  id INTEGER PRIMARY KEY,
  owner_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  visible INTEGER NOT NULL,
  FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE parts
(
  id INTEGER PRIMARY KEY,
  parent_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  kind TEXT NOT NULL,
  content TEXT,
  FOREIGN KEY (parent_id) REFERENCES notes(id)
);
