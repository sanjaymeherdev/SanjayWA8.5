#!/bin/bash
cd src
echo "Building whatsapp binary..."
CGO_ENABLED=1 go build -o whatsapp .
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi
echo "Starting whatsapp REST server on port 5000..."
./whatsapp rest --port=5000
