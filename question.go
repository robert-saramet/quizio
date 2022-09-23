package main

import (
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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

func (q *Questions) restrain(topics []string, difficulty string) *Questions {
	minDif, err := strconv.Atoi(difficulty[:1])
	handle(err, "strconv.Atoi")
	maxDif, err := strconv.Atoi(difficulty[len(difficulty)-1 : len(difficulty)])
	handle(err, "strconv.Atoi")
	newQuestions := new(Questions)
	for _, question := range *q {
		for _, topic := range topics {
			if question.Topic == topic {
				dif := question.Difficulty
				if dif >= minDif && dif <= maxDif {
					*newQuestions = append(*newQuestions, question)
				}
			}
		}
	}
	return newQuestions
}

func (q *Questions) shuffle() *Questions {
	qD := *q
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(qD), func(i, j int) {
		qD[i], qD[j] = qD[j], qD[i]
	})
	return &qD
}

func (q *Question) checkAnswer(input string) bool {
	return strings.ToLower(q.Answer) == strings.ToLower(input)
}
