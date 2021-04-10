FROM golang:1.16-alpine AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o feature-service

FROM scratch
WORKDIR /app
COPY --from=build /build/feature-service .
ENTRYPOINT ["/app/feature-service"]
