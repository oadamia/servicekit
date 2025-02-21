package nanoid

import "fmt"

type mockGenerator struct {
	id     int
	prefix string
	length int
}

func (m *mockGenerator) Generate() (string, error) {
	id := fmt.Sprintf("%v", m.id)

	for len(id) < m.length {
		id = fmt.Sprintf("0%s", id)
	}

	id = fmt.Sprintf("%s%v", m.prefix, id)
	return id, nil
}

func MockGenerator(prefix string, length int) Generator {
	return &mockGenerator{
		id:     1,
		prefix: prefix,
		length: length,
	}
}
