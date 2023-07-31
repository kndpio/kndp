FROM golang:1.20-alpine as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /operator

FROM alpine:latest
WORKDIR /
COPY --from=build /operator /
RUN apk --no-cache add curl
RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/local/bin/
ENTRYPOINT [ "/operator" ]
EXPOSE 8080
