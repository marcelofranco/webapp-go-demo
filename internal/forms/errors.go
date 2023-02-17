package forms

type errors map[string][]string

// Add adds an error message for a giving field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message
func (e errors) Get(field string) string {
	es, exists := e[field]
	if !exists {
		return ""
	}

	return es[0]
}
