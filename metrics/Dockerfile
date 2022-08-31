FROM golang:alpine AS build

WORKDIR /mmts

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o mtm-metrics

FROM alpine:3.16

WORKDIR /home/mmts

RUN addgroup --gid 1000 mmts
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/home/mmts" \
    --ingroup mmts \
    --no-create-home \
    --uid 1000 \
    mmts

USER mmts

COPY --from=build /mmts/mtm-metrics ./mtm-metrics

EXPOSE 8080
ENTRYPOINT ["./mtm-metrics"]