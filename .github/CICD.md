# CI/CD Pipeline for Blockchain Client

This document describes the continuous integration and continuous deployment (CI/CD) pipeline implemented for the blockchain client application using GitHub Actions.

## Pipeline Overview

The CI/CD pipeline automates the following processes:

1. **Code Quality Checks**: Linting and formatting validation
2. **Deployment**: Automated deployment to staging and production environments
3. **Release Management**: Automated GitHub releases

## Workflow Triggers

The pipeline is triggered by the following events:

- **Push to `develop` branch**: Runs the full pipeline and deploys to staging
- **Push to `main` branch**: Runs the full pipeline and deploys to production (after approval)
- **Pull Requests**: Runs linting without deployment
- **Manual Trigger**: Can be manually triggered with environment selection

## Pipeline Stages

### 1. Code Linting

- Uses `golangci-lint` to enforce code quality standards
- Verifies proper code formatting with `gofmt`
- Fails the pipeline if code doesn't meet standards

### 2. Docker Image Building

- Builds a Docker image using multi-stage builds
- Pushes the image to GitHub Container Registry (ghcr.io)
- Tags images with git SHA, branch name, and semantic versions

### 3. Staging Deployment

- Triggered automatically for pushes to the `develop` branch
- Uses Terraform to deploy to AWS
- Sets up proper infrastructure with appropriate variables for staging
- Verifies the deployment by checking the health endpoint

### 4. Production Approval

- Requires manual approval from authorized team members
- Uses GitHub Environments protection rules for approval workflow

### 5. Production Deployment

- Deploys to production environment with appropriate variables
- Enables additional security measures specific to production
- Performs post-deployment verification
- Creates a GitHub release for successful deployments

## Environment Configuration

The pipeline uses GitHub Environments to manage different deployment targets:

- **staging**: Development/test environment
- **production-approval**: Gating environment for approvals
- **production**: Production environment

## Secrets Management

The following secrets need to be configured in GitHub:

### AWS Credentials
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

### Terraform Configuration
- `TF_STATE_BUCKET` - S3 bucket for Terraform state

### Production Configuration
- `ALLOWED_CIDR_BLOCKS` - Comma-separated list of allowed CIDR blocks
- `DOMAIN_NAME` - Production domain name
- `ROUTE53_ZONE_ID` - AWS Route53 zone ID
- `ALARM_EMAIL` - Email for CloudWatch alarms

## Usage Instructions

### Running the Pipeline Manually

1. Go to the "Actions" tab in the GitHub repository
2. Select "CI/CD Pipeline" from the workflows list
3. Click "Run workflow"
4. Select the branch and environment to deploy to
5. Click "Run workflow"

### Monitoring Deployments

1. Review the outputs from the Terraform apply step for endpoint URLs
2. CloudWatch dashboards are automatically created for monitoring
3. Alarms will be sent to the configured email address

## Troubleshooting

### Common Issues

#### Linting Failures
- Review the linting output in the GitHub Actions logs
- Fix formatting or code quality issues as identified

#### Deployment Failures
- Check the Terraform output for specific errors
- Verify AWS credentials and permissions
- Review CloudWatch logs for the application

## Extending the Pipeline

### Adding Custom Steps
Modify the pipeline YAML file to include additional jobs or steps as needed.

### Adding New Environments
1. Create a new GitHub Environment in repository settings
2. Add required secrets for the environment
3. Duplicate and modify the deployment job in the workflow file
