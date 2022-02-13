package main

import (
	"fmt"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// Window size
const ws = 400.

// Circle size
const cs = 10.

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type PlayerStats struct {
	X float64
	Y float64
}

type ObstaclesData struct {
	Radius float64
	X      float64
	Y      float64
}

func makeWindow() *pixelgl.Window {
	// Make new window
	windowConfig := pixelgl.WindowConfig{
		Title:  "Game",
		Bounds: pixel.R(0, 0, ws, ws),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(windowConfig)
	check(err)
	return win
}

// Check if the player is colliding with anything inside the obstacles slice
func isCollision(obstacles []ObstaclesData, playerStats PlayerStats) bool {
	for _, obstacle := range obstacles {
		// Get distance between the obstacle and player
		distance := math.Sqrt(
			math.Pow(obstacle.X-playerStats.X, 2) +
				math.Pow(obstacle.Y-playerStats.Y, 2))
		// Check is the distance minus obstacle radiuses is less than zero
		if distance-cs-obstacle.Radius <= 0 {
			return true
		}
	}
	return false
}

// Check if player is inside the bounds of the window
func isInsideBounds(playerStats PlayerStats) bool {
	// Top bottom
	if playerStats.Y-cs <= 0 || playerStats.Y+cs >= ws {
		return false
	}
	// Left right
	if playerStats.X-cs <= 0 || playerStats.X+cs >= ws {
		return false
	}
	return true
}

func run() {
	win := makeWindow()

	// Create obstacles and the obstaclesData slice
	obstaclesData := []ObstaclesData{
		{Radius: 100, X: 50, Y: 100},
		{Radius: 150, X: 200, Y: ws},
	}
	obstacles := []*imdraw.IMDraw{
		imdraw.New(nil),
		imdraw.New(nil),
	}

	// Initiate player
	player := imdraw.New(nil)

	playerStats := PlayerStats{
		X: ws - 50,
		Y: 50,
	}

	// Draw our rectangles
	for i, data := range obstaclesData {
		obstacles[i].Push(pixel.V(data.X, data.Y))
		obstacles[i].Circle(data.Radius, 1)
		obstacles[i].Draw(win)
	}

	wiho := func(key pixelgl.Button) bool {
		if win.JustPressed(key) || win.Repeated(key) {
			return true
		}
		return false
	}

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) {
			break
		}

		win.Clear(pixel.RGB(0, 0, 0))

		// Redraw our rectangles
		for i := range obstacles {
			obstacles[i].Draw(win)
		}

		// Store prev position in case we moved out of bounds
		prev := playerStats

		// Movements keys
		if wiho(pixelgl.KeyW) {
			playerStats.Y += 10
		}
		if wiho(pixelgl.KeyS) {
			playerStats.Y -= 10
		}
		if wiho(pixelgl.KeyA) {
			playerStats.X -= 10
		}
		if wiho(pixelgl.KeyD) {
			playerStats.X += 10
		}

		// Check collisions
		if isInsideBounds(playerStats) && !isCollision(obstaclesData, playerStats) {
			player.Clear()
			player.Push(pixel.V(playerStats.X, playerStats.Y))
			player.Circle(cs, 1)
			player.Draw(win)
		} else {
			// Reset to the previous position
			playerStats = prev
		}

		// Update window
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
