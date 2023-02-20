package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var th myHandler
	h := NoSurf(&th)
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not http.Handler but is %t", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var th myHandler
	h := SessionLoad(&th)
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not http.Handler but is %t", v)
	}
}
