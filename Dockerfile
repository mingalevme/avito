FROM golang:alpine AS builder
WORKDIR /avito
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o avito .

FROM alpine AS avito
COPY --from=builder /avito/avito /usr/local/bin/avito
ENTRYPOINT [ "/usr/local/bin/avito" ]