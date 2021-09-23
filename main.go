package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/manifoldco/promptui"
)

func main() {

	// Define flags
	useProfile := flag.String("profile", "default", "The aws profile you want to use.")
	useRegion := flag.String("region", "eu-west-1", "The aws region you want to use")
	useBucket := flag.String("buckets", "", "The bucket(s) you want to destroy")
	verboseMode := flag.Bool("verbose", false, "Verbose mode")
	dryRunMode := flag.Bool("dry-run", false, "A mode to see what would happen")
	flag.Parse()

	if *verboseMode {
		fmt.Printf("Verbose: Profile = %q \n", *useProfile)
		fmt.Printf("Verbose: Region = %q \n", *useRegion)
		fmt.Printf("Verbose: Bucket(s) = %q \n", *useBucket)
	}

	// Check to see if the bucket search term is empty, if so tell the user and exit
	if len(*useBucket) == 0 {
		fmt.Println("--buckets is required but found to be empty")
		os.Exit(1)
	}

	// credentials from the shared credentials file ~/.aws/credentials.
	err := os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	if *verboseMode {
		fmt.Println("Verbose: Setting ENV AWS_SDK_LOAD_CONFIG=True")
	}
	if err != nil {
		return
	}
	// Set the AWS ENV VAR to use the profile we want
	err = os.Setenv("AWS_PROFILE", *useProfile)
	if *verboseMode {
		fmt.Printf("Verbose: Setting AWS_PROFILE=%q \n", *useProfile)
	}
	if err != nil {
		return
	}

	// If the region flag is not set just

	// Connect to AWS
	if *verboseMode {
		fmt.Println("Verbose: Attempting to connect to AWS")
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(*useRegion),
			CredentialsChainVerboseErrors: aws.Bool(true)},
		Profile: *useProfile,
	})
	if err != nil {
		exitErrorf("Unable to connect to AWS", err)
	}
	// Create S3 service client
	svc := s3.New(sess)
	// Try and list all the buckets
	if *verboseMode {
		fmt.Println("Verbose: Attempting to list all S3 buckets")
	}
	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	// Create a slice to hold the found buckets
	var foundBuckets []string

	// Loop over every bucket and if they match our search term add it to the slice.
	for _, b := range result.Buckets {
		bucketName := aws.StringValue(b.Name)
		// Case-insensitive search for buckets
		if strings.Contains(bucketName, strings.ToLower(*useBucket)) || strings.Contains(bucketName, strings.ToUpper(*useBucket)) {
			// Append to slice
			foundBuckets = append(foundBuckets, bucketName)
			if *verboseMode {
				fmt.Printf("Verbose: Found bucket %q \n", bucketName)
			}
		}
	}

	// If no buckets are found tell the user and exit
	if len(foundBuckets) == 0 {
		fmt.Printf("No buckets found matching the search term: %q \n", *useBucket)
		os.Exit(0)
	}

	// Ask if the user wants to delete the buckets
	fmt.Println("Would you like to delete the following buckets? ")
	for _, i := range foundBuckets {
		fmt.Printf("* %s\n", i)
	}
	userConfirmation := yesNo()
	if userConfirmation {
		if *verboseMode {
			fmt.Println("Verbose: User Confirmation = Yes")
		}
		for _, bucketName := range foundBuckets {
			if !*dryRunMode {
				fmt.Printf("Attempting to delete: %s\n", bucketName)
			}
			// If we are not in dry-run mode actually attempt to empty the bucket
			if !*dryRunMode {
				if *verboseMode {
					fmt.Printf("Verbose: Attempt to empty %q \n", bucketName)
				}
				emtpyBucket(svc, bucketName)
			} else {
				fmt.Printf("Dry Run: Would have attempted to empty %q \n", bucketName)
			}
			// If we are not in dry-run mode actually attempt to delete the bucket
			if !*dryRunMode {
				if *verboseMode {
					fmt.Printf("Verbose: Attempt to delete %q \n", bucketName)
				}
				deleteBucket(svc, bucketName)
			} else {
				fmt.Printf("Dry Run: Would have attempted to delete %q \n", bucketName)
			}

		}
		os.Exit(0)
	} else {
		if *verboseMode {
			fmt.Println("Verbose: User Confirmation = No")
		}
		os.Exit(0)
	}
}

// Error handling
func exitErrorf(msg string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, msg+"\n", args...)
	if err != nil {
		return
	}
	os.Exit(1)
}

// Fancy prompt for users to get a bool value
func yesNo() bool {
	prompt := promptui.Select{
		Label: "Select[Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

func emtpyBucket(svc *s3.S3, bucketName string) {
	// Create a list iterator to iterate through the list of bucket objects, deleting each object. If an error occurs, call exitErrorf.
	// Try and list all the buckets
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})

	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		exitErrorf("Unable to delete objects from bucket %q, %v \n", bucketName, err)
	}
	// Once all the items in the bucket have been deleted, inform the user that the objects were deleted.
	fmt.Printf("Deleted object(s) from bucket: %s \n", bucketName)
}

func deleteBucket(svc *s3.S3, bucketName string) {
	var err error
	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		exitErrorf("Unable to delete bucket %q, %v \n", bucketName, err)
	}

	// Wait until bucket is deleted before finishing
	fmt.Printf("Waiting for bucket %q to be deleted...\n", bucketName)
	fmt.Printf("Deleted Bucket %q \n", bucketName)

	_ = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
}
