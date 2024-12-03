package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	raw, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile(`mul\(([1-9][0-9]?[0-9]?),([1-9][0-9]?[0-9]?)\)`)
	matches := r.FindAllStringSubmatch(string(raw), -1)

	res := 0
	for _, match := range matches {
		l, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}
		r, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}
		res += l * r
	}

	log.Println(res)
}
