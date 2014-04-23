package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

type ValueMap struct {
	m         map[string]string
	numSolved int
	flag      bool // might be able to get rid of this
}

func (vm ValueMap) GetMinSquare(squares []string) (minSquare string) {
	minLength := size + 1
	for square, value := range vm.M() {
		length := len(value)
		if length == 2 {
			return square
		} else if length > 1 && length < minLength {
			minSquare = square
		}
	}
	return minSquare
}

func (vm ValueMap) Values(square string) string {
	return vm.m[square]
}

func (vm ValueMap) M() map[string]string {
	return vm.m
}

func (vm ValueMap) IsSolved() bool {
	return vm.numSolved == size*size
}

func (vm *ValueMap) Fail() {
	vm.flag = false
}

func (vm *ValueMap) IncrementNumSolved() {
	vm.numSolved++
}

func (vm *ValueMap) RemoveValue(square string, digit string) {
	vals := strings.Replace(vm.m[square], digit, "", 1)
	vm.m[square] = vals
}

func (vm *ValueMap) SetValues(square string, digits string) {
	vm.m[square] = digits
}

func crossProduct(s1, s2 string) []string {
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

func makeUnits(s1, s2 string, reverse bool) [][]string {
	units := make([][]string, size)
	for i, c := range s2 {
		if reverse {
			units[i] = crossProduct(s1, string(c))
		} else {
			units[i] = crossProduct(string(c), s1)
		}
	}
	return units
}

func makeBoxUnits(sl1, sl2 []string) [][]string {
	units := make([][]string, size)
	index := 0
	for _, s1 := range sl1 {
		for _, s2 := range sl2 {
			units[index] = crossProduct(s1, s2)
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

func parseGrid(grid map[string]string, values *ValueMap, units map[string][][]string, peers map[string][]string, digits string) bool {
	for sq, digit := range grid {
		if digit != "." {
			if !assign(values, sq, digit, units, peers, digits) {
				values.Fail()
				return false
			}
		}
	}
	return true
}

func assign(values *ValueMap, square, digit string, units map[string][][]string, peers map[string][]string, digits string) bool {
	for _, d := range digits {
		if string(d) != digit {
			if !eliminate(values, square, string(d), units, peers, digits) {
				values.Fail()
				return false
			}
		}
	}
	return true
}

func eliminate(values *ValueMap, square, digit string, units map[string][][]string, peers map[string][]string, digits string) bool {
	if !strings.Contains(values.Values(square), digit) {
		return true
	}
	values.RemoveValue(square, digit)
	vals := values.Values(square)
	length := len(vals)

	if length == 0 {
		values.Fail()
		return false
	} else if length == 1 {
		values.IncrementNumSolved()
		solutionDigit := vals // might be able to use vals directly without need to reassign
		for _, peerSquare := range peers[square] {
			if !eliminate(values, peerSquare, solutionDigit, units, peers, digits) {
				values.Fail()
				return false
			}
		}
	}

	for _, unit := range units[square] {
		digitPlaces := []string{}
		for _, unitSquare := range unit {
			if strings.Contains(values.Values(unitSquare), digit) {
				digitPlaces = append(digitPlaces, unitSquare)
			}
		}
		if len(digitPlaces) == 0 {
			return false
		} else if len(digitPlaces) == 1 {
			if !assign(values, digitPlaces[0], digit, units, peers, digits) {
				values.Fail()
				return false
			}
		}
	}
	return true
}

func displayGrid(values *ValueMap) {
	vals := values.M()
	box_sep := int(math.Sqrt(float64(size)))
	line_sep := size
	section_sep := line_sep * box_sep
	var squares []string
	for sq, _ := range vals {
		squares = append(squares, sq)
	}
	sort.Strings(squares)
	i := 1
	for _, sq := range squares {
		digit := vals[sq]
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

func CopyMap(oldmap map[string]string) map[string]string {
	newmap := make(map[string]string)
	for k, v := range oldmap {
		newmap[k] = v
	}
	return newmap
}

func (self *ValueMap) Copy() *ValueMap {
	vm := &ValueMap{}
	*vm = *self
	vm.m = CopyMap(self.M())
	return vm
}

func search(values *ValueMap, squares []string, units map[string][][]string, peers map[string][]string, digits string) {
	if !values.flag {
		return
	} else if values.IsSolved() {
		displayGrid(values)
		os.Exit(0)
	} else {
		sq := values.GetMinSquare(squares)
		for _, digit := range values.Values(sq) {
			copyVals := values.Copy()
			assign(copyVals, sq, string(digit), units, peers, digits)
			search(copyVals, squares, units, peers, digits)
		}
	}
	return
}

var size int

func main() {
	// Read Sudoku grid from file and get the dimension of the grid.
	grid := readInput()
	size = getSize(grid)

	// Represent the Sudoku grid the same way as a chess board.
	// Each square in the grid is represented as 'A1', 'A2', etc.
	rows := "ABCDEFGHI"
	cols := "123456789"
	rowBoxes := []string{"ABC", "DEF", "GHI"}
	colBoxes := []string{"123", "456", "789"}

	// The possible values at a square.
	digits := "123456789"
	// if size == 16 {
	// 	rows = "ABCDEFGHIJKLMNOP"
	// 	cols = "123456789ABCDEFG"
	// 	digits = "123456789ABCDEFG"
	// 	rowBoxes = []string{"ABCD", "EFGH", "IJKL", "MNOP"}
	// 	colBoxes = []string{"1234", "5678", "9ABC", "DEFG"}
	// } else if size == 25 {
	// 	rows = "ABCDEFGHIJKLMNOPQRSTUVWXY"
	// 	cols = "123456789ABCDEFGHIJKLMNOP"
	// 	digits = "123456789ABCDEFGHIJKLMNOP"
	// 	rowBoxes = []string{"ABCDE", "FGHIJ", "KLMNO", "PQRST", "UVWXY"}
	// 	colBoxes = []string{"12345", "6789A", "BCDEF", "GHIJK", "LMNOP"}
	// }
	squares := crossProduct(rows, cols)

	// Generate slices of row, column, and box units.
	// A unit is the set of all squares in a row, column, or box.
	rowUnits := makeUnits(rows, cols, true)
	colUnits := makeUnits(cols, rows, false)
	boxUnits := makeBoxUnits(rowBoxes, colBoxes)

	// Create a map of squares to units and peers associated with each square.
	// A square's peer is every unique unit associated with that square
	// except for the square itself.
	unitList := append(append(rowUnits, colUnits...), boxUnits...)
	units, peers := getUnitsAndPeers(squares, unitList)

	// Create a map of squares to possible values.
	// Initially, each value can be any digit.
	values := &ValueMap{
		m:         make(map[string]string),
		numSolved: 0,
		flag:      true,
	}

	for _, square := range squares {
		values.SetValues(square, digits)
	}

	// Parse the input grid into a map of squares to values.
	m := makeGrid(grid, squares)

	parseGrid(m, values, units, peers, digits)
	search(values, squares, units, peers, digits)
}
