#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Directory containing the .proto files
PROTO_DIR="./proto"

# Directory to output the generated Go files
OUTPUT_DIR="./proto"

# Create the output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# Generate Go code
protoc -I=$PROTO_DIR --go_out=$OUTPUT_DIR --go_opt=paths=source_relative \
    --go-grpc_out=$OUTPUT_DIR --go-grpc_opt=paths=source_relative \
    $PROTO_DIR/*.proto

echo "Proto files generated successfully."