package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Robot struct {
	position [2]int
	velocity [2]int
}

func readInput(raw string) []Robot {
	var bots []Robot
	for _, line := range strings.Split(raw, "\n") {
		i := strings.Split(line, " ")
		rawPos := strings.Split(i[0][2:], ",")
		posX, err := strconv.Atoi(rawPos[0])
		if err != nil {
			panic(err)
		}
		posY, err := strconv.Atoi(rawPos[1])
		if err != nil {
			panic(err)
		}
		rawVel := strings.Split(i[1][2:], ",")
		velX, err := strconv.Atoi(rawVel[0])
		if err != nil {
			panic(err)
		}
		velY, err := strconv.Atoi(rawVel[1])
		if err != nil {
			panic(err)
		}

		bots = append(bots, Robot{
			position: [2]int{posX, posY},
			velocity: [2]int{velX, velY},
		})
	}
	return bots
}

func renderMap(w, h int, bots []Robot) string {
	var m string
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			botCpt := 0
			for bot := range bots {
				if bots[bot].position[0] == i && bots[bot].position[1] == j {
					botCpt++
				}
			}
			if botCpt > 0 {
				m += strconv.Itoa(botCpt)
				continue
			}
			m += "."
		}
		m += "\n"
	}
	return m
}
func renderMapWOQuadrant(w, h int, bots []Robot) string {
	var m string
	for j := 0; j < h; j++ {
		if j == h/2 {
			m += "\n"
			continue
		}
		for i := 0; i < w; i++ {
			if i == w/2 {
				m += " "
				continue
			}

			botCpt := 0
			for bot := range bots {
				if bots[bot].position[0] == i && bots[bot].position[1] == j {
					botCpt++
				}
			}
			if botCpt > 0 {
				m += strconv.Itoa(botCpt)
				continue
			}
			m += "."
		}
		m += "\n"
	}
	return m
}

func getQuadrantBotCounts(w, h int, bots []Robot) [4]int {
	var botCounts [4]int
	for bot := range bots {
		if bots[bot].position[0] < w/2 {
			if bots[bot].position[1] < h/2 {
				botCounts[0]++
			} else if bots[bot].position[1] > h/2 {
				botCounts[1]++
			}
		} else if bots[bot].position[0] > w/2 {
			if bots[bot].position[1] < h/2 {
				botCounts[2]++
			} else if bots[bot].position[1] > h/2 {
				botCounts[3]++
			}
		}
	}
	return botCounts
}

func moveBot(bot Robot, w, h, times int) Robot {
	x := (bot.position[0] + bot.velocity[0]*times) % w
	if x < 0 {
		x = w + x
	}
	y := (bot.position[1] + bot.velocity[1]*times) % h
	if y < 0 {
		y = h + y
	}
	return Robot{
		position: [2]int{x, y},
		velocity: bot.velocity,
	}
}

func moveBots(bots []Robot, w, h, times int) []Robot {
	movedBots := make([]Robot, len(bots))
	for i, bot := range bots {
		movedBots[i] = moveBot(bot, w, h, times)
	}
	return movedBots
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "p=0,4 v=3,-3\np=6,3 v=-1,-3\np=10,3 v=-1,2\np=2,0 v=2,-1\np=0,0 v=1,3\np=3,0 v=-2,-2\np=7,6 v=-1,-3\np=3,0 v=-1,-2\np=9,3 v=2,3\np=7,3 v=-1,2\np=2,4 v=2,-3\np=9,5 v=-3,-3"
	//raw = "p=2,4 v=2,-3"
	//raw = "p=0,0 v=4,0"

	bots := readInput(raw)
	/*for _, bot := range bots {
		fmt.Println(bot)
	}*/
	w := 101 //11
	h := 103 //7
	var movedBots []Robot
	/*
		for i := 0; i <= 100; i++ {
			movedBots = moveBots(bots, w, h, i)
			fmt.Println(renderMap(w, h, movedBots))
		}
		fmt.Println(renderMap(w, h, movedBots) == "......2..1.\n...........\n1..........\n.11........\n.....1.....\n...12......\n.1....1....\n")

		fmt.Println(renderMapWOQuadrant(w, h, movedBots))
	*/
	movedBots = moveBots(bots, w, h, 100)

	botCounts := getQuadrantBotCounts(w, h, movedBots)
	safetyFactor := 1
	for _, count := range botCounts {
		safetyFactor *= count
	}

	fmt.Println("safety factor:", safetyFactor)
}
