# AWSRM

- This is a tool for deleting things in AWS that either take a few clicks per item

> Why not just use the `aws-cli`?

Feel free to use it, this is a personal project to help me learn `golang`.

# Limitations

At the moment all it does is delete S3 buckets.





# AWS Credential Loading Order

Taken from the AWS Go SDK Docs [here](https://docs.aws.amazon.com/sdk-for-go/api/aws/session/)
```
* Shared Credentials file
* Shared Configuration file
```

I have chosen not to support AWS $ENV VARS as most of the tools will be used to work across diffrent profiles. So YMMV