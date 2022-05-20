package alien

import (
	"testing"

	"github.com/santihernandezc/alien-invasion/world"
	"github.com/stretchr/testify/assert"
)

var (
	testCity1 = world.City{
		Name: "Gerli",
	}
	testCity2 = world.City{
		Name:      "Laferrere",
		Neighbors: []*world.City{&testCity1},
	}
	testCity3 = world.City{
		Name:      "SanJusto",
		Neighbors: []*world.City{&testCity1, &testCity2},
	}
	testWorld = world.World{
		Cities: map[string]*world.City{
			testCity1.Name: &testCity1,
			testCity2.Name: &testCity2,
			testCity3.Name: &testCity3,
		},
	}
)

func TestMove(t *testing.T) {
	testCityNoRoads := world.City{
		Name: "Berazategui",
	}
	testCityOneRoad := world.City{
		Name:      "Calzada",
		Neighbors: []*world.City{&testCity1},
	}
	testCityTwoRoads := world.City{
		Name:      "Bernal",
		Neighbors: []*world.City{&testCity1, &testCity2},
	}
	tests := []struct {
		name           string
		alien          *Alien
		trapped        bool
		possibleCities []string
	}{
		{
			"trapped alien",
			&Alien{
				Position: &testCityNoRoads,
			},
			true,
			[]string{},
		},
		{
			"one option",
			&Alien{
				Position: &testCityOneRoad,
			},
			false,
			[]string{testCity1.Name},
		},
		{
			"two options",
			&Alien{
				Position: &testCityTwoRoads,
			},
			false,
			[]string{testCity1.Name, testCity2.Name},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			ok := test.alien.move()
			if test.trapped {
				assert.False(tt, ok, "Alien should be trapped")
				return
			}

			assert.Contains(tt, test.possibleCities, test.alien.Position.Name)
		})
	}

}
