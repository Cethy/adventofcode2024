package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type InputEquations struct {
	result     int
	testValues []int
}

func getInputEquations(raw string) []InputEquations {
	eqs := make([]InputEquations, 0)
	for _, rawEq := range strings.Split(raw, "\n") {
		els := strings.Split(rawEq, ": ")
		result, err := strconv.Atoi(els[0])
		if err != nil {
			panic(err)
		}
		var testValues []int
		for _, el := range strings.Split(els[1], " ") {
			num, err := strconv.Atoi(el)
			if err != nil {
				panic(err)
			}
			testValues = append(testValues, num)
		}
		eqs = append(eqs, InputEquations{result, testValues})
	}

	return eqs
}

var operators = []rune{'+', '*'}

func generateOperatorCombinations(testValueCpt int) [][]rune {
	// (n*operators)^(testValueCpt-1) possibilities
	// 2^(testValueCpt-1)
	opCombinations := make([][]rune, 0)
	combinationCount := int(math.Pow(2, float64(testValueCpt-1)))
	format := "%." + strconv.Itoa(testValueCpt-1) + "b"

	for i := 0; i < combinationCount; i++ {
		bin := fmt.Sprintf(format, i)
		var combination []rune
		for _, bit := range bin {
			sBit, err := strconv.Atoi(string(bit))
			if err != nil {
				panic(err)
			}
			combination = append(combination, operators[sBit])
		}
		opCombinations = append(opCombinations, combination)
	}
	return opCombinations
}

func isSolvableEquation(eq InputEquations) bool {
	opCombinations := generateOperatorCombinations(len(eq.testValues))
	for _, combination := range opCombinations {
		combResult := eq.testValues[0]
		for i := 1; i < len(eq.testValues); i++ {
			switch combination[i-1] {
			case '+':
				combResult += eq.testValues[i]
			case '*':
				combResult *= eq.testValues[i]
			}
		}
		if combResult == eq.result {
			return true
		}
	}
	return false
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "190: 10 19\n3267: 81 40 27\n83: 17 5\n156: 15 6\n7290: 6 8 6 15\n161011: 16 10 13\n192: 17 8 14\n21037: 9 7 18 13\n292: 11 6 16 20"
	//raw = "190: 10 19\n3267: 81 40 27\n292: 11 6 16 20"

	inputEqs := getInputEquations(raw)

	solvedEqs := make([]InputEquations, 0)
	for _, eq := range inputEqs {
		if isSolvableEquation(eq) {
			solvedEqs = append(solvedEqs, eq)
		}
	}

	var sum int
	for _, eq := range solvedEqs {
		sum += eq.result
	}
	fmt.Println(sum)
}
