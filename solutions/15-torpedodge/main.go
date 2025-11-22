package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/rebirth-in-ruins/torpedodge/server/game"
	"github.com/sanity-io/litter"
)

func decideMove(state game.GameStateResponse) string {
	litter.Dump(state)
	return "LEFT"
}

const (
	gameserverURL = "wss://gameserver.resamvi.io/play"
	playerName    = "GolangBot"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	conn, _, err := websocket.Dial(context.Background(), gameserverURL, nil)
	if err != nil {
		return fmt.Errorf("could not dial server: %w (url: %v)", err, gameserverURL)
	}
	defer conn.CloseNow()

	// Send initial JOIN message with your name
	err = wsjson.Write(context.Background(), conn, "JOIN "+playerName+".go")
	if err != nil {
		return fmt.Errorf("could not join server: %w", err)
	}

	i := 0

	for {
		// RECEIVE NEXT STATE
		var state game.GameStateResponse
		err := wsjson.Read(context.Background(), conn, &state)
		if err != nil {
			return fmt.Errorf("could not read from conn: %w", err)
		}

		// Sail in a circle
		action := decideMove(state)
		slog.Info(action)

		// SEND ACTION
		err = wsjson.Write(context.Background(), conn, action)
		if err != nil {
			return fmt.Errorf("could not send next action: %w", err)
		}

		i++
	}
}
