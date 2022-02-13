package main

import (
	"fmt"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	ws = 800. // Windows size
	ps = 10.  // Player size
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type PlayerStats struct {
	X, Y float64
}

type ObstacleData struct {
	Radius, X, Y float64
}

type Ray struct {
	Pos         pixel.Vec
	HasCollided bool
	imd         *imdraw.IMDraw
}

type Game struct {
	ObstaclesData []ObstacleData
	Obstacles     []*imdraw.IMDraw
	Player        *imdraw.IMDraw
	PlayerStats   PlayerStats
	Win           *pixelgl.Window
	Rays          []Ray
	RaySpeed      float64
	GameOver      bool
}

// Check if the player is colliding with anything inside the obstacles slice
func (g *Game) isObstacleCollision() bool {
	for _, obstacle := range g.ObstaclesData {
		// Get distance between the obstacle and player
		distance := math.Sqrt(
			math.Pow(obstacle.X-g.PlayerStats.X, 2) +
				math.Pow(obstacle.Y-g.PlayerStats.Y, 2))

		// Check is the distance minus obstacle radiuses is less than zero
		if distance-ps-obstacle.Radius <= 0 {
			return true
		}
	}

	return false
}

// Check if ray is colliding with any of the obstacles
func (g *Game) isRayCollision(ray Ray) bool {
	for _, obstacle := range g.ObstaclesData {
		distance := math.Sqrt(
			math.Pow(obstacle.X-ray.Pos.X, 2) +
				math.Pow(obstacle.Y-ray.Pos.Y, 2))

		if distance-obstacle.Radius <= 0 {
			return true
		}
	}

	return false
}

func (g *Game) isDead(ray Ray) {
	distance := math.Sqrt(
		math.Pow(g.PlayerStats.X-ray.Pos.X, 2) +
			math.Pow(g.PlayerStats.Y-ray.Pos.Y, 2))

	if distance-ps <= 0 {
		g.GameOver = true
	}
}

// Check if player is inside the bounds of the window
func (g *Game) isInsideBounds() bool {
	// Top, bottom, left, right
	if g.PlayerStats.Y-ps <= 0 || g.PlayerStats.Y+ps >= ws ||
		g.PlayerStats.X-ps <= 0 || g.PlayerStats.X+ps >= ws {
		return false
	} else {
		return true
	}
}

func (g *Game) Run() {
	// Draw our obstacles
	for i, data := range g.ObstaclesData {
		g.Obstacles[i].Push(pixel.V(data.X, data.Y))
		g.Obstacles[i].Circle(data.Radius, 5)
		g.Obstacles[i].Draw(g.Win)
	}

	// Draw our rays
	for i, ray := range g.Rays {
		g.Rays[i].imd.Push(pixel.V(ray.Pos.X, ray.Pos.Y))
		g.Rays[i].imd.Line(1)
		g.Rays[i].imd.Draw(g.Win)
	}

	registerKeyPress := func(key pixelgl.Button) bool {
		if g.Win.JustPressed(key) || g.Win.Repeated(key) {
			return true
		}
		return false
	}

	g.Win.SetSmooth(true)

	speed := 1.

	for !g.Win.Closed() {
		if registerKeyPress(pixelgl.KeyEscape) {
			break
		}

		g.Win.Clear(pixel.RGB(0, 0, 0))

		if !g.GameOver {
			// Redraw our rectangles
			for i := range g.Obstacles {
				g.Obstacles[i].Draw(g.Win)
			}

			// Redraw our rays
			for i, ray := range g.Rays {
				if ray.Pos.X > ws {
					g.GameOver = true
				}
				if ray.HasCollided == false {
					g.Rays[i].Pos.X += g.RaySpeed
					g.Rays[i].imd.Push(pixel.V(ray.Pos.X, ray.Pos.Y))
					g.Rays[i].imd.Push(pixel.V(ray.Pos.X+g.RaySpeed, ray.Pos.Y))
					g.Rays[i].imd.Line(1.0001) // The float fixed a weird issue
					g.Rays[i].imd.Draw(g.Win)

					if g.isRayCollision(ray) {
						g.Rays[i].HasCollided = true
					}
					g.isDead(ray)
				} else {
					g.Rays[i].imd.Line(1)
					g.Rays[i].imd.Draw(g.Win)
				}
			}

			// Store prev position in case we moved out of bounds
			prev := g.PlayerStats

			if registerKeyPress(pixelgl.KeyLeftSuper) {
				if speed == 1. {
					speed = 20.
				} else {
					speed = 1.
				}
			}

			// Movements keys
			if registerKeyPress(pixelgl.KeyW) {
				g.PlayerStats.Y += speed
			}
			if registerKeyPress(pixelgl.KeyS) {
				g.PlayerStats.Y -= speed
			}
			if registerKeyPress(pixelgl.KeyA) {
				g.PlayerStats.X -= speed
			}
			if registerKeyPress(pixelgl.KeyD) {
				g.PlayerStats.X += speed
			}

			// Check collisions
			if g.isInsideBounds() && !g.isObstacleCollision() {
				g.Player.Clear()
				g.Player.Push(pixel.V(g.PlayerStats.X, g.PlayerStats.Y))
				g.Player.Circle(ps, 1)
				g.Player.Draw(g.Win)
			} else {
				// Reset to the previous position
				g.PlayerStats = prev
			}
		} else {
			g.Player.Circle(ps, 1)
			g.Player.Draw(g.Win)
		}
		// Update window
		g.Win.Update()
	}
}

func main() {
	pixelgl.Run(func() {

		// Create window
		win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
			Title:  "Game",
			Bounds: pixel.R(0, 0, ws, ws),
			VSync:  true,
		})
		check(err)

		// Create the rays
		rays := []Ray{}
		for i := 10.; i < ws; i += 20 {
			rays = append(
				rays,
				Ray{
					Pos:         pixel.V(0, i),
					HasCollided: false,
					imd:         imdraw.New(nil),
				},
			)
		}

		// Initiate game
		game := Game{
			ObstaclesData: []ObstacleData{
				{Radius: 20, X: 200, Y: 100},
				{Radius: 50, X: 200, Y: ws - 100},
			},
			Obstacles: []*imdraw.IMDraw{
				imdraw.New(nil),
				imdraw.New(nil),
			},
			Player: imdraw.New(nil),
			PlayerStats: PlayerStats{
				X: ws - 50,
				Y: ws / 2,
			},
			Win:      win,
			Rays:     rays,
			RaySpeed: 2,
			GameOver: false,
		}

		game.Run()
	})
}
