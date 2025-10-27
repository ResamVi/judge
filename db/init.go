package db

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func Init(url string) (*Queries, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, "postgres://"+url)
	if err != nil {
		return nil, fmt.Errorf("connecting: %w", err)
	}

	queries := New(conn)

	// Initialize Admin
	encrypted, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost) // TODO: lol
	if err != nil {
		return nil, fmt.Errorf("generating admin password: %w", err)
	}

	err = queries.CreateUser(ctx, CreateUserParams{
		Username: "julienministrator",
		Password: string(encrypted),
		Token:    "abc",
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

	return queries, nil
}
