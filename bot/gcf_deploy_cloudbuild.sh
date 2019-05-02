#!/usr/bin/env bash

gcloud beta functions deploy cloudbuild \
--project=your-project-id \
--region=asia-northeast1 \
--entry-point=CloudBuild \
--stage-bucket=your-project-id_cloudbuild \
--runtime=go111 \
--env-vars-file=.env.yaml \
--trigger-topic=cloud-builds
