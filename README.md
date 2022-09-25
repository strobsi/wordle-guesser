# WordleGuesser

This small programm plays the wordle game for you.

## Usage

To run the program just compile and run the code:

```sh
go run pkg/main.go -silent <run in silent mode> -simulation <run in simulation mode>
```

The program runs in simulation mode by default. If you want to play for yourself, change the `Simulation` and `Silent` mode.
To change the word input, change line 11 in main.go.

## Strategies

the program runs with three slightly different strategies. Each of them have different average results.

### Native

the native stratgey runs the program and simply sorts out the words which do not pass the filter.

**GreenFilter:**

Simply checks all words for the specific character at the spcecific position.

**YellowFilter:**

Checks if words contain the character - not at the given position.

**GrayFilrer:**

Checks if the words do NOT contain the character within the gray filter

> Average score at 100000 games: **4.49743**

### ImprovedStartWord

It turns out that based in the 'gainable' information from the startword, just changing the start word already improves the guesses. There are different strategies on why and how to get the best starting word. It also differentiates based on the given word list. What I found to be the best is the word 'slate'.

The rest of the strategy is exactly the same.

> Average score at 100000 games: **4.25725**

### ImprovedGuessing

Running the program in improved guessing mode is different from the methods above. The filter stay the same, but instead of randomly choosing a new word from the remaining words, before each guess the program scans for the frequency of the letters. The more a character at a certain position is found, the better. Using this strategy the average score could be decreased further.

> Average score at 100000 games: **4.09332**

## More

I am sure, there are strategies to further improve the guessing score. See 3Blue1Brown for example, which analyzed the issue on a whole new level. [https://www.youtube.com/watch?v=v68zYyaEmEA](https://www.youtube.com/watch?v=v68zYyaEmEA). However, before checking the internet, I wanted to find out a solution by myself and not just simply copying others.
