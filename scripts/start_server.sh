#!/bin/bash

main() {
    scripts/setup_dev_grpc_certs.sh
    echo -e "CA is available at $(pwd)/grpc-kms-certs/grpc_kms_ca_cert.pem"
    go run main.go /tmp/socket.sock grpc-kms-certs/grpc_kms_server_cert.pem grpc-kms-certs/grpc_kms_server_key.pem
}

main