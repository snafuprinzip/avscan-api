# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /avscan-api

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM clamav/clamav AS build-release-stage

WORKDIR /avscan-api/

COPY --from=build-stage /avscan-api /avscan-api/avscan-api
COPY conf /avscan-api/conf

RUN chown -R clamav:clamav /avscan-api

EXPOSE 8080

USER clamav:clamav
ENTRYPOINT ["/avscan-api/avscan-api"]
