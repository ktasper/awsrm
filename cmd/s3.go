package cmd

// TODO: Read in bucket names from a file (Add a sub command)
// TODO: GoRoutines
import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ktasper/awsrm/handlers"

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
			fmt.Printf("ℹ️ - Verbose: Bucket Name(s): %q \n", bucketNames)
		}

		// Connect to aws and create a session
		sess := handlers.AwsClient(awsRegion, awsProfile)

		// Create S3 service client
		svc := s3.New(sess)

		// Try and list all the buckets
		if verboseMode {
			fmt.Println("ℹ️ - Verbose: Attempting to list all S3 buckets")
		}

		result := handlers.ListBuckets(svc)

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
					fmt.Printf("ℹ️ - Verbose: Found bucket %q \n", bucketName)
				}
			}
		}

		// If no buckets are found tell the user and exit
		if len(foundBuckets) == 0 {
			fmt.Printf("❗️ - No buckets found matching the search term: %q \n", bucketNames)
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

	BUCKET_LOOP:
		for _, bucketName := range foundBuckets {
			// First we want to get the region of the bucket
			bucketRegion, err := s3manager.GetBucketRegion(aws.BackgroundContext(), sess, bucketName, awsRegion)
			if err != nil {
				handlers.ExitErrorf("❗️ - Unable to get bucket region, %v", err)
			}
			if verboseMode {
				fmt.Printf("ℹ️ - Verbose: Bucket %q is in region %q \n", bucketName, bucketRegion)
				fmt.Printf("ℹ️ - Verbose: Changing session region to match the bucket: UserSession=%q BucketSession=%q \n", awsRegion, bucketRegion)
			}
			// connect to the same region as the bucket
			sess = handlers.AwsClient(bucketRegion, awsProfile)
			// Create S3 service client in the same region as the bucket
			svc := s3.New(sess)
			// Also create a EC2 service client in the same region as the bucket
			vpcSvc := ec2.New(sess)
			// Get the names of all the vpc's in the region
			vpc_names := handlers.ListVPCNames(vpcSvc)
			for _, vpc_name := range vpc_names {
				// If the vpc_name is in the bucket name and skipVpcCheck is not set, then print a warning and skip the bucket
				if strings.Contains(bucketName, vpc_name) && !skipVpcCheck {
					fmt.Printf("❗️ - Bucket %q is in the same region as an active VPC %q, skipping \n", bucketName, vpc_name)
					continue BUCKET_LOOP
				}
			}
			// If we are not in Dry Run Mode empty the bucket
			if !dryRunMode {
				if verboseMode {
					fmt.Printf("🪣 - Attempting to empty: %s\n", bucketName)
				}
				// empty the bucket in the correct region
				handlers.EmptyBucket(svc, bucketName)
				if verboseMode {
					fmt.Printf("🪣 - Attempting to delete: %s\n", bucketName)
				}
				handlers.DeleteBucket(svc, bucketName)
			} else {
				// If we are in dry run mode just print that we WOULD have tried to empty the bucket
				fmt.Printf("😴 - Dry Run: Would have attempted to empty: %s\n", bucketName)
				fmt.Printf("😴 - Dry Run: Would have attempted to delete: %s\n", bucketName)
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
