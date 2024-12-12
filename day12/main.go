package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

/*
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
*/

// given an initial position, find all neighbours non-explored coords contained in same region
func exploreRegion(matrix [][]rune, traversedMatrix *[][2]int, x, y int) [][2]int {
	var regionCoords [][2]int

	regionCoords = append(regionCoords, [2]int{x, y})
	*traversedMatrix = append(*traversedMatrix, [2]int{x, y})

	// for each direction
	// if neighbour not OOB AND neighbour is same rune AND neighbour not explored
	ru := matrix[y][x]
	// top
	if y > 0 && matrix[y-1][x] == ru && !slices.Contains(*traversedMatrix, [2]int{x, y - 1}) {
		//fmt.Println("top")
		regionCoords = append(regionCoords, exploreRegion(matrix, traversedMatrix, x, y-1)...)
	}
	// down
	if y < len(matrix)-1 && matrix[y+1][x] == ru && !slices.Contains(*traversedMatrix, [2]int{x, y + 1}) {
		//fmt.Println("down")
		regionCoords = append(regionCoords, exploreRegion(matrix, traversedMatrix, x, y+1)...)
	}
	// left
	if x > 0 && matrix[y][x-1] == ru && !slices.Contains(*traversedMatrix, [2]int{x - 1, y}) {
		//fmt.Println("left")
		regionCoords = append(regionCoords, exploreRegion(matrix, traversedMatrix, x-1, y)...)
	}
	// right
	if x < len(matrix[0])-1 && matrix[y][x+1] == ru && !slices.Contains(*traversedMatrix, [2]int{x + 1, y}) {
		//fmt.Println("right")
		regionCoords = append(regionCoords, exploreRegion(matrix, traversedMatrix, x+1, y)...)
	}

	return regionCoords
}

// extractRegionMeasurements returns measures (area, perimeter) and add explored coordinates to traversedMatrix
func extractRegionMeasurements(matrix [][]rune, traversedMatrix *[][2]int, x, y int) [2]int {
	regionCoordinates := exploreRegion(matrix, traversedMatrix, x, y)
	// The perimeter of a region is
	// the number of sides of garden plots in the region
	// that do not touch another garden plot in the same region.
	perimeter := 0
	for _, r := range regionCoordinates {
		peri := 4
		for _, neighbour := range [][2]int{
			{r[0] - 1, r[1]},
			{r[0] + 1, r[1]},
			{r[0], r[1] - 1},
			{r[0], r[1] + 1},
		} {
			if slices.Contains(regionCoordinates, neighbour) {
				peri -= 1
			}
		}
		perimeter += peri
	}
	//fmt.Println(string(matrix[y][x]), regionCoordinates, "area:", len(regionCoordinates), "perimeter:", perimeter)

	return [2]int{len(regionCoordinates), perimeter}
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "AAAA\nBBCD\nBBCC\nEEEC"
	//raw = "OOOOO\nOXOXO\nOOOOO\nOXOXO\nOOOOO"
	//raw = "RRRRIICCFF\nRRRRIICCCF\nVVRRRCCFFF\nVVRCCCJFFF\nVVVVCJJCFE\nVVIVCCJJEE\nVVIIICJJEE\nMIIIIIJJEE\nMIIISIJEEE\nMMMISSJEEE"

	var matrix [][]rune
	for _, line := range strings.Split(raw, "\n") {
		var l []rune
		for _, ru := range line {
			l = append(l, ru)
		}
		matrix = append(matrix, l)
	}
	for _, l := range matrix {
		for _, ru := range l {
			fmt.Print(string(ru))
		}
		fmt.Println()
	}
	fmt.Println()

	// list regions [][area, perimeter]
	var regions [][2]int
	// list already explored coordinates in the matrix [][x,y]
	traversedMatrix := make([][2]int, 0)
	for y, l := range matrix {
		for x := range l {
			// @todo might be more convenient to create a flat array of all coords and remove them to avoid looping on nothing ?
			// if already explored, ignore the cell
			if slices.Contains(traversedMatrix, [2]int{x, y}) {
				continue
			}
			// new cell = new region
			region := extractRegionMeasurements(matrix, &traversedMatrix, x, y)
			regions = append(regions, region)
			//traversedMatrix = append(traversedMatrix, traversedRegion...)
		}
	}

	// pricing
	var pricing int
	for _, re := range regions {
		pricing += re[0] * re[1]
	}
	fmt.Println("pricing:", pricing)
}
