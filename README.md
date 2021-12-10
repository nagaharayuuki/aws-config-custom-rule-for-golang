## AWS Config Custom Rule Sample for Golang

### s3-lifecycle
Checks if "S3 Lifecycle configuration" is set for s3 buckets

Params:
- excludeBuckets: Skip compliance check for listed buckets
    - key: excludeBuckets
    - value: Comma separated list of buckets, e.g "bucket1, bucket2"

### Deployment

Tests config rules:
```go test -v -cover ./rules/s3-lifecycle```

Build config rules:
```make all```

Apply config rules:
```
export AWS_PROFILE=xxxxx
export AWS_REGION=xxxxx
cd terraform
terraform init
terraform apply
```
