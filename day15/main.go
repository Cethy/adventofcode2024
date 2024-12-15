package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

func widerArea(rawArea string) [][]string {
	var area [][]string

	for _, line := range strings.Split(rawArea, "\n") {
		var l []string
		for _, cell := range line {
			c := string(cell)
			if c == "@" {
				l = append(l, "@", ".")
				continue
			}
			if c == "O" {
				l = append(l, "[", "]")
				continue
			}
			l = append(l, c, c)
		}
		area = append(area, l)
	}

	return area
}

func readInput(raw string, wide bool) ([][]string, []string) {
	ra := strings.Split(raw, "\n\n")
	var area [][]string
	var instructions []string
	if wide {
		area = widerArea(ra[0])
	} else {
		for _, line := range strings.Split(ra[0], "\n") {
			var l []string
			for _, cell := range line {
				l = append(l, string(cell))
			}
			area = append(area, l)
		}
	}

	for _, i := range strings.Replace(ra[1], "\n", "", -1) {
		instructions = append(instructions, string(i))
	}

	return area, instructions
}

func renderMap(area [][]string) string {
	var o string
	for _, line := range area {
		for _, cell := range line {
			o += cell
		}
		o += "\n"
	}
	return o
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

func mBot(area [][]string, vector [2]int) [][]string {
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
	// box or wide box
	if area[nextPos[1]][nextPos[0]] == "O" ||
		area[nextPos[1]][nextPos[0]] == "[" || area[nextPos[1]][nextPos[0]] == "]" {
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

func allMoveableBoxes(area [][]string, vector [2]int, boxSidesCoords [][2]int) ([][2]int, error) {
	//var o [][2]int
	// add boxes sides if any missed
	// and remove "." ? why is there any in the first place ?
	for _, coords := range boxSidesCoords {
		v := area[coords[1]][coords[0]]
		if v == "[" {
			if !slices.Contains(boxSidesCoords, [2]int{coords[0] + 1, coords[1]}) {
				boxSidesCoords = append(boxSidesCoords, [2]int{coords[0] + 1, coords[1]})
			}
		}
		if v == "]" {
			if !slices.Contains(boxSidesCoords, [2]int{coords[0] - 1, coords[1]}) {
				boxSidesCoords = append(boxSidesCoords, [2]int{coords[0] - 1, coords[1]})
			}
		}
	}

	var nextPosValues []string
	for _, coords := range boxSidesCoords {
		nextPosValues = append(nextPosValues, area[coords[1]+vector[1]][coords[0]+vector[0]])
	}
	if slices.Contains(nextPosValues, "#") {
		// cant move
		return boxSidesCoords, errors.New("cant move")
	}
	if slices.Contains(nextPosValues, "[") || slices.Contains(nextPosValues, "]") {
		var nextBoxesPos [][2]int
		for _, coords := range boxSidesCoords {
			if slices.Contains([]string{"[", "]"}, area[coords[1]][coords[0]]) {
				nextPos := [2]int{coords[0] + vector[0], coords[1] + vector[1]}
				if area[nextPos[1]][nextPos[0]] == "." {
					continue
				}
				nextBoxesPos = append(nextBoxesPos, nextPos)
			}
		}

		nBoxSidesCoords, err := allMoveableBoxes(area, vector, nextBoxesPos)
		if err != nil {
			return boxSidesCoords, err
		}
		return append(boxSidesCoords, nBoxSidesCoords...), nil
	}
	return boxSidesCoords, nil
}

func mBotWide(area [][]string, vector [2]int) [][]string {
	botPos := findBotPos(area)

	nextPos := [2]int{botPos[0] + vector[0], botPos[1] + vector[1]}

	// if nextPos not an obstacle
	if area[nextPos[1]][nextPos[0]] != "[" && area[nextPos[1]][nextPos[0]] != "]" {
		return mBot(area, vector)
	}

	// simple case, going < or >
	if vector[0] != 0 {
		// needs only 1 space
		firstEmptySpacePos := [2]int{-1, -1}
		for x := nextPos[0]; x > 0 && x < len(area[0]); x = x + vector[0] {
			// no empty space: bot not moving
			if area[nextPos[1]][x] == "#" {
				return area
			}
			// empty space found: bot not moving
			if area[nextPos[1]][x] == "." {
				firstEmptySpacePos = [2]int{x, nextPos[1]}
				break
			}
		}

		// move boxes
		for i := abs(nextPos[0] - firstEmptySpacePos[0]); i > 0; i-- {
			area[nextPos[1]][nextPos[0]+i*vector[0]] = area[nextPos[1]][nextPos[0]+i*vector[0]-1*vector[0]]
		}
		area[nextPos[1]][nextPos[0]] = "."
		// move bot
		area[nextPos[1]][nextPos[0]], area[botPos[1]][botPos[0]] = area[botPos[1]][botPos[0]], area[nextPos[1]][nextPos[0]]
	} else {
		// complex case, when bot going ^ or v
		// needs <= 2*n empty spaces on last row where n is nb of boxes involved on the last row

		boxSidesCoords, err := allMoveableBoxes(area, vector, [][2]int{nextPos})
		if err != nil {
			return area
		}

		// move boxes
		slices.Reverse(boxSidesCoords)
		for _, coords := range boxSidesCoords {
			area[coords[1]+vector[1]][coords[0]+vector[0]], area[coords[1]][coords[0]] = area[coords[1]][coords[0]], area[coords[1]+vector[1]][coords[0]+vector[0]]
		}
		// move bot
		area[nextPos[1]][nextPos[0]], area[botPos[1]][botPos[0]] = area[botPos[1]][botPos[0]], area[nextPos[1]][nextPos[0]]
	}

	return area
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func moveBot(area [][]string, instruction string, wide bool) [][]string {
	var vector [2]int
	switch instruction {
	case "^":
		vector = [2]int{0, -1}
	case ">":
		vector = [2]int{1, 0}
	case "v":
		vector = [2]int{0, 1}
	case "<":
		vector = [2]int{-1, 0}
	}

	if wide {
		return mBotWide(area, vector)
	}
	return mBot(area, vector)
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
	//raw = "#######\n#...#.#\n#.....#\n#..OO@#\n#..O..#\n#.....#\n#######\n\n<vv<<^^<<^^"
	//raw = "#######\n#...#.#\n#.....#\n#..OO@#\n#..O..#\n#.....#\n#######\n\n<"
	//raw = "#######\n#..OO@#\n#..O..#\n#.....#\n#######\n\n<<<<<<<<"
	//raw = "#######\n#...#.#\n#.....#\n#..OO@#\n#..O..#\n#.....#\n#######\n\nv<^^^"
	//raw = "#######\n#...#.#\n#.....#\n#.@OO.#\n#..O..#\n#.....#\n#######\n\n>>>>"
	//raw = "#######\n#...#.#\n#.....#\n#..OO@#\n#..O..#\n#.....#\n#######\n\n<<<<<"
	//raw = "#######\n#...#.#\n#.....#\n#..O..#\n#..OO@#\n#.....#\n#.....#\n#######\n\n<^^<<vvvvvv"

	//raw = "################\n#..............#\n#.....O........#\n#.....OO@......#\n#....O.O.......#\n#..............#\n#..............#\n#..............#\n################\n\n<^^<<vvvv"
	//raw = "################\n#..............#\n#.....O........#\n#.....OO@......#\n#....O.O.......#\n#..............#\n#.....O........#\n#..............#\n#..............#\n#..............#\n################\n\n<^^<<vvvv"
	//raw = "################\n#..............#\n#.....O........#\n#.....OO@......#\n#....O.O.......#\n#..............#\n#.....#........#\n#..............#\n#..............#\n#..............#\n################\n\n<^^<<vvvv"
	area, instructions := readInput(raw, false)

	//fmt.Println(renderMap(area))

	var nArea [][]string

	for _, instruction := range instructions {
		nArea = moveBot(area, instruction, false)
		//fmt.Println(instruction)
		//fmt.Println(renderMap(nArea))
		//fmt.Println()
	}
	fmt.Println(renderMap(nArea))

	var sumGPS int
	for y, line := range nArea {
		for x, cell := range line {
			if cell == "O" {
				sumGPS += x + y*100
			}
		}
	}

	fmt.Println("GPS:", sumGPS)
	fmt.Println()

	area, instructions = readInput(raw, true)

	//fmt.Print("\033[2J") // clear screen
	for _, instruction := range instructions {
		nArea = moveBot(area, instruction, true)
		//fmt.Print("\033[H") // move cursor to top-left corner
		//fmt.Println(instruction)
		//fmt.Println(renderMap(nArea))
		//fmt.Println()
		//time.Sleep(1 * time.Second / 100)
	}
	//os.Exit(0)
	fmt.Println(renderMap(nArea))

	sumGPS = 0
	for y, line := range nArea {
		for x, cell := range line {
			if cell == "[" {
				sumGPS += x + y*100
			}
		}
	}

	fmt.Println("GPSwide:", sumGPS)
}
