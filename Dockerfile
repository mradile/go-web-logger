#backend build
FROM golang:1.16-buster as gobuilder

RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates git

#run as unpriviledged user
RUN addgroup --gid 1100 app && adduser --disabled-password --uid 1100 --gid 1100 --gecos '' app

WORKDIR /build

#copy go.mod and go.sum and download all modules to cache this docker layer
ADD go.mod .
ADD go.sum .
RUN go mod download

#copy the rest and build the binary
COPY . /build/
RUN CGO_ENABLED=0 go build -v ./...

#production image
FROM scratch
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuilder /etc/passwd /etc/passwd
COPY --chown=1100:1100 --from=gobuilder /build/go-web-logger /app

USER 1100:1100
ENTRYPOINT ["./app"]
