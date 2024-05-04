FROM golang:1.22.2-bookworm as builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o quickstub.bin cmd/quickstub/main.go

FROM alpine:3.19.1 as runtime
RUN apk add --no-cache gcompat

COPY --from=builder /build/quickstub.bin /bin/quickstub
COPY --from=builder /build/cmd/quickstub/sample_config.yaml /etc/quickstub.yaml

ENTRYPOINT ["/bin/quickstub"]
CMD ["-conf", "/etc/quickstub.yaml"] 

