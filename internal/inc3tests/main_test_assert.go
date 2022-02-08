package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbsAssert(t *testing.T) {
	tests := []struct { // добавился слайс тестов
		name  string
		value float64
		want  float64
	}{
		{
			name:  "simple test #1",
			value: 3,
			want:  3,
		},
		{
			name:  "simple test #2",
			value: 0,
			want:  0,
		},
		{
			name:  "simple test #3",
			value: -3,
			want:  3,
		},
		{
			name:  "simple test #4",
			value: -2.000000001,
			want:  2.000000001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Abs(tt.value))
		})
	}
}

func TestFullnameAssert(t *testing.T) {
	tests := []struct { // добавился слайс тестов
		name  string
		value User
		want  string
	}{
		{
			name: "simple test #1",
			value: User{
				FirstName: "Asd",
				LastName:  "Asd",
			},
			want: "Asd Asd",
		},
		{
			name: "simple test #2",
			value: User{
				FirstName: "123",
				LastName:  "123",
			},
			want: "123 123",
		},
		{
			name: "simple test #3",
			value: User{
				FirstName: "",
				LastName:  "",
			},
			want: " ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.value.FullName())
		})
	}
}

func TestFamilyAssert(t *testing.T) {
	family := Family{}
	emptyPerson := Person{FirstName: "", LastName: "", Age: 0}
	t.Run("Father In family", func(t *testing.T) {
		err := family.AddNew(Father, Person{
			FirstName: "Al",
			LastName:  "Pachino",
			Age:       133,
		})
		assert.NotEqual(t, emptyPerson, family.Members[Father])
		assert.NotNil(t, err)
	})
	t.Run("Mother In family", func(t *testing.T) {
		err := family.AddNew(Mother, Person{
			FirstName: "Some",
			LastName:  "name",
			Age:       131,
		})
		assert.NotEqual(t, emptyPerson, family.Members[Mother])
		assert.NotNil(t, err)
	})
	t.Run("Second Father In family", func(t *testing.T) {
		err := family.AddNew(Father, Person{
			FirstName: "Some",
			LastName:  "name",
			Age:       131,
		})
		assert.Nil(t, err)
	})
}
