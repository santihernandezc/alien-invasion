package world

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected *cityDefinition
		err      string
	}{
		{
			"empty string",
			"",
			nil,
			"invalid input: empty string",
		},
		{
			"white space",
			"		 ",
			nil,
			"invalid input: empty string",
		},
		{
			"invalid input, to many directions",
			"Gerli west=Burzaco east=Gerli north=DockSud south=Quilmes southwest=Ezeiza",
			nil,
			"invalid number of segments: 6",
		},
		{
			"invalid input, invalid road format",
			"Gerli west=Burzaco south=Gerli=something",
			nil,
			`invalid road definition: "south=Gerli=something"`,
		},
		{
			"invalid input, invalid direction",
			"Gerli west=Burzaco southeast=Gerli",
			nil,
			`error converting to direction: cannot convert string "southeast" to direction type`,
		},
		{
			"no neighbors",
			"Gerli",
			&cityDefinition{name: "Gerli"},
			"",
		},
		{
			"one neighbor",
			"Gerli south=DockSud",
			&cityDefinition{
				name:        "Gerli",
				neighbors:   []string{"DockSud"},
				neighborMap: map[string]direction{"DockSud": south},
			},
			"",
		},
		{
			"two neighbors",
			"Gerli2 south=Burzãco west=DockSud",
			&cityDefinition{
				name:        "Gerli2",
				neighbors:   []string{"Burzãco", "DockSud"},
				neighborMap: map[string]direction{"Burzãco": south, "DockSud": west},
			},
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			cityDef, err := parseLine(test.line)
			if test.err != "" {
				assert.EqualError(tt, err, test.err)
				return
			}

			assert.NoError(tt, err)

			if ok := assert.NotNil(tt, cityDef); !ok {
				return
			}

			assert.Equal(tt, cityDef.name, test.expected.name)
			assert.ElementsMatch(tt, test.expected.neighbors, cityDef.neighbors)

			assert.Equal(tt, len(test.expected.neighborMap), len(cityDef.neighborMap))
			for k, v := range test.expected.neighborMap {
				assert.Equal(tt, v, cityDef.neighborMap[k])
			}
		})
	}
}
