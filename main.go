package main

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/DevPedroHB/hangman-game/hangman"
	"github.com/labstack/echo/v4"
)

var game hangman.Hangman

type GuessRequest struct {
	Letter string `json:"letter"`
}

func main() {
	// Inicializa o servidor e o jogo
	rand.Seed(time.Now().UnixNano())

	game = hangman.Hangman{
		Words: []string{"apple", "banana", "cherry", "peteca", "luz", "janela", "abacaxi", "suquinho"},
	}

	game.InitGame()

	e := echo.New()

	// Rotas
	e.PATCH("/game", startNewGame)
	e.GET("/game", getGameState)
	e.POST("/game/guess", makeGuess)

	e.Logger.Fatal(e.Start(":1323"))
}

// startNewGame inicia um novo jogo.
func startNewGame(c echo.Context) error {
	game.InitGame()

	return c.String(http.StatusOK, "Novo jogo iniciado!")
}

// getGameState retorna o estado atual do jogo.
func getGameState(c echo.Context) error {
	result := hangman.GuessResult{
		ShowWord: game.ShowWord,
		Lives:    game.Lives,
		Status:   "playing",
    Attempts: game.GetAttempts(),
	}

	if !strings.Contains(game.ShowWord, "_") {
		result.Status = "won"
		result.ChosenWord = game.ChosenWord
	} else if game.Lives <= 0 {
		result.Status = "lost"
		result.ChosenWord = game.ChosenWord
	}

	return c.JSON(http.StatusOK, result)
}

// makeGuess processa uma tentativa do jogador.
func makeGuess(c echo.Context) error {
	var req GuessRequest

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Entrada inválida.")
	}

	result := game.GuessLetter(req.Letter)

	return c.JSON(http.StatusOK, result)
}
