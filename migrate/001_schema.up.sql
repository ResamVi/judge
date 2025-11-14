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

-- INFO: Break second normal form for better reading experience when querying in psql
CREATE TABLE solved (
    user_id     BIGSERIAL   NOT NULL    REFERENCES users(id),
    username    text        NOT NULL    REFERENCES users(username),
    exercise_id text        NOT NULL    REFERENCES exercises(id),
    title       text        NOT NULL    REFERENCES exercises(title),
    PRIMARY KEY (user_id, exercise_id)
);

CREATE TABLE submissions (
    user_id     BIGSERIAL   REFERENCES users(id),
    exercise_id text        REFERENCES exercises(id),
    code        text        NOT NULL,
    output      text        NOT NULL,
    evaluation  text        NOT NULL,
    solved      boolean     NOT NULL,
    PRIMARY KEY (user_id, exercise_id)
);