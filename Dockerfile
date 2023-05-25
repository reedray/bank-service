FROM golang:1.19-alpine as base

WORKDIR /app

ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./converter ./cmd/converter/main.go
RUN go build -o ./gateway ./cmd/gateway/main.go
RUN go build -o ./transact ./cmd/transact/main.go

FROM base as converter
COPY --from=base /app/converter /bin/converter
CMD ["/bin/converter"]

FROM base as gateway
EXPOSE 8080
COPY --from=base /app/gateway /bin/gateway
CMD ["/bin/gateway"]

FROM base as transact
COPY --from=base /app/transact /bin/transact
CMD ["/bin/transact"]




