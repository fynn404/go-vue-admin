# Build stage for backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY . .
RUN go mod download
RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -o ../bin/server

# Build stage for frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY frontend/ .
RUN npm install
RUN npm run build

# Final stage
FROM alpine:latest
WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the built binary from backend builder
COPY --from=backend-builder /app/bin/server .
COPY --from=backend-builder /app/configs ./configs

# Copy the built frontend from frontend builder
COPY --from=frontend-builder /app/dist ./frontend/dist

# Expose necessary ports
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release

# Command to run the application
CMD ["./server"] 