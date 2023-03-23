package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/term"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

type CLI struct {
	temp  string
	score int
}

func (c *CLI) init() {
	// Handle SIGINT gracefully
	ch := make(chan os.Signal, 3)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ch, os.Kill, syscall.SIGKILL)
	go func() {
		<-ch
		c.exit()
	}()
	c.clear()
}

func (c *CLI) clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	handle(err, "cmd.Run")
	fmt.Print("Press 'CTRL+C' to quit at any time\n")
}

func (c *CLI) read() (input string) {
	_, err := fmt.Scanln(&input)
	if err != nil {
		message := ""
		if err.Error() == "expected newline" {
			message = "Try again without spaces please: "
			// flush stdin before retrying
			stdin := bufio.NewReader(os.Stdin)
			_, err = stdin.ReadString('\n')
			handle(err, "ReadString")
		} else if err.Error() == "unexpected newline" {
			message = "Try again please: "
		}
		if message != "" {
			fmt.Print(message)
			input = c.read()
			return
		} else {
			handle(err, "Scanln()")
		}
	}
	return
}

func (c *CLI) readSpaces() (input string) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		handle(err, "scanner.Text()")
	}
	return
}

func (c *CLI) readKey() (char string) {
	t, _ := term.Open("/dev/tty")
	err := term.CBreakMode(t)
	handle(err, "term.CBreakMode")
	bytes := make([]byte, 3)
	numRead, err := t.Read(bytes)
	handle(err, "term.Read")
	if numRead == 1 {
		char = string(bytes[0])
	}
	err = t.Restore()
	handle(err, "term.Restore")
	err = t.Close()
	handle(err, "term.Close")
	return
}

func (c *CLI) cursor(on bool) {
	target := bufio.NewWriter(os.Stdout)
	defer func(target *bufio.Writer) {
		err := target.Flush()
		handle(err, "Flush")
	}(target)

	var cursor string
	if on {
		cursor = "\033[?25h"
	} else {
		cursor = "\033[?25l"
	}
	_, err := fmt.Fprintf(target, cursor)
	handle(err, "Fprintf")
}

func (c *CLI) newPlayer(username string) *Player {
	newPlayer := player.make()
	player = *newPlayer
	player.Name = username
	for {
		c.printMenu("Base")
		if player.Mode != "" && player.Difficulty != "" {
			return &player
		} else {
			c.temp = "Please configure all settings"
			c.printMenu("Base")
		}
	}
}

func (c *CLI) exit() {
	c.clear()
	fmt.Println("See you around!")
	c.cursor(true)
	os.Exit(0)
}

func (c *CLI) printKeyMenu(submenu []string, property string) {
	if property != "base" {
		fmt.Println(c.temp)
	}
	for _, item := range submenu {
		fmt.Println(item)
	}
	c.cursor(false)
	key := c.readKey()
	id, _ := strconv.Atoi(key)
	if id == 0 {
		if property == "base" {
			c.temp = "EXIT"
			return
		}
	} else if id < len(submenu) {
		if property == "mode" {
			player.Mode = submenu[id][3:]
		} else if property == "learning" {
			val, _ := strconv.ParseBool(submenu[id][3:])
			player.Learning = val
		} else if property == "base" {
			c.temp = (submenu[id])[3:]
			paths := []string{"Base", "Name", "Mode", "Learning", "Difficulty"}
			c.printMenu(paths[id])
		}
	} else {
		// print current menu
		c.printMenu(fmt.Sprint(strings.ToUpper(property[:1]), property[1:]))
	}
	c.printMenu("Base")
}

func (c *CLI) printConfirmMenu(submenu []string, property string) {
	fmt.Print("Confirm ", property, ": ", c.temp, "\n")
	for _, item := range submenu {
		fmt.Println(item)
	}
	key := c.readKey()
	if property == "name" {
		if key == "\n" {
			player.Name = c.temp
		} else if key == "0" {
			c.printMenu("Base")
		} else {
			c.printMenu("ConfirmName")
		}
	} else if property == "difficulty" {
		if key == "\n" {
			player.Difficulty = c.temp
		} else if key == "0" {
			c.printMenu("Base")
		} else {
			c.printMenu("ConfirmDifficulty")
		}
	}
	c.printMenu("Base")
}

func (c *CLI) printTextMenu(submenu []string) {
	fmt.Print(submenu[0])
	c.cursor(true)
	c.temp = c.read()
	c.cursor(false)
}

func (c *CLI) printMenu(propName string) {
	submenu := reflect.ValueOf(menu).
		FieldByName(propName).Interface().([]string)
	c.clear()
	if c.temp == "EXIT" {
		c.cursor(true)
		return
	}
	switch propName {
	case "Base":
		c.printKeyMenu(submenu, "base")
	case "Name":
		c.printTextMenu(submenu)
		c.printMenu("ConfirmName")
	case "ConfirmName":
		c.printConfirmMenu(submenu, "name")
	case "Mode":
		c.printKeyMenu(submenu, "mode")
	case "Learning":
		c.printKeyMenu(submenu, "learning")
	case "Difficulty":
		c.printTextMenu(submenu)
		c.printMenu("ConfirmDifficulty")
	case "ConfirmDifficulty":
		c.printConfirmMenu(submenu, "difficulty")
	}
}

func handle(err error, function string) {
	if err != nil {
		log.Fatalln(function, "failed:", err)
	}
	return
}
