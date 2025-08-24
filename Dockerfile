ARG BUILDPLATFORM=linux/amd64

FROM --platform=$BUILDPLATFORM golang:1.24-alpine as get-deps
LABEL authors="svck"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

FROM get-deps as build
WORKDIR /app

COPY . .

FROM build AS amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/sentinel -tags prod -mod=readonly -ldflags "-s -w" ./main.go


FROM alpine:3.21.3

COPY --from=amd64 /out/sentinel /usr/local/bin/sentinel
CMD ["sentinel"]
