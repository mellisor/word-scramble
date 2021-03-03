package scramble

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/fatih/color"
)

var words []string
var baseValue = byte(0)

// Indicates the direction a word is running
type direction int

const (
	horizontal direction = 0
	vertical             = 1
	diagonal             = 2
)

// Options contains options for board configuration
type Options struct {
	WordCount      int
	AllowBackwards bool
	Height         int
	Width          int
	MaxWordLength  int
	MinWordLength  int
	AllowDiagonals bool
	Seed           int64
	Cheat          bool
}

type space struct {
	letter byte
	used   bool
}

// Puzzle contains state of the game
type Puzzle struct {
	board   [][]space
	Options *Options
	words   []string
}

// Print prints the board
func (p Puzzle) Print() {
	cheatPrinter := color.New(color.FgYellow)
	for y := 0; y < len(p.board[0]); y++ {
		for x := range p.board {
			str := string(p.board[x][y].letter) + " "
			if p.Options.Cheat && p.board[x][y].used {
				cheatPrinter.Print(str)
			} else {
				print(str)
			}
		}
		println()
	}
}

func pad(word string, size int) string {
	padded := word
	for len(padded) < size {
		padded = padded + " "
	}
	return padded
}

// PrintWords handles the printing of words nicely
func (p Puzzle) PrintWords(columns int) {
	c := color.New(color.Underline)

	colWidth := p.Options.MaxWordLength + 1
	boardWords := p.GetWords()
	c.Printf("\n%d Words\n", len(boardWords))
	n := 0
	for _, v := range boardWords {
		if n >= columns {
			println();
			n = 0
		}

		print(pad(v, colWidth))
		n++
	}

	color.Yellow("\nSeed: %d", p.Options.Seed)
}

// GetWords gets the list of available words
func (p Puzzle) GetWords() []string {
	sort.Strings(p.words)
	return p.words
}

// Get all available spots a word can be placed
func (p Puzzle) getOpenSpaces(word string) ([]int, []int, []direction) {
	xChoices := make([]int, 0)
	yChoices := make([]int, 0)
	dChoices := make([]direction, 0)
	directions := []direction { horizontal, vertical }
	if (p.Options.AllowDiagonals) {
		directions = append(directions, diagonal)
	}
	for _, d := range directions {
		// Get the maximum x and y coordinate this word can be placed
		maxX := p.Options.Width - 1
		maxY := p.Options.Height - 1
		if d != vertical {
			maxX = maxX - len(word) + 1
		}
		if d != horizontal {
			maxY = maxY - len(word) + 1
		}

		// Evaluate each location to see if there is a fit
		for x := 0; x <= maxX; x++ {
			for y := 0; y <= maxY; y++ {
				tempX := x
				tempY := y
				success := true
				for curr := range word {
					boardValue := p.board[tempX][tempY].letter
					char := word[curr]

					// If the current value is non-default and doesn't match this letter, not a match
					if boardValue != baseValue && char != boardValue {
						success = false
						break
					}

					// Iterate indices
					if d != vertical {
						tempX++
					}
					if d != horizontal {
						tempY++
					}
				}

				// Append to choices if successful
				if success {
					xChoices = append(xChoices, x)
					yChoices = append(yChoices, y)
					dChoices = append(dChoices, d)
				}

			}
		}
	}
	return xChoices, yChoices, dChoices
}

func (p *Puzzle) populateBoard() error {
	// Initialize things
	if p.Options.Seed == 0 {
		p.Options.Seed = time.Now().UnixNano()
	}
	rand.Seed(p.Options.Seed)

	// Get the list of available words
	availableWords := make([]string, 0)
	for i := range words {
		if len(words[i]) > p.Options.MaxWordLength || len(words[i]) < p.Options.MinWordLength {
			continue
		}

		availableWords = append(availableWords, words[i])
	}

	if len(availableWords) == 0 {
		return errors.New("No available words")
	}

	for j := 0; j < p.Options.WordCount; j++ {
		// Break if there are no available words left
		if len(availableWords) == 0 {
			break
		}

		// Determine the word
		index := rand.Intn(len(availableWords))
		word := availableWords[index]

		// Remove the chosen word from available words
		availableWords[index], availableWords[len(availableWords)-1] = availableWords[len(availableWords)-1], availableWords[index]
		availableWords = availableWords[:len(availableWords)-1]

		// Determine whether the word is forwards or backwards
		boardWord := word
		if p.Options.AllowBackwards && rand.Intn(2) == 1 {
			boardWord = ""
			for _, v := range word {
				boardWord = string(v) + boardWord
			}
		}

		// Get the available spots for this word. If there are none, continue
		xChoices, yChoices, directions := p.getOpenSpaces(boardWord)
		if len(xChoices) == 0 {
			continue
		}
		p.words = append(p.words, word)

		// Pick a choice and place the word
		choice := rand.Intn(len(xChoices))
		x := xChoices[choice]
		y := yChoices[choice]
		d := directions[choice]
		
		for i := range boardWord {
			newSpace := space{
				letter: boardWord[i],
				used:   true,
			}
			p.board[x][y] = newSpace
			if d != vertical {
				x++
			}
			if d != horizontal {
				y++
			}
		}
	}

	// Fill in the rest of the spaces
	for i := 0; i < p.Options.Width; i++ {
		for j := 0; j < p.Options.Height; j++ {
			if p.board[i][j].letter == baseValue {
				newSpace := space{
					letter: byte(rand.Intn(26) + 97),
					used:   false,
				}
				p.board[i][j] = newSpace
			}
		}
	}

	return nil
}

// New returns a board using given the options
func New(options Options) (Puzzle, error) {

	var puzzle Puzzle
	var e error

	// Evaluate options
	if options.MaxWordLength > options.Height && options.MaxWordLength > options.Width {
		e = errors.New("Max word length exceeds board dimensions")
	} else if options.Height < 1 || options.Width < 1 {
		e = errors.New("Invalid board dimensions")
	} else if options.WordCount < 1 {
		e = errors.New("Invalid word count")
	} else if options.MinWordLength > options.MaxWordLength {
		e = errors.New("Min word length is greater than max word length")
	}

	// Abort if there is an error
	if e != nil {
		return puzzle, e
	}

	// Make the board
	a := make([][]space, options.Width)
	for i := range a {
		a[i] = make([]space, options.Height)
	}

	// Declare the puzzle
	puzzle = Puzzle{
		Options: &options,
		board:   a,
	}

	// Populate the board's values
	e = puzzle.populateBoard()

	return puzzle, e
}

// LoadWords loads words for the puzzle from the specified file
func LoadWords(file string) {

	f, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()
	data, _ := ioutil.ReadAll(f)
	json.Unmarshal([]byte(data), &words)
}
