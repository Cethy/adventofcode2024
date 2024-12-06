package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func getMap(raw string) [][]rune {
	var area [][]rune
	for _, line := range strings.Split(raw, "\n") {
		var row []rune
		for i := 0; i < len(line); i++ {
			row = append(row, rune(line[i]))
		}
		area = append(area, row)
	}
	return area
}

func getGuardPosition(area [][]rune) (int, int) {
	for i, line := range area {
		for j, cell := range line {
			if cell != '.' && cell != '#' && cell != 'X' {
				return i, j
			}
		}
	}
	return -1, -1
}

// tick makes the guard move by one cell
func tick(area [][]rune) ([][]rune, error) {
	guardY, guardX := getGuardPosition(area)
	guardRune := area[guardY][guardX]

	area[guardY][guardX] = 'X'
	switch guardRune {
	case '^': // "up"
		if guardY <= 0 {
			return area, errors.New("out of bound")
		}
		if area[guardY-1][guardX] == '#' {
			area[guardY][guardX] = '>'
			return tick(area)
		}
		area[guardY-1][guardX] = '^'
	case '>':
		if guardX >= len(area[guardY])-1 {
			return area, errors.New("out of bound")
		}
		if area[guardY][guardX+1] == '#' {
			area[guardY][guardX] = 'v'
			return tick(area)
		}
		area[guardY][guardX+1] = '>'
	case 'v': // down
		if guardY >= len(area)-1 {
			return area, errors.New("out of bound")
		}
		if area[guardY+1][guardX] == '#' {
			area[guardY][guardX] = '<'
			return tick(area)
		}
		area[guardY+1][guardX] = 'v'
	case '<':
		if guardX == 0 {
			return area, errors.New("out of bound")
		}
		if area[guardY][guardX-1] == '#' {
			area[guardY][guardX] = '^'
			return tick(area)
		}
		area[guardY][guardX-1] = '<'
	}

	return area, nil
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "....#.....\n.........#\n..........\n..#.......\n.......#..\n..........\n.#..^.....\n........#.\n#.........\n......#..."
	//raw = "....#.....\n.........#\n....>.....\n..#.......\n.......#..\n..........\n.#........\n........#.\n#.........\n......#..."
	//raw = "....#.....\n.........#\n..........\n..#.......\n.......#..\n..........\n.#........\n........#.\n#.........\n.....<#..."
	//raw = "....#....v\n.........#\n..........\n..#.......\n.......#..\n..........\n.#........\n........#.\n#.........\n......#..."

	area := getMap(raw)

	var oob error
	for oob == nil {
		/*for _, line := range area {
			for _, cell := range line {
				fmt.Print(string(cell))
			}
			fmt.Println("")
		}*/

		area, oob = tick(area)
	}

	fmt.Println("final area")
	var patrolDistinctPositions int
	for _, line := range area {
		for _, cell := range line {
			fmt.Print(string(cell))
			if cell == 'X' {
				patrolDistinctPositions++
			}
		}
		fmt.Println("")
	}
	fmt.Println("patrolDistinctPositions:", patrolDistinctPositions)
}
