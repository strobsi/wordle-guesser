package game

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// New initializes new game with mode given by parameter
func New(wordListPath string, mode GameMode) *Game {

	g := Game{}
	g.Score = 0
	// Read out words from given path to word file.
	readFile, err := os.Open(wordListPath)
	if err != nil {
		g.Print("Unable to read words file: " + err.Error())
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		word := Word{
			Characters: []byte(fileScanner.Text()),
		}
		g.WordList.Words = append(g.WordList.Words, word)
	}
	readFile.Close()
	g.Mode = mode

	return &g
}

// Play starts the game
func (g *Game) Play() {

	// When in improved guessing mode, calculate frequency of letters at each position
	if g.Mode == ImprovedGuessing {
		g.LetterFrequency = CalculateFreq(g.WordList.Words)
	}
	// Check if we are in simulation mode and generate random target word, if so.
	if g.Simulation {
		target, _ := g.GetNewWord()
		g.Target = target
	}

	// Get new start word.
	startWord := Word{}
	if g.Mode == Native {
		startWord, _ = g.GetNewWord()
	} else {
		// This is simply looked up in the internet - no brain work here.
		// Alternatively, also the word with the highest frequency score can be taken
		startWord = Word{Characters: []byte("slate")}
	}
	g.CurrentWord = startWord.String()

	// Print output to user if not in silent mode
	g.Print("-----------------------------------")
	if g.Simulation {
		g.Print("Searching: " + g.Target.String())
		r := g.AnalyzeWords(g.Target, startWord)
		// Automatically guess when in simulation mode
		g.Guess(r)
	} else {
		// Ask user for result.
		g.Print("Enter start word: " + startWord.String())
		g.Print("What's the result? Enter comma separated ( green = 2, yellow = 1, gray = 0)")
		g.WaitForInput()
	}
}

// WaitForInput simply reads out the input of the user
// Has to be 5 digits, separated by comma
func (g *Game) WaitForInput() {
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	if !g.SilentMode {
		fmt.Print("> ")
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		g.Print("An error occured while reading input. Please try again" + err.Error())
		return
	}
	// remove the delimeter from the string
	inputRaw := strings.Split(strings.TrimSuffix(input, "\n"), ",")
	guesses := []int{}
	for _, v := range inputRaw {
		intValue, err := strconv.Atoi(v)
		if err != nil {
			g.Print("Error parsing the results, please enter again")
		}
		guesses = append(guesses, int(intValue))
	}
	// Start guess, when input is available
	g.Guess(guesses)
}

// Guess is the central function to calculate the new word
func (g *Game) Guess(results []int) {
	g.Score++
	// Check if input is 2,2,2,2,2
	if isCorrect(results) {
		g.Print("Congrats, you won!")
	} else {
		// Append results to filter slices
		for index, v := range results {

			// This can be further improved ( codewise ) by combining slices into one filter with separate field for color
			switch v {
			case 0:
				g.GrayChars = append(g.GrayChars, Result{Position: index, Char: g.CurrentWord[index]})
			case 1:
				g.YellowChars = append(g.YellowChars, Result{Position: index, Char: g.CurrentWord[index]})
			case 2:
				g.GreenChars = append(g.GreenChars, Result{Position: index, Char: g.CurrentWord[index]})
			}
		}

		// Remove words, which do not pass the filter slices
		g.WordList.Words = g.removeWords(g.WordList.Words, g.GreenChars, g.YellowChars, g.GrayChars)

		// When in native or improves start word mode, just take a random word out of the remaining ones
		nW := Word{}
		if g.Mode != ImprovedGuessing {
			w, err := g.GetNewWord()
			if err != nil {
				g.Print("No word found")
			}
			nW = w
		} else {
			// Get new word based on frequency score.
			w, err := g.GetMostLikelyWord(results)
			if err != nil {
				g.Print("No word found")
			}
			nW = w
		}
		g.CurrentWord = nW.String()
		g.Print("Guessed new word: " + g.CurrentWord)
		// Check if in simulation and continue, if so. Otherwise, wait for input from user
		if !g.Simulation {
			g.WaitForInput()
		} else {
			result := g.AnalyzeWords(g.Target, Word{Characters: []byte(g.CurrentWord)})
			g.Print("result: ", result)
			g.Guess(result)
		}

	}
}

// GetNewWord simply returns a random word from wordlist
func (g *Game) GetNewWord() (Word, error) {
	rand.Seed(time.Now().UnixNano())
	if (len(g.WordList.Words)) > 0 {
		return g.WordList.Words[rand.Intn(len(g.WordList.Words))], nil
	} else {
		return Word{}, errors.New("no word found")
	}
}

// GetMostLikelyWord returns the word with the least distance to the optimal frequency
func (g *Game) GetMostLikelyWord(results []int) (Word, error) {

	// First, check out characters we have to check
	// Example: we need to get first letter
	type ScoredWord struct {
		Word  Word
		Score int
	}

	// Crate slice of words with additional score to sort it later on
	scoredWordList := []ScoredWord{}
	if (len(g.WordList.Words)) > 0 {
		// Get new frequency of characters
		g.LetterFrequency = CalculateFreq(g.WordList.Words)
		for _, v := range g.WordList.Words {
			// Append scored word to slice
			scoredWordList = append(scoredWordList, ScoredWord{Word: v, Score: g.getScoreForWord(v)})
		}
	}
	// Sort slice by least score. The less, the better - See getScoreForLetter
	sort.Slice(scoredWordList, func(i, j int) bool {
		return scoredWordList[i].Score < scoredWordList[j].Score
	})
	if len(scoredWordList) > 0 {
		return scoredWordList[0].Word, nil
	}
	return Word{}, errors.New("no word available")
}

// isCorrect checks, if all the results are 2
func isCorrect(results []int) bool {
	isCorrect := true
	for i := 0; i < len(results); i++ {
		if results[i] != 2 {
			isCorrect = false
		}
	}
	return isCorrect
}

// AnalyzeWords generates the 'user input' in simulation mode
func (g *Game) AnalyzeWords(target Word, guess Word) []int {
	g.Print("Comparing: ", target.String(), " and ", guess.String())
	points := []int{}
	for i, v := range guess.Characters {
		if v == target.Characters[i] {
			points = append(points, 2)
		} else if strings.Contains(target.String(), string(v)) {
			points = append(points, 1)
		} else {
			points = append(points, 0)
		}
	}
	g.Print("Input: ", points)
	return points
}

// removeWords removes the words based on the filter slices
func (g *Game) removeWords(words []Word, green []Result, yellow []Result, gray []Result) []Word {

	greenFilter := []Word{}
	// Apply green filter
	for _, v := range words {
		isCorrectGreen := true
		for _, g := range green {
			if v.Characters[g.Position] != g.Char {
				isCorrectGreen = false
			}
		}
		if isCorrectGreen {
			greenFilter = append(greenFilter, v)
		}
	}
	// Apply gray filter
	grayFilter := []Word{}
	for _, v := range greenFilter {
		if !v.HasOneOfChars(gray) {
			grayFilter = append(grayFilter, v)
		}
	}

	// Apply yellow filter
	yellowFilter := []Word{}
	for _, v := range grayFilter {
		if v.HasAllCharacters(yellow) {
			if v.HasRightPositions(yellow) {
				yellowFilter = append(yellowFilter, v)
			}
		}
	}
	return yellowFilter
}

// getScoreForWord iterates through each character and sums up the score of the characters
func (g *Game) getScoreForWord(w Word) int {
	score := 0
	for i, v := range w.Characters {
		charScore := g.getScoreForCharacter(i, v)
		score += charScore
	}
	return score
}

// getScoreForCharacter calculates the distance between the given character at the given position
// and the index 0 element of the frequency slice at this position
func (g *Game) getScoreForCharacter(index int, b byte) int {
	score := 0
	for i := 0; i < len(g.LetterFrequency.Characters[index].Frequency); i++ {
		el := g.LetterFrequency.Characters[index].Frequency[i]
		if el.Key == b {
			score = i
		}
	}
	return score
}
