package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func splitRaw(raw string) (string, string) {
	pageOrderingRules := ""
	updates := ""
	half := 0
	for _, line := range strings.Split(raw, "\n") {
		if line == "" {
			half++
			continue
		}
		switch half {
		case 0:
			pageOrderingRules += line + "\n"
		case 1:
			updates += line + "\n"
		}
	}
	return pageOrderingRules[:len(pageOrderingRules)-1], updates[:len(updates)-1]
}

func extractPageOrderingRules(raw string) [][2]int {
	var o [][2]int
	for _, line := range strings.Split(raw, "\n") {
		xy := strings.Split(line, "|")
		if len(xy) != 2 {
			panic("xy size is wrong for line: " + line)
		}
		x, err := strconv.Atoi(xy[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(xy[1])
		if err != nil {
			panic(err)
		}

		o = append(o, [2]int{x, y})
	}
	return o
}

func extractUpdates(raw string) [][]int {
	var o [][]int
	for _, line := range strings.Split(raw, "\n") {
		var s []int
		for _, cell := range strings.Split(line, ",") {
			x, err := strconv.Atoi(cell)
			if err != nil {
				panic(err)
			}
			s = append(s, x)
		}
		o = append(o, s)
	}
	return o
}

func fixUpdate(update []int, rule [2]int) []int {
	o := make([]int, len(update))
	for i := 0; i < len(update); i++ {
		o[i] = update[i]
	}
	xi := slices.Index(update, rule[0])
	yi := slices.Index(update, rule[1])
	o[xi] = update[yi]
	o[yi] = update[xi]
	return o
}

func fixUpdateOrder(update []int, rules [][2]int) []int {
	err := checkUpdateOrder(update, rules)
	if err, ok := err.(*UpdateOrderError); ok {
		return fixUpdateOrder(fixUpdate(update, err.rule), rules)
	}

	return update
}

type UpdateOrderError struct {
	rule [2]int
}

func (e *UpdateOrderError) Error() string {
	return "UpdateOrderError"
}

func checkUpdateOrder(update []int, rules [][2]int) error {
	for _, pageNumber := range update {
		// find all rules related to pageNumber
		// and keep only ones related to other numbers in the update
		var relatedRules [][2]int
		for _, pageOrderingRule := range rules {
			if slices.Contains(relatedRules, pageOrderingRule) {
				continue
			}
			if pageOrderingRule[0] == pageNumber && slices.Contains(update, pageOrderingRule[1]) {
				relatedRules = append(relatedRules, pageOrderingRule)
			}
			if pageOrderingRule[1] == pageNumber && slices.Contains(update, pageOrderingRule[0]) {
				relatedRules = append(relatedRules, pageOrderingRule)
			}
		}

		// check the rules
		for _, rule := range relatedRules {
			x := rule[0]
			y := rule[1]
			xi := slices.Index(update, x)
			yi := slices.Index(update, y)
			if xi == -1 || yi == -1 {
				panic("x or y not found in the update")
			}
			if xi >= yi {
				fmt.Println("update", update, "is not in order ; rule", rule)
				return &UpdateOrderError{rule}
			}
		}
	}
	return nil
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "47|53\n97|13\n97|61\n97|47\n75|29\n61|13\n75|53\n29|13\n97|29\n53|29\n61|53\n97|53\n61|29\n47|13\n75|47\n97|75\n47|61\n75|61\n47|29\n75|13\n53|13\n\n75,47,61,53,29\n97,61,53,29,13\n75,29,13\n75,97,47,61,53\n61,13,29\n97,13,75,29,47"
	//raw = "47|53\n97|13\n97|61\n97|47\n75|29\n61|13\n75|53\n29|13\n97|29\n53|29\n61|53\n97|53\n61|29\n47|13\n75|47\n97|75\n47|61\n75|61\n47|29\n75|13\n53|13\n\n75,97,47,61,53\n61,13,29\n97,13,75,29,47"

	rawPageOrderingRules, rawUpdates := splitRaw(raw)
	pageOrderingRules := extractPageOrderingRules(rawPageOrderingRules)
	updates := extractUpdates(rawUpdates)

	//fmt.Println(pageOrderingRules)
	//fmt.Println(updates)

	var middlePagesOfCorrectUpdates []int
	var middlePagesOfCorrectedUpdates []int
	// determine if updates are in the right order
	for _, update := range updates {
		err := checkUpdateOrder(update, pageOrderingRules)
		if err != nil {
			fixedUpdate := fixUpdateOrder(update, pageOrderingRules)
			middlePagesOfCorrectedUpdates = append(middlePagesOfCorrectedUpdates, fixedUpdate[len(fixedUpdate)/2])
			continue
		}
		middlePagesOfCorrectUpdates = append(middlePagesOfCorrectUpdates, update[len(update)/2])
	}

	var sumCorrectUpdates int
	for _, pageNumber := range middlePagesOfCorrectUpdates {
		sumCorrectUpdates += pageNumber
	}
	fmt.Println("correct:", sumCorrectUpdates)

	var sumFixedUpdates int
	for _, pageNumber := range middlePagesOfCorrectedUpdates {
		sumFixedUpdates += pageNumber
	}
	fmt.Println("fixed:  ", sumFixedUpdates)
}
