package app_test

import (
	"github.com/AtCliffUnderline/url-shortener/internal/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Run("add to storage and read successfully", func(t *testing.T) {
		id := app.ShortRoute("some route")
		route, err := app.GetRouteById(id)
		assert.Nil(t, err)
		assert.Equal(t, "some route", route)
	})
	t.Run("read unexciting element", func(t *testing.T) {
		_, err := app.GetRouteById(123)
		assert.NotNil(t, err)
	})
}
