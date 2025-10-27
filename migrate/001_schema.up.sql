CREATE TABLE users (
  id        BIGSERIAL   PRIMARY KEY,
  username  text        NOT NULL UNIQUE,
  password  text        NOT NULL,
  token     text        NOT NULL,
  approved  boolean     NOT NULL DEFAULT FALSE
);

CREATE TABLE exercises (
   id        char(24)    PRIMARY KEY,
   title     text        NOT NULL UNIQUE
);

-- INFO: Break second normal form for better reading experience when querying in psql
CREATE TABLE user_solved_exercise (
    user_id     BIGSERIAL   REFERENCES users(id),
    username    text        REFERENCES users(username),
    exercise_id char(24)    REFERENCES exercises(id),
    title       text        REFERENCES exercises(title),
    PRIMARY KEY (user_id, exercise_id)
);

------ Test data ------

INSERT INTO users (username, token, password, approved) VALUES ('admin', 'abc0', '$2a$10$aIX0H/Wpntz7VAHJ3rWs1OKlMPVStaG1FZn25hdsvdnLmNq2/SITy', true);
INSERT INTO users (username, token, password, approved) VALUES ('l0rdpwned', 'abc1', '$2a$10$aIX0H/Wpntz7VAHJ3rWs1OKlMPVStaG1FZn25hdsvdnLmNq2/SITy', true);
INSERT INTO users (username, token, password, approved) VALUES ('mschmidt', 'abc2', '$2a$10$aIX0H/Wpntz7VAHJ3rWs1OKlMPVStaG1FZn25hdsvdnLmNq2/SITy', true);
INSERT INTO users (username, token, password, approved) VALUES ('deeeeeeeeeeeeeez3', 'abc3', '$2a$10$aIX0H/Wpntz7VAHJ3rWs1OKlMPVStaG1FZn25hdsvdnLmNq2/SITy', true);

INSERT INTO exercises(id, title) VALUES ('01-compiler', 'Der Compiler');
INSERT INTO exercises(id, title) VALUES ('02-hello-world', 'Das erste Programm');
INSERT INTO exercises(id, title) VALUES ('04-booleans', 'Lasagne kochen');
INSERT INTO exercises(id, title) VALUES ('03-funktionen', 'Rettungsaktion');
INSERT INTO exercises(id, title) VALUES ('05-mathe', 'Autofabrik');
INSERT INTO exercises(id, title) VALUES ('06-strings', 'Willkommensnachricht');
INSERT INTO exercises(id, title) VALUES ('07-if', 'Autokauf');
INSERT INTO exercises(id, title) VALUES ('08-switch', 'Blackjack');
INSERT INTO exercises(id, title) VALUES ('09-structs', 'Rennfahren');
INSERT INTO exercises(id, title) VALUES ('10-slices', 'Kartentricks');

INSERT INTO user_solved_exercise (user_id, username, exercise_id, title) VALUES (1, 'admin', '01-compiler', 'Der Compiler');
INSERT INTO user_solved_exercise (user_id, username, exercise_id, title) VALUES (2, 'l0rdpwned', '01-compiler', 'Der Compiler');
INSERT INTO user_solved_exercise (user_id, username, exercise_id, title) VALUES (2, 'l0rdpwned', '02-hello-world', 'Das erste Programm');
