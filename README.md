# Doomsday Algorithm

Guess (calculate) the day of the week using the Doomsday Algorithm.

## How to Run

In the project root, run the Go server

```
go run main.go
```

Now, the application is available at `http://localhost:8080`

## API Details

This Go server supports both `text/html` and `application/json` responses.

For `GET /date`, the **response** struct is `DateDisplay`:

```go
type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type DateDisplay struct {
	Date    Date   `json:"date"`
	Month_S string `json:"month_s"`
}
```

For `POST /date`, the **request** struct is `PostDateRequest`:

```go
type PostDateRequest struct {
	Day   int64  `json:"day"`
	Month int64  `json:"month"`
	Year  int64  `json:"year"`
	Guess string `json:"guess"`
}
```

For `POST /date`, the **response** struct is `DateAnswer`:

```go
type DateAnswer struct {
	Correct bool   `json:"correct"`
	Answer  string `json:"answer"`
}
```

## Tech Stack

- Go (1.25.4)
- HTMX

## Motivation

This is a toy project for me to learn Go and HTMX.
