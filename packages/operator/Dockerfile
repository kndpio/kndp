FROM golang:1.20-alpine as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /operator
FROM scratch
WORKDIR /
COPY --from=build /operator /
ENTRYPOINT [ "/operator" ]
EXPOSE 8080