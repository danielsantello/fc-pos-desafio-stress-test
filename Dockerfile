FROM golang:1.26-alpine AS build

WORKDIR /src
RUN apk add --no-cache ca-certificates

COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /out/stress-test .

FROM scratch

COPY --from=build /out/stress-test /stress-test
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 65532:65532
ENTRYPOINT ["/stress-test"]
