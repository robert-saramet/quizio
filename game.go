package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func initPlayer() {
	fmt.Print("Welcome, please enter your username: ")
	username := cli.read()
	filename := username + ".yaml"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Print("Player data not found, let's create it now")
		for i := 0; i < 3; i++ {
			fmt.Print(".")
			time.Sleep(time.Second / 3)
		}
		player = *cli.newPlayer(username)
	} else {
		fmt.Println("Player data found, loading...")
		player = *player.load(filename)
		cli.cursor(false)
		fmt.Print("Would you like to change your settings? [y/n] ")
		for {
			key := cli.readKey()
			key = strings.ToLower(key)
			if key == "y" {
				cli.printMenu("Base")
				break
			} else if key == "n" {
				break
			}
		}
	}
	player.write(player.Name + ".yaml")
	cli.cursor(true)
	cli.clear()
}

func selectQuestions() {
	cli.clear()
	cli.cursor(true)
	fmt.Println("Please select the topics you're interested in")
	fmt.Println("Enter their indexes on a single line, separated by commas")
	var topics, allTopics []string
	for _, question := range questions {
		newQuestion := true
		for _, topic := range allTopics {
			if question.Topic == topic {
				newQuestion = false
			}
		}
		if newQuestion {
			allTopics = append(allTopics, question.Topic)
		}
	}
	for i, topic := range allTopics {
		fmt.Print(i, ") ", topic, "\n")
	}
	fmt.Print("Topics: ")
	selected := strings.Split(cli.read(), ",")
	for _, item := range selected {
		index, err := strconv.Atoi(item)
		handle(err, "strconv.Atoi")
		topics = append(topics, allTopics[index])
	}
	questions = *questions.restrain(topics, player.Difficulty)
}

func runGame() {
	fmt.Print("Starting game")
	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(time.Second / 3)
	}
	cli.clear()
	switch player.Mode {
	case "Classic":
		runClassicGame()
	case "Endless":
		runEndlessGame()
	}
	cli.clear()
	cli.cursor(false)
	fmt.Println("Congratulations!")
	fmt.Println("You earned", cli.score, "points")
	if cli.score > player.Highscore {
		fmt.Println("This is your new highscore")
		fmt.Println("The old one was", player.Highscore)
		player.Highscore = cli.score
	} else {
		fmt.Println("Your highscore is", player.Highscore)
	}
	player.XP += cli.score
	fmt.Println("Your total XP is", player.XP)
	player.write(player.Name + ".yaml")
	fmt.Println("Press any key to exit")
	cli.readKey()
}

func runClassicGame() {
	lives := 5
	for _, question := range questions {
		fmt.Print(question.Text, ": ")
		input := cli.readSpaces()
		incorrectModif := ""
		if !question.checkAnswer(input) {
			incorrectModif = "in"
			lives -= 1
		} else {
			cli.score += 10 * question.Difficulty
		}
		if player.Learning {
			fmt.Print("Your answer is ", incorrectModif, "correct\n")
			if incorrectModif == "in" {
				fmt.Println("Correct answer:", question.Answer)
				fmt.Println("You have", lives, "lives left")
			}
		}
		fmt.Println()
		if lives < 1 {
			break
		}
	}
}

func runEndlessGame() {
	questions = *questions.shuffle()
	for _, question := range questions {
		fmt.Print(question.Text, ": ")
		input := cli.readSpaces()
		incorrectModif := ""
		if !question.checkAnswer(input) {
			incorrectModif = "in"
		} else {
			cli.score += 10 * question.Difficulty
		}
		fmt.Print("Your answer is ", incorrectModif, "correct\n")
		if incorrectModif == "in" && player.Learning {
			fmt.Println("Correct answer:", question.Answer)
		}
		cli.cursor(false)
		fmt.Println("Press 'Enter' to continue or 'T' to stop")
		for {
			key := cli.readKey()
			if key == "\n" {
				break
			} else if strings.ToUpper(key) == "T" {
				return
			}
		}
		cli.cursor(true)
		fmt.Println()
	}
}
