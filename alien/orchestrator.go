package alien

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/santihernandezc/alien-invasion/world"
)

type AlienOrchestrator struct {
	Aliens    []*Alien
	positions map[string][]*Alien
	world     *world.World
	log       *log.Logger
}

func NewOrchestrator(amount int, rngSeed int64, w *world.World, log *log.Logger) (*AlienOrchestrator, error) {
	// Prevent panics
	if w == nil {
		return nil, fmt.Errorf("invalid World value: %v", w)
	}
	if len(w.Cities) < 1 {
		return nil, fmt.Errorf("invalid World value: %d cities", len(w.Cities))
	}

	alienOrchestrator := AlienOrchestrator{
		Aliens:    make([]*Alien, 0, amount),
		positions: make(map[string][]*Alien),
		world:     w,
		log:       log,
	}

	var cities []*world.City
	for _, city := range w.Cities {
		cities = append(cities, city)
	}

	rand.Seed(rngSeed)
	// Choose a random city for each alien
	for i := 0; i < amount; i++ {
		city := cities[rand.Intn(len(cities))]
		alien := Alien{
			ID:       i + 1,
			Position: city,
		}
		alienOrchestrator.Aliens = append(alienOrchestrator.Aliens, &alien)
		alienOrchestrator.positions[city.Name] = append(alienOrchestrator.positions[city.Name], &alien)
	}

	return &alienOrchestrator, nil
}

func (ao *AlienOrchestrator) UnleashAliens(maxMovements int) {
	for i := 0; i < maxMovements; i++ {
		// If there are no aliens left, the simulation is over
		if len(ao.Aliens) < 1 {
			return
		}

		for _, alien := range ao.Aliens {
			// Check if the alien was killed or stuck in the current loop
			if alien.isDeleted {
				continue
			}

			// Make the alien move
			prevPos := alien.Position.Name
			if ok := alien.move(); !ok {
				ao.log.Printf("🚷 Alien %d is trapped forever in %s", alien.ID, alien.Position.Name)
				ao.deleteAliens([]*Alien{alien})
				continue
			}

			// Remove it from the city it was previously in
			newPos := alien.Position.Name
			ao.removeAlienFromCity(prevPos, alien)
			ao.log.Printf("👾 Alien %d moved from %s to %s", alien.ID, prevPos, newPos)

			// Check if there's another alien in the new position
			rivalAliens, ok := ao.positions[newPos]
			if ok && len(rivalAliens) > 0 {
				// If two aliens find each other, the city gets destroyed and the aliens die.
				ao.log.Printf("👀 Alien %d found Alien %d in %s", alien.ID, rivalAliens[0].ID, newPos)
				ao.world.DeleteCityAndRoads(alien.Position)
				ao.log.Printf("💥 %s has been destroyed by Alien %d and Alien %d", newPos, alien.ID, rivalAliens[0].ID)

				// Since the city is destroyed, other aliens can't go to or through it
				aliensToEliminate := []*Alien{alien}
				aliensToEliminate = append(aliensToEliminate, rivalAliens...)

				if len(rivalAliens) > 1 {
					for _, ra := range rivalAliens[1:] {
						ao.log.Printf("🚷 Alien %d is trapped forever in the ruins of %s", ra.ID, alien.Position.Name)
					}
				}

				ao.deleteCityAndAliens(aliensToEliminate, newPos)
			}

			// After checking for other aliens, add alien to city
			ao.addAlienToCity(newPos, alien)
		}
	}

}

func (ao *AlienOrchestrator) deleteAliens(aliens []*Alien) {
	aliensToDelete := make(map[*Alien]struct{}, len(aliens))
	for _, alien := range aliens {
		// First, check if the alien is currently deleted
		if !alien.isDeleted {
			aliensToDelete[alien] = struct{}{}
			alien.isDeleted = true
		}
	}

	remainingAliens := make([]*Alien, 0, len(ao.Aliens)-len(aliensToDelete))
	for _, alien := range ao.Aliens {
		if _, ok := aliensToDelete[alien]; !ok {
			remainingAliens = append(remainingAliens, alien)
		}
	}

	ao.Aliens = remainingAliens
}

func (ao *AlienOrchestrator) removeAlienFromCity(prevCity string, alien *Alien) {
	var newPositionSlice []*Alien
	for _, a := range ao.positions[prevCity] {
		if a.ID != alien.ID {
			newPositionSlice = append(newPositionSlice, a)
		}
	}
	ao.positions[prevCity] = newPositionSlice
}

func (ao *AlienOrchestrator) addAlienToCity(newCity string, alien *Alien) {
	ao.positions[newCity] = append(ao.positions[newCity], alien)
}

func (ao *AlienOrchestrator) deleteCityAndAliens(alien []*Alien, cityName string) {
	ao.deleteAliens(alien)
	delete(ao.positions, cityName)
}
