FROM golang:1.13-alpine AS build

WORKDIR /src/

# Copy go mod and sum files
COPY . /src/

# Build the Go app
RUN CGO_ENABLED=0 go build -o /bin/free-create

FROM scratch
COPY --from=build /bin/free-create /bin/free-create
ENTRYPOINT ["/bin/free-create"]