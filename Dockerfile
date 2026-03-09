FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libchromaprint-tools ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the pre-compiled Go binary from the host
COPY tagSonic .

COPY .env .

ENTRYPOINT ["./tagSonic"]
CMD ["/music"]
