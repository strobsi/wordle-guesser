package game

import (
	"fmt"
	"sort"
	"strings"
)

// Returns the string representation of the word
func (w Word) String() string {
	return string(w.Characters)
}

// Checks if all the characters given have the right positions
func (w Word) HasRightPositions(chars []Result) bool {
	hasRightPositions := true
	for _, v := range chars {
		for index, x := range w.Characters {
			if x == v.Char && index == v.Position {
				hasRightPositions = false
			}
		}
	}
	return hasRightPositions
}

// HasAllCharacters checks if all characters given as parameter are within the word
func (w Word) HasAllCharacters(chars []Result) bool {
	hasAll := true
	for _, v := range chars {
		if !strings.Contains(w.String(), string(v.Char)) {
			hasAll = false
		}
	}
	return hasAll
}

// HasOneOfChars checks if one of the given characters are within the word
func (w Word) HasOneOfChars(chars []Result) bool {
	contains := false
	for _, v := range chars {
		if strings.Contains(w.String(), string(v.Char)) {
			contains = true
		}
	}
	return contains
}

// Print prints out - if not in silent mode
func (g *Game) Print(msg ...interface{}) {
	if !g.SilentMode {
		fmt.Println(msg...)
	}
}

// CalculateFreq calculates the frequency of the letters in an slice of words
func CalculateFreq(words []Word) FrequencySlice {
	// Get frequency for first character
	characterFreq := FrequencyMap{}
	// Get length of words to stay dynamic
	for i := 0; i < len(words[0].Characters); i++ {
		posFreq := map[byte]int{}
		characterFreq.Characters = append(characterFreq.Characters, PositionFrequency{Frequency: posFreq})
	}
	for _, v := range words {
		for i := 0; i < len(characterFreq.Characters); i++ {
			characterFreq.Characters[i].Frequency[v.Characters[i]]++
		}
	}
	freqArr := FrequencySlice{}
	for _, v := range characterFreq.Characters {

		posFreqArr := PositionFrequencySlice{}

		var ss []kv
		for k, val := range v.Frequency {
			ss = append(ss, kv{k, val})
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})
		posFreqArr.Frequency = ss
		freqArr.Characters = append(freqArr.Characters, posFreqArr)
	}
	return freqArr
}
