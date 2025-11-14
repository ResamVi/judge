-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserFromId :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserFromToken :one
SELECT * FROM users
WHERE token = $1;

-- name: GetUsers :many
SELECT * FROM users
WHERE approved = true
ORDER BY id ASC;

-- name: GetExercises :many
SELECT * FROM exercises
ORDER BY id ASC;

-- name: GetStatus :many
SELECT
    u.id AS user_id,
    COALESCE(s.solved, 0) AS solved
FROM users u
    LEFT JOIN submissions s
    ON s.user_id = u.id
    AND s.exercise_id = $1
ORDER BY u.id ASC;

-- name: GetSubmission :one
SELECT * FROM submissions
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
    user_id, exercise_id, code, output, evaluation, solved
) VALUES (
     $1, $2, $3, $4, $5, $6
) ON CONFLICT (user_id, exercise_id) DO UPDATE
SET
    code = EXCLUDED.code,
    output = EXCLUDED.output,
    evaluation = EXCLUDED.evaluation,
    solved = EXCLUDED.solved;