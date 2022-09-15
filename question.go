package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type Question struct {
	Topic      string `yaml:"topic"`
	Text       string `yaml:"text"`
	Answer     string `yaml:"answer"`
	Difficulty int    `yaml:"difficulty"`
}

type Questions []Question

func newQuestion(topic string, text string, answer string, difficulty int) *Question {
	if difficulty < 1 || difficulty > 5 {
		log.Fatal("NewQuestion: difficulty must be int 1 to 5")
	}
	return &Question{Topic: topic, Text: text, Answer: answer, Difficulty: difficulty}
}

func (q *Questions) parseData() *Questions {
	file, err := os.ReadFile("data.yaml")
	if err != nil {
		log.Fatal("ReadFile failed: ", err)
	}
	err = yaml.Unmarshal(file, q)
	if err != nil {
		log.Fatal("Unmarshal failed: ", err)
	}
	return q
}

func (q *Question) checkAnswer(input string) bool {
	return strings.ToLower(q.Answer) == strings.ToLower(input)
}
