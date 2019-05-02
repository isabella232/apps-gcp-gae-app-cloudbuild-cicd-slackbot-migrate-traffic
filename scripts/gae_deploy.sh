#!/usr/bin/env bash

# ブランチ名の取得
bn=$(git rev-parse --abbrev-ref HEAD)
gcloud app deploy --project=your-project-id --version=$bn ./appengine/app.yaml
