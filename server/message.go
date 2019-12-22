package server

import "strings"

type serverMessage string

func (s serverMessage) String() string {
	return strings.TrimSpace(string(s))
}

func (s serverMessage) Valid() bool {
	msg := s.String()
	splited := strings.Split(msg, " ")
	if len(splited) != 2 {
		return false
	}

	if splited[0] != "lock" && splited[0] != "unlock" {
		return false
	}

	return true
}

func (s serverMessage) Action() string {
	msg := s.String()
	splited := strings.Split(msg, " ")
	return splited[0]
}

func (s serverMessage) Resource() string {
	msg := s.String()
	splited := strings.Split(msg, " ")
	return splited[1]
}
