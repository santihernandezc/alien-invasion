package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/santihernandezc/alien-invasion/alien"
	"github.com/santihernandezc/alien-invasion/world"
)

var (
	path      = flag.String("path", "config.txt", "path to the config file")
	n         = flag.Int("n", 5, "number of aliens for the simulation")
	movements = flag.Int("movements", 10000, "how many iterations this simulation is going to run")
	directed  = flag.Bool("directed", true, "use a directed graph")
	o         = flag.String("o", "", "path to the output file (optional)")
)

func init() {
	flag.Parse()
}

func main() {
	log := log.New(os.Stdout, "", 0)

	// Read and parse file into World map.
	log.Printf("Initializing world from file with path %q", *path)
	file, err := os.Open(*path)
	if err != nil {
		log.Fatalf("Error opening file in path %s: %v", *path, err)
	}
	defer file.Close()

	worldMap, err := world.NewFromReader(file, *directed)
	if err != nil {
		log.Fatalf("Error reading and parsing file: %v", err)
	}

	// Make aliens
	log.Printf("Initializing %d aliens", *n)
	rngSeed := time.Now().UnixNano()
	ao, err := alien.NewOrchestrator(*n, rngSeed, worldMap, log)
	if err != nil {
		log.Fatalf("error creating aliens: %v", err)
	}

	// Start the simulation
	log.Printf("Unleashing %d aliens\n\n", *n)
	ao.UnleashAliens(*movements)

	worldString := worldMap.String()
	log.Print("Simulation done, this is what's left of the world...\n", worldString)

	if *o != "" {
		os.WriteFile(*o, []byte(worldMap.String()), 0777)
	}
}
