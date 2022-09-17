package main

import (
	"fmt"
)

var input, username, incorrectModif, clearCmd string
var questions Questions
var player Player
var cli CLI

func main() {
	questions.load()

	defer cli.exit()
	cli.init()
	cli.initPlayer()

	// loop through the questions
	for _, question := range questions {
		fmt.Print(question.Text, ": ")

		input = cli.read()

		incorrectModif = ""
		if !question.checkAnswer(input) {
			incorrectModif = "in"
		}

		fmt.Print("Your answer is ", incorrectModif, "correct\n\n")
	}
}
