#!/bin/bash
# Creates a bunch of S3 buckets for testing

# Count to 10 and create a bucket for each
for i in {1..5}; do
  aws s3 --profile $1 mb s3://kwtestbucket$i
done