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
	menu := new(Menu)
	c.menu = menu.load()
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

func (c *CLI) initPlayer() {
	fmt.Print("Welcome, please enter your username: ")
	username = c.read()
	filename := username + ".yaml"
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
	c.player.write(c.player.Name + ".yaml")
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
			c.player.DefaultMode = submenu[id][3:]
		} else if property == "learning" {
			val, _ := strconv.ParseBool(submenu[id][3:])
			c.player.DefaultLearning = val
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
			c.player.Name = c.temp
		} else if key == "0" {
			c.printMenu("Base")
		} else {
			c.printMenu("ConfirmName")
		}
	} else if property == "difficulty" {
		if key == "\n" {
			c.player.DefaultDifficulty = c.temp
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
	submenu := reflect.ValueOf(c.menu).Elem().
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
