#!/usr/bin/env bash

gcloud kms decrypt \
--keyring=keyring-name \
--key=key-name \
--location=asia-northeast1 \
--project=your-project-id \
--ciphertext-file=./bot/.env.yaml.enc \
--plaintext-file=./bot/.env.yaml
