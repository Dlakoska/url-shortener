package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 30",
			size: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewRandomString(tt.size)

			// Проверяем длину
			assert.Len(t, result, tt.size)

			// Проверяем, что все символы — из допустимого набора
			validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
			for _, r := range result {
				assert.Contains(t, validChars, string(r), "символ %q не входит в разрешённый набор", r)
			}
		})
	}
}
