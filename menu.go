package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Menu struct {
	Base              []string `yaml:"base"`
	Name              []string `yaml:"name"`
	ConfirmName       []string `yaml:"confirmName"`
	Mode              []string `yaml:"mode"`
	Learning          []string `yaml:"learning"`
	Difficulty        []string `yaml:"difficulty"`
	ConfirmDifficulty []string `yaml:"confirmDifficulty"`
}

func (m *Menu) load() *Menu {
	file, err := os.ReadFile("menu.yaml")
	handle(err, "os.ReadFile")
	err = yaml.Unmarshal(file, m)
	handle(err, "yaml.Unmarshal")
	return m
}
