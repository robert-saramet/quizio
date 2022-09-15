package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	var input, incorrectModif, clearCmd string
	var questions Questions

	// load questions from yaml
	questions.parseData()

	// run os specific clear command
	if runtime.GOOS == "windows" {
		clearCmd = "cls"
	} else {
		clearCmd = "clear"
	}
	cmd := exec.Command(clearCmd)
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Print("Press 'CTRL+C' to quit\n\n")

	// Handle SIGTERM gracefully
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	// loop through the questions
	for _, question := range questions {
		fmt.Print(question.Text, ": ")

		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal("Scanln failed: ", err)
		}

		incorrectModif = ""
		if !question.checkAnswer(input) {
			incorrectModif = "in"
		}

		fmt.Print("Your answer is ", incorrectModif, "correct\n\n")
	}
}
