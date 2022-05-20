package world

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// World is a graph with interconnected cities.
type World struct {
	Cities   map[string]*City
	directed bool
}

// City is an edge on the graph.
type City struct {
	Name        string
	Neighbors   []*City
	neighborMap map[*City]direction
}

// NewFromReader returns a new World based on the contents of an io.Reader.
func NewFromReader(reader io.Reader, isDirected bool) (*World, error) {
	world := World{
		Cities:   make(map[string]*City),
		directed: isDirected,
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		// Parse city and roads from each file line
		cityDef, err := parseLine(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %w", err)
		}

		world.addCityAndRoads(cityDef)
	}

	return &world, nil
}

func (w *World) addCityAndRoads(cityDef *cityDefinition) {
	// Create or retrieve city, name must be unique.
	cityFrom, ok := w.Cities[cityDef.name]
	if !ok {
		cityFrom = &City{
			Name:        cityDef.name,
			neighborMap: make(map[*City]direction, maxRoads),
		}
	}

	// Add roads to the city
	for _, neighborName := range cityDef.neighbors {
		cityTo, ok := w.Cities[neighborName]
		if !ok {
			// If the neighbor city hasn't been created yet,
			// create it and add it to the list.
			cityTo = &City{
				Name:        neighborName,
				neighborMap: make(map[*City]direction, maxRoads),
			}
			w.Cities[neighborName] = cityTo
			// If the graph is non-directed, make the connection bi-directional
			if !w.directed {
				cityTo.Neighbors = append(cityTo.Neighbors, cityFrom)
				cityTo.neighborMap[cityFrom] = oppositeDirectionMap[cityDef.neighborMap[cityTo.Name]]
			}
			cityFrom.Neighbors = append(cityFrom.Neighbors, cityTo)
			cityFrom.neighborMap[cityTo] = cityDef.neighborMap[cityTo.Name]
			continue
		}

		// If the city already exists,
		// append to neighbors list only if it's not already there
		var alreadyNeighbor bool
		for _, n := range cityFrom.Neighbors {
			if n == cityTo {
				alreadyNeighbor = true
			}
		}

		if !alreadyNeighbor {
			if !w.directed {
				cityTo.Neighbors = append(cityTo.Neighbors, cityFrom)
				cityTo.neighborMap[cityFrom] = oppositeDirectionMap[cityDef.neighborMap[cityTo.Name]]
			}
			cityFrom.Neighbors = append(cityFrom.Neighbors, cityTo)
			cityFrom.neighborMap[cityTo] = cityDef.neighborMap[cityTo.Name]
		}
	}

	// Add newly created City
	w.Cities[cityFrom.Name] = cityFrom
}

func (w *World) DeleteCityAndRoads(city *City) {
	// Delete City from world map
	delete(w.Cities, city.Name)

	// If it's not a directed graph, delete all roads to the city
	// from its neighbors' adjacency list
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
		delete(c.neighborMap, city)
		newNeighbors := make([]*City, 0, len(c.Neighbors))
		for _, n := range c.Neighbors {
			if n != city {
				newNeighbors = append(newNeighbors, n)
			}
		}

		if len(newNeighbors) != len(c.Neighbors) {
			c.Neighbors = newNeighbors
		}
	}

}

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
