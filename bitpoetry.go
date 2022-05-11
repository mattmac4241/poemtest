package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
Finite Poetry Protocol (from the docs, PDF, page 42)

| Poem (variable size)                                                       |
+-------------+-------------+-------------+-----------------+----------------+
| PoetryLine1 | PoetryLine2 | PoetryLine3 | PoetryLineN ... | EndPoem (0x00) |

| PoetryLine  (3 bytes)                                    |
+--------------------+-------------+-----------------------+
| PartOfSpeech uint8 | Count uint8 | DictionaryIndex uint8 |

| PartOfSpeech (1 byte)  |
+------------------------+
| Part      |   uint8    |
+-----------+------------+
| Verb      |    0x01    |
| Noun      |    0x02    |
| Adjective |    0x03    |

Pseudo-BNF:
PoetryLine = PartOfSpeech uint8 | Count uint8 | DictionaryIndex uint8
Poem = PoetryLine* | endOfPoem
*/

type partOfSpeech int

const (
	endOfPoem partOfSpeech = iota
	verb
	noun
	adjective
)

var verbs = []string{"jump", "dance", "scream"}
var nouns = []string{"fish", "bear", "taco"}
var adjectives = []string{"blue", "tasty", "smelly"}
var parts = [][]string{verbs, nouns, adjectives}

var poem = []byte{0x01, 0xA0, 0x00, 0x03, 0x02, 0x02, 0x02, 0x01, 0x00, 0x00}

func main() {
	blocks, err := getBlocks(poem)
	if err != nil {
		panic(fmt.Sprintf("failed to convert blocks: %s\n", err))
	}

	byteBlocks := []byte(blocks)
	file, err := os.Create("./output.txt")
	if err != nil {
		panic(fmt.Sprintf("failed to create file: %s\n", err))
	}
	defer file.Close()

	_, err = file.Write(byteBlocks)
	if err != nil {
		panic(fmt.Sprintf("failed to write file: %s\n", err))
	}
}

func getBlocks(poem []byte) (string, error) {
	startBlock := 0
	endBlock := 2
	blocks := ""

	if len(poem) < 3 {
		return "", errors.New("invalid blocks size")
	}

	for endBlock < len(poem) {
		block := poem[startBlock : endBlock+1]
		convertedBlock, err := convertBlock(block)
		if err != nil {
			return "", err
		}

		blocks += fmt.Sprintf("%s\n", convertedBlock)

		startBlock += 3
		endBlock += 3
	}

	return blocks, nil
}

// convert byte block into poem
func convertBlock(block []byte) (string, error) {
	if len(block) != 3 {
		return "", errors.New("invalid block size, must be exactly 3 bytes")
	}

	partOfSpeech := block[0] - 1 // because count starts at 1 for byte
	index := block[2]
	countByte := block[1]

	// because count is a byte we need to convert to a
	// string and then finally to an int
	countString := fmt.Sprintf("%v", countByte)
	count, err := strconv.Atoi(countString)
	if err != nil {
		return "", err
	}

	wordType := parts[partOfSpeech]
	word := wordType[index]
	words := []string{}

	for i := 0; i < count; i++ {
		words = append(words, word)
	}

	return strings.Join(words, ", "), nil
}
