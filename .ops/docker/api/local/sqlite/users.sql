CREATE TABLE IF NOT EXISTS users (login VARCHAR, token TEXT, activity INTEGER, UNIQUE(login));
INSERT INTO users VALUES ('sqlite_user', 'test_token', 1);
