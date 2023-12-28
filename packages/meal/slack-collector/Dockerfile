FROM golang:1.20-alpine as build

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -o slack-deployment

FROM golang:1.20-alpine
WORKDIR /app
COPY --from=build /build/slack-deployment /app/
EXPOSE 3000
CMD [ "./slack-deployment" ]
