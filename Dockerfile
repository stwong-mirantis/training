ARG GOLANG_ALPINE
ARG ALPINE
FROM ${GOLANG_ALPINE} as builder

WORKDIR /build
COPY . .

RUN go build -o server server.go

FROM ${ALPINE}

COPY --from=builder /build/server /app/

EXPOSE 8080

ENTRYPOINT ["/app/server"]
