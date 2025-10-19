package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func Init(url string) (*Queries, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://"+url)
	if err != nil {
		return nil, fmt.Errorf("connecting: %w", err)
	}
	defer conn.Close(ctx)

	queries := New(conn)

	encrypted, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost) // TODO: lol
	if err != nil {
		return nil, fmt.Errorf("generating admin password: %w", err)
	}

	err = queries.UpsertUser(ctx, UpsertUserParams{
		Username: "admin",
		Password: string(encrypted),
		Approved: true,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert user: %w", err)
	}

	return queries, nil
}
