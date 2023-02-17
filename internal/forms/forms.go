package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initialize a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks if form fields exists in post request
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field exists in post request
func (f *Form) Has(field string, r *http.Request) bool {
	exists := r.Form.Get(field)
	return exists != ""
}

// MinLenght checks for string minimal lenght
func (f *Form) MinLenght(field string, lenght int, r *http.Request) bool {
	valid := r.Form.Get(field)
	if len(valid) < lenght {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", lenght))
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

// Valid returns true if form has no error, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
