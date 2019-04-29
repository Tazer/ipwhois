#first stage - builder
FROM golang:1.12.4-stretch as builder
COPY . /ipwhois
WORKDIR /ipwhois
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o ipwhois
#second stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /ipwhois .
EXPOSE 8080
CMD ["./ipwhois"]