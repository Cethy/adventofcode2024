package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Area [][]rune

type Antenna struct {
	x, y int
	code rune
}

func getArea(raw string) Area {
	var area Area
	for _, line := range strings.Split(raw, "\n") {
		var larea []rune
		for _, r := range line {
			larea = append(larea, r)
		}
		area = append(area, larea)
	}
	return area
}

func getAntennas(area Area) []Antenna {
	var antennas []Antenna
	for y, a := range area {
		for x, ru := range a {
			if ru != '.' {
				antennas = append(antennas, Antenna{
					code: ru,
					x:    x,
					y:    y,
				})
			}
		}
	}

	return antennas
}

type AntennaGroups map[rune][]Antenna

// group antennas per code
func getAntennaGroups(antennas []Antenna) AntennaGroups {
	groups := make(map[rune][]Antenna)
	for _, a := range antennas {
		groups[a.code] = append(groups[a.code], a)
	}
	return groups
}

type AntennaPairs [][2]Antenna

// list all possible pair of same `code` antenna
func getAntennaPairs(groups AntennaGroups) AntennaPairs {
	var pairs [][2]Antenna
	for _, group := range groups {
		for _, a := range group {
			for _, b := range group {
				pair := [2]Antenna{a, b}
				if a == b || slices.ContainsFunc(pairs, func(p [2]Antenna) bool {
					return (pair[0] == p[0] && pair[1] == p[1]) || (pair[0] == p[1] && pair[1] == p[0])
				}) {
					continue
				}
				pairs = append(pairs, [2]Antenna{a, b})
			}
		}
	}
	return pairs
}

// using midpoint formula in reverse
// given (x1,y1), (x2,y2), (xm,ym) : x2 = 2*xm - x1 ; y2 = 2*ym - y1
func getAntiNodes(pairs AntennaPairs) [][2]int {
	var antiNodes [][2]int
	for _, pair := range pairs {
		// pair direction
		x2 := 2*pair[1].x - pair[0].x
		y2 := 2*pair[1].y - pair[0].y
		antiNodes = append(antiNodes, [2]int{x2, y2})
		// inverse
		x2 = 2*pair[0].x - pair[1].x
		y2 = 2*pair[0].y - pair[1].y
		antiNodes = append(antiNodes, [2]int{x2, y2})
	}

	return antiNodes
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "............\n........0...\n.....0......\n.......0....\n....0.......\n......A.....\n............\n............\n........A...\n.........A..\n............\n............"

	//raw = "..........\n..........\n..........\n....a.....\n........a.\n.....a....\n..........\n......A...\n..........\n.........."

	area := getArea(raw)
	antennas := getAntennas(area)

	groups := getAntennaGroups(antennas)
	pairs := getAntennaPairs(groups)
	antinodes := getAntiNodes(pairs)

	fmt.Println("area")
	for _, a := range area {
		for _, ru := range a {
			fmt.Print(string(ru))
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Println("antennas")
	for _, a := range antennas {
		fmt.Println(string(a.code), a.x, a.y)
	}
	fmt.Println()
	fmt.Println("groups")
	for ru, g := range groups {
		fmt.Print(string(ru))
		for _, a := range g {
			fmt.Print(" ", a.x, " ", a.y, " ")
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Println("pairs")
	for _, p := range pairs {
		fmt.Print(string(p[0].code))
		for _, a := range p {
			fmt.Print(" ", a.x, " ", a.y, " ")
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Println("antinodes")
	for _, p := range antinodes {
		fmt.Println("(", p[0], ";", p[1], ")")
	}
	fmt.Println("area w/ antinodes")

	var antinodesInArea [][2]int
	for y, a := range area {
		for x, ru := range a {
			if slices.Contains(antinodes, [2]int{x, y}) {
				antinodesInArea = append(antinodesInArea, [2]int{x, y})
			}
			if ru == '.' && slices.Contains(antinodes, [2]int{x, y}) {
				fmt.Print("#")
				continue
			}
			fmt.Print(string(ru))
		}
		fmt.Println()
	}

	fmt.Println("antinodesInArea:", len(antinodesInArea))
}
