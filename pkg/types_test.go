package pkg

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestInstantToDuration(t *testing.T) {
	instant := Instant("5m")
	d, err := instant.ToDuration()
	assert.Nil(t, err)
	assert.Equal(t, 5 * time.Minute, d)

	instant = Instant("2h")
	d, err = instant.ToDuration()
	assert.Nil(t, err)
	assert.Equal(t, 2 * time.Hour, d)

	instant = Instant("ab")
	_, err = instant.ToDuration()
	assert.NotNil(t, err)
}

