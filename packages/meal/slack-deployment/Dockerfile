FROM golang:1.20-alpine as build

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -o slack-deployment

FROM golang:1.20-alpine
WORKDIR /app
COPY --from=build /build/slack-deployment /app/
COPY --from=alpine/k8s:1.24.16 /usr/bin/kubectl /usr/local/bin/kubectl
EXPOSE 3000
CMD [ "./slack-deployment" ]
