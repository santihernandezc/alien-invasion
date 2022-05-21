package world

import (
	"fmt"
	"strings"
)

type cityDefinition struct {
	name        string
	neighbors   []string
	neighborMap map[string]direction
}

// parseLine takes a string defining a city and its roads.
// Each line consists of a city name and up to four roads going to neighbor cities.
func parseLine(line string) (*cityDefinition, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("invalid input: empty string")
	}

	elements := strings.Split(line, " ")
	// If there are more than 5 segments, it's an invalid definition.
	numElements := len(elements)
	if numElements > 5 {
		return nil, fmt.Errorf("invalid number of segments: %d", numElements)
	}

	var neighbors []string
	neighborMap := make(map[string]direction, numElements-1)
	for _, dirCityPair := range elements[1:] {
		// Edges should be defined in "direction=city" pairs.
		edgeSlice := strings.Split(dirCityPair, "=")
		if len(edgeSlice) != 2 {
			return nil, fmt.Errorf("invalid road definition: %q", dirCityPair)
		}

		// Convert to direction type and check if it's valid
		dir, err := stringToDirection(edgeSlice[0])
		if err != nil {
			return nil, fmt.Errorf("error converting to direction: %w", err)
		}

		neighbors = append(neighbors, edgeSlice[1])
		// Add the city name as key and the direction as value.
		neighborMap[edgeSlice[1]] = dir
	}

	return &cityDefinition{
		name:        elements[0],
		neighbors:   neighbors,
		neighborMap: neighborMap,
	}, nil
}
