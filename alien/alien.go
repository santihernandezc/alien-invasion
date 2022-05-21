package alien

import (
	"math/rand"

	"github.com/santihernandezc/alien-invasion/world"
)

type Alien struct {
	ID        int
	Position  *world.City
	isDeleted bool
}

func (a *Alien) move() (ok bool) {
	// Check whether the alien is trapped
	neighbors := a.Position.Neighbors
	if len(neighbors) < 1 {
		return false
	}

	// Choose a random city
	a.Position = neighbors[rand.Intn(len(neighbors))]
	return true
}
