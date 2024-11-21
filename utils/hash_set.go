package utils

type HashSet struct {
	elements map[string]bool
}

func NewHashSet() *HashSet {
	return &HashSet{elements: make(map[string]bool)}
}

func (s *HashSet) Add(element string) {
	s.elements[element] = true
}

func (s *HashSet) AddList(elements []string) {
	for _, element := range elements {
		s.elements[element] = true
	}
}

func (s *HashSet) Remove(element string) {
	delete(s.elements, element)
}

func (s *HashSet) RemoveList(elements []string) {
	for _, element := range elements {
		delete(s.elements, element)
	}
}

func (s *HashSet) Contains(element string) bool {
	return s.elements[element]
}

func (s *HashSet) Size() int {
	return len(s.elements)
}

func (s *HashSet) Elements() []string {
	keys := make([]string, 0, len(s.elements))
	for key := range s.elements {
		keys = append(keys, key)
	}
	return keys
}