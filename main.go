package main

import (
	"flag"
	"os"
	"scramble/scramble"
	"github.com/fatih/color"
)

func main() {
	// Parse cl args
	count := flag.Int("words", 5, "Number of words")
	height := flag.Int("height", 10, "Grid height")
	width := flag.Int("width", 10, "Grid width")
	maxLength := flag.Int("maxLength", 7, "Max word length")
	minLength := flag.Int("minLength", 3, "Min word length")
	backwards := flag.Bool("noBackwards", true, "Prevent backwards words")
	diagonals := flag.Bool("noDiagonals", true, "Prevent diagonals")
	input := flag.String("input", "words.json", "Input word file")
	seed := flag.Int64("seed", 0, "Puzzle seed")
	flag.Parse()

	// Create the game's options
	options := scramble.Options{
		WordCount:      *count,
		AllowBackwards: *backwards,
		Height:         *height,
		Width:          *width,
		MaxWordLength:  *maxLength,
		MinWordLength:  *minLength,
		AllowDiagonals: *diagonals,
		Seed:           *seed,
	}

	// Load words and generate the board
	scramble.LoadWords(*input)
	board, err := scramble.New(options)

	// Print error and exit if something went wrong
	if err != nil {
		println(err.Error())
		os.Exit(3)
	}

	// Print board and words
	board.Print()
	c := color.New(color.Underline)
	c.Println("\nWords")
	for _, v := range board.GetWords() {
		println(v)
	}

	color.Yellow("\nSeed: %d", board.Options.Seed)
}
