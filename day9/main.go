package main

import (
	"fmt"
	"os"
	"strconv"
)

func diskMapToBlocks(raw string) []int {
	var blocks []int
	var fileID int
	for i, ru := range raw {
		l, err := strconv.Atoi(string(ru))
		if err != nil {
			panic(err)
		}

		switch i % 2 {
		case 0: // file
			for j := 0; j < l; j++ {
				blocks = append(blocks, fileID)
			}
			fmt.Println("file : length:", l, " ID: ", fileID)
			fileID++
		case 1: // free space
			for j := 0; j < l; j++ {
				// "." = -1
				blocks = append(blocks, -1)
			}

			fmt.Println("space: length:", l)
		}
	}

	return blocks
}

func getSpacesCount(blocks []int) int {
	var sCount int
	for _, b := range blocks {
		if b == -1 { // free space
			sCount++
		}
	}

	return sCount
}

func lastFileBlockIndex(blocks []int) int {
	for i := len(blocks) - 1; i >= 0; i-- {
		if blocks[i] != -1 {
			return i
		}
	}
	return -1
}
func firstSpaceBlockIndex(blocks []int) int {
	for i := 0; i < len(blocks); i++ {
		if blocks[i] == -1 {
			return i
		}
	}
	return -1
}

// don't work for ids > 9
func displayBlocks(blocks []int) string {
	var output string
	for _, block := range blocks {
		if block == -1 {
			output += "."
			continue
		}
		if block > 9 {
			panic("cant use displayBlocks with blockIds > 9")
		}
		output += strconv.Itoa(block)
	}

	return output
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	//raw = "2333133121414131402"
	//raw = "12345"

	blocks := diskMapToBlocks(raw)

	//fmt.Println(displayBlocks(blocks), "00...111...2...333.44.5555.6666.777.888899" == displayBlocks(blocks))
	//fmt.Println(displayBlocks(blocks), "0..111....22222" == displayBlocks(blocks))

	//moving file blocks

	// maximum permutations
	sCount := getSpacesCount(blocks)
	fmt.Println(sCount)

	for i := 0; i < sCount; i++ {
		lfi := lastFileBlockIndex(blocks)
		fsi := firstSpaceBlockIndex(blocks)

		if lfi < fsi {
			break
		}

		blocks[lfi], blocks[fsi] = blocks[fsi], blocks[lfi]

		//fmt.Println(displayBlocks(blocks), "0099811188827773336446555566.............." == displayBlocks(blocks))
		//fmt.Println(displayBlocks(blocks), "022111222......" == displayBlocks(blocks))
	}

	// checksum
	var checksum int
	for i, fileID := range blocks {
		if fileID == -1 {
			continue
		}
		checksum += fileID * i
	}
	fmt.Println("checksum: ", checksum)
}
