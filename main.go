package main

import (
	"fmt"
	"log"
)

func main() {
	var input, incorrectModif string
	question := newQuestion("math", "5 + 5", "10", 1)

	fmt.Print(question.Text, ": ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal("Scanln failed: ", err)
	}
	if !question.checkAnswer(input) {
		incorrectModif = "in"
	}

	fmt.Print("Your answer is ", incorrectModif, "correct")
}
