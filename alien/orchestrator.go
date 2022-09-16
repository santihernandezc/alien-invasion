package alien

import (
	"fmt"
	"log"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/santihernandezc/alien-invasion/world"
)

// AlienOrchestrator is in charge of managing the state and behavior of aliens.
// It contains the main logic to run and stop the simulation.
type AlienOrchestrator struct {
	Aliens []*Alien

	// positions maps a city name with the aliens in that city
	world     *world.World
	positions map[string][]*Alien
	log       *log.Logger
}

func NewOrchestrator(amount int, rngSeed int64, w *world.World, alienTexture rl.Texture2D, log *log.Logger) (*AlienOrchestrator, error) {
	// Prevent panics
	if w == nil {
		return nil, fmt.Errorf("invalid World value: <nil>")
	}
	if log == nil {
		return nil, fmt.Errorf("invalid value for logger: <nil>")
	}
	if len(w.Cities) == 0 {
		return nil, fmt.Errorf("invalid World value: 0 cities")
	}

	alienOrchestrator := AlienOrchestrator{
		Aliens:    make([]*Alien, 0, amount),
		positions: make(map[string][]*Alien, len(w.Cities)),
		world:     w,
		log:       log,
	}

	// Make a slice to choose random cities as starting positions
	var cities []*world.City
	for _, city := range w.Cities {
		cities = append(cities, city)
	}

	// To make things more random-like, use the seed
	rand.Seed(rngSeed)

	// Choose a random city for each alien.
	// Start from 1 instead of 0 to use the same value for the alien's ID.
	for i := 1; i <= amount; i++ {
		city := cities[rand.Intn(len(cities))]
		alien := Alien{
			ID:           i,
			City:         city,
			Position:     rl.NewVector2(float32(city.Position.X), float32(city.Position.Y)),
			NextPosition: rl.NewVector2(float32(city.Position.X), float32(city.Position.Y)),
			Texture:      alienTexture,
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
			prevPos := alien.City.Name
			if ok := alien.move(); !ok {
				ao.log.Printf("🚷 Alien %d is trapped forever in %s", alien.ID, alien.City.Name)
				ao.deleteAliens([]*Alien{alien})
				continue
			}

			// Remove it from the city it was previously in
			newPos := alien.City.Name
			ao.removeAlienFromCity(prevPos, alien)
			ao.log.Printf("👾 Alien %d moved from %s to %s", alien.ID, prevPos, newPos)

			// Check if there's another alien in the new position
			rivalAliens, ok := ao.positions[newPos]
			if ok && len(rivalAliens) > 0 {
				// If two aliens find each other, the city gets destroyed and the aliens die.
				ao.log.Printf("👀 Alien %d found Alien %d in %s", alien.ID, rivalAliens[0].ID, newPos)
				ao.world.DeleteCityAndRoads(alien.City)
				ao.log.Printf("💥 %s has been destroyed by Alien %d and Alien %d", newPos, alien.ID, rivalAliens[0].ID)

				// Since the city is destroyed, other aliens can't go to or through it
				aliensToEliminate := append(rivalAliens, alien)

				if len(rivalAliens) > 1 {
					for _, ra := range rivalAliens[1:] {
						ao.log.Printf("🚷 Alien %d is trapped forever in the ruins of %s", ra.ID, alien.City.Name)
					}
				}

				ao.deleteCityAndAliens(aliensToEliminate, newPos)
			}

			// After checking for other aliens, add alien to city
			ao.addAlienToCity(newPos, alien)
		}
	}
}

func (ao *AlienOrchestrator) Step(alien *Alien) {
	// Check if the alien was killed or stuck in the current loop
	if alien.isDeleted {
		return
	}

	// Make the alien move
	prevPos := alien.City.Name
	if ok := alien.move(); !ok {
		ao.log.Printf("🚷 Alien %d is trapped forever in %s", alien.ID, alien.City.Name)
		ao.deleteAliens([]*Alien{alien})
		return
	}

	// Remove it from the city it was previously in
	newPos := alien.City.Name
	ao.removeAlienFromCity(prevPos, alien)
	ao.log.Printf("👾 Alien %d moved from %s to %s", alien.ID, prevPos, newPos)

	// Check if there's another alien in the new position
	rivalAliens, ok := ao.positions[newPos]
	if ok && len(rivalAliens) > 0 {
		// If two aliens find each other, the city gets destroyed and the aliens die.
		ao.log.Printf("👀 Alien %d found Alien %d in %s", alien.ID, rivalAliens[0].ID, newPos)
		ao.world.DeleteCityAndRoads(alien.City)
		ao.log.Printf("💥 %s has been destroyed by Alien %d and Alien %d", newPos, alien.ID, rivalAliens[0].ID)

		// Since the city is destroyed, other aliens can't go to or through it
		aliensToEliminate := append(rivalAliens, alien)

		if len(rivalAliens) > 1 {
			for _, ra := range rivalAliens[1:] {
				ao.log.Printf("🚷 Alien %d is trapped forever in the ruins of %s", ra.ID, alien.City.Name)
			}
		}

		ao.deleteCityAndAliens(aliensToEliminate, newPos)
	}

	// After checking for other aliens, add alien to city
	ao.addAlienToCity(newPos, alien)
}

func (ao *AlienOrchestrator) deleteAliens(aliens []*Alien) {
	aliensToDelete := make(map[*Alien]struct{}, len(aliens))

	// First, check if the aliens are already in deleted state
	for _, alien := range aliens {
		if !alien.isDeleted {
			aliensToDelete[alien] = struct{}{}
			alien.isDeleted = true
		}
	}

	// Make a slice with the aliens that are still active
	remainingAliens := make([]*Alien, 0, len(ao.Aliens)-len(aliensToDelete))
	for _, alien := range ao.Aliens {
		if _, ok := aliensToDelete[alien]; !ok {
			remainingAliens = append(remainingAliens, alien)
		}
	}

	ao.Aliens = remainingAliens
}

func (ao *AlienOrchestrator) removeAlienFromCity(prevCity string, alien *Alien) {
	// Filter out the alien from the slice corresponding to the previous city
	newPositionSlice := make([]*Alien, 0, len(ao.positions)-1)
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
