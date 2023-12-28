FROM golang:1.20-alpine as build

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -o slack-cronjob

FROM golang:1.20-alpine
WORKDIR /app
COPY --from=build /build/slack-cronjob /app/
CMD [ "./slack-cronjob" ]
