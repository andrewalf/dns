FROM golang:1.16 AS build
WORKDIR /dns_app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cmd/dns ./cmd

FROM alpine:3.14 AS exec
COPY --from=build /dns_app/cmd/dns .
CMD ["/dns"]