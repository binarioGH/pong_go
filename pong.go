package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"golang.org/x/term"
)

type Paddle struct {
	X      int
	Y      int
	Height int
	Width  int
	Sprite string
}

func (p Paddle) Draw() {
	gotoxy(p.X, p.Y)
	var i, total int
	total = p.Y + p.Height
	for i = p.Y; i < total; i++ {
		gotoxy(p.X, i)
		fmt.Print(p.Sprite)
	}
}

func (p *Paddle) Clear() {
	oldSprite := p.Sprite
	p.Sprite = " "
	p.Draw()
	p.Sprite = oldSprite
}

func (p Paddle) IsInHitBox(x int, y int) bool {
	if (x >= p.X && x <= (p.X+p.Width)) && (y >= p.Y && y <= (p.Y+p.Height)) {
		return true
	}
	return false

}

type Ball struct {
	X          int
	Y          int
	XDirection int
	YDirection int
	Sprite     string
	lastMove   int64
}

func (b *Ball) ChangeDirection() {
	b.XDirection *= -1
	b.YDirection *= -1
}

func (b *Ball) ChangeVerticalDirection() {
	b.YDirection *= -1
}

func (b *Ball) ChangeHorizontalDirection() {
	b.XDirection *= -1
}

func (b *Ball) Move() {
	t := time.Unix(0, b.lastMove*int64(time.Millisecond))
	elapsed := time.Since(t)
	ms := elapsed.Milliseconds()
	if ms >= 50 {
		b.Clear()
		b.X += b.XDirection
		b.Y += b.YDirection
		b.lastMove = int64(time.Now().UnixMilli())
		b.Draw()
	}

}

func (b Ball) Draw() {
	gotoxy(b.X, b.Y)
	fmt.Print(b.Sprite)
}

func (b *Ball) Clear() {
	oldSprite := b.Sprite
	b.Sprite = " "
	b.Draw()
	b.Sprite = oldSprite

}

func main() {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	clearScreen()
	gotoxy(width, height)
	var score int
	player := Paddle{}
	player.Sprite = "█"
	player.X = 5
	player.Y = 5
	player.Width = 1
	player.Height = 3

	robot := Paddle{}
	robot.Sprite = "█"
	robot.X = width - 6
	robot.Y = height - 6
	robot.Width = 1
	robot.Height = 3

	ball := Ball{}

	ball.X = width / 2
	ball.Y = height / 2
	ball.XDirection = 1
	ball.YDirection = 1
	ball.Sprite = "O"
	ball.lastMove = int64(time.Now().UnixMilli())

	score = 0

	ch := make(chan string)

	var robotPoints, playerPoints int
	robotPoints = 0
	playerPoints = 0

	go func(ch chan string) {
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

	//handle exiting the program with ctrl c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			exec.Command("stty", "-F", "/dev/tty", "echo").Run()
			os.Exit(0)
		}
	}()

	clearScreen()
	var halfWidth int = width / 2
	for score < 5 {

		// Robot logic
		if ball.X >= halfWidth {
			robot.Clear()
			if ball.Y >= (robot.Y+robot.Height) && (robot.Y+robot.Height) < height {
				robot.Y += 1
			} else if ball.Y < robot.Y && robot.Y > 1 {
				robot.Y -= 1
			}
			robot.Draw()

		}

		// Player input
		select {
		case stdin, _ := <-ch:
			if stdin == "w" {

				if player.Y > 1 {
					player.Clear()
					player.Y -= 1
				}

			} else if stdin == "s" {

				if (player.Y + player.Height) <= height {
					player.Clear()
					player.Y += 1
				}

			}

		default:

			player.Draw()

		}

		// ball logic
		nextX := ball.X + ball.XDirection
		nextY := ball.Y + ball.YDirection
		if player.IsInHitBox(nextX, nextY) || robot.IsInHitBox(nextX, nextY) {
			ball.ChangeDirection()
		}
		if ball.Y <= 1 || ball.Y >= height {
			ball.ChangeVerticalDirection()
		}

		if ball.X <= 1 || ball.X >= width {
			clearScreen()
			score += 1
			if ball.X <= 1 {
				robotPoints += 1
			} else if ball.X >= width {
				playerPoints += 1
			}
			ball.Clear()
			ball.X = width / 2
			ball.Y = height / 2
			ball.ChangeHorizontalDirection()

		}

		ball.Move()

		drawScore(playerPoints, robotPoints, width)
		waitMil(10)

	}
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	os.Exit(0)

}

func drawScore(playerPoints int, robotPoints int, width int) {
	scoreString := fmt.Sprintf("Human Score: %d    ||    Robot Score %d", playerPoints, robotPoints)
	scoreLengthHalf := len(scoreString) / 2
	newX := (width / 2) - scoreLengthHalf
	gotoxy(newX, 2)
	fmt.Print(scoreString)

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
