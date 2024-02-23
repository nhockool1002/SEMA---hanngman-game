package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

// Profile struct represents the player's profile including their name, high score, and pet.
type Profile struct {
	Name      string
	HighScore int
	Pet       Pet
}

// Pet struct represents the virtual pet in the game.
type Pet struct {
	Name      string
	Happiness int
	Hunger    int
}

// Game struct represents the Hangman game.
type Game struct {
	word          string
	guessedLetters []string
	score         int
	maxAttempts   int
	attempts      int
	profile       Profile
}

// Function to select a random word from a list
func randomWord() string {
	words := []string{"apple", "banana", "cherry", "orange", "grape", "kiwi", "melon", "peach"}
	rand.Seed(time.Now().UnixNano())
	return words[rand.Intn(len(words))]
}

// Function to display the current state of the word with guessed letters
func (g *Game) displayWord() string {
	display := ""
	for _, letter := range g.word {
		if contains(g.guessedLetters, string(letter)) {
			display += string(letter)
		} else {
			display += "_"
		}
		display += " "
	}
	return display
}

// Function to check if a slice contains a specific element
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// Function to read the profile from a file
func readProfile() Profile {
	file, err := os.Open("profile.gob")
	if err != nil {
		return Profile{}
	}
	defer file.Close()

	var profile Profile
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&profile)
	if err != nil {
		return Profile{}
	}
	return profile
}

// Function to write the profile to a file
func writeProfile(profile Profile) {
	file, err := os.Create("profile.gob")
	if err != nil {
		fmt.Println("Error writing profile:", err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(profile)
	if err != nil {
		fmt.Println("Error writing profile:", err)
	}
}

func main() {
	// Initialize the game
	game := Game{
		word:          randomWord(),
		guessedLetters: []string{},
		score:         0,
		maxAttempts:   6,
		attempts:      0,
		profile:       readProfile(),
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Hangman!")
	fmt.Println("Try to guess the word.")

	// Main game loop
	for {
		// Display the word with guessed letters
		fmt.Println(game.displayWord())

		// Prompt the player for a guess
		fmt.Print("Enter a letter: ")
		guess, _ := reader.ReadString('\n')
		guess = strings.TrimSpace(strings.ToLower(guess))

		// Check if the guess is a single letter
		if len(guess) != 1 || !('a' <= guess[0] && guess[0] <= 'z') {
			fmt.Println("Please enter a single letter.")
			continue
		}

		// Check if the letter has already been guessed
		if contains(game.guessedLetters, guess) {
			fmt.Println("You already guessed that letter.")
			continue
		}

		// Add the guess to the list of guessed letters
		game.guessedLetters = append(game.guessedLetters, guess)

		// Check if the guess is correct
		if strings.Contains(game.word, guess) {
			fmt.Println("Correct guess!")
			game.score += 10 // Increase score for correct guess
		} else {
			fmt.Println("Incorrect guess.")
			game.attempts++
			if game.attempts >= game.maxAttempts {
				fmt.Println("You've run out of attempts. The word was:", game.word)
				break
			}
			game.score -= 5 // Decrease score for incorrect guess
		}

		// Check if the player has guessed all the letters
		if strings.ReplaceAll(game.displayWord(), " ", "") == game.word {
			fmt.Println("Congratulations! You've guessed the word:", game.word)
			break
		}
	}

	// Update profile with high score
	if game.score > game.profile.HighScore {
		game.profile.HighScore = game.score
	}
	game.profile.Pet.Happiness += 10
	game.profile.Pet.Hunger += 5

	// Display final score and pet status
	fmt.Println("Final Score:", game.score)
	fmt.Println("Pet Happiness:", game.profile.Pet.Happiness)
	fmt.Println("Pet Hunger:", game.profile.Pet.Hunger)

	// Save profile
	writeProfile(game.profile)
}
