# AWSRM

- This is a tool for deleting things in AWS that either take a few clicks per item

> Why not just use the `aws-cli`?

Feel free to use it, this is a personal project to help me learn `golang`.

# Limitations

At the moment all it does is delete S3 buckets.





# AWS Credential Loading Order

Taken from the AWS Go SDK Docs [here](https://docs.aws.amazon.com/sdk-for-go/api/aws/session/)
```
* Environment Variables
* Shared Credentials file
* Shared Configuration file (if SharedConfig is enabled)
* EC2 Instance Metadata (credentials only)
```