package game

// Result is an element of a filter slice
type Result struct {
	Char     byte
	Position int
}

// Word represents the word with all of its characters
type Word struct {
	Characters []byte
}

// Word list simply wraps the given wordlist as word slice
type WordList struct {
	Words []Word
}

// Enum for game mode
type GameMode int64

const (
	Native           GameMode = 0
	ImprovedStart    GameMode = 1
	ImprovedGuessing GameMode = 2
)

// Game is the central struct of the program and contains all information for processing the guesses
type Game struct {
	Mode            GameMode
	WordList        WordList       // The wordlist to start with
	CurrentWord     string         // current word of the game
	GrayChars       []Result       // Filter for gray characters
	YellowChars     []Result       // Filter for yellow characters
	GreenChars      []Result       // Filter for green characters
	LetterFrequency FrequencySlice // Contains the character slices with calculated frequencies
	// Followng only in simulation mode
	Target     Word // Target word if in simulation mode
	Simulation bool
	Score      int
	SilentMode bool
}

// Simple wrapper to sort better
type kv struct {
	Key   byte `json:"key"`
	Value int  `json:"value"`
}

// PositionFrequencySlice contains the sortable frequency of letters
type PositionFrequencySlice struct {
	Frequency []kv `json:"frequency"`
}

// FrequencySlice wraps the sortable frequency of letters
type FrequencySlice struct {
	Characters []PositionFrequencySlice `json:"characters"`
}

// FrequencySlice contains the mapped frequency of letters
type PositionFrequency struct {
	Frequency map[byte]int `json:"frequency"`
}

// FrequencySlice contains the slice of the mapped frequency of letters
type FrequencyMap struct {
	Characters []PositionFrequency `json:"characters"`
}
