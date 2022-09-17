package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pkg/term"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type CLI struct {
	menu   *Menu
	player *Player
	temp   string
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
	c.menu = c.initMenu()
}

func (c *CLI) clear() {
	if runtime.GOOS == "windows" {
		clearCmd = "cls"
	} else {
		clearCmd = "clear"
	}
	cmd := exec.Command(clearCmd)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	handle(err, "cmd.Run")
	fmt.Print("Press 'CTRL+C' to quit\n")
}

func (c *CLI) read() (input string) {
	_, err := fmt.Scanln(&input)
	if err != nil {
		if err.Error() == "expected newline" {
			stdin := bufio.NewReader(os.Stdin)
			// flush stdin before retrying
			_, err = stdin.ReadString('\n')
			handle(err, "ReadString")
			fmt.Print("Try again without spaces please: ")
			input = c.read()
			return
		} else {
			log.Fatal("Scanln failed: ", err)
		}
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
		ascii := int(bytes[0])
		char = string(ascii)
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

func (c *CLI) initPlayer() {
	fmt.Print("Welcome, please enter your username: ")
	username = c.read()
	filename := username + ".conf"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Print("Player data not found, let's create it now")
		for i := 0; i < 3; i++ {
			fmt.Print(".")
			time.Sleep(time.Second / 3)
		}
		c.player = c.newPlayer()
	} else {
		fmt.Println("Player data found, loading...")
		c.player = player.load(filename)
		c.cursor(false)
		fmt.Print("Would you like to change your settings? [y/n] ")
		for {
			key := c.readKey()
			key = strings.ToLower(key)
			if key == "y" {
				c.printMenu("Base")
				break
			} else if key == "n" {
				break
			}
		}
	}
	c.player.write(filename)
	c.cursor(true)
	c.clear()
}

func (c *CLI) newPlayer() *Player {
	newPlayer := player.make()
	c.player = newPlayer
	for {
		c.printMenu("Base")
		if c.player.Name != "" && c.player.DefaultMode != "" && c.player.DefaultDifficulty != "" {
			return c.player
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

type Menu struct {
	Base              []string
	Name              []string
	ConfirmName       []string
	Mode              []string
	Learning          []string
	Difficulty        []string
	ConfirmDifficulty []string
}

func (c *CLI) initMenu() *Menu {
	base := []string{
		"0) Return",
		"1) Change name",
		"2) Select classic/endless mode",
		"3) Toggle learning mode",
		"4) Change difficulty",
	}
	name := []string{
		"Enter new name: ",
	}
	confirmName := []string{
		"0) Return",
		"Enter: Confirm",
	}
	mode := []string{
		"0) Return",
		"1) Classic",
		"2) Endless",
	}
	learning := []string{
		"0) Return",
		"1) On",
		"2) Off",
	}
	difficulty := []string{
		"Enter difficulty range [1-5]: ",
	}
	confirmDifficulty := []string{
		"0) Return",
		"Enter: confirm",
	}
	return &Menu{base, name, confirmName,
		mode, learning, difficulty, confirmDifficulty}
}

func (c *CLI) printMenu(propName string) {
	submenu := reflect.ValueOf(c.menu).Elem().
		FieldByName(propName).Interface().([]string)

	var key string

	if propName != "Name" && propName != "Difficulty" {
		c.clear()
		if propName != "Base" {
			fmt.Println(c.temp)
		}
		for _, item := range submenu {
			fmt.Println(item)
		}
		key = c.readKey()
		if propName == "ConfirmName" || propName == "ConfirmDifficulty" {
			//fmt.Println(propName[8:], c.temp)
		} else {
			id, _ := strconv.Atoi(key)
			if id >= len(submenu) {
				c.printMenu(propName)
				return
			}
			c.temp = (submenu[id])[3:]
		}
	}
	c.cursor(false)
	if c.temp == "Return" {
		c.temp = ""
	}

	switch propName {
	case "Base":
		switch key {
		case "0":
			return
		case "1":
			c.printMenu("Name")
		case "2":
			c.printMenu("Mode")
		case "3":
			c.printMenu("Learning")
		case "4":
			c.printMenu("Difficulty")
		}
	case "Name":
		c.clear()
		fmt.Print(submenu[0])
		c.cursor(true)
		c.temp = c.read()
		c.cursor(false)
		c.printMenu("ConfirmName")
	case "ConfirmName":
		if key == "\n" {
			c.player.Name = c.temp
		} else if key == "0" {
			c.printMenu("Base")
		} else {
			c.printMenu(propName)
		}
		c.temp = ""
		c.printMenu("Base")
	case "Mode":
		switch key {
		case "0":
			c.printMenu("Base")
		case "1":
			c.player.DefaultMode = "Classic"
		case "2":
			c.player.DefaultMode = "Endless"
		}
		c.printMenu("Base")
	case "Learning":
		switch key {
		case "0":
			c.printMenu("Base")
		case "1":
			c.player.DefaultLearning = true
		case "2":
			c.player.DefaultLearning = false
		}
		c.printMenu("Base")
	case "Difficulty":
		c.clear()
		fmt.Print(submenu[0])
		c.cursor(true)
		c.temp = c.read()
		c.cursor(false)
		c.printMenu("ConfirmDifficulty")
	case "ConfirmDifficulty":
		fmt.Println(fmt.Sprint("New difficulty: ", c.temp))
		if key == "\n" {
			c.player.DefaultDifficulty = c.temp
		} else if key == "0" {
			c.printMenu("Base")
		} else {
			c.printMenu(propName)
		}
		c.temp = ""
		c.printMenu("Base")
	}
	c.cursor(true)
	return
}

func handle(err error, function string) {
	if err != nil {
		log.Fatalln(function, "failed:", err)
	}
	return
}
