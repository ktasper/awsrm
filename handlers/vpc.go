package handlers

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// create a function that searches all vpcs in the account and returns the vpc names
func ListVPCs(svc *ec2.EC2) *ec2.DescribeVpcsOutput {
	result, err := svc.DescribeVpcs(nil)
	if err != nil {
		ExitErrorf("❗️ Unable to list VPCs, %v", err)
	}
	return result
}

// create a function that takes the list of vpcs and returns the vpc names
func ListVPCNames(svc *ec2.EC2) []string {
	result := ListVPCs(svc)
	var vpcNames []string
	for _, vpc := range result.Vpcs {
		// Iterate over the VPCs and append the VPC names to the vpcNames array
		for _, tag := range vpc.Tags {
			if *tag.Key == "Name" {
				vpcNames = append(vpcNames, *tag.Value)
			}
		}
	}
	// We need to trim the "-vpc" from the end of each item in the vpcNames array
	for i, vpcName := range vpcNames {
		vpcNames[i] = strings.TrimSuffix(vpcName, "-vpc")
	}
	return vpcNames
}
