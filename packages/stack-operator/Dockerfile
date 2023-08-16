FROM golang:1.20-alpine as build
WORKDIR /app
ENV STACKS_DIRECTORY=""
COPY . .
RUN CGO_ENABLED=0 go build -o /stack-operator

FROM scratch
WORKDIR /
COPY --from=build /operator /
COPY --from=alpine/k8s:1.24.16 /usr/bin/kubectl /usr/local/bin/kubectl
ENTRYPOINT [ "/stack-operator" ]
EXPOSE 8080
 