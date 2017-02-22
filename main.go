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
const mayhem = -1
const ticker_sleep = 200;
const monsterUpdateModulo = 2;

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
	z := -1;
	for {
		select {
		case <-tickChan:
			fmt.Println("Updating the world! ")
			gameOver := updateWorld(world, z)
			if gameOver {
				itsGameOver(world)
				drawWorld(x, world, fb)
				time.Sleep(time.Second * 3)
				screen.Clear()
				return
			}
			z = z + 1
			if z % 8 == 0 {

				monsterX := rand.Int() % 8
				tmp := []int{0,0,0,0,0,0,0,0};
				copy(tmp, world[monsterX]);
				tmp[0] = monster
				world[monsterX] = tmp
				fmt.Println("added a monster: ", world[monsterX], " MonsterX, " , monsterX  )
			}

		case <-signals:
			fmt.Println("EXIT")
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

func updateWorld(world [][]int, z int) bool {

	for x := 0; x <= 7; x ++  { // iterate the x-axis
 		oldY := world[x];

		if oldY[7] == monster {
			return true  // GAME OVER
		}

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
			if x == 0{
				fmt.Println("Moved the monsters  ", oldY, " newY ", newY)
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
				fmt.Println("Detected explosion oldY ", oldY, " newY ", newY)
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
				fmt.Println("1. Detected Explosion!! newMonsters: ", newMonsters , " oldY ", oldY )

				newY[y] = explosion

			} else if oldY[y] == bullet && newMonsters[y-1] == monster{
				fmt.Println("2. Detected Explosion!! newMonsters: ", newMonsters , " oldY ", oldY )
				newY[y-1] = explosion
			} else {
				if oldY[y] == bullet {
					newY[y-1] = bullet
				}
			}
		}
		if x == 0 {
			fmt.Println("Moved the Bullets:  ", oldY, " newY ", newY, " newMonsters ", newMonsters)
		}

		world[x] = newY;

	}
	return false
}


func drawWorld(spaceship int, world [][]int, fb *screen.FrameBuffer) {
	for x := 0; x <= 7; x ++ {

		for  y:= 0 ; y <=  7; y++ {
			if world[x][y]  == monster {
				fb.SetPixel(x, y, color.Red)
			} else if world[x][y]  == bullet {
				fb.SetPixel(x, y, color.Blue)
			} else if world[x][y]  == explosion {
				fb.SetPixel(x, y, color.New(255,255,0))

			} else if world[x][y]  == mayhem {
				fb.SetPixel(x, y, color.New(244,66,241))

			} else {
				fb.SetPixel(x, y, color.New(0, 0, 0))
			}
		}
	}
	fb.SetPixel(spaceship, 7, color.Green)
	screen.Draw(fb)
}

func itsGameOver(world [][]int) {
	for x := 0; x <= 7; x ++ {
		world[x] =  []int{mayhem,mayhem,mayhem,mayhem,mayhem,mayhem,mayhem,mayhem};
	}
}
