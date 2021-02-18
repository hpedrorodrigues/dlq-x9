FROM gcr.io/distroless/base-debian10
COPY dlq-x9 /
ENTRYPOINT ["/dlq-x9"]
