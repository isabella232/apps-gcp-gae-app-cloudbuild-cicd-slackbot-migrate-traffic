#!/usr/bin/env bash

gcloud beta functions deploy challenge \
--project=your-project-id \
--region=asia-northeast1 \
--entry-point=Challenge \
--runtime=go111 \
--env-vars-file=.env.yaml \
--trigger-http
