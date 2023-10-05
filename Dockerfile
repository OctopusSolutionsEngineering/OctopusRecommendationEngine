# syntax=docker/dockerfile:1

FROM golang:1.20 as build

ARG Version=development

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . /app

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.Version=${Version}'" -o /octolint cmd/octolint.go

# Create the execution image
FROM alpine:latest

COPY --from=build /octolint /octolint

# Run
ENTRYPOINT ["/octolint"]
