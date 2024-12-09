package main

import (
	"fmt"
	"os"
	"slices"
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
			//fmt.Println("file : length:", l, " ID: ", fileID)
			fileID++
		case 1: // free space
			for j := 0; j < l; j++ {
				// "." = -1
				blocks = append(blocks, -1)
			}

			//fmt.Println("space: length:", l)
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

func part1(raw string) []int {
	blocks := diskMapToBlocks(raw)

	// maximum permutations
	sCount := getSpacesCount(blocks)

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
	return blocks
}

func part2displayBlocks(diskmap []struct {
	size int
	id   int
}) string {
	var output string
	for _, item := range diskmap {
		if item.id == -1 {
			for j := 0; j < item.size; j++ {
				output += "."
			}
		} else {
			for j := 0; j < item.size; j++ {
				output += strconv.Itoa(item.id)
			}
		}
	}

	return output
}

func part2(raw string) []int {
	var diskmap []struct {
		//index int
		size int
		id   int // if id=-1 => space
	}

	var fileID int
	for i, ru := range raw {
		l, err := strconv.Atoi(string(ru))
		if err != nil {
			panic(err)
		}

		currID := -1 // identify space
		if i%2 == 0 {
			currID = fileID
			fileID++
		}
		diskmap = append(diskmap, struct{ size, id int }{size: l, id: currID})
	}
	//fmt.Println(diskmap)

	var maxFileID int
	if diskmap[len(diskmap)-1].id > maxFileID {
		maxFileID = diskmap[len(diskmap)-1].id
	}
	if diskmap[len(diskmap)-2].id > maxFileID {
		maxFileID = diskmap[len(diskmap)-2].id
	}

	//fmt.Println(part2displayBlocks(diskmap))
	//fmt.Println()

	// for each file, last to first
	for curFileID := maxFileID; curFileID >= 0; curFileID-- {
		var file struct{ size, id int }
		fileIndex := -1
		for i, item := range diskmap {
			if item.id == curFileID {
				file = item
				fileIndex = i
			}
		}

		// for each space, first to "file index"
		for spaceIndex := 0; spaceIndex < fileIndex; spaceIndex++ {
			if diskmap[spaceIndex].id != -1 {
				continue
			}

			space := diskmap[spaceIndex]

			// if space can fit file
			if space.size >= file.size {
				// replace file with a space of same size
				diskmap[fileIndex].id = -1

				// put file before space
				diskmap[spaceIndex].id = file.id
				diskmap[spaceIndex].size = file.size
				//fmt.Println(part2displayBlocks(diskmap), " move file", file.id)

				// insert new space if needed
				if space.size != file.size {
					diskmap = slices.Insert(diskmap, spaceIndex+1, struct {
						size int
						id   int
					}{size: space.size - file.size, id: -1})
				}
				//fmt.Println(part2displayBlocks(diskmap), " fix space")

				break
			}
		}
	}

	//fmt.Println("\n00992111777.44.333....5555.6666.....8888..")
	//fmt.Println(part2displayBlocks(diskmap), "00992111777.44.333....5555.6666.....8888.." == part2displayBlocks(diskmap))

	var blocks []int
	for _, item := range diskmap {
		for i := 0; i < item.size; i++ {
			blocks = append(blocks, item.id)
		}
	}

	return blocks
}

func main() {
	r, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	raw := string(r)

	raw = "2333133121414131402"
	//raw = "12345"

	blocks := part1(raw)
	// checksum
	var checksum int
	for i, fileID := range blocks {
		if fileID == -1 {
			continue
		}
		checksum += fileID * i
	}
	fmt.Println("checksum (part1): ", checksum)

	blocks = part2(raw)

	// checksum
	checksum = 0
	for i, fileID := range blocks {
		if fileID == -1 {
			continue
		}
		checksum += fileID * i
	}
	fmt.Println("checksum (part2): ", checksum)
}
