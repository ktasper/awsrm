package handlers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

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

// Create a function that returns all objects keys on a bucket
func ListObjects(svc *s3.S3, bucketName string) ([]string, error) {
	var keys []string
	err := svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			keys = append(keys, *obj.Key)
		}
		return !lastPage
	})
	if err != nil {
		ExitErrorf("❗️ Unable to list objects, %v", err)
	}

	return keys, nil
}

func EmptyBucketV2(svc *s3.S3, bucketName string, debugMode bool) (ret bool) {
	// TODO dump the iters out of the function

	// Get a list of all the objects in the bucket
	objects, err := ListObjects(svc, bucketName)
	// Iterate over each object and delete it
	for _, object_key := range objects {
		result, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(object_key),
		})
		// Print the result of the deletion if debug mode flag is set
		if debugMode {
			fmt.Printf("🗑️ Deleted object %q \n", result)
		}
		if err != nil {
			ExitErrorf("❗️ Unable to delete object, %v", err)
		}
	}
	// We need to get the object versions
	object_versions, _ := svc.ListObjectVersions(&s3.ListObjectVersionsInput{
		Bucket: aws.String(bucketName),
	})

	// For each object version delete it
	for _, current_object_version_metadata := range object_versions.Versions {
		result, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket:    aws.String(bucketName),
			Key:       current_object_version_metadata.Key,
			VersionId: current_object_version_metadata.VersionId,
		})
		// Print the result of the deletion
		if debugMode {
			fmt.Printf("🗑️ Deleted object %q \n", result)
		}
		if err != nil {
			ExitErrorf("❗️ Unable to delete object, %v", err)
		}
	}

	// for each object version delete marker, delete it so the object is permanently deleted
	for _, delete_marker := range object_versions.DeleteMarkers {
		result, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket:    aws.String(bucketName),
			Key:       delete_marker.Key,
			VersionId: delete_marker.VersionId,
		})
		// Print the result of the deletion
		if debugMode {
			fmt.Printf("🗑️ Deleted object %q \n", result)
		}
		if err != nil {
			ExitErrorf("❗️ Unable to delete object, %v", err)
		}
	}

	if err != nil {
		ExitErrorf("❗️ Unable to list objects, %v", err)
	}

	return true
}
