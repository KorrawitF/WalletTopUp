FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app .

FROM alpine:3.24

RUN apk add --no-cache ca-certificates tzdata \
 && addgroup -S app && adduser -S -G app app

COPY --from=build /out/app /app

EXPOSE 8000

USER app:app

ENTRYPOINT ["/app"]
