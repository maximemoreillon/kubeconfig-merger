#!/bin/bash

outDir="./dist/"
outFileName="kubeconfig_merger"
outFilePath="$outDir/$outFileName"


GOOS=linux GOARCH=amd64 go build -o "${outFilePath}_linux_amd64" .
GOOS=windows GOARCH=amd64 go build -o "${outFilePath}_windows_amd64.exe" .