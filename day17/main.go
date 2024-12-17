package main

import "C"
import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// readInput returns the registers and instruction list
func readInput(raw string) ([3]int, []int) {
	ra := strings.Split(raw, "\n\n")

	regRaw := strings.Split(ra[0], "\n")
	regA, err := strconv.Atoi(strings.Split(regRaw[0], ": ")[1])
	if err != nil {
		panic(err)
	}
	regB, err := strconv.Atoi(strings.Split(regRaw[1], ": ")[1])
	if err != nil {
		panic(err)
	}
	regC, err := strconv.Atoi(strings.Split(regRaw[2], ": ")[1])
	if err != nil {
		panic(err)
	}

	var instructions []int
	for _, v := range strings.Split(strings.Split(ra[1], ": ")[1], ",") {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		instructions = append(instructions, i)
	}

	return [3]int{
		regA,
		regB,
		regC,
	}, instructions
}

func operandComboValue(operand int, registers [3]int) int {
	if operand == 7 {
		panic("\nCombo operand 7 is reserved and will not appear in valid programs.")
	}

	// Combo operands 0 through 3 represent literal values 0 through 3.
	if operand <= 3 {
		return operand
	}

	// Combo operand 4 represents the value of register A.
	// Combo operand 5 represents the value of register B.
	// Combo operand 6 represents the value of register C.
	operand -= 4
	return registers[operand]
}

func adv(operand int, registers [3]int) int {
	// operand is combo
	operand = operandComboValue(operand, registers)

	// division
	// numerator is registerA value
	// denominator = 2^operand
	return registers[0] / int(math.Pow(2, float64(operand)))
}

func main() {

	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "Register A: 729\nRegister B: 0\nRegister C: 0\n\nProgram: 0,1,5,4,3,0"

	registers, instructions := readInput(raw)
	//fmt.Println(registers, instructions)

	var output []int
	instructionPointer := 0
	for instructionPointer < len(instructions) {
		instruction := instructions[instructionPointer]
		operand := instructions[instructionPointer+1]

		//fmt.Println(instructionPointer, registers, instruction, operand)
		//time.Sleep(1 * time.Second)

		switch instruction {
		case 1: // bxl
			// operand is a literal
			// bitwise XOR regB and literal operand
			registers[1] = registers[1] ^ operand
		case 2: // bst
			// operand is combo
			operand = operandComboValue(operand, registers)
			// modulo 8 of combo operand
			registers[1] = operand % 8
		case 3: // jnz
			// do nothing if regA is 0
			if registers[0] != 0 {
				// GOTO
				// change instructionPointer value
				instructionPointer = operand
				// and do not increase pointer for this iteration
				continue
			}

		case 4: // bxc
			// dont read operand
			// bitwise XOR regB and regC
			registers[1] = registers[1] ^ registers[2]
		case 5: // out
			// operand is combo
			operand = operandComboValue(operand, registers)
			// output
			output = append(output, operand%8)

		case 0: //adv
			registers[0] = adv(operand, registers)
		case 6: // bdv
			registers[1] = adv(operand, registers)
		case 7: // cdv
			registers[2] = adv(operand, registers)
		}

		// move to next instruction
		instructionPointer += 2
	}

	var fOutput string
	for i, v := range output {
		if i != 0 {
			fOutput += ","
		}
		fOutput += strconv.Itoa(v)
	}

	fmt.Println(fOutput)
	//fmt.Println(fOutput == "4,6,3,5,6,3,5,2,1,0")
}
