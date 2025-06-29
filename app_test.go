package main

import (
	"reflect"
	"testing"
)

func TestNewGame(t *testing.T) {
	app := NewApp() // In a real scenario, startup would be called.
	// For unit tests, explicitly initialize or rely on NewApp's behavior if it sets up gameState.
	// Based on app.go, NewApp just returns &App{}. startup initializes gameState.
	// So, we should call NewGame() to get a game state.
	// However, the App methods like NewGame, MakeMove, ResetGame operate on app.gameState.
	// So, app.gameState needs to be set.

	// The methods like NewGame(), MakeMove(), ResetGame() are methods of App struct.
	// NewGame() itself sets and returns app.gameState.
	gameState := app.NewGame()

	if gameState.CurrentPlayer != "X" {
		t.Errorf("Expected CurrentPlayer to be 'X', got '%s'", gameState.CurrentPlayer)
	}
	if gameState.GameOver != false {
		t.Errorf("Expected GameOver to be false, got '%t'", gameState.GameOver)
	}
	if gameState.Winner != "" {
		t.Errorf("Expected Winner to be '', got '%s'", gameState.Winner)
	}
	expectedBoard := [9]string{}
	if !reflect.DeepEqual(gameState.Board, expectedBoard) {
		t.Errorf("Expected empty board, got '%v'", gameState.Board)
	}

	// Also check if app.gameState is what NewGame returned
	if !reflect.DeepEqual(app.gameState, gameState) {
		t.Errorf("app.gameState was not updated by NewGame()")
	}
}

func TestMakeMove(t *testing.T) {
	app := NewApp()

	t.Run("Valid moves", func(t *testing.T) {
		app.gameState = app.NewGame() // Ensure fresh state

		gs, err := app.MakeMove(0)
		if err != nil {
			t.Fatalf("MakeMove(0) returned error: %v", err)
		}
		if gs.Board[0] != "X" {
			t.Errorf("Expected Board[0] to be 'X', got '%s'", gs.Board[0])
		}
		if gs.CurrentPlayer != "O" {
			t.Errorf("Expected CurrentPlayer to be 'O', got '%s'", gs.CurrentPlayer)
		}
		if !reflect.DeepEqual(app.gameState, gs) {
			t.Errorf("app.gameState was not updated by MakeMove")
		}

		gs, err = app.MakeMove(1)
		if err != nil {
			t.Fatalf("MakeMove(1) returned error: %v", err)
		}
		if gs.Board[1] != "O" {
			t.Errorf("Expected Board[1] to be 'O', got '%s'", gs.Board[1])
		}
		if gs.CurrentPlayer != "X" {
			t.Errorf("Expected CurrentPlayer to be 'X', got '%s'", gs.CurrentPlayer)
		}
	})

	t.Run("Cell already occupied", func(t *testing.T) {
		app.gameState = app.NewGame()
		_, err := app.MakeMove(0) // X moves at 0
		if err != nil {
			t.Fatalf("Initial MakeMove(0) failed: %v", err)
		}

		originalGameState := *app.gameState // Capture state before invalid move

		gs, err := app.MakeMove(0) // O tries to move to the same cell 0
		if err == nil {
			t.Fatalf("Expected error for occupied cell, got nil")
		}
		expectedError := "Cell already occupied"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
		}

		// Check that game state was not altered by the failed move
		if gs.Board[0] != "X" { // Board should remain as X
			t.Errorf("Expected Board[0] to still be 'X', got '%s'", gs.Board[0])
		}
		if gs.CurrentPlayer != "O" { // Current player should still be O (who attempted the invalid move)
			t.Errorf("Expected CurrentPlayer to still be 'O', got '%s'", gs.CurrentPlayer)
		}
		if gs.Winner != "" { // No winner should be set
			t.Errorf("Expected Winner to be empty, got '%s'", gs.Winner)
		}
		if gs.GameOver != false { // Game should not be over
			t.Errorf("Expected GameOver to be false, got '%t'", gs.GameOver)
		}
		// Check if the returned state `gs` from the failed MakeMove is the same as before the call
		if !reflect.DeepEqual(gs, &originalGameState) {
			t.Errorf("GameState was unexpectedly altered by failed MakeMove. Got: %v, Expected: %v", gs, originalGameState)
		}

	})

	t.Run("Game is over", func(t *testing.T) {
		app.gameState = app.NewGame()
		app.gameState.GameOver = true // Manually set game to over

		_, err := app.MakeMove(0) // Try to move on an empty cell
		if err == nil {
			t.Fatalf("Expected error for game over, got nil")
		}
		expectedError := "Game is over"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("Invalid move index", func(t *testing.T) {
		app.gameState = app.NewGame()

		_, err := app.MakeMove(9) // Index out of bounds
		if err == nil {
			t.Fatalf("Expected error for index 9, got nil")
		}
		expectedError := "Invalid move index"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s' for index 9, got '%s'", expectedError, err.Error())
		}

		_, err = app.MakeMove(-1) // Index out of bounds
		if err == nil {
			t.Fatalf("Expected error for index -1, got nil")
		}
		if err.Error() != expectedError { // Should be the same error message
			t.Errorf("Expected error '%s' for index -1, got '%s'", expectedError, err.Error())
		}
	})
}

func TestWinConditions(t *testing.T) {
	app := NewApp()

	// Helper to run a sequence of moves and test win
	testWinScenario := func(t *testing.T, player string, moves []int, winningMove int, description string) {
		t.Run(description, func(t *testing.T) {
			app.gameState = app.NewGame() // Fresh board

			// Apply initial moves
			for i, move := range moves {
				// Determine current player for this move based on i
				// This setup assumes alternating players starting with X
				// If player for scenario is "O", first move is by X, then O, etc.
				// This part might need adjustment if a specific player needs to make all setup moves
				// For now, let's assume the `moves` array represents sequential plays by alternating players
				_, err := app.MakeMove(move)
				if err != nil {
					t.Fatalf("Setup move %d (%d) failed: %v", i+1, move, err)
				}
			}

			// At this point, app.gameState.CurrentPlayer should be the one to make the winning move.
			// We need to ensure this is the 'player' parameter.
			// This requires careful setup of the 'moves' array.
			// For simplicity, let's ensure CurrentPlayer is set correctly before the winning move.
			app.gameState.CurrentPlayer = player

			gs, err := app.MakeMove(winningMove)
			if err != nil {
				t.Fatalf("Winning MakeMove(%d) for %s failed: %v", winningMove, player, err)
			}
			if !gs.GameOver {
				t.Errorf("Expected GameOver to be true for %s win (%s)", player, description)
			}
			if gs.Winner != player {
				t.Errorf("Expected Winner to be '%s', got '%s' (%s)", player, gs.Winner, description)
			}
		})
	}

	// --- Test cases for Player X ---
	// Row wins
	testWinScenario(t, "X", []int{0, 3, 1, 4}, 2, "X row 1 (0,1,2)") // X:0,1,2 O:3,4
	testWinScenario(t, "X", []int{3, 0, 4, 1}, 5, "X row 2 (3,4,5)") // X:3,4,5 O:0,1
	testWinScenario(t, "X", []int{6, 0, 7, 1}, 8, "X row 3 (6,7,8)") // X:6,7,8 O:0,1

	// Column wins
	testWinScenario(t, "X", []int{0, 1, 3, 2}, 6, "X col 1 (0,3,6)") // X:0,3,6 O:1,2
	testWinScenario(t, "X", []int{1, 0, 4, 2}, 7, "X col 2 (1,4,7)") // X:1,4,7 O:0,2
	testWinScenario(t, "X", []int{2, 0, 5, 1}, 8, "X col 3 (2,5,8)") // X:2,5,8 O:0,1

	// Diagonal wins
	testWinScenario(t, "X", []int{0, 1, 4, 2}, 8, "X diag 1 (0,4,8)") // X:0,4,8 O:1,2
	testWinScenario(t, "X", []int{2, 1, 4, 0}, 6, "X diag 2 (2,4,6)") // X:2,4,6 O:1,0

	// --- Test cases for Player O ---
	// To test for O, X must make the first move.
	// Row wins for O
	// X O .
	// X O .
	// . O . (winning O move)
	testWinScenario(t, "O", []int{0, 3, 1, 4, 2}, 5, "O row 2 (3,4,5)") // X:0,1,2 O:3,4,5
	// Example: O wins on top row (0,1,2)
	// X plays 3, O plays 0
	// X plays 4, O plays 1
	// X plays 6 (arbitrary other cell), O plays 2 (winning)
	testWinScenario(t, "O", []int{3, 0, 4, 1, 6}, 2, "O row 1 (0,1,2)")

	// Column wins for O
	// X O .
	// X O .
	// . O . (winning O move for col 2: 1,4,7)
	testWinScenario(t, "O", []int{0, 1, 2, 4, 3}, 7, "O col 2 (1,4,7)") // X:0,2,3 O:1,4,7

	// Diagonal wins for O
	// X . O
	// . O .
	// O X X (winning O move for diag 1: 0,4,8 - no, this is X if O is 0,4,8)
	// O . X
	// . O X
	// X . O (winning O move for diag 1: 0,4,8)
	testWinScenario(t, "O", []int{2, 0, 1, 4, 5}, 8, "O diag 1 (0,4,8)") // X:2,1,5 O:0,4,8
	// Example: O wins on diag 2 (2,4,6)
	// X plays 0, O plays 2
	// X plays 1, O plays 4
	// X plays 3 (arbitrary), O plays 6 (winning)
	testWinScenario(t, "O", []int{0, 2, 1, 4, 3}, 6, "O diag 2 (2,4,6)")
}

func TestDrawCondition(t *testing.T) {
	app := NewApp()
	app.gameState = app.NewGame()

	// Sequence leading to a draw:
	// X O X
	// X O X
	// O X O (O makes the last move at index 8)
	moves := []struct {
		player string
		index  int
	}{
		{"X", 0}, {"O", 1}, {"X", 2},
		{"X", 3}, {"O", 4}, {"X", 5},
		{"O", 6}, {"X", 7}, // O will make the last move at 8
	}

	for _, move := range moves {
		app.gameState.CurrentPlayer = move.player // Explicitly set player for test predictability
		_, err := app.MakeMove(move.index)
		if err != nil {
			t.Fatalf("MakeMove for draw setup failed at index %d for player %s: %v", move.index, move.player, err)
		}
	}

	// Last move by O at index 8
	app.gameState.CurrentPlayer = "O"
	gs, err := app.MakeMove(8)
	if err != nil {
		t.Fatalf("Final MakeMove for draw failed: %v", err)
	}

	if !gs.GameOver {
		t.Errorf("Expected GameOver to be true for a draw, got false")
	}
	if gs.Winner != "draw" {
		t.Errorf("Expected Winner to be 'draw', got '%s'", gs.Winner)
	}

	// Verify board is full
	boardIsFull := true
	for _, cell := range gs.Board {
		if cell == "" {
			boardIsFull = false
			break
		}
	}
	if !boardIsFull {
		t.Errorf("Expected board to be full for a draw, got %v", gs.Board)
	}
}

func TestResetGame(t *testing.T) {
	app := NewApp()
	app.gameState = app.NewGame() // Initial state

	// Make some moves to change the state
	_, _ = app.MakeMove(0) // X at 0
	_, _ = app.MakeMove(1) // O at 1

	// Manually set a win condition to ensure ResetGame clears it
	app.gameState.Winner = "X"
	app.gameState.GameOver = true
	app.gameState.Board[2] = "Something" // Put some garbage data

	resetState := app.ResetGame()

	if resetState.CurrentPlayer != "X" {
		t.Errorf("Expected CurrentPlayer to be 'X' after reset, got '%s'", resetState.CurrentPlayer)
	}
	if resetState.GameOver != false {
		t.Errorf("Expected GameOver to be false after reset, got '%t'", resetState.GameOver)
	}
	if resetState.Winner != "" {
		t.Errorf("Expected Winner to be '' after reset, got '%s'", resetState.Winner)
	}
	expectedBoard := [9]string{} // All empty
	if !reflect.DeepEqual(resetState.Board, expectedBoard) {
		t.Errorf("Expected empty board after reset, got '%v'", resetState.Board)
	}

	// Also verify that app.gameState is the same as the resetState
	if !reflect.DeepEqual(app.gameState, resetState) {
		t.Errorf("app.gameState was not correctly updated by ResetGame()")
	}
}
