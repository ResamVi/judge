CREATE TABLE users (
  id        BIGSERIAL   PRIMARY KEY,
  username  text        NOT NULL UNIQUE,
  password  text        NOT NULL,
  token     text        NOT NULL,
  approved  boolean     NOT NULL DEFAULT FALSE
);

CREATE TABLE exercises (
   id        text        PRIMARY KEY,
   title     text        NOT NULL UNIQUE
);

CREATE TABLE submissions (
    user_id     BIGSERIAL   REFERENCES users(id),
    exercise_id text        REFERENCES exercises(id),
    code        text        NOT NULL,
    output      text        NOT NULL,
    evaluation  text        NOT NULL,
    solved      int         NOT NULL DEFAULT 0,
    attempts    int         NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, exercise_id)
);