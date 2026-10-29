FROM golang:1.27.3-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api .

FROM alpine:3.24
COPY --from=build /api /api
EXPOSE 8181
ENTRYPOINT ["/api"]
