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

	z := 0

	for {
		select {
		case <-tickChan:
			updateBullets(&world, fb)
			detectExplosions(&world, fb)
			z = z + 1
			if ((z % 3) == 0) {
				updateMonsters(&world, fb)
				detectExplosions(&world, fb)
			}
			if ((z % 8) == 0) {
				createMonster(&world, fb)
			}

		case <-signals:
			fmt.Println("")
			screen.Clear()

			// Exit the loop
			return
		case e := <-input.Events:
			switch e.Code {

			case stick.Up:
				fmt.Println("Shoot!")
				b := Bullet{X: world.X, Y: 7, Color: color.Blue, Explode:false}
				world.Bullets = append(world.Bullets, b)

			case stick.Right:
				//fmt.Println("1. L: ", x, " y ", y)
				fb.SetPixel(x, y, color.New(0, 0, 0))
				x = (x + 1) % 8
				if x < 0 {
					x = x - 8
				}
				world.X = x
			case stick.Left:
				//fmt.Println("1. R: ", x, " y ", y)
				fb.SetPixel(x, y, color.New(0, 0, 0))
				x = (x - 1) % 8
				if x < 0 {
					x = x + 8
				}
				world.X = x
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
		if b.Explode {
			fb.SetPixel(x, y, color.New(255,255,0))
		} else if x >= 0 && yy >= 0 {
			bullet := Bullet{X: x, Y: yy, Color: color.Blue}
			remainingBullets = append(remainingBullets, bullet)
			fb.SetPixel(x, yy, b.Color)
		}
	}
	world.Bullets = remainingBullets
}

func detectExplosions(world *World, fb *screen.FrameBuffer) {
	for i := 0; i < len(world.Bullets); i++ {
		bullet := &world.Bullets[i]

		for x := 0;  x < len(world.Monsters); x++ {
			monster := &world.Monsters[x]
			if (bullet.Y == monster.Y && bullet.X == monster.X)  {
				//fmt.Println("detected explosion! Monster.Y:", monster.Y, " monster.X ", monster.X,
				//	" bullet.Y ", bullet.Y, " bullet.X  ", bullet.X)

				monster.Explode = true
				bullet.Explode = true
			}
		}
	}
}


func updateMonsters(world *World, fb *screen.FrameBuffer) {
	remainingMonsters := []Monster{}

	for i := 0; i < len(world.Monsters); i++ {
		m := world.Monsters[i]
		x := m.X
		y := m.Y

		counter := m.Counter

		if m.Explode  && counter >= 0{

			//fmt.Println("12. Render explosion! monsterY:", y, " M.X ", x, " counter: ", counter)

			expColor := color.New(255, 255, 0)

			if (m.Counter <= 10) {
				expColor = color.New(255, 153, 0)
			}

			fb.SetPixel(x, y, expColor) // yellow
			m := Monster{X: x, Y: y, Color: expColor, Explode: true, Counter: (counter - 1) }
			remainingMonsters = append(remainingMonsters, m)

		} else {
			fb.SetPixel(x, y, color.New(0, 0, 0))
		}


		yy := y + 1

		if !m.Explode && x < 8 && yy < 8 {
			m := Monster{X: x, Y: yy, Color: color.Red}
			remainingMonsters = append(remainingMonsters, m)
			fb.SetPixel(x, yy, m.Color)
		}
	}
	world.Monsters = remainingMonsters
}

func drawWorld(world World, fb *screen.FrameBuffer) {
	fb.SetPixel(world.X, 7, color.Green)
	screen.Draw(fb)
}

func createMonster(world *World, fb *screen.FrameBuffer) {
	//fmt.Println("Monster!")
	x := rand.Int() % 8
	y := 0

	for  {
		if ( world.PreviousMonsterX == x) {
			x = rand.Int() % 8
		} else {
			break
		}
	}

	world.PreviousMonsterX = x;

	m := Monster{X: x, Y: y, Color: color.Red, Explode: false, Counter: 20}
	world.Monsters = append(world.Monsters, m)

	fb.SetPixel(m.X, m.Y, m.Color)
}

type Bullet struct {
	X     int
	Y     int
	Color color.Color
	Explode bool
}

type Monster struct {
	X     int
	Y     int
	Explode bool
	Counter int
	Color color.Color
}

type World struct {
	X       int
	Bullets []Bullet
	Monsters []Monster
	PreviousMonsterX  int
}


