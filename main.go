package main

var questions Questions
var menu Menu
var player Player
var cli CLI

func main() {
	defer cli.exit()
	questions.load()
	menu.load()
	cli.init()
	initPlayer()
	selectQuestions()
	runGame()
}
