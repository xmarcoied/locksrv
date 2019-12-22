package utils

type Set map[interface{}]bool

func NewSet() Set {
	return make(map[interface{}]bool)
}

// Add ...
func (s Set) Add(item interface{}) Set {
	s[item] = true
	return s
}

// Remove ...
func (s Set) Remove(item interface{}) Set {
	delete(s, item)
	return s
}
