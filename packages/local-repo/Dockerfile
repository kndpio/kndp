FROM golang:1.20-alpine as build
WORKDIR /app
COPY packages/local-repo/src .
RUN CGO_ENABLED=0 go build -o /micro

FROM scratch
WORKDIR /
COPY --from=build /micro /
COPY --from=alpine/k8s:1.24.16 /usr/bin/helm /usr/local/bin/helm
ENTRYPOINT [ "/micro" ] 