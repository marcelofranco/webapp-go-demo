package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/testvalid", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/testvalid", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("did not validated required fields when it should")
	}

	postedValues := url.Values{}
	postedValues.Add("a", "a")
	postedValues.Add("b", "a")
	postedValues.Add("c", "a")

	r = httptest.NewRequest("POST", "/testvalid", nil)
	r.PostForm = postedValues
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("errors on required field when shoud have none")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/testvalid", nil)
	form := New(r.PostForm)
	if form.Has("a") {
		t.Error("expected to return false when form does not have field but returned true")
	}

	postedValues := url.Values{}
	postedValues.Add("a", "a")

	r = httptest.NewRequest("POST", "/testvalid", nil)
	r.PostForm = postedValues
	form = New(r.PostForm)
	if !form.Has("a") {
		t.Error("expected to return true when form does have field but returned false")
	}
}

func TestForm_MinLenght(t *testing.T) {
	r := httptest.NewRequest("POST", "/testvalid", nil)
	form := New(r.PostForm)
	if form.MinLenght("a", 2) {
		t.Error("expected to return false when form does not have field to validate lenght")
	}

	postedValues := url.Values{}
	postedValues.Add("a", "a")
	r = httptest.NewRequest("POST", "/testvalid", nil)
	r.PostForm = postedValues
	form = New(r.PostForm)
	if form.MinLenght("a", 2) {
		t.Error("expected to return false when field does not have minimal lenght")
	}

	postedValues = url.Values{}
	postedValues.Add("a", "ab")
	r = httptest.NewRequest("POST", "/testvalid", nil)
	r.PostForm = postedValues
	form = New(r.PostForm)
	if !form.MinLenght("a", 2) {
		t.Error("expected to return true when field have minimal lenght but returned false")
	}

	postedValues = url.Values{}
	postedValues.Add("a", "abc")
	r = httptest.NewRequest("POST", "/testvalid", nil)
	r.PostForm = postedValues
	form = New(r.PostForm)
	if !form.MinLenght("a", 2) {
		t.Error("expected to return true when field have more then minimal lenght but returned false")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/testvalid", nil)
	form := New(r.PostForm)
	if form.IsEmail("a") {
		t.Error("expected to return false when form does not have field to validate if is email")
	}

	postedValues := url.Values{}
	postedValues.Add("a", "a")
	r = httptest.NewRequest("POST", "/testvalid", nil)
	r.PostForm = postedValues
	form = New(r.PostForm)
	if form.IsEmail("a") {
		t.Error("expected to return false when field is not a valid email")
	}

	postedValues = url.Values{}
	postedValues.Add("a", "me@here.com")
	r = httptest.NewRequest("POST", "/testvalid", nil)
	r.PostForm = postedValues
	form = New(r.PostForm)
	if !form.IsEmail("a") {
		t.Error("expected to return true when field is a valid email but returned false")
	}
}
