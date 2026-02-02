package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

func (d Date) String() string {
	return fmt.Sprintf("%d-%d-%d", d.Day, d.Month, d.Year)
}

var monthDays map[int][]int
var templates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/date.html",
))

func main() {
	monthDays = map[int][]int{
		1:  {31},
		2:  {28, 29},
		3:  {31},
		4:  {30},
		5:  {31},
		6:  {30},
		7:  {31},
		8:  {31},
		9:  {30},
		10: {31},
		11: {30},
		12: {31},
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", rootHandler)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /date", getDateHandler)
	mux.HandleFunc("POST /date", postDateHandler)

	protocol := "http"
	serverAddr := "localhost:8080"

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: loggingMiddleware(mux),
	}
	log.Printf("Server starting on %s://%s", protocol, serverAddr)
	server.ListenAndServe()
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		log.Printf("%d %s %s; %v", wrapped.statusCode, r.Method, r.URL.RequestURI(), time.Since(start))
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func getDateHandler(w http.ResponseWriter, r *http.Request) {
	date := newDate(2000, 2026, monthDays)
	if r.Header.Get("HX-Request") != "true" {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
		return
	}
	templates.ExecuteTemplate(w, "date", date)
}

type DateAnswer struct {
	Correct bool   `json:"correct"`
	Answer  string `json:"answer"`
}

type PostDateRequest struct {
	Day   int64  `json:"day"`
	Month int64  `json:"month"`
	Year  int64  `json:"year"`
	Guess string `json:"guess"`
}

func postDateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var day int64
	var month int64
	var year int64
	var guess string

	if r.Header.Get("Content-Type") == "application/json" {
		defer r.Body.Close()
		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}

		var postDateRequst PostDateRequest
		err := json.NewDecoder(r.Body).Decode(&postDateRequst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		day = postDateRequst.Day
		month = postDateRequst.Month
		year = postDateRequst.Year
		guess = postDateRequst.Guess
	} else {
		var err error
		day, err = strconv.ParseInt(r.FormValue("day"), 10, 0)
		if err != nil {
			http.Error(w, "Error parsing 'day'", http.StatusBadRequest)
			return
		}
		month, err = strconv.ParseInt(r.FormValue("month"), 10, 0)
		if err != nil {
			http.Error(w, "Error parsing 'month'", http.StatusBadRequest)
			return
		}
		year, err = strconv.ParseInt(r.FormValue("year"), 10, 0)
		if err != nil {
			http.Error(w, "Error parsing 'year'", http.StatusBadRequest)
			return
		}
		guess = r.FormValue("guess")
	}

	if !(isValidDate(int(day), int(month), int(year))) {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	date := Date{
		Day:   int(day),
		Month: int(month),
		Year:  int(year),
	}

	dateAnswer := checkWeekday(date, guess)

	if r.Header.Get("HX-Request") != "true" {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(dateAnswer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if dateAnswer.Correct {
		fmt.Fprint(w, "<p class='text-green-600'>Correct!</p>")
	} else {
		fmt.Fprint(w, "<p class='text-red-600'>Wrong! "+date.String()+" was a "+dateAnswer.Answer+"</p>")
	}
}

func checkWeekday(date Date, guess string) DateAnswer {
	t := time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.UTC)
	weekday := t.Weekday().String()
	return DateAnswer{weekday == guess, weekday}
}

func newDate(startYear int, endYear int, monthDays map[int][]int) Date {
	// choose year
	year := rand.IntN(endYear+1-startYear) + startYear
	leapYear := isLeapYear(year)
	// choose month
	month := rand.IntN(12) + 1
	// choose day
	var day int
	if leapYear && month == 2 {
		day = rand.IntN(monthDays[month][1]) + 1
	} else {
		day = rand.IntN(monthDays[month][0]) + 1
	}

	return Date{day, month, year}
}

func isLeapYear(year int) bool {
	if year%100 == 0 && year%400 != 0 {
		return false
	}
	return year%4 == 0
}

func isValidDate(day int, month int, year int) bool {
	leapYear := isLeapYear(year)
	// check if month is valid
	if !(1 <= month && month <= 12) {
		return false
	}
	// check if day is valid
	var maxDays int
	if leapYear && month == 2 {
		maxDays = monthDays[month][1] + 1
	} else {
		maxDays = monthDays[month][0] + 1
	}
	if !(1 <= day && day <= maxDays) {
		return false
	}
	return true
}
