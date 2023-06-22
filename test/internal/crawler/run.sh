#!/bin/bash

# Run Go tests in the current folder
go test

# Push all code changes to remote repository
git add .
git commit -m "Code updates"
git push origin master
