package main

import (
	"fmt"
	"strings"
)

type Question struct {
	Topic      string `yaml:"topic"`
	Text       string `yaml:"text"`
	Answer     string `yaml:"answer"`
	Difficulty int    `yaml:"difficulty"`
}

func newQuestion(topic string, text string, answer string, difficulty int) *Question {
	if difficulty < 1 || difficulty > 5 {
		panic("NewQuestion: difficulty must be int 1 to 5")
	}
	return &Question{Topic: topic, Text: text, Answer: answer, Difficulty: difficulty}
}

func (q *Question) print() {
	fmt.Print(q.Text, ": ")
}

func (q *Question) checkAnswer(input string) bool {
	return strings.ToLower(q.Answer) == strings.ToLower(input)
}
