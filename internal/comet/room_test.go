package comet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoom_MasterId(t *testing.T) {
	room := NewRoom("1000")
	assert.Equal(t, "", room.MasterId())

	room.Put(&Channel{Key: "1"})
	room.Put(&Channel{Key: "2"})
	assert.Equal(t, "1", room.MasterId())
}

func TestRoom_Users(t *testing.T) {
	room := NewRoom("1000")
	assert.Equal(t,[]string{}, room.Users())

	room.Put(&Channel{Key: "1"})
	room.Put(&Channel{Key: "2"})
	room.Put(&Channel{Key: "3"})
	assert.Equal(t, []string{"3", "2", "1"}, room.Users())
}
