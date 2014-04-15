package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

func crossProduct(s1, s2 string, size int) []string {
	arr := make([]string, size*size)
	index := 0
	for _, c1 := range s1 {
		for _, c2 := range s2 {
			arr[index] = string(c1) + string(c2)
			index++
		}
	}
	return arr
}

func makeUnits(s1, s2 string, reverse bool, size int) [][]string {
	units := make([][]string, size)
	for i, c := range s2 {
		if reverse {
			units[i] = crossProduct(s1, string(c), size)
		} else {
			units[i] = crossProduct(string(c), s1, size)
		}
	}
	return units
}

func makeBoxUnits(sl1, sl2 []string, size int) [][]string {
	units := make([][]string, size)
	index := 0
	for _, s1 := range sl1 {
		for _, s2 := range sl2 {
			units[index] = crossProduct(s1, s2, size)
			index++
		}
	}
	return units
}

func exists(sl []string, s1 string) bool {
	for _, s2 := range sl {
		if s1 == s2 {
			return true
		}
	}
	return false
}

func getUnitsAndPeers(squares []string, unitList [][]string) (map[string][][]string, map[string][]string) {
	units := map[string][][]string{}
	peers := map[string][]string{}
	for _, s := range squares {
		for _, u := range unitList {
			if exists(u, s) {
				units[s] = append(units[s], u)
				for _, sq := range u {
					if sq != s && !exists(peers[s], sq) {
						peers[s] = append(peers[s], sq)
					}
				}
			}
		}
	}
	return units, peers
}

func readInput() string {
	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getSize(grid string) int {
	for i, c := range grid {
		if c == '\n' {
			return i
		}
	}
	return -1
}

func makeGrid(grid string, squares []string) map[string]string {
	entries := map[string]string{}
	grid = strings.Replace(grid, "\n", "", -1)
	for i, c := range grid {
		entries[squares[i]] = string(c)
	}
	return entries
}

func parseGrid(grid, values map[string]string, units map[string][][]string, peers map[string][]string) map[string]string {
	for sq, digit := range grid {
		if digit != "." {
			tmp_values := assign(values, sq, digit, units, peers)
			if tmp_values == nil {
				return nil
			} else {
				values = tmp_values
			}
		}
	}
	return values
}

func assign(values map[string]string, square, digit string, units map[string][][]string, peers map[string][]string) map[string]string {
	other_digits := strings.Replace(values[square], digit, "", 1)
	for _, other_digit := range other_digits {
		tmp_values := eliminate(values, square, string(other_digit), units, peers)
		if tmp_values == nil {
			return nil
		} else {
			values = tmp_values
		}
	}
	return values
}

func eliminate(values map[string]string, square, digit string, units map[string][][]string, peers map[string][]string) map[string]string {
	if !strings.Contains(values[square], digit) {
		return values
	}
	values[square] = strings.Replace(values[square], digit, "", 1)

	if len(values[square]) == 0 {
		return nil
	} else if len(values[square]) == 1 {
		other_digit := values[square]
		for _, other_square := range peers[square] {
			tmp_values := eliminate(values, other_square, other_digit, units, peers)
			if tmp_values == nil {
				return nil
			} else {
				values = tmp_values
			}
		}
	}

	for _, unit := range units[square] {
		digit_places := []string{}
		for _, sq := range unit {
			if strings.Contains(values[sq], digit) {
				digit_places = append(digit_places, sq)
			}
		}
		if len(digit_places) == 0 {
			return nil
		} else if len(digit_places) == 1 {
			tmp_values := assign(values, digit_places[0], digit, units, peers)
			if tmp_values == nil {
				return nil
			} else {
				values = tmp_values
			}
		}
	}
	return values
}

func displayGrid(values map[string]string, size int) {
	box_sep := int(math.Sqrt(float64(size)))
	line_sep := size
	section_sep := line_sep * box_sep
	var squares []string
	for sq, _ := range values {
		squares = append(squares, sq)
	}
	sort.Strings(squares)
	i := 1
	for _, sq := range squares {
		digit := values[sq]
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

// It might be possible to modify this to range over values without using squares at all.
func isSolved(values map[string]string, squares []string) bool {
	for _, sq := range squares {
		if len(values[sq]) != 1 {
			return false
		}
	}
	return true
}

func getMin(values map[string]string, squares []string) string {
	minVal := 10
	var minSq string
	for _, sq := range squares {
		l := len(values[sq])
		if l > 1 && l < minVal {
			minSq = sq
		}
	}
	return minSq
}

func copy(oldmap map[string]string) map[string]string {
	newmap := make(map[string]string)
	for k, v := range oldmap {
		newmap[k] = v
	}
	return newmap
}

func search(values map[string]string, squares []string, units map[string][][]string, peers map[string][]string) map[string]string {
	if values == nil {
		return nil
	}
	if isSolved(values, squares) {
		return values
	}
	sq := getMin(values, squares)
	sequence := make([]map[string]string, len(values))
	for i, digit := range values[sq] {
		copyVals := copy(values)
		newVals := assign(copyVals, sq, string(digit), units, peers)
		sequence[i] = search(newVals, squares, units, peers)
	}
	return some(sequence)
}

func some(sequence []map[string]string) map[string]string {
	for _, elem := range sequence {
		if elem != nil {
			return elem
		}
	}
	return nil
}

func main() {
	// Read Sudoku grid from file and get the dimension of the grid.
	grid := readInput()
	size := getSize(grid)

	// Represent the Sudoku grid the same way as a chess board.
	// Each square in the grid is represented as 'A1', 'A2', etc.
	rows := "ABCDEFGHI"
	cols := "123456789"
	squares := crossProduct(rows, cols, size)

	// The possible values at a square.
	digits := "123456789"

	// Generate slices of row, column, and box units.
	// A unit is the set of all squares in a row, column, or box.
	rowBoxes := []string{"ABC", "DEF", "GHI"}
	colBoxes := []string{"123", "456", "789"}
	rowUnits := makeUnits(rows, cols, true, size)
	colUnits := makeUnits(cols, rows, false, size)
	boxUnits := makeBoxUnits(rowBoxes, colBoxes, size)

	// Create a map of squares to units and peers associated with each square.
	// A square's peer is every unique unit associated with that square
	// except for the square itself.
	unitList := append(append(rowUnits, colUnits...), boxUnits...)
	units, peers := getUnitsAndPeers(squares, unitList)

	// Create a map of squares to possible values.
	// Initially, each value can be any digit.
	values := map[string]string{}
	for _, sq := range squares {
		values[sq] = digits
	}

	// Parse the input grid into a map of squares to values.
	m := makeGrid(grid, squares)

	newVals := parseGrid(m, values, units, peers)
	finalVals := search(newVals, squares, units, peers)
	displayGrid(finalVals, size)
}
