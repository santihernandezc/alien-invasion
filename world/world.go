package world

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// World is a graph with interconnected cities.
type World struct {
	Cities          map[string]*City
	DestroyedCities []*City
	directed        bool
}

type position struct {
	X int32
	Y int32
}

type cityDefinition struct {
	Name        string   `json:"name"`
	Neighbors   []string `json:"neighbors"`
	neighborMap map[string]direction
}

// City is an edge on the graph.
type City struct {
	Name        string
	Neighbors   []*City
	neighborMap map[*City]direction
	Position    rl.Vector2
}

// NewFromBytes returns a new World based on raw bytes.
func NewFromBytes(b []byte, isDirected bool, width int32, height int32) (*World, error) {
	world := World{
		Cities:   make(map[string]*City),
		directed: isDirected,
	}

	// Parse city and roads from each file line
	var cityDefs []cityDefinition
	if err := json.Unmarshal(b, &cityDefs); err != nil {
		return nil, fmt.Errorf("error unmarshaling bytes: %w", err)
	}

	for _, cityDef := range cityDefs {
		world.addCityAndRoads(&cityDef, width, height)
	}

	return &world, nil
}

func (w *World) addCityAndRoads(cityDef *cityDefinition, width int32, height int32) {
	// Create or retrieve city, name must be unique
	cityFrom, ok := w.Cities[cityDef.Name]
	if !ok {
		cityFrom = &City{
			Position:    rl.NewVector2(rand.Float32()*float32(width-100)+50, rand.Float32()*float32(height-100)+50),
			Name:        cityDef.Name,
			neighborMap: make(map[*City]direction, maxRoads),
		}
		// Add newly created city to World
		w.Cities[cityFrom.Name] = cityFrom
	}

	// Add roads to neighbor cities
	for _, neighborName := range cityDef.Neighbors {
		// Define directions both ways
		cityToNeighborDir := cityDef.neighborMap[neighborName]
		neighborToCityDir := oppositeDirectionMap[cityToNeighborDir]

		cityTo, ok := w.Cities[neighborName]
		if !ok {
			// If the neighbor city hasn't been created yet,
			// create it and add it to the World before proceeding.
			cityTo = &City{
				Position:    rl.NewVector2(rand.Float32()*float32(width-100)+50, rand.Float32()*float32(height-100)+50),
				Name:        neighborName,
				neighborMap: make(map[*City]direction, maxRoads),
			}
			w.Cities[neighborName] = cityTo
		}

		// Append to city neighbors only if it's not already there
		if _, ok := cityFrom.neighborMap[cityTo]; !ok {
			cityFrom.Neighbors = append(cityFrom.Neighbors, cityTo)
			cityFrom.neighborMap[cityTo] = cityToNeighborDir

			// If the graph is non-directed, make the connection bi-directional
			if !w.directed {
				cityTo.Neighbors = append(cityTo.Neighbors, cityFrom)
				cityTo.neighborMap[cityFrom] = neighborToCityDir
			}
		}
	}
}

// DeleteCityAndRoads removes a city and all its edges from the World.
func (w *World) DeleteCityAndRoads(city *City) {
	// Delete City from the World's City map
	w.DestroyedCities = append(w.DestroyedCities, city)
	delete(w.Cities, city.Name)

	// If it's not a directed graph, delete all roads to the city
	// from its neighbors' adjacency lists
	if !w.directed {
		for _, neighbor := range city.Neighbors {
			newNeighbors := make([]*City, 0, len(neighbor.Neighbors)-1)
			for _, nn := range neighbor.Neighbors {
				if nn != city {
					newNeighbors = append(newNeighbors, nn)
				}
			}

			neighbor.Neighbors = newNeighbors
			delete(neighbor.neighborMap, city)
		}
		return
	}

	// If it's directed, we have to check node by node
	// and delete every reference to the city
	for _, c := range w.Cities {
		newNeighbors := make([]*City, 0, len(c.Neighbors))
		for _, n := range c.Neighbors {
			if n != city {
				newNeighbors = append(newNeighbors, n)
			}
		}

		if len(newNeighbors) != len(c.Neighbors) {
			c.Neighbors = newNeighbors
			delete(c.neighborMap, city)
		}
	}
}

// String returns the string representation of the World
// using the same format as the input file.
func (w *World) String() string {
	var builder strings.Builder
	for _, city := range w.Cities {
		fmt.Fprintf(&builder, "%s", city.Name)
		for n, dir := range city.neighborMap {
			fmt.Fprintf(&builder, " %s=%s", dir, n.Name)
		}
		fmt.Fprintln(&builder)
	}

	return builder.String()
}
