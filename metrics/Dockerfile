FROM golang:alpine AS build

WORKDIR /mmts

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o mmts-metrics

FROM alpine:3.16

WORKDIR /mmts
COPY --from=build /mmts/mmts-metrics ./mmts-metrics

EXPOSE 8080
ENTRYPOINT ["./mmts-metrics"]