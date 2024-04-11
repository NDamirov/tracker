CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  login VARCHAR(20) NOT NULL UNIQUE,
  phash VARCHAR(64) NOT NULL,
  name VARCHAR(50),
  surname VARCHAR(50),
  bdate DATE,
  email VARCHAR(100),
  phoneno VARCHAR(15)
);

CREATE TABLE user_tokens (
    login VARCHAR(20) NOT NULL PRIMARY KEY,
    token text
);