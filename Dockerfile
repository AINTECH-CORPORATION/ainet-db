FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o /out/ainet-db-explore ./cmd/ainet-db-explore && \
    go build -o /out/ainet-db-sync ./cmd/ainet-db-sync

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/ainet-db-explore /usr/local/bin/ainet-db-explore
COPY --from=builder /out/ainet-db-sync /usr/local/bin/ainet-db-sync
EXPOSE 8099
ENTRYPOINT ["ainet-db-explore"]
