package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func generateArea(raw string) [][]int {
	var area [][]int
	for _, line := range strings.Split(raw, "\n") {
		var l []int
		for _, c := range line {
			if c == '.' {
				l = append(l, -1)
				continue
			}
			cell, err := strconv.Atoi(string(c))
			if err != nil {
				panic(err)
			}
			l = append(l, cell)
		}
		area = append(area, l)
	}
	return area
}

func getTrailheadsCoords(area [][]int) [][2]int {
	var coords [][2]int
	for y, line := range area {
		for x, c := range line {
			if c == 0 {
				coords = append(coords, [2]int{x, y})
			}
		}
	}
	return coords
}

type pathTree struct {
	x, y  int
	leafs []*pathTree
}

func newPathTree(x, y int) *pathTree {
	return &pathTree{x, y, []*pathTree{}}
}

func buildPathTreeLeafs(area [][]int, x, y int) *pathTree {
	branch := newPathTree(x, y)
	currHeight := area[y][x]
	// we're at the top
	if currHeight == 9 {
		return branch
	}

	// top
	if y > 0 && area[y-1][x] == currHeight+1 {
		//fmt.Println("top")
		branch.leafs = append(branch.leafs, buildPathTreeLeafs(area, x, y-1))
	}
	// down
	if y < len(area)-1 && area[y+1][x] == currHeight+1 {
		//fmt.Println("down")
		branch.leafs = append(branch.leafs, buildPathTreeLeafs(area, x, y+1))
	}
	// left
	if x > 0 && area[y][x-1] == currHeight+1 {
		//fmt.Println("left")
		branch.leafs = append(branch.leafs, buildPathTreeLeafs(area, x-1, y))
	}
	// right
	if x < len(area[0])-1 && area[y][x+1] == currHeight+1 {
		//fmt.Println("right")
		branch.leafs = append(branch.leafs, buildPathTreeLeafs(area, x+1, y))
	}

	return branch
}

func explorePaths(branch *pathTree) [][][2]int {
	var paths [][][2]int

	for _, b := range branch.leafs {
		lpaths := explorePaths(b)
		for _, lpath := range lpaths {
			paths = append(paths, slices.Insert(lpath, 0, [2]int{branch.x, branch.y}))
		}
	}
	if len(branch.leafs) == 0 {
		paths = append(paths, [][2]int{{branch.x, branch.y}})
	}

	return paths
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "...0...\n...1...\n...2...\n6543456\n7.....7\n8.....8\n9.....9"
	//raw = "..90..9\n...1.98\n...2..7\n6543456\n765.987\n876....\n987...."
	//raw = "10..9..\n2...8..\n3...7..\n4567654\n...8..3\n...9..2\n.....01"
	//raw = "89010123\n78121874\n87430965\n96549874\n45678903\n32019012\n01329801\n10456732"

	//raw = ".....0.\n..4321.\n..5..2.\n..6543.\n..7..4.\n..8765.\n..9...."
	//raw = "..90..9\n...1.98\n...2..7\n6543456\n765.987\n876....\n987...."
	//raw = "012345\n123456\n234567\n345678\n4.6789\n56789."

	area := generateArea(raw)

	for _, l := range area {
		fmt.Println(l)
	}
	fmt.Println()

	// find trailheads coords
	trailheads := getTrailheadsCoords(area)
	fmt.Println("trailheads:", trailheads)

	var trailsTree []*pathTree
	for _, coords := range trailheads {
		trailsTree = append(trailsTree, buildPathTreeLeafs(area, coords[0], coords[1]))
	}

	// all path to any summit
	var pathsToSummit [][][2]int
	for _, b := range trailsTree {
		tmp := explorePaths(b)
		for _, p := range tmp {
			if len(p) == 10 {
				pathsToSummit = append(pathsToSummit, p)
			}
		}
	}

	// trailheads to unique summits
	trailheadsToSummits := make([][][2][2]int, len(trailheads))
	for i, th := range trailheads {
		for _, p := range pathsToSummit {
			if p[0] == th {
				tuple := [2][2]int{th, p[len(p)-1]}
				if !slices.Contains(trailheadsToSummits[i], tuple) {
					trailheadsToSummits[i] = append(trailheadsToSummits[i], tuple)
				}
			}
		}
	}

	var score int
	for i := range trailheadsToSummits {
		score += len(trailheadsToSummits[i])
	}
	fmt.Println("score", score)

	// rating
	fmt.Println("rating", len(pathsToSummit))
}
