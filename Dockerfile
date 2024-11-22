FROM golang:latest AS builder

WORKDIR /src
RUN apt-get update && apt-get install -y ca-certificates openssl
ARG cert_location=/usr/local/share/ca-certificates
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
# Update certificates
RUN update-ca-certificates


# Download dependencies
COPY go.mod go.sum /
RUN go mod download

# Add source code
COPY . .
RUN CGO_ENABLED=0 go build -o main .

# Multi-Stage production build
FROM alpine AS production
RUN apk --no-cache add ca-certificates

WORKDIR /app
# Retrieve the binary from the previous stage
COPY --from=builder /src/main .
# Copy static template files
COPY templates templates
# Copy frontend
COPY public public
# Expose port
EXPOSE 3000
# Set the binary as the entrypoint of the container
CMD ["./main", "serve"]