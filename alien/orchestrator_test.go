package alien

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/santihernandezc/alien-invasion/world"
	"github.com/stretchr/testify/assert"
)

var (
	nopLogger = log.New(ioutil.Discard, "", 0)
)

func TestNewOrchestrator(t *testing.T) {
	tests := []struct {
		name   string
		w      *world.World
		logger *log.Logger
		n      int
		err    string
	}{
		{
			"nil world",
			nil,
			nopLogger,
			0,
			"invalid World value: <nil>",
		},
		{
			"nil logger",
			&testWorld,
			nil,
			0,
			"invalid value for logger: <nil>",
		},
		{
			"empty world",
			&world.World{},
			nopLogger,
			0,
			"invalid World value: 0 cities",
		},
		{
			"5 aliens",
			&testWorld,
			nopLogger,
			5,
			"",
		},
		{
			"100 aliens",
			&testWorld,
			nopLogger,
			100,
			"",
		},
		{
			"10000 aliens",
			&testWorld,
			nopLogger,
			10000,
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			ao, err := NewOrchestrator(test.n, 0, test.w, test.logger)
			if test.err != "" {
				assert.EqualError(tt, err, test.err)
				return
			}

			if ok := assert.NotNil(tt, ao); !ok {
				return
			}

			assert.Equal(tt, len(ao.Aliens), test.n)

			ids := make(map[int]struct{}, len(ao.Aliens))
			for _, alien := range ao.Aliens {
				// IDs should be unique
				_, ok := ids[alien.ID]
				if !assert.False(tt, ok, fmt.Sprintf("ID %d is duplicated", alien.ID)) {
					return
				}
				ids[alien.ID] = struct{}{}

				// Check they have an assigned position
				_, ok = testWorld.Cities[alien.Position.Name]
				if !assert.True(tt, ok, fmt.Sprintf("Position %s for Alien %d not found in test World", alien.Position.Name, alien.ID)) {
					return
				}

			}
		})
	}
}

func TestDeleteAlien(t *testing.T) {
	tests := []struct {
		name           string
		n              int
		toDelete       int
		alreadyDeleted int
	}{
		{
			"5 aliens, zero to delete",
			5,
			0,
			0,
		},
		{
			"5 aliens, one to delete",
			5,
			1,
			0,
		},
		{
			"5 aliens, two to delete",
			5,
			2,
			0,
		},
		{
			"10 aliens, two to delete, two already deleted",
			10,
			2,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			ao, err := NewOrchestrator(test.n, 0, &testWorld, nopLogger)
			assert.NoError(tt, err)

			toDelete := ao.Aliens[:test.toDelete]
			ao.deleteAliens(toDelete)

			alreadyDeleted := ao.Aliens[test.toDelete:(test.toDelete + test.alreadyDeleted)]
			for _, a := range alreadyDeleted {
				a.isDeleted = true
			}

			expectedLenght := test.n - test.toDelete
			assert.Equal(tt, expectedLenght, len(ao.Aliens))

			assert.NotContains(tt, ao.Aliens, toDelete)
		})
	}
}

func TestUnleashAliens(t *testing.T) {
	t.Run("when two aliens encounter, they kill each other and destroy the city", func(tt *testing.T) {
		worldDef := "City1 south=City2\nCity2 north=City1"
		w, err := world.NewFromReader(strings.NewReader(worldDef), false)
		if !assert.NoError(tt, err) {
			return
		}

		ao, err := NewOrchestrator(2, 0, w, nopLogger)
		assert.NoError(tt, err)

		ao.UnleashAliens(1)

		assert.Equal(tt, 0, len(ao.Aliens))
		assert.Equal(tt, 1, len(w.Cities))

		var remainingCity *world.City
		if city, ok := w.Cities["City1"]; ok {
			remainingCity = city
		}
		if city, ok := w.Cities["City2"]; ok {
			remainingCity = city
		}

		assert.Equal(tt, 0, len(remainingCity.Neighbors))
	})

	t.Run("when a city is destroyed, aliens can no longer travel to or through it", func(tt *testing.T) {
		worldDef := "City1 south=City2\nCity2 north=City1"
		w, err := world.NewFromReader(strings.NewReader(worldDef), false)
		if !assert.NoError(tt, err) {
			return
		}

		ao, err := NewOrchestrator(5, 0, w, nopLogger)
		assert.NoError(tt, err)

		ao.UnleashAliens(2)

		// Two aliens are killed, the other ones are trapped
		assert.Equal(tt, 0, len(ao.Aliens))
		assert.Equal(tt, 1, len(w.Cities))
	})
}
