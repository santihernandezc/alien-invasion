package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/santihernandezc/alien-invasion/alien"
	"github.com/santihernandezc/alien-invasion/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	path      = flag.String("path", "config.json", "path to the json config file")
	n         = flag.Int("n", 5, "number of aliens for the simulation")
	movements = flag.Int("movements", 10000, "how many iterations this simulation is going to run")
	directed  = flag.Bool("directed", false, "use a directed graph")
	o         = flag.String("o", "", "path to the output file (optional)")
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

	// Instantiate aliens and seed randomness
	log.Printf("Initializing %d aliens", *n)
	rngSeed := time.Now().UnixNano()
	ao, err := alien.NewOrchestrator(*n, rngSeed, worldMap, log)
	if err != nil {
		log.Fatalf("error creating aliens: %v", err)
	}
	var counter int

	// Init window
	rl.InitWindow(800, 450, "Alien Invasion")
	rl.SetTargetFPS(60)

	// Load textures
	alienImg := rl.LoadImage("./assets/alien.png")
	alienTexture := rl.LoadTextureFromImage(alienImg)
	alienTexture.Width = int32(float32(alienTexture.Width) * 0.2)
	alienTexture.Height = int32(float32(alienTexture.Height) * 0.2)
	rl.UnloadImage(alienImg)

	// Draw
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawText(fmt.Sprintf("%f", rl.GetFPS()), 10, 10, 20, rl.Black)
		drawMap(worldMap)
		drawAliens(ao.Aliens, alienTexture)
		rl.EndDrawing()

		if rl.IsKeyReleased(32) {
			if len(ao.Aliens) >= 1 {
				fmt.Println("counter:", counter)
				ao.Step(ao.Aliens[counter%len(ao.Aliens)])
				counter++
			}
		}
	}

	rl.UnloadTexture(alienTexture)
	rl.CloseWindow()

	// Start the simulation
	// log.Printf("Unleashing %d aliens\n\n", *n)
	// ao.UnleashAliens(*movements)

	// worldString := worldMap.String()
	// log.Print("Simulation done, this is what's left of the world...\n\n", worldString)

	// if *o != "" {
	// 	os.WriteFile(*o, []byte(worldMap.String()), 0777)
	// }
}

func drawAliens(aliens []*alien.Alien, texture rl.Texture2D) {
	for _, alien := range aliens {
		x := alien.Position.Position.X
		y := alien.Position.Position.Y
		rl.DrawTexture(texture, x-texture.Width/2, y-texture.Height/2, rl.White)
	}
}

func drawMap(worldMap *world.World) {
	for _, city := range worldMap.Cities {
		rl.DrawText(city.Name, city.Position.X+10, city.Position.Y+10, 10, rl.Black)
		rl.DrawCircleLines(city.Position.X, city.Position.Y, 10, rl.Black)
		for _, neighbor := range city.Neighbors {
			rl.DrawLine(city.Position.X, city.Position.Y, neighbor.Position.X, neighbor.Position.Y, rl.Gray)
		}
	}
}
