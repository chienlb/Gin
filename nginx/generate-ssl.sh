#!/bin/bash

# Generate self-signed SSL certificates for development
# In production, use Let's Encrypt or your CA

mkdir -p ssl

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout ssl/key.pem \
    -out ssl/cert.pem \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

echo "SSL certificates generated in nginx/ssl/"
echo "Note: These are self-signed certificates for development only"
echo "For production, use Let's Encrypt: https://letsencrypt.org/"
