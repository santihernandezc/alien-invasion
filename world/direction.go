package world

import "fmt"

type direction string

const (
	north direction = "north"
	south direction = "south"
	east  direction = "east"
	west  direction = "west"
)

var (
	// Useful for knowing all the valid directions and the opposite value of each one.
	oppositeDirectionMap = map[direction]direction{
		north: south,
		south: north,
		east:  west,
		west:  east,
	}
	// The maximum amount of neighbors a city can have is derived from the number of available directions.
	maxRoads = len(oppositeDirectionMap)
)

func stringToDirection(str string) (direction, error) {
	dir := direction(str)
	if _, ok := oppositeDirectionMap[dir]; !ok {
		return "", fmt.Errorf("cannot convert string %q to direction type", str)
	}

	return dir, nil
}
