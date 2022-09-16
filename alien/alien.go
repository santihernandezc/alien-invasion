package alien

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/santihernandezc/alien-invasion/world"
)

type Alien struct {
	ID           int
	City         *world.City
	Position     rl.Vector2
	NextPosition rl.Vector2
	Texture      rl.Texture2D
	isDeleted    bool
}

func (a *Alien) move() (ok bool) {
	// Check whether the alien is trapped
	neighbors := a.City.Neighbors
	if len(neighbors) < 1 {
		return false
	}

	// Move to a random city
	a.City = neighbors[rand.Intn(len(neighbors))]
	a.NextPosition = rl.NewVector2(float32(a.City.Position.X), float32(a.City.Position.Y))

	return true
}

func (a *Alien) Draw() {
	rl.DrawTexture(a.Texture, int32(a.Position.X)-a.Texture.Width/2, int32(a.Position.Y)-a.Texture.Height/2, rl.White)
}
