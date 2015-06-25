package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const gridSize = 32
const nAgents = 10

// create an agent struct with location and orientation
type Agent struct {
	X       int // X position
	Y       int // Y position
	num     int // its number
	bearing int // its bearing 0 is left, 1 is up, 2 is right, 3 is down
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	timeSteps := 100
	pObs := 0.05
	// Environment is going to have -2 unoccupied, -1 obstacle and index of agent
	var Env [gridSize][gridSize]int

	// Fill with objects at rate pObs
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			if rand.Intn(100) <= int(pObs*100.0) {
				Env[x][y] = -1
			} else {
				Env[x][y] = -2
			}
		}
	}
	var agents [nAgents]Agent
	// place the agents
	for i := 0; i < nAgents; i++ {
		agents[i].X = rand.Intn(gridSize)
		agents[i].Y = rand.Intn(gridSize)
		agents[i].num = i
		agents[i].bearing = rand.Intn(4)
	}

	/*	for i := 0; i < nAgents; i++ {
			fmt.Printf("%d ", agents[i].X)
		}
	*/

	// do the simulation
	for t := 0; t < timeSteps; t++ {
		for i := 0; i < nAgents; i++ {
			// mark loc unoccupied
			Env[agents[i].X][agents[i].Y] = -2
			updateAgent(i, &agents[i], Env)
			// remark occupied
			Env[agents[i].X][agents[i].Y] = i
		}

		writeToJP(Env, t)
	}

}

func writeToJP(env [gridSize][gridSize]int, t int) {

	outFile := fmt.Sprint("/Users/muralib/go/output/output", strconv.Itoa(t), ".jpg")
	if t < 10 {
		outFile = fmt.Sprint("/Users/muralib/go/output/output00", strconv.Itoa(t), ".jpg")
	} else if t < 100 {
		outFile = fmt.Sprint("/Users/muralib/go/output/output0", strconv.Itoa(t), ".jpg")
	}

	out, err := os.Create(outFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	imgRect := image.Rect(0, 0, gridSize*10, gridSize*10)
	img := image.NewGray(imgRect)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			fill := &image.Uniform{color.White}
			if env[x][y] == -1 {
				fill = &image.Uniform{color.Black}
			} else if env[x][y] > 1 {
				c := color.Gray{uint8(env[x][y] * 20)}
				fill = &image.Uniform{c}
			}

			draw.Draw(img, image.Rect((x-1)*10, (y-1)*10, x*10, y*10), fill, image.ZP, draw.Src)
		}
	}

	var opt jpeg.Options

	opt.Quality = 80
	// ok, write out the data into the new JPEG file
	err = jpeg.Encode(out, img, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Generated image to output.jpg \n")
}

func updateAgent(i int, agent *Agent, env [gridSize][gridSize]int) {
	// pick a direction at random
	count := 0
	MAXCOUNT := 20
	flag := true
	for flag {
		dirX := 2 * (float64(rand.Intn(2)) - 0.5)
		dirY := 2 * (float64(rand.Intn(2)) - 0.5)
		newX := int(math.Min(math.Max(float64(agent.X)+dirX, 0), float64(gridSize-1)))
		newY := int(math.Min(math.Max(float64(agent.Y)+dirY, 0), float64(gridSize-1)))
		// if unoccupied move there
		if env[newX][newY] == -2 {
			agent.X = newX
			agent.Y = newY
			flag = false
		}
		if count > MAXCOUNT {
			flag = false
		}
		count = count + 1

	}
}
