FROM golang:1.18 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine

RUN apk add asciiquarium

COPY --from=build /go/bin/app /
CMD ["/app"]