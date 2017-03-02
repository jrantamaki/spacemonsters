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

const bullet = 1;
const monster = 2;
const explosion = 5;
const update_world_interval = 230;
const render_interval = 50;
const monsterUpdateModulo = 2; // speed of monsters in relation to bullets
const randomSleepBetweenMonsters = 400;
const constantSleepBetweenMonsters = 650;

func main() {

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

	// Tick channel to handle periodic updating.
	tickChan := time.NewTicker(time.Millisecond * update_world_interval).C

	renderChan := time.NewTicker(time.Millisecond * render_interval).C

	zero := []int{0,0,0,0,0,0,0,0};
	world := [][]int{zero, zero, zero, zero, zero, zero, zero, zero}

	x := 7
	z := -1;


	go func()  {
		fmt.Println("Starting the update the world ")
		for {
			select {
			case <-tickChan:
				z = z + 1
				updateWorld(world, z)
			}
		}
	}()

	go func()  {
		fmt.Println("Starting the rendering channel!")
		for {
			select {
			case <-renderChan:
				drawWorld(x, world, fb)
			}
		}
	}()

	go func()  {
		fmt.Println("Starting the monster engine")
		oldX := -1
		monsterX := rand.Int() % 8
		for {
			for {
				monsterX = (rand.Int() % 8)
				if monsterX != oldX {
					break;
				}

			}

			oldX = monsterX

			tmp := []int{0, 0, 0, 0, 0, 0, 0, 0};
			copy(tmp, world[monsterX]);
			tmp[0] = monster
			world[monsterX] = tmp
			sleep := rand.Int31n(randomSleepBetweenMonsters) + constantSleepBetweenMonsters

			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}()


	for {
		select {

		case <-signals:
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
			case stick.Right:
				//fmt.Println("1. L: ", x, " y ", y)
				x = (x + 1) % 8
				if x < 0 {
					x = x - 8
				}

			case stick.Left:
				//fmt.Println("1. R: ", x, " y ", y)
				x = (x - 1) % 8
				if x < 0 {
					x = x + 8
				}
			}
		}
	}
}

func updateWorld(world [][]int, z int) {

	for x := 0; x <= 7; x ++  { // iterate the x-axis
 		oldY := world[x];

		newY := []int{0,0,0,0,0,0,0,0};
		if z % monsterUpdateModulo == 0 { // hack to manage the monster vs bullet speed.
			for  y := 0 ;  y <= 7 ; y ++ { // move down the monsters in the matrix.
				if y == 7 {
					break;
				}
				if oldY[y] == monster {
					newY[y+1] = monster
				}
			}
		} else {
			copy(newY, oldY)
			for y := 0 ; y <= 7 ; y ++ {
				if newY[y] == bullet {
					newY[y] = 0;
				}
			}
		}

		// Was there any collisions with the bullets?
		for  y := 0 ;  y <= 7 ; y ++ {

			if newY[y] == monster && oldY[y] == bullet {
				newY[y] = explosion
				oldY[y] = 0;
			}
		}
		newMonsters := []int{0,0,0,0,0,0,0,0};
		copy(newMonsters, newY)

		for  y := 6 ; y > 0; y-- { // move the bullets
			if y == 0 {
				break;
			}
			if oldY[y]  == bullet && newMonsters[y] == monster {
				newY[y] = explosion

			} else if oldY[y] == bullet && newMonsters[y-1] == monster{
				newY[y-1] = explosion
			} else {
				if oldY[y] == bullet {
					newY[y-1] = bullet
				}
			}
		}

		world[x] = newY;

	}
}


func drawWorld(spaceship int, world [][]int, fb *screen.FrameBuffer) {
	for x := 0; x <= 7; x ++ {

		for  y:= 0 ; y <=  7; y++ {
			if world[x][y]  == monster {
				fb.SetPixel(x, y, color.Red)
			} else if world[x][y]  == bullet {
				fb.SetPixel(x, y, color.Blue)
			} else if world[x][y]  == explosion {
				fb.SetPixel(x, y, color.New(255, 255, 0))
			} else {
				fb.SetPixel(x, y, color.New(0, 0, 0))
			}
		}
	}
	fb.SetPixel(spaceship, 7, color.Green)
	screen.Draw(fb)
}