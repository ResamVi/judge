-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserFromToken :one
SELECT * FROM users
WHERE token = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
WHERE approved = true
ORDER BY id ASC;

-- name: GetExercises :many
SELECT * FROM exercises
ORDER BY id ASC;

-- name: GetSolvers :many
SELECT
    u.id,
    u.username,
    (usex.user_id IS NOT NULL)::boolean AS solved
FROM users u
LEFT JOIN user_solved_exercise usex ON usex.user_id = u.id
AND usex.exercise_id = $1
ORDER BY u.id;


-- name: GetCode :one
SELECT code FROM submissions
WHERE user_id = $1 AND exercise_id = $2;

-- name: CreateUser :exec
INSERT INTO users (
  username, password, token, approved
) VALUES (
  $1, $2, $3, $4
) ON CONFLICT DO NOTHING;

-- name: CreateExercise :exec
INSERT INTO exercises (
    id, title
) VALUES (
  $1, $2
) ON CONFLICT DO NOTHING;

-- name: CreateSubmission :exec
INSERT INTO submissions (
    user_id, exercise_id, code
) VALUES (
     $1, $2, $3
) ON CONFLICT (user_id, exercise_id) DO UPDATE
SET code = EXCLUDED.code;

-- name: UserSolvedExercise :exec
INSERT INTO user_solved_exercise (
    user_id, username, exercise_id, title
) VALUES (
    $1, $2, $3, $4
);
