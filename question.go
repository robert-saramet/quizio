package main

import (
	"gopkg.in/yaml.v3"
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

func (q *Questions) load() *Questions {
	file, err := os.ReadFile("questions.yaml")
	handle(err, "os.ReadFile")
	err = yaml.Unmarshal(file, q)
	handle(err, "yaml.Unmarshal")
	return q
}

func (q *Question) checkAnswer(input string) bool {
	return strings.ToLower(q.Answer) == strings.ToLower(input)
}
