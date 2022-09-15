package main

import (
	"fmt"
	"log"
)

func main() {
	var input, incorrect string
	question := newQuestion("math", "5 + 5", "10", 1)

	fmt.Print(question.Text, ": ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal("Scanln failed: ", err)
	}
	if !question.checkAnswer(input) {
		incorrect = "in"
	}

	fmt.Print("Your answer is ", incorrect, "correct")
}
