FROM golang:1.25.6 AS build-stage

WORKDIR /go/src/app

# copy all the files because main.go requires static and templates
# therefore, can't just do `*.go ./`
COPY . ./

# build the binary and keep it project root. binary is called `app`
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app

# distroless image to minimise build size
FROM gcr.io/distroless/base-debian13 AS build-release-stage

WORKDIR /

COPY static /static
COPY templates /templates
COPY --from=build-stage /go/bin/app /

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app"]
