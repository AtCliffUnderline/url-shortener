package app_test

import (
	"github.com/AtCliffUnderline/url-shortener/internal/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	storage := app.DefaultRouteStorage{}
	t.Run("add to storage and read successfully", func(t *testing.T) {
		id := storage.ShortRoute("some route")
		route, err := storage.GetRouteByID(id)
		assert.NoError(t, err)
		assert.Equal(t, "some route", route)
	})
	t.Run("read unexciting element", func(t *testing.T) {
		_, err := storage.GetRouteByID(123)
		assert.Error(t, err)
	})
}
