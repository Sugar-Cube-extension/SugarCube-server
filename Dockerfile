FROM golang:1.24.2-alpine as builder

WORKDIR /build
COPY . .

RUN go mod tidy
RUN go build -o sugarcubed .

FROM golang:1.24.2-alpine

WORKDIR /

COPY --from=builder /build/sugarcubed /usr/bin/

EXPOSE 80

CMD ["sugarcubed"]

