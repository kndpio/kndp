FROM golang:1.20-alpine as build

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -o slack-notify

FROM golang:1.20-alpine
WORKDIR /app
COPY --from=build /build/slack-notify /app/
CMD [ "./slack-notify" ]
