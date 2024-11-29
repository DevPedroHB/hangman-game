package hangman

import (
	"math/rand"
	"strings"
	"unicode"
)

type Hangman struct {
	Words      []string
	ChosenWord string
	ShowWord   string
	Lives      int
	Attempts   map[string]bool
}

type GuessResult struct {
	ShowWord   string   `json:"showWord"`
	Lives      int      `json:"lives"`
	Status     string   `json:"status"`
	ChosenWord string   `json:"chosenWord,omitempty"`
	Attempts   []string `json:"attempts"`
}

// InitGame inicializa um novo jogo da forca.
func (h *Hangman) InitGame() {
	h.ChosenWord = h.Words[rand.Intn(len(h.Words))]
	h.ShowWord = strings.Repeat("_", len(h.ChosenWord))
	h.Lives = 6
	h.Attempts = make(map[string]bool)
}

// GuessLetter processa uma letra enviada pelo jogador.
func (h *Hangman) GuessLetter(letter string) GuessResult {
	if h.Lives <= 0 || !strings.Contains(h.ShowWord, "_") {
		status := "playing"

		if h.Lives <= 0 {
			status = "lost"
		} else if !strings.Contains(h.ShowWord, "_") {
			status = "won"
		}

		return GuessResult{
			ShowWord:   h.ShowWord,
			Lives:      h.Lives,
			Status:     status,
			ChosenWord: h.ChosenWord,
			Attempts:   h.GetAttempts(),
		}
	}

	letter = strings.ToLower(letter)

	if len(letter) != 1 || !unicode.IsLetter([]rune(letter)[0]) {
		return GuessResult{
			ShowWord: h.ShowWord,
			Lives:    h.Lives,
			Status:   "invalid_input",
			Attempts: h.GetAttempts(),
		}
	}

	if h.Attempts[letter] {
		return GuessResult{
			ShowWord: h.ShowWord,
			Lives:    h.Lives,
			Status:   "already_guessed",
			Attempts: h.GetAttempts(),
		}
	}

	h.Attempts[letter] = true

	if strings.Contains(h.ChosenWord, letter) {
		for i, r := range h.ChosenWord {
			if string(r) == letter {
				h.ShowWord = replaceAt(h.ShowWord, i, r)
			}
		}

		if !strings.Contains(h.ShowWord, "_") {
			return GuessResult{
				ShowWord:   h.ShowWord,
				Lives:      h.Lives,
				Status:     "won",
				ChosenWord: h.ChosenWord,
				Attempts:   h.GetAttempts(),
			}
		}
	} else {
		h.Lives--

		if h.Lives <= 0 {
			return GuessResult{
				ShowWord:   h.ShowWord,
				Lives:      0,
				Status:     "lost",
				ChosenWord: h.ChosenWord,
				Attempts:   h.GetAttempts(),
			}
		}
	}

	return GuessResult{
		ShowWord: h.ShowWord,
		Lives:    h.Lives,
		Status:   "playing",
		Attempts: h.GetAttempts(),
	}
}

// getAttempts retorna as letras já tentadas em forma de lista.
func (h *Hangman) GetAttempts() []string {
	attempts := make([]string, 0, len(h.Attempts))

	for letter := range h.Attempts {
		attempts = append(attempts, letter)
	}

	return attempts
}

// replaceAt substitui o caractere na posição específica.
func replaceAt(s string, index int, char rune) string {
	runes := []rune(s)
	runes[index] = char

	return string(runes)
}
