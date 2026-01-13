package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Run("should return a new server", func(t *testing.T) {
		server := NewServer()
		assert.NotNil(t, server)
		assert.NotNil(t, server.Handler)
	})
}
