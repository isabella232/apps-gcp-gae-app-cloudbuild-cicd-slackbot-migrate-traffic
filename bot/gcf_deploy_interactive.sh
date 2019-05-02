#!/usr/bin/env bash

gcloud beta functions deploy interactive \
--project=your-project-id \
--region=asia-northeast1 \
--entry-point=Interactive \
--runtime=go111 \
--env-vars-file=.env.yaml \
--trigger-http
