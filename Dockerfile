FROM golang:1.19-bullseye
WORKDIR /opt
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app

FROM gcr.io/distroless/static-debian11
COPY --from=0 /opt/app /
CMD ["/app"]