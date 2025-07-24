package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAwsCodePipelineProduction(t *testing.T) {
	t.Parallel()

	// Give this project a unique name so we can run tests in parallel
	uniqueId := random.UniqueId()
	projectName := fmt.Sprintf("test-prod-app-%s", uniqueId)
	awsRegion := "ap-south-1" // The test will run in this region

	// Configure the options for Terraform
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"project_name":      projectName,
			"aws_region":        awsRegion,
			"github_repo_owner": "gruntwork-io",
			"github_repo_name":  "terratest-hello-world",
		},
		// NOTE: We are NOT using the WorkspaceName feature in this simplified version
		// to ensure compatibility with your Go environment.
	}

	// At the end of the test, run `terraform destroy` to clean up all resources.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`.
	terraform.InitAndApply(t, terraformOptions)

	// Get the website URL output from Terraform
	websiteURL := terraform.Output(t, terraformOptions, "website_url")
	assert.NotEmpty(t, websiteURL, "Website URL output is empty")

	// Validate that the correct number of EC2 instances were launched by the ASG.
	// We find them by looking for the tag we assigned in the Launch Template.
	// This is a stable function that should work with your library version.
	instanceIds := aws.GetEc2InstanceIdsByTag(t, awsRegion, "Name", fmt.Sprintf("%s-dev-web-server", projectName))
	assert.Equal(t, 2, len(instanceIds), "Expected to find exactly 2 EC2 instances with the correct tag")

	t.Logf("Successfully validated that 2 instances were created for the ASG.")
}
