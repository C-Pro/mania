FROM golang:latest AS builder

RUN apt update && apt install -y ca-certificates
COPY . /src
WORKDIR /src

RUN CGO_ENABLED=0 go build -installsuffix cgo -o /mania .

FROM scratch

COPY --from=builder /mania /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

VOLUME /firebase
ENV GOOGLE_APPLICATION_CREDENTIALS /firebase/credentials.json

EXPOSE 8080:8080

CMD ["/mania"]
