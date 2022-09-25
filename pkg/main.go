package main

import (
	"flag"
	"fmt"

	"github.com/strobsi/wordleguessr/pkg/game"
)

const TOTAL_ROUNDS = 5000
const WORDS_PATH = "words.txt"

func main() {

	var mode string
	var simulation bool
	var silent bool

	flag.StringVar(&mode, "m", "improvedGuessing", "The mode you want to play. Can be 'native', 'improvedStart' and  'improvedGuessing'")
	flag.BoolVar(&simulation, "simulation", false, "Play in simulation mode")
	flag.BoolVar(&silent, "silent", false, "Play in silent mode")
	flag.Parse()

	if mode == "native" {
		playNative(simulation, silent)
	} else if mode == "improvedStart" {
		playImprovedStart(simulation, silent)
	} else {
		playImprovedGuessing(simulation, silent)
	}
}

func playNative(simulation bool, silent bool) {

	fmt.Println("Play game in Native mode")

	if simulation {

		totalScore := 0
		rounds := 0
		for rounds < TOTAL_ROUNDS {
			rounds++
			g := game.New(WORDS_PATH, game.Native)
			g.Simulation = simulation
			g.SilentMode = silent
			g.Play()
			totalScore += g.Score
		}
		var average = 0.0
		average = float64(totalScore) / TOTAL_ROUNDS
		fmt.Println("Average score in", rounds, " rounds:", average)
	} else {
		g := game.New(WORDS_PATH, game.Native)
		g.SilentMode = silent
		g.Play()
	}
}

func playImprovedStart(simulation bool, silent bool) {

	fmt.Println("Play game in ImprovedStart mode")

	if simulation {
		totalScore := 0
		rounds := 0
		for rounds < TOTAL_ROUNDS {
			rounds++
			g := game.New(WORDS_PATH, game.ImprovedStart)
			g.Simulation = simulation
			g.SilentMode = silent
			g.Play()
			totalScore += g.Score
		}
		var average = 0.0
		average = float64(totalScore) / TOTAL_ROUNDS
		fmt.Println("Average score in", rounds, " rounds:", average)
	} else {
		g := game.New(WORDS_PATH, game.ImprovedStart)
		g.SilentMode = silent
		g.Play()
	}
}

func playImprovedGuessing(simulation bool, silent bool) {

	fmt.Println("Play game in ImprovedGuessing mode")
	if simulation {
		totalScore := 0
		rounds := 0
		for rounds < TOTAL_ROUNDS {
			rounds++
			g := game.New(WORDS_PATH, game.ImprovedGuessing)
			g.Simulation = simulation
			g.SilentMode = silent
			g.Play()
			totalScore += g.Score
		}
		var average = 0.0
		average = float64(totalScore) / TOTAL_ROUNDS
		fmt.Println("Average score in", rounds, " rounds:", average)
	} else {
		g := game.New(WORDS_PATH, game.ImprovedGuessing)
		g.SilentMode = silent
		g.Play()
	}
}
