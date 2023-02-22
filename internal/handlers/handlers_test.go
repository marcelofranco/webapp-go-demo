package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/marcelofranco/webapp-go-demo/internal/models"
)

var tests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"get-search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range tests {
		if e.method == "GET" {
			res, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected status code %d but got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}
		}
	}
}

var getReservationTests = []struct {
	name               string
	reservation        models.Reservation
	expectedStatusCode int
	expectedLocation   string
	expectedHTML       string
}{
	{
		name: "reservation-in-session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				RoomName: "General's Quarters",
			},
		},
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `action="/make-reservation"`,
	},
	{
		name:               "reservation-not-in-session",
		reservation:        models.Reservation{},
		expectedStatusCode: http.StatusTemporaryRedirect,
		expectedLocation:   "/",
		expectedHTML:       "",
	},
	{
		name: "non-existent-room",
		reservation: models.Reservation{
			RoomID: 100,
			Room: models.Room{
				RoomName: "General's Quarters",
			},
		},
		expectedStatusCode: http.StatusTemporaryRedirect,
		expectedLocation:   "/",
		expectedHTML:       "",
	},
}

func TestReservation(t *testing.T) {
	for _, e := range getReservationTests {
		req, _ := http.NewRequest("GET", "/make-reservation", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.Reservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var postReservationTests = []struct {
	name                 string
	reservation          models.Reservation
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name: "valid-data",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				RoomName: "General's Quarters",
			},
		},
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedHTML:         "",
		expectedLocation:     "/reservation-summary",
	},
	{
		name: "empty-reservation-session",
		reservation: models.Reservation{
			RoomID: 0,
		},
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedHTML:         "",
		expectedLocation:     "/",
	},
	{
		name: "empty-form-data",
		reservation: models.Reservation{
			RoomID: 1,
		},
		postedData:           nil,
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedHTML:         "",
		expectedLocation:     "/",
	},
	{
		name: "invalid-form-data",
		reservation: models.Reservation{
			RoomID: 1,
		},
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"notanemail"},
			"phone":      {"555-555-5555"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusOK,
		expectedHTML:         `action="/make-reservation"`,
		expectedLocation:     "",
	},
	{
		name: "error-inserting-reservation",
		reservation: models.Reservation{
			RoomID: 2,
		},
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
			"room_id":    {"2"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedHTML:         "",
		expectedLocation:     "/",
	},
	{
		name: "error-inserting-room-restriction",
		reservation: models.Reservation{
			RoomID: 3,
		},
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
			"room_id":    {"3"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedHTML:         "",
		expectedLocation:     "/",
	},
}

func TestRepository_PostReservation(t *testing.T) {
	for _, e := range postReservationTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/make-reservation", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.PostReservation)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedResponseCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var getReservationSummaryTests = []struct {
	name                 string
	reservation          models.Reservation
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name: "reservation-in-session",
		reservation: models.Reservation{
			RoomID:    1,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			Room: models.Room{
				RoomName: "General's Quarters",
			},
		},
		expectedResponseCode: http.StatusOK,
		expectedLocation:     "",
	},
	{
		name: "reservation-not-in-session",
		reservation: models.Reservation{
			RoomID:    0,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			Room: models.Room{
				RoomName: "General's Quarters",
			},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
}

func TestRepository_ReservationSummary(t *testing.T) {
	for _, e := range getReservationSummaryTests {
		req, _ := http.NewRequest("GET", "/reservation-summary", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.ReservationSummary)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedResponseCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var postPostAvailability = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
}{
	{
		name: "valid-data",
		postedData: url.Values{
			"start_date": {"2049-11-30"},
			"end_date":   {"2049-11-30"},
		},
		expectedResponseCode: http.StatusOK,
	},
	{
		name: "invalid-start-date",
		postedData: url.Values{
			"start_date": nil,
			"end_date":   {"2050-01-02"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
	},
	{
		name: "invalid-end-date",
		postedData: url.Values{
			"start_date": {"2050-01-02"},
			"end_date":   nil,
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
	},
	{
		name: "failed-date",
		postedData: url.Values{
			"start_date": {"2060-01-01"},
			"end_date":   {"2060-01-01"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
	},
	{
		name: "not-available-room",
		postedData: url.Values{
			"start_date": {"2051-01-02"},
			"end_date":   {"2051-01-03"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
	{
		name:                 "invalid-form",
		postedData:           nil,
		expectedResponseCode: http.StatusTemporaryRedirect,
	},
}

func TestRepository_PostAvailability(t *testing.T) {
	for _, e := range postPostAvailability {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/search-availability", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.PostAvailability)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedResponseCode)
		}
	}
}

var getChooseRoom = []struct {
	name               string
	reservation        models.Reservation
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name: "reservation-in-session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				RoomName: "General's Quarters",
			},
		},
		url:                "/choose-room/1",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/make-reservation",
	},
	{
		name:               "reservation-not-in-session",
		reservation:        models.Reservation{},
		url:                "/choose-room/1",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name:               "malformed-url",
		reservation:        models.Reservation{},
		url:                "/choose-room/fish",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
}

func TestChooseRoom(t *testing.T) {
	for _, e := range getChooseRoom {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		// set the RequestURI on the request so that we can grab the ID from the URL
		req.RequestURI = e.url

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.ChooseRoom)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var getBookRoom = []struct {
	name                 string
	url                  string
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name:                 "valid-data",
		url:                  "/book-room?id=1&s=2050-01-02&e=2050-01-02",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/make-reservation",
	},
	{
		name:                 "invalid-id",
		url:                  "/book-room?id=something&s=2050-01-02&e=2050-01-02",
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
	{
		name:                 "invalid-start-date",
		url:                  "/book-room?id=1&s=2050-01-32&e=2050-01-02",
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
	{
		name:                 "invalid-end-date",
		url:                  "/book-room?id=1&s=2050-01-02&e=2050-01-32",
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
}

func TestRepository_BookRoom(t *testing.T) {
	for _, e := range getBookRoom {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		// set the RequestURI on the request so that we can grab the ID from the URL
		req.RequestURI = e.url

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.BookRoom)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedResponseCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var postAvailabilityJSON = []struct {
	name            string
	postedData      url.Values
	expectedOK      bool
	expectedMessage string
}{
	{
		name: "valid-data",
		postedData: url.Values{
			"start_date": {"2049-11-30"},
			"end_date":   {"2049-11-30"},
			"room_id":    {"1"},
		},
		expectedOK: true,
	},
	{
		name:            "invalid-form",
		postedData:      nil,
		expectedOK:      false,
		expectedMessage: "Internal server error",
	},
	{
		name: "not-available",
		postedData: url.Values{
			"start_date": {"2049-11-30"},
			"end_date":   {"2049-11-30"},
			"room_id":    {"2"},
		},
		expectedOK: false,
	},
	{
		name: "error-database",
		postedData: url.Values{
			"start_date": {"2049-11-30"},
			"end_date":   {"2049-11-30"},
			"room_id":    {"3"},
		},
		expectedOK:      false,
		expectedMessage: "Error connecting to database",
	},
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	for _, e := range postAvailabilityJSON {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/search-availability-json", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AvailabilityJSON)

		handler.ServeHTTP(rr, req)

		var j jsonResponse
		err := json.Unmarshal([]byte(rr.Body.String()), &j)
		if err != nil {
			t.Error("failed to parse json!")
		}

		if j.OK != e.expectedOK {
			t.Errorf("%s: expected %v but got %v", e.name, e.expectedOK, j.OK)
		}
	}
}

var postSignUp = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name: "valid-form",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"test@here.com"},
			"password":   {"12345Q@e"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/",
		expectedHTML:         "",
	},
	{
		name:                 "empty-form",
		postedData:           nil,
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
		expectedHTML:         "",
	},
	{
		name: "invalid-email-form",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"notanemail"},
			"password":   {"12345Q@e"},
		},
		expectedResponseCode: http.StatusOK,
		expectedLocation:     "",
		expectedHTML:         `action="/sign-up"`,
	},
	{
		name: "already-exist-user",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"me@here.com"},
			"password":   {"12345Q@e"},
		},
		expectedResponseCode: http.StatusOK,
		expectedLocation:     "",
		expectedHTML:         `action="/sign-up"`,
	},
	{
		name: "database-error",
		postedData: url.Values{
			"first_name": {"Error"},
			"last_name":  {"Smith"},
			"email":      {"test@here.com"},
			"password":   {"12345Q@e"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
		expectedHTML:         "",
	},
}

func TestRepository_PostSignUp(t *testing.T) {
	for _, e := range postSignUp {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/sign-up", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/sign-up", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.PostSignUp)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedResponseCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var getSignUp = []struct {
	name               string
	expectedStatusCode int
	expectedLocation   string
	expectedHTML       string
}{
	{
		name:               "enter-signup",
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `action="/sign-up"`,
	},
}

func TestSignUp(t *testing.T) {
	for _, e := range getSignUp {
		req, _ := http.NewRequest("GET", "/sign-up", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.SignUp)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var postSignIn = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name: "valid-form",
		postedData: url.Values{
			"username_login": {"me@here.com"},
			"password_login": {"12345Q@e"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/",
	},
	{
		name:                 "empty-form",
		postedData:           nil,
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
	{
		name: "unauthorized-form",
		postedData: url.Values{
			"username_login": {"me@here.com"},
			"password_login": {"unauthorized"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
}

func TestRepository_PostSignIn(t *testing.T) {
	for _, e := range postSignIn {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/sign-up", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/sign-up", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.Signin)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedResponseCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var getLogout = []struct {
	name               string
	expectedStatusCode int
	expectedLocation   string
	expectedHTML       string
}{
	{
		name:               "enter-logout",
		expectedStatusCode: http.StatusSeeOther,
	},
}

func TestLogout(t *testing.T) {
	for _, e := range getLogout {
		req, _ := http.NewRequest("GET", "/logout", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.Logout)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var getBookedRooms = []struct {
	name               string
	userID             int
	expectedStatusCode int
	expectedLocation   string
	expectedHTML       string
}{
	{
		name:               "valid-user-and-reservations",
		userID:             1,
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "user-not-in-session",
		userID:             0,
		expectedStatusCode: http.StatusTemporaryRedirect,
		expectedLocation:   "/",
	},
	{
		name:               "invalid-user",
		userID:             2,
		expectedStatusCode: http.StatusTemporaryRedirect,
	},
	{
		name:               "valid-user-no-reservations",
		userID:             3,
		expectedStatusCode: http.StatusTemporaryRedirect,
		expectedLocation:   "/",
	},
}

func TestBookedRooms(t *testing.T) {
	for _, e := range getBookedRooms {
		req, _ := http.NewRequest("GET", "/booked-rooms", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		if e.userID > 0 {
			session.Put(ctx, "user_id", e.userID)
		}

		handler := http.HandlerFunc(Repo.BookedRooms)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
