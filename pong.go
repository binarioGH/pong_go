package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"golang.org/x/term"
)

type Paddle struct {
	X      int
	Y      int
	Height int
	Width  int
}

func (p Paddle) Draw() {
	gotoxy(p.X, p.Y)
	var i, total int
	total = p.Y + p.Height
	for i = p.Y; i < total; i++ {
		gotoxy(p.X, i)
		fmt.Print("█")
	}
}

func (p Paddle) IsInHitBox(x int, y int) bool {
	if (x >= p.X && x <= (p.X+p.Width)) && (y >= p.Y && y <= (p.Y+p.Height)) {
		return true
	}
	return false

}

type EnemyBrain struct {
}

func main() {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	clearScreen()
	gotoxy(width, height)
	var score int
	player := Paddle{}
	player.X = 5
	player.Y = 5
	player.Width = 1
	player.Height = 3
	score = 0

	ch := make(chan string)

	go func(ch chan string) {
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

	for score < 5 {
		clearScreen()

		select {
		case stdin, _ := <-ch:
			if stdin == "w" {
				player.Y -= 1
			} else if stdin == "s" {
				player.Y += 1

			}

		default:

			player.Draw()
			waitMil(30)
			gotoxy(width, height)

		}
	}

}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func gotoxy(x int, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func waitMil(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}
