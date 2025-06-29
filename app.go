package main

import (
	"context"
	"fmt"
)

// GameState represents the state of the Tic-tac-toe game
type GameState struct {
	Board         [9]string `json:"board"`
	CurrentPlayer string    `json:"currentPlayer"`
	Winner        string    `json:"winner"`
	GameOver      bool      `json:"gameOver"`
}

// App struct
type App struct {
	ctx       context.Context
	gameState *GameState
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.gameState = a.NewGame()
}

// GetInitialState returns the initial game state
func (a *App) GetInitialState() *GameState {
	return a.gameState
}

// NewGame initializes a new game
func (a *App) NewGame() *GameState {
	a.gameState = &GameState{
		Board:         [9]string{},
		CurrentPlayer: "X",
		Winner:        "",
		GameOver:      false,
	}
	return a.gameState
}

// MakeMove handles a player's move
func (a *App) MakeMove(index int) (*GameState, error) {
	if a.gameState.GameOver {
		return nil, fmt.Errorf("game is over")
	}

	if index < 0 || index >= 9 {
		return nil, fmt.Errorf("invalid move index")
	}

	if a.gameState.Board[index] != "" {
		return nil, fmt.Errorf("cell already occupied")
	}

	a.gameState.Board[index] = a.gameState.CurrentPlayer

	if checkWin(a.gameState.Board, a.gameState.CurrentPlayer) {
		a.gameState.Winner = a.gameState.CurrentPlayer
		a.gameState.GameOver = true
		return a.gameState, nil
	}

	if checkDraw(a.gameState.Board) {
		a.gameState.Winner = "draw"
		a.gameState.GameOver = true
		return a.gameState, nil
	}

	if a.gameState.CurrentPlayer == "X" {
		a.gameState.CurrentPlayer = "O"
	} else {
		a.gameState.CurrentPlayer = "X"
	}

	return a.gameState, nil
}

// ResetGame resets the game to its initial state
func (a *App) ResetGame() *GameState {
	return a.NewGame()
}

// checkWin checks if the current player has won
func checkWin(board [9]string, player string) bool {
	// Check rows
	for i := 0; i < 3; i++ {
		if board[i*3] == player && board[i*3+1] == player && board[i*3+2] == player {
			return true
		}
	}
	// Check columns
	for i := 0; i < 3; i++ {
		if board[i] == player && board[i+3] == player && board[i+6] == player {
			return true
		}
	}
	// Check diagonals
	if board[0] == player && board[4] == player && board[8] == player {
		return true
	}
	if board[2] == player && board[4] == player && board[6] == player {
		return true
	}
	return false
}

// checkDraw checks if the game is a draw
func checkDraw(board [9]string) bool {
	for _, cell := range board {
		if cell == "" {
			return false // If there's an empty cell, game is not a draw
		}
	}
	return true // All cells are filled
}
