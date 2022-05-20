package world

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cityDetails struct {
	name        string
	neighbors   []string
	neighborMap map[string]direction
}

func TestAddCityAndRoads(t *testing.T) {
	tests := []struct {
		name           string
		cityDefs       []cityDefinition
		expectedCities []cityDetails
		err            string
		directed       bool
	}{
		{
			"nil cityDefinition",
			nil,
			nil,
			"invalid city definition: <nil>",
			true,
		},
		{
			"one city, no neighbors",
			[]cityDefinition{
				{
					name: "Hurlingham",
				},
			},
			[]cityDetails{{name: "Hurlingham"}},
			"",
			true,
		},
		{
			"one city, one neighbor, directed",
			[]cityDefinition{
				{name: "Hurlingham",
					neighbors:   []string{"Morón"},
					neighborMap: map[string]direction{"Morón": west},
				},
			},
			[]cityDetails{
				{
					name:        "Hurlingham",
					neighbors:   []string{"Morón"},
					neighborMap: map[string]direction{"Morón": west},
				},
				{
					name: "Morón",
				},
			},
			"",
			true,
		},
		{
			"two cities, two neighbors each, directed",
			[]cityDefinition{
				{
					name:      "Hurlingham",
					neighbors: []string{"Morón", "Bernal"},
					neighborMap: map[string]direction{
						"Morón":  west,
						"Bernal": east,
					},
				},
				{
					name:      "Bernal",
					neighbors: []string{"Quilmes", "Hurlingham"},
					neighborMap: map[string]direction{
						"Hurlingham": west,
						"Quilmes":    south,
					},
				},
			},
			[]cityDetails{
				{
					name:      "Hurlingham",
					neighbors: []string{"Morón", "Bernal"},
					neighborMap: map[string]direction{
						"Morón":  west,
						"Bernal": east,
					},
				},
				{
					name:      "Bernal",
					neighbors: []string{"Quilmes", "Hurlingham"},
					neighborMap: map[string]direction{
						"Hurlingham": west,
						"Quilmes":    south,
					},
				},
				{
					name: "Morón",
				},
				{
					name: "Quilmes",
				},
			},
			"",
			true,
		},
		{
			"one city, one neighbor, non-directed",
			[]cityDefinition{
				{name: "Hurlingham",
					neighbors:   []string{"Morón"},
					neighborMap: map[string]direction{"Morón": west},
				},
			},
			[]cityDetails{
				{
					name:        "Hurlingham",
					neighbors:   []string{"Morón"},
					neighborMap: map[string]direction{"Morón": west},
				},
				{
					name:        "Morón",
					neighbors:   []string{"Hurlingham"},
					neighborMap: map[string]direction{"Hurlingham": east},
				},
			},
			"",
			false,
		},
		{
			"two defined cities, two neighbors each, non-directed",
			[]cityDefinition{
				{
					name:      "Hurlingham",
					neighbors: []string{"Morón", "Gerli"},
					neighborMap: map[string]direction{
						"Morón": west,
						"Gerli": south,
					},
				},
				{
					name:      "Bernal",
					neighbors: []string{"Gerli", "Hurlingham"},
					neighborMap: map[string]direction{
						"Hurlingham": west,
						"Gerli":      north,
					},
				},
			},
			[]cityDetails{
				{
					name:      "Hurlingham",
					neighbors: []string{"Morón", "Bernal", "Gerli"},
					neighborMap: map[string]direction{
						"Morón":  west,
						"Bernal": east,
						"Gerli":  south,
					},
				},
				{
					name:      "Bernal",
					neighbors: []string{"Gerli", "Hurlingham"},
					neighborMap: map[string]direction{
						"Hurlingham": west,
						"Gerli":      north,
					},
				},
				{
					name:      "Morón",
					neighbors: []string{"Hurlingham"},
					neighborMap: map[string]direction{
						"Hurlingham": east,
					},
				},
				{
					name:      "Gerli",
					neighbors: []string{"Bernal", "Hurlingham"},
					neighborMap: map[string]direction{
						"Bernal":     south,
						"Hurlingham": north,
					},
				},
			},
			"",
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			w := &World{
				Cities:   make(map[string]*City),
				directed: test.directed,
			}

			for _, cityDef := range test.cityDefs {
				w.addCityAndRoads(&cityDef)
			}

			assert.Equal(tt, len(test.expectedCities), len(w.Cities))
			for _, expectedCity := range test.expectedCities {
				actualCity, ok := w.Cities[expectedCity.name]
				assert.True(tt, ok)
				assert.NotNil(tt, actualCity)
				assert.Equal(tt, expectedCity.name, actualCity.Name)

				// Check neighbors
				assert.Equal(tt, len(expectedCity.neighbors), len(actualCity.Neighbors))
				assert.Equal(tt, len(expectedCity.neighborMap), len(actualCity.neighborMap))

				for _, n := range actualCity.Neighbors {
					assert.Contains(tt, expectedCity.neighbors, n.Name)
				}

				for c, dir := range actualCity.neighborMap {
					assert.Equal(tt, expectedCity.neighborMap[c.Name], dir)
				}
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name                string
		world               World
		expectedDefinitions map[string]cityDefinition
	}{
		{
			"empty world",
			World{},
			nil,
		},
		{
			"one city, no neighbors",
			World{
				Cities: map[string]*City{
					"Laferrere": {
						Name: "Laferrere",
					},
				},
			},
			map[string]cityDefinition{
				"Laferrere": {name: "Laferrere"},
			},
		},
		{
			"one city, two neighbors",
			World{
				Cities: map[string]*City{
					"Laferrere": {
						Name: "Laferrere",
						neighborMap: map[*City]direction{
							{Name: "ValentínAlsina"}: east,
							{Name: "Mataderos"}:      north,
						},
					},
				},
			},
			map[string]cityDefinition{
				"Laferrere": {
					name:      "Laferrere",
					neighbors: []string{"ValentínAlsina", "Mataderos"},
					neighborMap: map[string]direction{
						"ValentínAlsina": east,
						"Mataderos":      north,
					},
				},
			},
		},
		{
			"two cities, one neighbor",
			World{
				Cities: map[string]*City{
					"Calzada": {
						Name: "Calzada",
						neighborMap: map[*City]direction{
							{Name: "Gerli"}:    north,
							{Name: "Claypole"}: south,
						},
					},
					"Gerli": {
						Name: "Gerli",
						neighborMap: map[*City]direction{
							{Name: "Calzada"}: south,
							{Name: "Sarandí"}: east,
						},
					},
				},
			},
			map[string]cityDefinition{
				"Calzada": {
					name:      "Calzada",
					neighbors: []string{"Claypole", "Gerli"},
					neighborMap: map[string]direction{
						"Claypole": south,
						"Gerli":    north,
					},
				},
				"Gerli": {
					name:      "Gerli",
					neighbors: []string{"Calzada", "Sarandí"},
					neighborMap: map[string]direction{
						"Calzada": south,
						"Sarandí": east,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			res := test.world.String()
			resSlice := strings.Split(res, "\n")

			if len(resSlice) < 1 {
				return
			}
			for _, line := range resSlice {
				if line == "" {
					continue
				}
				cityDef, err := parseLine(line)
				if ok := assert.NoError(tt, err); !ok {
					return
				}

				expected := test.expectedDefinitions[cityDef.name]
				assert.Equal(tt, expected.name, cityDef.name)
				assert.ElementsMatch(tt, expected.neighbors, cityDef.neighbors)

				assert.Equal(tt, len(expected.neighborMap), len(cityDef.neighborMap))
				for k, v := range expected.neighborMap {
					assert.Equal(tt, v, cityDef.neighborMap[k])
				}
			}
		})
	}
}

func TestWorld(t *testing.T) {
	tests := []struct {
		name                        string
		input                       string
		isError                     bool
		expectedCities              map[string][]string
		cityToDelete                string
		expectedCitiesAfterDeletion map[string][]string
	}{
		{
			"one city, no neighbors",
			"Gerli",
			false,
			map[string][]string{"Gerli": {}},
			"Gerli",
			map[string][]string{},
		},
		{
			"two cities, no neighbors",
			"Ezeiza\n\nDockSud\n",
			false,
			map[string][]string{"Ezeiza": {}, "DockSud": {}},
			"DockSud",
			map[string][]string{"Ezeiza": {}},
		},
		{
			"one city, one neighbor",
			"Gerli south=DockSud",
			false,
			map[string][]string{"Gerli": {"DockSud"}, "DockSud": {}},
			"DockSud",
			map[string][]string{"Gerli": {}},
		},
		{
			"two cities, neighbors",
			"Gerli south=Burzaco west=DockSud\nDockSud east=Gerli",
			false,
			map[string][]string{"Gerli": {"Burzaco", "DockSud"}, "DockSud": {"Gerli"}, "Burzaco": {}},
			"Burzaco",
			map[string][]string{"Gerli": {"DockSud"}, "DockSud": {"Gerli"}},
		},
		{
			"invalid input, to many directions",
			"Gerli west=Burzaco east=Gerli north=DockSud south=Quilmes southwest=Ezeiza",
			true,
			nil,
			"",
			nil,
		},
		{
			"invalid input, invalid direction format",
			"Gerli west=Burzaco south=Gerli=something",
			true,
			nil,
			"",
			nil,
		},
		{
			"invalid input, invalid direction",
			"Gerli west=Burzaco southeast=Gerli",
			true,
			nil,
			"",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s directed", test.name), func(tt *testing.T) {
			w, err := NewFromReader(strings.NewReader(test.input), true)
			if test.isError {
				// If an error is expected we have nothing else to check, return.
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)

			// Check amount of cities.
			if ok := assert.Equal(tt, len(test.expectedCities), len(w.Cities)); !ok {
				return
			}

			// If lengths are equal, then check that every city is in place.
			for expCity, expNeighbors := range test.expectedCities {
				city, ok := w.Cities[expCity]
				assert.True(tt, ok)
				assert.NotNil(tt, city)
				assert.Equal(tt, len(expNeighbors), len(city.Neighbors))

				// Check that the neighbors are the same we expect
				neighborMap := make(map[string]struct{}, len(city.Neighbors))
				for _, neighbor := range city.Neighbors {
					neighborMap[neighbor.Name] = struct{}{}
				}

				for _, expNeighbor := range expNeighbors {
					_, ok := neighborMap[expNeighbor]
					assert.True(tt, ok, fmt.Sprintf("neighbor %q not found", expNeighbor))
				}
			}

			// Delete a city
			city := w.Cities[test.cityToDelete]
			w.DeleteCityAndRoads(city)

			// Check amount of cities.
			if ok := assert.Equal(tt, len(test.expectedCitiesAfterDeletion), len(w.Cities)); !ok {
				return
			}

			// If lengths are equal, then check that every city is in place.
			for expCity, expNeighbors := range test.expectedCitiesAfterDeletion {
				city, ok := w.Cities[expCity]
				assert.True(tt, ok)
				assert.NotNil(tt, city)
				assert.Equal(tt, len(expNeighbors), len(city.Neighbors))

				// Check that the neighbors are the same we expect
				neighborMap := make(map[string]struct{}, len(city.Neighbors))
				for _, neighbor := range city.Neighbors {
					neighborMap[neighbor.Name] = struct{}{}
				}

				for _, expNeighbor := range expNeighbors {
					_, ok := neighborMap[expNeighbor]
					assert.True(tt, ok, fmt.Sprintf("neighbor %q not found", expNeighbor))
				}
			}
		})

		t.Run(fmt.Sprintf("%s non-directed", test.name), func(tt *testing.T) {
			w, err := NewFromReader(strings.NewReader(test.input), false)
			if test.isError {
				// If an error is expected we have nothing else to check, return.
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)

			// Check amount of cities.
			if ok := assert.Equal(tt, len(test.expectedCities), len(w.Cities)); !ok {
				return
			}

			// If lengths are equal, then check that every city is in place.
			for expCity, expNeighbors := range test.expectedCities {
				city, ok := w.Cities[expCity]
				assert.True(tt, ok)
				assert.NotNil(tt, city)

				// Check that the neighbors are in place
				neighborMap := make(map[string]struct{}, len(city.Neighbors))
				for _, neighbor := range city.Neighbors {
					neighborMap[neighbor.Name] = struct{}{}
				}

				for _, expNeighbor := range expNeighbors {
					_, ok := neighborMap[expNeighbor]
					assert.True(tt, ok, fmt.Sprintf("neighbor %q not found", expNeighbor))

					// Check that the connection is bi-directional.
					neighbor, ok := w.Cities[expNeighbor]
					if !assert.True(tt, ok, fmt.Sprintf("city %q not found", expNeighbor)) {
						return
					}

					var found bool
					for _, n := range neighbor.Neighbors {
						if n.Name == city.Name {
							found = true
						}
					}

					assert.True(tt, found)
				}
			}

			// Delete a city
			city := w.Cities[test.cityToDelete]
			w.DeleteCityAndRoads(city)

			// Check amount of cities.
			if ok := assert.Equal(tt, len(test.expectedCitiesAfterDeletion), len(w.Cities)); !ok {
				return
			}

			// If lengths are equal, then check that every city is in place.
			for expCity, expNeighbors := range test.expectedCitiesAfterDeletion {
				city, ok := w.Cities[expCity]
				assert.True(tt, ok)
				assert.NotNil(tt, city)
				assert.Equal(tt, len(expNeighbors), len(city.Neighbors))

				// Check that the neighbors are the same we expect
				neighborMap := make(map[string]struct{}, len(city.Neighbors))
				for _, neighbor := range city.Neighbors {
					neighborMap[neighbor.Name] = struct{}{}
				}

				for _, expNeighbor := range expNeighbors {
					_, ok := neighborMap[expNeighbor]
					assert.True(tt, ok, fmt.Sprintf("neighbor %q not found", expNeighbor))
				}
			}
		})
	}
}

// Benchmarks
func BenchmarkAddCityAndRoads_Directed(b *testing.B) {
	cityDefs := []cityDefinition{
		{
			name:      "Hurlingham",
			neighbors: []string{"Morón", "Bernal"},
			neighborMap: map[string]direction{
				"Morón":  west,
				"Bernal": east,
			},
		},
		{
			name:      "Bernal",
			neighbors: []string{"Quilmes", "Hurlingham"},
			neighborMap: map[string]direction{
				"Hurlingham": west,
				"Quilmes":    south,
			},
		},
		{
			name:      "Gerli",
			neighbors: []string{"Lanús", "Avellaneda"},
			neighborMap: map[string]direction{
				"Lanús":      south,
				"Avellaneda": north,
			},
		},
		{
			name:      "Escalada",
			neighbors: []string{"Lanús", "Banfield"},
			neighborMap: map[string]direction{
				"Banfield": south,
				"Lanús":    north,
			},
		},
	}

	b.Run("non-directed", func(b *testing.B) {
		w := &World{
			Cities: make(map[string]*City),
		}
		for i := 0; i < b.N; i++ {
			for _, cityDef := range cityDefs {
				w.addCityAndRoads(&cityDef)
			}
		}
	})

	b.Run("directed", func(b *testing.B) {
		w := &World{
			Cities:   make(map[string]*City),
			directed: true,
		}

		for i := 0; i < b.N; i++ {
			for _, cityDef := range cityDefs {
				w.addCityAndRoads(&cityDef)
			}
		}
	})

}
