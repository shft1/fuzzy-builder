#!/usr/bin/env bash
set -euo pipefail

# Backend tests
echo "==> Backend tests"
go test ./...

# Frontend tests
echo "==> Frontend tests"
cd frontend
npm i --silent > /dev/null 2>&1 || true
npm run test --silent
