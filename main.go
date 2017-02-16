package main

import (
	"flag"
	"fmt"
	"github.com/nathany/bobblehat/sense/screen"
	"github.com/nathany/bobblehat/sense/screen/color"
	"github.com/nathany/bobblehat/sense/stick"
	"os"
	"os/signal"
	"time"
	"math/rand"
)

func main() {
	tickChan := time.NewTicker(time.Millisecond * 200).C

	var path string
	flag.StringVar(&path, "path", "/dev/input/event2", "path to the event device")

	// Parse command line flags
	flag.Parse()

	fb := screen.NewFrameBuffer()
	// Open the input device (and defer closing it)
	input, err := stick.Open(path)
	if err != nil {
		fmt.Printf("Unable to open input device: %s\nError: %v\n", path, err)
		os.Exit(1)
	}

	// Clear the screen
	screen.Clear()

	// Print the name of the input device
	fmt.Println(input.Name())

	// Set up a signals channel (stop the loop using Ctrl-C)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	world := World{X: 7, Bullets: []Bullet{}}
	x := 7
	y := 7

	for {
		select {
		case <-tickChan:
			//fmt.Println("upate the world! ")
			updateBullets(&world, fb)
			//Update the world for monsters and bullets.
			//fb.SetPixel(x, y, color.New(0,0,0))
			//fb.SetPixel(x, y, color.Green)

			//screen.Draw(fb)

		case <-signals:
			fmt.Println("")
			screen.Clear()

			// Exit the loop
			return
		case e := <-input.Events:
			switch e.Code {

			case stick.Up:
				fmt.Println("up")
				b := Bullet{X: world.X, Y: 7, Color: color.Blue}
				world.Bullets = append(world.Bullets, b)

			case stick.Down:
				fmt.Println("Monster!")
				monster := rand.Int() % 8
				monstery := 0
				fb.SetPixel(monster, monstery, color.Red)

			case stick.Right:
				fmt.Println("1. L: ", x, " y ", y)
				fb.SetPixel(x, y, color.New(0, 0, 0))
				x = (x + 1) % 8
				if x < 0 {
					x = x - 8
				}
				world.X = x
				fmt.Println("L: ", x, " y ", y)
			case stick.Left:
				//fmt.Println("right")
				fmt.Println("1. R: ", x, " y ", y)
				fb.SetPixel(x, y, color.New(0, 0, 0))
				x = (x - 1) % 8
				if x < 0 {
					x = x + 8
				}
				world.X = x
				fmt.Println("R: ", x, " y ", y)

			}
		}
		drawWorld(world, fb)
	}
}

func updateBullets(world *World, fb *screen.FrameBuffer) {

	remainingBullets := []Bullet{}

	for i := 0; i < len(world.Bullets); i++ {
		b := world.Bullets[i]
		x := (b.X)
		y := (b.Y)
		fb.SetPixel(x, y, color.New(0, 0, 0))

		yy := y - 1

		if x >= 0 && yy >= 0 {
			bullet := Bullet{X: x, Y: yy, Color: color.Blue}
			//world.Bullets[i] = Bullet{X: x, Y: yy, Color: color.Blue}
			remainingBullets = append(remainingBullets, bullet)
			fb.SetPixel(x, yy, b.Color)
		}
	}
	world.Bullets = remainingBullets

}

func drawWorld(world World, fb *screen.FrameBuffer) {
	fb.SetPixel(world.X, 7, color.Green)
	screen.Draw(fb)
}

type Bullet struct {
	X     int
	Y     int
	Color color.Color
}

type World struct {
	X       int
	Bullets []Bullet
}
