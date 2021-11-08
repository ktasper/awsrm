package cmd

// TODO: Read in bucket names from a file (Add a sub command)
// TODO: GoRoutines
// TODO: Add a check to ensure that we are not trying to delete a bucket that is in the same region as an active vpc (Add a local flag to skip the safety check)

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/spf13/cobra"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Deletes an S3 buckets contents & the bucket",
	Long: `Deletes an S3 buckets contents & the bucket:

This will find the buckets you have provided (case insensitive)
find all the buckets that match, and then prompt you if you want to delete
them, if you say yes it will empty the bucket and then delete the bucket.
Unless you are in quiet mode, then it will just empty & delete with no prompt

WARNING: This will delete a bucket and its contents. Double check you
actually want to delete whatever you are using this tool with.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Set an easier var for readability, the args the user passes is the bucket name(s)
		// Take it as its currently in a slice
		bucketNames := args[0]
		if verboseMode {
			fmt.Printf("ℹ️ Verbose: Bucket Name(s): %q \n", bucketNames)
		}

		// Connect to aws and create a session
		sess := awsClient(awsRegion, awsProfile)

		// Create S3 service client
		svc := s3.New(sess)

		// Try and list all the buckets
		if verboseMode {
			fmt.Println("ℹ️ Verbose: Attempting to list all S3 buckets")
		}

		result := listBuckets(svc)

		// Create a slice to hold the found buckets
		var foundBuckets []string

		// Loop over every bucket and if they match our search term add it to the slice.
		for _, b := range result.Buckets {
			bucketName := aws.StringValue(b.Name)
			// Case-insensitive search for buckets
			if strings.Contains(bucketName, strings.ToLower(bucketNames)) || strings.Contains(bucketName, strings.ToUpper(bucketNames)) {
				// Append to slice
				foundBuckets = append(foundBuckets, bucketName)
				if verboseMode {
					fmt.Printf("ℹ️ Verbose: Found bucket %q \n", bucketName)
				}
			}
		}

		// If no buckets are found tell the user and exit
		if len(foundBuckets) == 0 {
			fmt.Printf("❗️ No buckets found matching the search term: %q \n", bucketNames)
			os.Exit(0)
		}

		if !quietMode {
			// Ask if the user wants to delete the buckets
			fmt.Println("Would you like to delete the following buckets? ")
			for _, i := range foundBuckets {
				fmt.Printf("* %s\n", i)
			}
			// Get user input
			userConfirmation := yesNo()
			if !userConfirmation {
				os.Exit(0)
			}
		}

		for _, bucketName := range foundBuckets {
			// First we want to get the region of the bucket
			bucketRegion, err := s3manager.GetBucketRegion(aws.BackgroundContext(), sess, bucketName, awsRegion)
			if err != nil {
				exitErrorf("❗️ Unable to get bucket region, %v", err)
			}
			if verboseMode {
				fmt.Printf("ℹ️ Verbose: Bucket %q is in region %q \n", bucketName, bucketRegion)
				fmt.Printf("ℹ️ Verbose: Changing session region to match the bucket: UserSession=%q BucketSession=%q \n", awsRegion, bucketRegion)
			}
			// connect to the same region as the bucket
			sess = awsClient(bucketRegion, awsProfile)
			// Create S3 service client
			svc := s3.New(sess)
			// If we are not in Dry Run Mode empty the bucket
			if !dryRunMode {
				if verboseMode {
					fmt.Printf("🪣 Attempting to empty: %s\n", bucketName)
				}
				// empty the bucket in the correct region
				emptyBucket(svc, bucketName)
				if verboseMode {
					fmt.Printf("🪣 Attempting to delete: %s\n", bucketName)
				}
				deleteBucket(svc, bucketName)
			} else {
				// If we are in dry run mode just print that we WOULD have tried to empty the bucket
				fmt.Printf("😴 Dry Run: Would have attempted to empty: %s\n", bucketName)
				fmt.Printf("😴 Dry Run: Would have attempted to delete: %s\n", bucketName)
				os.Exit(0)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func emptyBucket(svc *s3.S3, bucketName string) {
	// Create a list iterator to iterate through the list of bucket objects, deleting each object. If an error occurs, call exitErrorf.
	// Try and list all the buckets
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})

	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		exitErrorf("❗️ Unable to delete objects from bucket %q, %v \n", bucketName, err)
	}
	// Once all the items in the bucket have been deleted, inform the user that the objects were deleted.
	fmt.Printf("🪣 Deleted object(s) from bucket: %s \n", bucketName)
}

func deleteBucket(svc *s3.S3, bucketName string) {
	var err error
	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		exitErrorf("❗️ Unable to delete bucket %q, %v \n", bucketName, err)
	}

	// Wait until bucket is deleted before finishing
	fmt.Printf("🪣 Waiting for bucket %q to be deleted...\n", bucketName)
	fmt.Printf("🪣 Deleted Bucket %q \n", bucketName)

	_ = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
}

func awsClient(region string, profile string) *session.Session {
	// Connect to AWS
	sess, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: profile,
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(region),
		},
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		exitErrorf("❗️ Unable to create session to AWS: %v \n", err)
	}
	return sess

}

func listBuckets(svc *s3.S3) *s3.ListBucketsOutput {
	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("❗️ Unable to list buckets, %v", err)
	}
	return result
}
