package db

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ResamVi/judge/grading"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

func Init(url string, adminPassword string, environment string) (*Queries, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, "postgres://"+url)
	if err != nil {
		return nil, fmt.Errorf("connecting: %w", err)
	}

	queries := New(conn)

	// Initialize Admin
	encrypted, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating admin password: %w", err)
	}

	token := uuid.New().String()
	if environment == "development" {
		token = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	}

	err = queries.CreateUser(ctx, CreateUserParams{
		Username: "julien",
		Password: string(encrypted),
		Token:    token,
		Approved: true,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Initialize exercises
	entries, err := os.ReadDir("exercises")
	if err != nil {
		return nil, fmt.Errorf("read dir of exercises: %w", err)
	}

	for _, e := range entries {
		file, err := os.ReadFile("exercises/" + e.Name() + "/README.md")
		if err != nil {
			return nil, fmt.Errorf("read README.md of folder: %w", err)
		}
		title, _, found := bytes.Cut(file, []byte("\n"))
		if !found {
			return nil, fmt.Errorf("missing title in README.md of " + e.Name())
		}
		title = bytes.TrimPrefix(title, []byte("# "))

		number, _, found := strings.Cut(e.Name(), "-")
		if !found {
			return nil, fmt.Errorf("missing '-' in title of " + e.Name())
		}

		err = queries.CreateExercise(ctx, CreateExerciseParams{
			ID:    e.Name(),
			Title: fmt.Sprintf("Aufgabe %s: %s", number, string(title)),
		})
		if err != nil {
			return nil, fmt.Errorf("create exercise for %s: %w", e.Name(), err)
		}
	}

	if environment == "development" {
		err = initializeTestdata(ctx, queries)
		if err != nil {
			return nil, fmt.Errorf("could not initialize testdata: %w", err)
		}
	}

	return queries, nil
}

func initializeTestdata(ctx context.Context, queries *Queries) error {
	err := queries.CreateUser(ctx, CreateUserParams{
		Username: "lou",
		Password: "$2a$10$aIX0H/Wpntz7VAHJ3rWs1OKlMPVStaG1FZn25hdsvdnLmNq2/SITy",
		Token:    "3aef3950-a3b0-46ae-8fda-5d97d26740c4",
		Approved: true,
	})
	if err != nil {
		return fmt.Errorf("failed creating lou: %w", err)
	}

	err = queries.CreateUser(ctx, CreateUserParams{
		Username: "anna",
		Password: "$2a$10$aIX0H/Wpntz7VAHJ3rWs1OKlMPVStaG1FZn25hdsvdnLmNq2/SITy",
		Token:    "0d93769c-6aa8-4316-b9c4-dd14a61311de",
		Approved: true,
	})
	if err != nil {
		return fmt.Errorf("failed creating anna: %w", err)
	}
	err = queries.CreateSubmission(ctx, CreateSubmissionParams{
		UserID:     2,
		ExerciseID: "01-judge-einrichten",
		Code:       "package main",
		Output:     "Hello World!",
		Evaluation: "❌, Ist noch nicht am funktionieren",
		Solved:     int32(grading.Attempted),
	})
	if err != nil {
		return fmt.Errorf("failed creating lou solving: %w", err)
	}
	err = queries.CreateSubmission(ctx, CreateSubmissionParams{
		UserID:     3,
		ExerciseID: "01-judge-einrichten",
		Code:       "package main\n\nfunc main() {\n\tfmt.Println(\"Hello World!\")\n}",
		Output:     "Hello World!",
		Evaluation: "✅ Ist am funktionieren",
		Solved:     int32(grading.Solved),
	})
	if err != nil {
		return fmt.Errorf("failed creating anna solving: %w", err)
	}

	return nil
}
