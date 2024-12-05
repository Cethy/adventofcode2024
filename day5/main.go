package main

import (
	"log"
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

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "47|53\n97|13\n97|61\n97|47\n75|29\n61|13\n75|53\n29|13\n97|29\n53|29\n61|53\n97|53\n61|29\n47|13\n75|47\n97|75\n47|61\n75|61\n47|29\n75|13\n53|13\n\n75,47,61,53,29\n97,61,53,29,13\n75,29,13\n75,97,47,61,53\n61,13,29\n97,13,75,29,47"
	//raw = "47|53\n97|13\n97|61\n97|47\n75|29\n61|13\n75|53\n29|13\n97|29\n53|29\n61|53\n97|53\n61|29\n47|13\n75|47\n97|75\n47|61\n75|61\n47|29\n75|13\n53|13\n\n75,47,61,53,29\n97,61,53,29,13\n75,29,13"

	rawPageOrderingRules, rawUpdates := splitRaw(raw)
	pageOrderingRules := extractPageOrderingRules(rawPageOrderingRules)
	updates := extractUpdates(rawUpdates)

	//fmt.Println(pageOrderingRules)
	//fmt.Println(updates)

	var middlePagesOfCorrectUpdates []int
	// determine if updates are in the right order
mainLoop:
	for _, update := range updates {
		for _, pageNumber := range update {
			// find all rules related to pageNumber
			// and keep only ones related to other numbers in the update
			var relatedRules [][2]int
			for _, pageOrderingRule := range pageOrderingRules {
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
					//fmt.Println("update", update, "is not in order ; rule", rule)
					continue mainLoop
					//panic("err")
				}
			}
		}

		middlePagesOfCorrectUpdates = append(middlePagesOfCorrectUpdates, update[len(update)/2])
	}

	var sum int
	for _, pageNumber := range middlePagesOfCorrectUpdates {
		sum += pageNumber
	}

	log.Println(sum)
}
