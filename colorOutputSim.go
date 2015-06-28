package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
)

const gridSize = 32
const nAgents = 10
const destname = "/Users/muralib/go/output/output.gif"
const delay = 10

var Env [gridSize][gridSize]int
var frames []*image.Paletted
var colors [nAgents]color.Color

// create an agent struct with location and orientation
type Agent struct {
	X       int // X position
	Y       int // Y position
	num     int // its number
	bearing int // its bearing 0 is left, 1 is up, 2 is right, 3 is down
}

var agents [nAgents]Agent

func main() {
	runtime.GOMAXPROCS(10)
	rand.Seed(10)
	nWorkers := 4
	timeSteps := 100
	pObs := 0.05
	// Environment is going to have -2 unoccupied, -1 obstacle and index of agent

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
	// place the agents
	for i := 0; i < nAgents; i++ {
		colors[i] = color.RGBA{uint8(rand.Intn(128)), uint8(rand.Intn(128)), uint8(rand.Intn(128)), 10}
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
	// create one channel per worked

	jobs := make(chan int, nAgents)
	results := make(chan Agent, nAgents)

	for i := 0; i < nWorkers; i++ {
		go worker(i, jobs, results)
	}
	simulateLoop(timeSteps, jobs, results)
	createGIF()
	fmt.Println("Creating GIF")
}

func createGIF() {
	delays := make([]int, len(frames))
	for j, _ := range delays {
		delays[j] = delay
	}
	opfile, err := os.Create(destname)
	if err != nil {
		log.Fatalf("Error creating the destination file %s : %s", destname, err)
	}

	if err := gif.EncodeAll(opfile, &gif.GIF{frames, delays, 0}); err != nil {
		log.Printf("Error encoding output into animated gif :%s", err)
	}
	opfile.Close()
}

func worker(id int, jobs <-chan int, results chan Agent) {
	for j := range jobs {
		updateAgent(j, results)
		fmt.Println("worker", id, "processed agent", j)
	}
}

func simulateLoop(timeSteps int, jobs chan int, results chan Agent) {

	for t := 0; t < timeSteps; t++ {
		fmt.Println("Time Step :", t)
		for i := 0; i < nAgents; i++ {
			// add agent to jobs
			jobs <- i
		}

		for i := 0; i < nAgents; i++ {
			tmp := <-results
			j := tmp.num
			Env[agents[j].X][agents[j].Y] = -2
			//fmt.Println("Iter ", t, " Agent ", tmp.num, " X ", tmp.X, " Y ", tmp.Y)
			agents[j].X = tmp.X
			agents[j].Y = tmp.Y
			Env[agents[j].X][agents[j].Y] = j
		}
		writeToJP()
		if t+1 == timeSteps {
			close(jobs)
		}
	}
}

func updateAgent(i int, results chan Agent) {
	agent := agents[i]
	// pick a direction at random
	count := 0
	MAXCOUNT := 200
	flag := true
	var outAgent Agent
	outAgent.num = i
	for flag {
		dirX := 2 * (float64(rand.Intn(2)) - 0.5)
		dirY := 2 * (float64(rand.Intn(2)) - 0.5)
		newX := int(math.Min(math.Max(float64(agent.X)+dirX, 0), float64(gridSize-1)))
		newY := int(math.Min(math.Max(float64(agent.Y)+dirY, 0), float64(gridSize-1)))
		// if unoccupied move there
		if Env[newX][newY] == -2 {
			outAgent.X = newX
			outAgent.Y = newY
			flag = false
		}
		if count > MAXCOUNT {
			outAgent.X = agent.X
			outAgent.Y = agent.Y
			flag = false
		}
		count = count + 1
	}
	results <- outAgent
}

func writeToJP() {

	imgRect := image.Rect(0, 0, gridSize*10, gridSize*10)
	img := image.NewRGBA(imgRect)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			fill := &image.Uniform{color.White}
			if Env[x][y] == -1 {
				fill = &image.Uniform{color.Black}
			} else if Env[x][y] > 1 {
				//c := color.RGBA{uint8(Env[x][y] * 20), 10, 10, 10}
				c := colors[Env[x][y]-1]
				fill = &image.Uniform{c}
			}

			draw.Draw(img, image.Rect(x*10, y*10, (x+1)*10, (y+1)*10), fill, image.ZP, draw.Src)
		}
	}
	buf := bytes.Buffer{}
	// ok, write out the data into the new JPEG file
	err := gif.Encode(&buf, img, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tmpimg, err := gif.Decode(&buf)
	if err != nil {
		log.Printf("Skipping frame due to weird error reading the temporary gif :%s", err)
	}
	frames = append(frames, tmpimg.(*image.Paletted))
}
