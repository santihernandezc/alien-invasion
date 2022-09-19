package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/santihernandezc/alien-invasion/alien"
	"github.com/santihernandezc/alien-invasion/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	path     = flag.String("path", "config.json", "path to the json config file")
	n        = flag.Int("n", 5, "number of aliens for the simulation")
	directed = flag.Bool("directed", false, "use a directed graph")
)

func init() {
	flag.Parse()
}

func main() {
	rand.Seed(time.Now().Unix())
	log := log.New(os.Stdout, "", 0)

	// Read and parse file into World map.
	log.Printf("Initializing world from file %q", *path)
	b, err := os.ReadFile(*path)
	if err != nil {
		log.Fatalf("Error opening file in path %s: %v", *path, err)
	}

	worldMap, err := world.NewFromBytes(b, *directed, 800, 450)
	if err != nil {
		log.Fatalf("Error reading and parsing file: %v", err)
	}

	// Init window
	rl.InitWindow(800, 450, "Alien Invasion")
	rl.SetTargetFPS(60)

	// Load textures
	alienImg := rl.LoadImage("./assets/alien.png")
	alienTexture := rl.LoadTextureFromImage(alienImg)
	alienTexture.Width = int32(float32(alienTexture.Width) * 0.2)
	alienTexture.Height = int32(float32(alienTexture.Height) * 0.2)
	rl.UnloadImage(alienImg)
	explosionImg := rl.LoadImage("./assets/explosion.png")
	explosionTexture := rl.LoadTextureFromImage(explosionImg)
	explosionTexture.Width = int32(float32(explosionTexture.Width) * 0.1)
	explosionTexture.Height = int32(float32(explosionTexture.Height) * 0.1)
	rl.UnloadImage(explosionImg)

	// Instantiate aliens and seed randomness
	log.Printf("Initializing %d aliens", *n)
	rngSeed := time.Now().UnixNano()
	ao, err := alien.NewOrchestrator(*n, rngSeed, worldMap, alienTexture, log)
	if err != nil {
		log.Fatalf("error creating aliens: %v", err)
	}

	var counter int

	stepSignal := make(chan bool)
	// go tick(stepSignal)
	go listener(stepSignal)

	// Draw
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		drawMap(worldMap, explosionTexture)

		for _, alien := range ao.Aliens {
			alien.Draw()
		}
		rl.EndDrawing()

		select {
		case <-stepSignal:
			if len(ao.Aliens) > 0 {
				ao.Step(ao.Aliens[counter%len(ao.Aliens)])
				counter++
			}
		default:
			continue
		}
	}

	rl.UnloadTexture(alienTexture)
	rl.CloseWindow()
}

func drawMap(worldMap *world.World, explosionTexture rl.Texture2D) {
	for _, city := range worldMap.DestroyedCities {
		rl.DrawText(city.Name, int32(city.Position.X)+10, int32(city.Position.Y)+10, 10, rl.Black)
		rl.DrawTexture(explosionTexture, int32(city.Position.X)-10, int32(city.Position.Y)-10, rl.White)
	}

	for _, city := range worldMap.Cities {
		rl.DrawText(city.Name, int32(city.Position.X)+10, int32(city.Position.Y)+10, 10, rl.Black)
		rl.DrawCircleLines(int32(city.Position.X), int32(city.Position.Y), 10, rl.Black)
		for _, neighbor := range city.Neighbors {
			rl.DrawLine(int32(city.Position.X), int32(city.Position.Y), int32(neighbor.Position.X), int32(neighbor.Position.Y), rl.Gray)
		}
	}
}

func listener(c chan bool) {
	for {
		if rl.IsKeyReleased(32) {
			c <- true
		}
	}
}

// func tick(c chan bool) {
// 	for {
// 		c <- true
// 		time.Sleep(1 * time.Second)
// 	}
// }
