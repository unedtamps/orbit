FROM golang:1.25-alpine AS build

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /api .

FROM alpine:3.17
COPY templates /templates
COPY static /static
COPY --from=build /api /api
EXPOSE 9999

ENTRYPOINT ["/api"]

