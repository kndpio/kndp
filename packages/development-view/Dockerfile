FROM golang:1.20-alpine as build
WORKDIR /app
COPY packages/development-view/src .
RUN CGO_ENABLED=0 go build -o /development-view

FROM scratch
WORKDIR /
COPY --from=build /development-view /
COPY --from=alpine/k8s:1.24.16 /usr/bin/helm /usr/local/bin/helm
ENTRYPOINT [ "/development-view" ] 