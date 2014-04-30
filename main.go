package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

type Puzzle struct {
	grid      map[string]string // Grid of possible values at each square in the puzzle.
	numSolved int               // Number of solved squares (only 1 possible value for a square).
	alive     bool              // This puzzle is still a possible solution.
}

// Getter function for grid field.
func (p Puzzle) Grid() map[string]string {
	return p.grid
}

// Looks up the possible values for a given square.
func (p Puzzle) Get(square string) string {
	return p.grid[square]
}

// Sets the possible values for a given square.
func (p *Puzzle) Set(square string, digits string) {
	p.grid[square] = digits
}

// Removes a character from the possible values for a given square.
func (p *Puzzle) Remove(square string, char string) {
	vals := strings.Replace(p.grid[square], char, "", 1)
	p.grid[square] = vals
}

// Called when there is only one possible value for a square.
func (p *Puzzle) Increment() {
	p.numSolved++
}

// Called when there is only one possible value for all squares in a puzzle.
func (p Puzzle) IsSolved() bool {
	return p.numSolved == size*size
}

// Called when a contradiction is discovered during propogation and is
// is no longer possible for an instance of a Puzzle to be a solution.
func (p *Puzzle) Fail() {
	p.alive = false
}

// Finds the square with the fewest possible values in order to
// intelligently choose a starting point for the next search attempt.
func (p Puzzle) GetBestSquare(squares []string) (minSquare string) {
	minLength := size + 1
	for square, values := range p.Grid() {
		length := len(values)
		if length == 2 {
			return square
		} else if length > 1 && length < minLength {
			minSquare = square
		}
	}
	return minSquare
}

// Returns a copy of a Puzzle to be used in the next search attempt.
func (self *Puzzle) Copy() *Puzzle {
	p := &Puzzle{}
	*p = *self
	p.grid = CopyMap(self.Grid())
	return p
}

// Copies all elements in a map into a new map.
func CopyMap(oldMap map[string]string) map[string]string {
	newMap := make(map[string]string)
	for key, val := range oldMap {
		newMap[key] = val
	}
	return newMap
}

// Computes the cross product of two strings to create a
// combination of all squares associated with those strings.
func allSquares(str1, str2 string) []string {
	squares := make([]string, size*size)
	index := 0
	for _, char1 := range str1 {
		for _, char2 := range str2 {
			squares[index] = string(char1) + string(char2)
			index++
		}
	}
	return squares
}

// Computes slices of units for rows and columns.
func makeUnits(str1, str2 string, reverse bool) [][]string {
	units := make([][]string, size)
	for index, char := range str2 {
		if reverse {
			units[index] = allSquares(str1, string(char))
		} else {
			units[index] = allSquares(string(char), str1)
		}
	}
	return units
}

// Computes slices of units for boxes.
func makeBoxUnits(slice1, slice2 []string) [][]string {
	units := make([][]string, size)
	index := 0
	for _, str1 := range slice1 {
		for _, str2 := range slice2 {
			units[index] = allSquares(str1, str2)
			index++
		}
	}
	return units
}

// Checks if a slice contains a given string.
func contains(slice []string, toMatch string) bool {
	for _, str := range slice {
		if toMatch == str {
			return true
		}
	}
	return false
}

// Creates maps of squares to units and peers associated with each square.
func getUnitsAndPeers(squares []string, allUnits [][]string) (map[string][][]string, map[string][]string) {
	units := make(map[string][][]string)
	peers := make(map[string][]string)
	for _, square := range squares {
		for _, unit := range allUnits {
			if contains(unit, square) {
				units[square] = append(units[square], unit)
				for _, usquare := range unit {
					if usquare != square && !contains(peers[square], usquare) {
						peers[square] = append(peers[square], usquare)
					}
				}
			}
		}
	}
	return units, peers
}

// Places given input values into a puzzle and propagates constraints.
func parseInput(p *Puzzle, grid map[string]string) bool {
	for square, char := range grid {
		if char != "." && !place(p, square, char) {
			p.Fail()
			return false
		}
	}
	return true
}

// Assigns a given character as the solution value at a given square.
func place(p *Puzzle, square, char string) bool {
	for _, val := range allValues {
		if string(val) != char && !propagate(p, square, string(val)) {
			p.Fail()
			return false
		}
	}
	return true
}

// Removes a given character as a possible solution value at a given square and uses
// the new values at that square to eliminate more possibilities from other squares.
func propagate(p *Puzzle, square, char string) bool {
	// The value is not present in the possibilites; already propagated.
	if !strings.Contains(p.Get(square), char) {
		return true
	}

	// Else, remove that character from the possibilites and begin propogation.
	p.Remove(square, char)
	vals := p.Get(square)
	length := len(vals)

	// If a square has a solution value, remove that value from the square's peers.
	if length == 0 { // None of the possible values for this square result in a solution; fail.
		p.Fail()
		return false
	} else if length == 1 { // There is only one possible value for this square; propagate.
		p.Increment()
		for _, peerSquare := range peers[square] {
			if !propagate(p, peerSquare, vals) { // Propogation does not return a valid solution; fail.
				p.Fail()
				return false
			}
		}
	} else if length == 2 { // Check if any squares can be eliminated using the "naked twins" strategy.
		if int(square[1])+1 < size {
			adjacent := string(square[0]) + string(square[1]+1)
			if vals == p.Get(adjacent) {
				for _, unitSquare := range units[square][1] {
					if unitSquare != adjacent { // Remove "twins" and propagate.
						if !propagate(p, unitSquare, string(vals[0])) || !propagate(p, unitSquare, string(vals[1])) {
							p.Fail()
							return false
						}
					}
				}
			}
		}
	}

	// If there is only one possible place for a value in a unit, place that value accordingly.
	for _, unit := range units[square] {
		locations := make([]string, size) // Valid locations for a value in a unit.
		index := 0
		for _, unitSquare := range unit {
			if strings.Contains(p.Get(unitSquare), char) {
				locations[index] = unitSquare
				index++
			}
		}
		if len(locations) == 0 { // The value can't be placed at any of the possible locations; fail.
			return false
		} else if len(locations) == 1 && !place(p, locations[0], char) {
			// There is only one possible location for a value; place and propagate.
			p.Fail()
			return false
		}
	}
	return true
}

// Uses depth-first search to look for solutions to a puzzle.
func search(p *Puzzle) {
	if !p.alive { // There was a contradiction during propogation; this is not a solution.
		return
	} else if p.IsSolved() { // All squares have only one possible value; display the solution.
		showPuzzle(p)
		os.Exit(0)
	} else { // This could be a solution; pick a good square and keep searching.
		square := p.GetBestSquare(squares)
		for _, val := range p.Get(square) {
			newP := p.Copy()
			place(newP, square, string(val))
			search(newP)
		}
	}
	return
}

// Takes in a filename as a command-line argument and returns the contents as a string.
func readInput() string {
	filename := os.Args[1]
	content, _ := ioutil.ReadFile(filename)
	return string(content)
}

// Scans the first line of the input grid to determine the dimensions of the grid.
func getSize(grid string) int {
	for index, char := range grid {
		if char == '\n' {
			return index
		}
	}
	return -1
}

// Converts the contents of a grid file into a map of squares to
// blank '.' characters and "hint" values given in the input grid.
func makeGrid(grid string, squares []string) map[string]string {
	entries := make(map[string]string)
	grid = strings.Replace(grid, "\n", "", -1)
	for index, char := range grid {
		entries[squares[index]] = string(char)
	}
	return entries
}

// Formats and displays a puzzle.
func showPuzzle(p *Puzzle) {
	grid := p.Grid()
	box_sep := int(math.Sqrt(float64(size)))
	line_sep := size
	section_sep := line_sep * box_sep
	var squares []string
	for sq, _ := range grid {
		squares = append(squares, sq)
	}
	sort.Strings(squares)
	i := 1
	for _, sq := range squares {
		digit := grid[sq]
		fmt.Print(digit + " ")
		if i%section_sep == 0 {
			fmt.Println("\n")
		} else if i%line_sep == 0 {
			fmt.Println()
		} else if i%box_sep == 0 {
			fmt.Print("| ")
		}
		i++
	}
}

var size int                    // Dimensions of a grid.
var allValues string            // All possible solution values for a square.
var squares []string            // All squares in a grid.
var units map[string][][]string // A unit is the set of all squares in a row, column, or box.
var peers map[string][]string   // A peer is the set of all unique unit squares for a given square.

func main() {
	// Read Sudoku from file and get the size of the grid.
	rawGrid := readInput()
	size = getSize(rawGrid)

	// Represent the Sudoku grid the same way as a chess board.
	// Each square in the grid is represented as 'A1', 'A2', etc.
	rows := "ABCDEFGHI"
	cols := "123456789"
	allValues = "123456789"
	rowBoxes := []string{"ABC", "DEF", "GHI"}
	colBoxes := []string{"123", "456", "789"}

	// Adjust values for bigger puzzles.
	if size == 16 {
		rows = "ABCDEFGHIJKLMNOP"
		cols = "123456789ABCDEFG"
		allValues = "123456789ABCDEFG"
		rowBoxes = []string{"ABCD", "EFGH", "IJKL", "MNOP"}
		colBoxes = []string{"1234", "5678", "9ABC", "DEFG"}
	} else if size == 25 {
		rows = "ABCDEFGHIJKLMNOPQRSTUVWXY"
		cols = "123456789ABCDEFGHIJKLMNOP"
		allValues = "123456789ABCDEFGHIJKLMNOP"
		rowBoxes = []string{"ABCDE", "FGHIJ", "KLMNO", "PQRST", "UVWXY"}
		colBoxes = []string{"12345", "6789A", "BCDEF", "GHIJK", "LMNOP"}
	}

	// Get a slice of all squares in the Sudoku grid.
	squares = allSquares(rows, cols)

	// For each square, create a map of units and peers associated with that square.
	rowUnits := makeUnits(rows, cols, true)
	colUnits := makeUnits(cols, rows, false)
	boxUnits := makeBoxUnits(rowBoxes, colBoxes)
	allUnits := append(append(rowUnits, colUnits...), boxUnits...)
	units, peers = getUnitsAndPeers(squares, allUnits)

	// Represent a Sudoku puzzle as a map of squares to possible values at each square.
	puzzle := &Puzzle{
		grid:      make(map[string]string),
		numSolved: 0,
		alive:     true,
	}

	// Initially, each square can take on any value.
	for _, square := range squares {
		puzzle.Set(square, allValues)
	}

	// Parse the input grid into a map of squares to values.
	grid := makeGrid(rawGrid, squares)

	// Assign the values from the input grid to their corresponding squares in the solution grid.
	// propagate constraints to solve as much of the puzzle as possible before searching.
	parseInput(puzzle, grid)

	// Use depth-first search to try possible values until a solution is found.
	search(puzzle)
}
