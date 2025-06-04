package wire

import (
	"testing"
)

func TestBase(t *testing.T) {
	monster := NewMonster()
	player := NewPlayer("dj")
	mission := NewMission(player, monster)

	mission.Start() // dj defeats kitty, world peace!
}
