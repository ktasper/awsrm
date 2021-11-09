package handlers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func EmptyBucket(svc *s3.S3, bucketName string) {
	// Create a list iterator to iterate through the list of bucket objects, deleting each object. If an error occurs, call handlers.ExitErrorF.
	// Try and list all the buckets
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})

	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		ExitErrorf("❗️ Unable to delete objects from bucket %q, %v \n", bucketName, err)
	}
	// Once all the items in the bucket have been deleted, inform the user that the objects were deleted.
	fmt.Printf("🪣 Deleted object(s) from bucket: %s \n", bucketName)
}

func DeleteBucket(svc *s3.S3, bucketName string) {
	var err error
	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		ExitErrorf("❗️ Unable to delete bucket %q, %v \n", bucketName, err)
	}

	// Wait until bucket is deleted before finishing
	fmt.Printf("🪣 Waiting for bucket %q to be deleted...\n", bucketName)
	fmt.Printf("🪣 Deleted Bucket %q \n", bucketName)

	_ = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
}

func ListBuckets(svc *s3.S3) *s3.ListBucketsOutput {
	result, err := svc.ListBuckets(nil)
	if err != nil {
		ExitErrorf("❗️ Unable to list buckets, %v", err)
	}
	return result
}
