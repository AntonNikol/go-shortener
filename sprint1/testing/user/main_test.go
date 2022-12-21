package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFullName(t *testing.T) {
	type fields struct {
		FirstName string
		LastName  string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Тест Иванова Ивана",
			fields: fields{
				FirstName: "Иван",
				LastName:  "Иванов",
			},
			want: "Иван Иванов",
		},
		{
			name: "long name",
			fields: fields{
				FirstName: "Pablo Diego KHoze Frantsisko de Paula KHuan Nepomukeno Krispin Krispiano de la Santisima Trinidad Ruiz",
				LastName:  "Picasso",
			},
			want: "Pablo Diego KHoze Frantsisko de Paula KHuan Nepomukeno Krispin Krispiano de la Santisima Trinidad Ruiz Picasso",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
			}

			//if got := u.FullName(); got != tt.want {
			//	t.Errorf("FullName() = %v, want %v", got, tt.want)
			//}
			v := u.FullName()
			assert.Equal(t, v, tt.want)
		})
	}
}
