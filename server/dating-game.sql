CREATE DATABSE `test-dating-game`;

CREATE TABLE users (
  id INT64 NOT NULL,
  username STRING(64) NOT NULL,
  matches ARRAY<INT64>,
) PRIMARY KEY(id);

CREATE TABLE chat (
  id INT64 NOT NULL,
  mark STRING(64) NOT NULL,
  sender STRING(64) NOT NULL,
  message STRING(MAX) NOT NULL,
  timestamp_created INT64 NOT NULL,
) PRIMARY KEY(id);

