package server

import "testing"

func Test_serverMessage_Valid(t *testing.T) {
	tests := []struct {
		name string
		s    serverMessage
		want bool
	}{
		{
			name: "Valid message",
			s:    "lock disk1",
			want: true,
		},
		{
			name: "Valid message",
			s:    "unlock disk1",
			want: true,
		},
		{
			name: "Valid message",
			s:    "check disk1",
			want: false,
		},
		{
			name: "Invalid message",
			s:    "marco",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Valid(); got != tt.want {
				t.Errorf("serverMessage.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
