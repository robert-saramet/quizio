package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Player struct {
	Name              string `yaml:"name"`
	Theme             string `yaml:"theme"`
	DefaultMode       string `yaml:"defaultMode"`
	DefaultLearning   bool   `yaml:"defaultLearning"`
	DefaultDifficulty string `yaml:"defaultDifficulty"`
	Highscore         int    `yaml:"highscore"`
	XP                int    `yaml:"xp"`
}

func (p *Player) make() *Player {
	p.Name = ""
	p.Theme = "default"
	p.DefaultMode = ""
	p.DefaultLearning = false
	p.DefaultDifficulty = ""
	p.Highscore = 0
	p.XP = 0
	return p
}

func (p *Player) load(filename string) *Player {
	file, err := os.ReadFile(filename)
	handle(err, "os.ReadFile")
	err = yaml.Unmarshal(file, p)
	handle(err, "yaml.Unmarshal")
	return p
}

func (p *Player) write(filename string) {
	data, err := yaml.Marshal(p)
	handle(err, "yaml.Marshal")
	err = os.WriteFile(filename, data, 0666)
	handle(err, "WriteFile")
}
