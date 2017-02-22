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

const empty = 0;
const bullet = 1;
const monster = 2;
const explosion = 5;

const ticker_sleep = 400;

func main() {
	tickChan := time.NewTicker(time.Millisecond * ticker_sleep).C

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

	zero := []int{0,0,0,0,0,0,0,0};
	world := [][]int{zero, zero, zero, zero, zero, zero, zero, zero}
	x := 7
	y := 7

	z := 0;
	for {
		select {
		case <-tickChan:
			fmt.Println("Updating the world! ")
			updateWorld(world, fb)
			z = z + 1
			if z % 5 == 0 {

				monsterX := rand.Int() % 8
				monsterX = 0
				tmp := []int{0,0,0,0,0,0,0,0};
				copy(tmp, world[monsterX]);
				tmp[0] = monster
				world[monsterX] = tmp
				fmt.Println("added a monster: ", world[monsterX], " MonsterX, " , monsterX  )
			}

		case <-signals:
			fmt.Println("")
			screen.Clear()

			// Exit the loop
			return
		case e := <-input.Events:
			switch e.Code {

			case stick.Up:
				tmp := []int{0,0,0,0,0,0,0,0};

				copy(tmp, world[x]);
				tmp[6] = bullet
				world[x] = tmp
				fmt.Println("ADDING A BULLET: ", world[x], " bullet x, " , x  )

			case stick.Right:
				fmt.Println("1. L: ", x, " y ", y)
				x = (x + 1) % 8
				if x < 0 {
					x = x - 8
				}

			case stick.Left:
				fmt.Println("1. R: ", x, " y ", y)
				x = (x - 1) % 8
				if x < 0 {
					x = x + 8
				}
			}
		}
		drawWorld(x, world, fb)
	}
}

func updateWorld(world [][]int, fb *screen.FrameBuffer) {

	for x := 0; x <= 7; x ++  { // iterate the x-axis
 		oldY := world[x];
		newY := []int{0,0,0,0,0,0,0,0};
		for  y := 0 ;  y <= 7 ; y ++ { // move down the monsters in the matrix.
			if y == 7 {
				break;
			}
			if oldY[y] == monster {
				newY[y+1] = monster
			}
		}
		if x == 0{
			fmt.Println("Moved the monsters  ", oldY, " newY ", newY)
		}

		// Was there any collisions with the bullets?
		for  y := 0 ;  y <= 7 ; y ++ {

			if newY[y] == monster && oldY[y] == bullet {
				fmt.Println("Detected explosion oldY ", oldY, " newY ", newY)
				newY[y] = explosion
				oldY[y] = 0;
			}
		}


		for  y := 6 ; y > 0; y-- { // move the bullets
			if y == 0 {
				break;
			}
			if oldY[y]  == bullet && newY[y] == monster {
				newY[y] = explosion

			} else if oldY[y] == bullet {
				newY[y-1] = bullet
			}
		}
		if x == 0 {
			fmt.Println("Moved the Bullets:  ", oldY, " newY ", newY)
		}
		world[x] = newY;

	}
}


func drawWorld(spaceship int, world [][]int, fb *screen.FrameBuffer) {
	for x := 0; x <= 7; x ++ {

		for  y:= 0 ; y <=  7; y++ {
			if world[x][y]  == monster {
				fmt.Println("Rendering a monster x:" ,x, " y: ",y)
				fb.SetPixel(x, y, color.Red)

			} else if world[x][y]  == bullet {
				fmt.Println("Rendering a bullet x:" ,x, " y: ",y)
				fb.SetPixel(x, y, color.Blue)
			} else if world[x][y]  == explosion {
				fmt.Println("Rendering a explosion x:" ,x, " y: ",y)
				fb.SetPixel(x, y, color.New(255,255,0))

			} else {
				fb.SetPixel(x, y, color.New(0, 0, 0))
			}

		}
	}
	fb.SetPixel(spaceship, 7, color.Green)
	screen.Draw(fb)
}

