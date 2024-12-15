package main

import (
	"fmt"
	"os"
	"strings"
)

func readInput(raw string) ([][]string, []string) {
	ra := strings.Split(raw, "\n\n")
	var area [][]string
	var instructions []string
	for _, line := range strings.Split(ra[0], "\n") {
		var l []string
		for _, cell := range line {
			l = append(l, string(cell))
		}
		area = append(area, l)
	}
	for _, i := range strings.Replace(ra[1], "\n", "", -1) {
		instructions = append(instructions, string(i))
	}

	return area, instructions
}

func displayMap(area [][]string) {
	for _, line := range area {
		for _, cell := range line {
			fmt.Print(cell)
		}
		fmt.Println()
	}
}

func findBotPos(area [][]string) [2]int {
	for y, line := range area {
		for x, cell := range line {
			if cell == "@" {
				return [2]int{x, y}
			}
		}
	}
	return [2]int{-1, -1}
}

func mBot(area [][]string, vector [2]int /*, limit [2]int*/) [][]string {
	botPos := findBotPos(area)

	nextPos := [2]int{botPos[0] + vector[0], botPos[1] + vector[1]}
	// empty tile
	if area[nextPos[1]][nextPos[0]] == "." {
		// simple move
		area[nextPos[1]][nextPos[0]], area[botPos[1]][botPos[0]] = area[botPos[1]][botPos[0]], area[nextPos[1]][nextPos[0]]
		return area
	}
	// wall
	if area[nextPos[1]][nextPos[0]] == "#" {
		// blocked
		// stay in place
		return area
	}
	// box
	if area[nextPos[1]][nextPos[0]] == "O" {
		//firstBoxPos := [2]int{nextPos[0], nextPos[1]}
		// look for some empty pos after the box
		firstEmptySpacePos := [2]int{-1, -1}
	mainLoop:
		for y := nextPos[1]; y > 0 && y < len(area); y = y + vector[1] {
			for x := nextPos[0]; x > 0 && x < len(area[0]); x = x + vector[0] {
				// no empty space: bot not moving
				if area[y][x] == "#" {
					return area
				}
				// empty space found: bot not moving
				if area[y][x] == "." {
					firstEmptySpacePos = [2]int{x, y}
					break mainLoop
				}

				if vector[0] == 0 {
					break
				}
			}

			if vector[1] == 0 {
				break
			}
		}

		// no empty space: bot not moving
		if firstEmptySpacePos[0] == -1 {
			return area
		}

		// exchange firstBox with firstEmptySpace pos
		area[nextPos[1]][nextPos[0]], area[firstEmptySpacePos[1]][firstEmptySpacePos[0]] = area[firstEmptySpacePos[1]][firstEmptySpacePos[0]], area[nextPos[1]][nextPos[0]]
		// move bot
		area[nextPos[1]][nextPos[0]], area[botPos[1]][botPos[0]] = area[botPos[1]][botPos[0]], area[nextPos[1]][nextPos[0]]

		return area
	}
	return area
}

func moveBot(area [][]string, instruction string) [][]string {
	switch instruction {
	case "^":
		return mBot(area, [2]int{0, -1})
	case ">":
		return mBot(area, [2]int{1, 0})
	case "v":
		return mBot(area, [2]int{0, 1})
	case "<":
		return mBot(area, [2]int{-1, 0})
	}

	return area
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "##########\n#..O..O.O#\n#......O.#\n#.OO..O.O#\n#..O@..O.#\n#O#..O...#\n#O..O..O.#\n#.OO.O.OO#\n#....O...#\n##########\n\n<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^\nvvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v\n><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<\n<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v^<<^^^^^\n^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><\n^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^\n>^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^\n<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>\n^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>\nv^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^"
	//raw = "########\n#..O.O.#\n##@.O..#\n#...O..#\n#.#.O..#\n#...O..#\n#......#\n########\n\n<^^>>>vv<v>>v<<"
	//raw = "########\n#..O.O.#\n##..O..#\n#...O..#\n#.#.O..#\n#...O..#\n#...@..#\n########\n\n^^^^^"
	//raw = "########\n#@.O.O.#\n##..O..#\n#...O..#\n#.#.O..#\n#...O..#\n#......#\n########\n\n>>>>"
	//raw = "########\n#..O.O@#\n##..O..#\n#...O..#\n#.#.O..#\n#...O..#\n#......#\n########\n\n<<<<"
	//raw = "########\n#..O@O.#\n##..O..#\n#...O..#\n#.#.O..#\n#...O..#\n#......#\n########\n\nvvvvv"

	area, instructions := readInput(raw)

	displayMap(area)

	var nArea [][]string
	for _, instruction := range instructions {
		nArea = moveBot(area, instruction)
		//fmt.Println(instruction)
		//displayMap(nArea)
		//fmt.Println()
	}

	var sumGPS int
	for y, line := range nArea {
		for x, cell := range line {
			if cell == "O" {
				sumGPS += x + y*100
			}
		}
	}

	fmt.Println("GPS:", sumGPS)
}
