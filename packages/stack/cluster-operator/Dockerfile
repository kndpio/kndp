FROM golang:1.20-alpine as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /cluster-operator

FROM alpine/k8s:1.24.16
WORKDIR /
COPY --from=build /cluster-operator /
ENTRYPOINT [ "/cluster-operator" ]