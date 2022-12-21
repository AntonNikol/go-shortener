package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbs(t *testing.T) {
	//v := rand.Float64()
	//absResult := math.Abs(v)
	//if res := Abs(v); res != absResult {
	//	t.Errorf("result expected to be %f, got %f", absResult, res)
	//}

	tests := []struct {
		name  string
		value float64
		want  float64
	}{
		{
			name:  "Тест отрицательного значения",
			value: -3,
			want:  3,
		}, {
			name:  "Тест отрицательного значения float",
			value: -0.000000005,
			want:  0.000000005,
		},
	}

	t.Logf("Начало запуска теста")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if result := Abs(tt.value); result != tt.want {
			//	t.Errorf("Abs() = %f, want = %f", result, tt.want)
			//}
			v := Abs(tt.value)
			assert.Equal(t, v, tt.want)
		})
	}
}
