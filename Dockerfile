# Stage 1: Build the Go agent
FROM golang:1.16-alpine AS go-builder
WORKDIR /app
# Copy the Go source code into the container
COPY go-agent/ ./
# Build the Go application
RUN go mod download
RUN go build -o go-agent .

# Stage 2: Set up the PHP Laravel application
FROM php:7.4-fpm AS laravel-setup
# Install system dependencies for Laravel
RUN apt-get update && apt-get install -y \
    git \
    curl \
    zip \
    unzip \
    quota \
    && apt-get clean
# Install PHP extensions
RUN docker-php-ext-install pdo pdo_mysql
# Set working directory
WORKDIR /var/www
# Get Composer
COPY --from=composer:latest /usr/bin/composer /usr/local/bin/composer
# Copy existing application directory contents
COPY laravel-app/ /var/www
# Copy the Go agent binary from the first stage
COPY --from=go-builder /app/go-agent /usr/local/bin/go-agent
# Install all PHP dependencies
RUN composer install
# Expose port 9000 and start php-fpm server
EXPOSE 9000
CMD ["php-fpm"]
