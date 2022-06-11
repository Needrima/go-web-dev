package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	body := `
	{
       "username":"user1",
	   "password":"password1"
	}
	`
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("wanted: %v, got: %v", http.StatusOK, rec.Code)
		t.Fail()
	}
}

func TestHome(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/home", nil)
	rec := httptest.NewRecorder()

	Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("wanted: %v, got: %v", http.StatusOK, rec.Code)
		t.Fail()
	}
}
