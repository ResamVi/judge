package db

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func Init(url string, adminPassword string) (*Queries, error) {
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

	err = queries.CreateUser(ctx, CreateUserParams{
		Username: "julienministrator",
		Password: string(encrypted),
		Token:    uuid.New().String(),
		Approved: true,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Initialize exercises
	entries, err := os.ReadDir("tasks")
	if err != nil {
		return nil, fmt.Errorf("read dir of tasks: %w", err)
	}

	for _, e := range entries {
		file, err := os.ReadFile("tasks/" + e.Name() + "/README.md")
		if err != nil {
			return nil, fmt.Errorf("read README.md of folder: %w", err)
		}
		title, _, found := bytes.Cut(file, []byte("\n"))
		if !found {
			return nil, fmt.Errorf("missing title in README.md of " + e.Name())
		}
		title = bytes.TrimPrefix(title, []byte("# "))

		err = queries.CreateExercise(ctx, CreateExerciseParams{
			ID:    e.Name(),
			Title: string(title),
		})
		if err != nil {
			return nil, fmt.Errorf("create exercise for %s: %w", e.Name(), err)
		}
	}

	// Testdaten
	err = testData(ctx, queries)
	if err != nil {
		return nil, fmt.Errorf("could not initialize testData: %w", err)
	}

	return queries, nil
}

func testData(ctx context.Context, queries *Queries) error {
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
	err = queries.UserSolvedExercise(ctx, UserSolvedExerciseParams{
		UserID: 2,
		Username: pgtype.Text{
			String: "lou",
			Valid:  true,
		},
		ExerciseID: "01-compiler",
		Title: pgtype.Text{
			String: "Der Compiler",
			Valid:  true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed creating lou solving: %w", err)
	}
	err = queries.UserSolvedExercise(ctx, UserSolvedExerciseParams{
		UserID: 3,
		Username: pgtype.Text{
			String: "anna",
			Valid:  true,
		},
		ExerciseID: "02-hello-world",
		Title: pgtype.Text{
			String: "Das erste Programm",
			Valid:  true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed creating anna solving: %w", err)
	}

	return nil
}
