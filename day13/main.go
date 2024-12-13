package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ClawMachine struct {
	buttonA, buttonB [2]int
	prize            [2]int
}

func readLineValues(line string) [2]int {
	splitTmp := strings.Split(line, ":")
	rawV := strings.TrimSpace(splitTmp[1])
	splitV := strings.Split(rawV, ", ")
	vX, err := strconv.Atoi(splitV[0][2:])
	if err != nil {
		panic(err)
	}
	vY, err := strconv.Atoi(splitV[1][2:])
	if err != nil {
		panic(err)
	}
	return [2]int{vX, vY}
}

func readInput(raw string) []ClawMachine {
	var claws []ClawMachine
	for _, block := range strings.Split(raw, "\n\n") {
		lines := strings.Split(block, "\n")
		claws = append(claws, ClawMachine{
			buttonA: readLineValues(lines[0]),
			buttonB: readLineValues(lines[1]),
			prize:   readLineValues(lines[2]),
		})
	}

	return claws
}

func gcd(a, b int) int {
	if a == 0 {
		return b
	}
	return gcd(b%a, a)
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "Button A: X+94, Y+34\nButton B: X+22, Y+67\nPrize: X=8400, Y=5400\n\nButton A: X+26, Y+66\nButton B: X+67, Y+21\nPrize: X=12748, Y=12176\n\nButton A: X+17, Y+86\nButton B: X+84, Y+37\nPrize: X=7870, Y=6450\n\nButton A: X+69, Y+23\nButton B: X+27, Y+71\nPrize: X=18641, Y=10279"
	//raw = "Button A: X+94, Y+34\nButton B: X+22, Y+67\nPrize: X=8400, Y=5400"

	var totalTokens int
	clawMachines := readInput(raw)
	for _, claw := range clawMachines {
		//fmt.Println(claw)
		gcdX := gcd(claw.buttonA[0], claw.buttonB[0])
		gcdY := gcd(claw.buttonA[1], claw.buttonB[1])

		// claw no solution
		if claw.prize[0]%gcdX+claw.prize[1]%gcdY > 0 {
			continue
		}

		if claw.buttonA[0]+claw.buttonB[0] == claw.prize[0] ||
			claw.buttonA[1]+claw.buttonB[1] == claw.prize[1] {
			panic("one only solution is not supported")
		}
		// just in case
		if claw.buttonA[0]+claw.buttonB[0] > claw.prize[0] ||
			claw.buttonA[1]+claw.buttonB[1] > claw.prize[1] {
			panic("A+B > C")
		}

		// now, bruteforce ! :D
		var solutions [][2]int
		for pushA := 0; pushA <= 100; pushA++ {
			for pushB := 0; pushB <= 100; pushB++ {
				if claw.prize[0] < claw.buttonA[0]*pushA+claw.buttonB[0]*pushB ||
					claw.prize[1] < claw.buttonA[1]*pushA+claw.buttonB[1]*pushB {
					// too far
					break
				}
				if claw.prize[0] == claw.buttonA[0]*pushA+claw.buttonB[0]*pushB &&
					claw.prize[1] == claw.buttonA[1]*pushA+claw.buttonB[1]*pushB {
					// solution
					solutions = append(solutions, [2]int{pushA, pushB})
				}
			}

			if claw.prize[0] < claw.buttonA[0]*pushA ||
				claw.prize[1] < claw.buttonA[1] {
				// too far
				break
			}
		}

		if len(solutions) == 0 {
			continue
		}

		bestSolution := solutions[0]
		for _, solution := range solutions {
			cost := solution[0]*3 + solution[1]
			//fmt.Println(solution, cost)

			if cost < bestSolution[0]*3+bestSolution[1] {
				bestSolution = solution
			}
		}
		//fmt.Println("bestSolution:", bestSolution)
		totalTokens += bestSolution[0]*3 + bestSolution[1]
	}

	fmt.Println("totalTokens:", totalTokens)
}
