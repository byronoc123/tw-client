# Blockchain Client

A lightweight and high-performance blockchain client that provides a RESTful API for interacting with EVM-compatible blockchains like Ethereum and Polygon.

## Project Description

This service acts as a simplified abstraction layer over blockchain JSON-RPC API endpoints. It provides a clean RESTful interface for:

- Fetching the latest block number
- Getting detailed block information by block number

The application is containerized using Docker and designed for production deployment on AWS ECS Fargate.

## Installation & Setup

### Prerequisites

- Go 1.19+
- Docker (for containerization)
- AWS CLI (for deployment)
- Terraform (for infrastructure provisioning)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/byronoc123/tw-client.git
cd tw-client
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file with required environment variables:
```bash
RPC_URL=https://polygon-rpc.com/
PORT=8080
TIMEOUT_SECONDS=10
```

4. Run the application:
```bash
go run main.go
```

5. Run tests:
```bash
go test ./...
```

### Docker Build

Build the Docker image locally:
```bash
docker build -t tw-client .
```

Run the containerized application:
```bash
docker run -p 8080:8080 --env-file .env tw-client
```

## API Documentation

### Health Check
```
GET /health
curl http://localhost:8080/health
```
Response:
```json
{
  "status": "ok"
}
```

### Get Latest Block Number
```
GET /api/v1/block/latest
curl http://localhost:8080/api/v1/block/latest
```
Response:
```json
{
  "blockNumber": "0x134e82a"
}
```

### Get Block By Number
```
GET /api/v1/block/:number
curl http://localhost:8080/api/v1/block/0xnumber
```
Parameters:
- `number`: Block number in decimal (e.g., `12345678`) or hexadecimal (e.g., `0xbc614e`) format

Response (example):
```json
{
  "number": "0x134e82a",
  "hash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
  "parentHash": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
  "nonce": "0x0000000000000000",
  "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
  "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
  "stateRoot": "0xdeed4b5727e421c66e2d6b9f12cf8d2cac432dce1c8429b9180891aa6b44b5d5",
  "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
  "miner": "0x0000000000000000000000000000000000000000",
  "difficulty": "0x0",
  "totalDifficulty": "0x0",
  "extraData": "0x",
  "size": "0x0",
  "gasLimit": "0x1c9c380",
  "gasUsed": "0x0",
  "timestamp": "0x6456a1d4",
  "transactions": [],
  "uncles": []
}

## Deployment Instructions

### AWS Deployment with Terraform

1. Configure AWS CLI with appropriate credentials:
```bash
aws configure
```

2. Navigate to the terraform directory:
```bash
cd terraform
```

3. Initialize Terraform:
```bash
terraform init
```

4. Plan the deployment:
```bash
terraform plan -out=tfplan
```

5. Apply the Terraform plan:
```bash
terraform apply tfplan
```

6. To clean up resources:
```bash
terraform destroy
```

#### Infrastructure Components

The Terraform configuration deploys:

- ECS Fargate tasks for running the containerized application
- Application Load Balancer with WAF protection
- CloudWatch alarms for monitoring application health
- SNS topic for alarm notifications

For details on the monitoring infrastructure, see `terraform/OBSERVABILITY_GUIDE.md`.

### Pushing to ECR

1. Authenticate with ECR:
```bash
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin YOUR_AWS_ACCOUNT_ID.dkr.ecr.us-west-2.amazonaws.com
```

2. Tag your Docker image:
```bash
docker tag tw-client:latest YOUR_AWS_ACCOUNT_ID.dkr.ecr.us-west-2.amazonaws.com/blockchain-client:latest
```

3. Push the image:
```bash
docker push YOUR_AWS_ACCOUNT_ID.dkr.ecr.us-west-2.amazonaws.com/blockchain-client:latest
```

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Port the server listens on | `8080` | No |
| `RPC_URL` | Blockchain RPC endpoint URL | `https://polygon-rpc.com/` | No |
| `TIMEOUT_SECONDS` | Timeout for RPC requests in seconds | `10` | No |
| `GIN_MODE` | Gin framework mode (debug/release) | `release` (in Docker) | No |

## Production Considerations

For a production-ready application, consider implementing:

- Authentication and authorization
- Additional rate limiting strategies
- Enhanced monitoring and alerting
- Caching layer for frequently requested blocks
- Circuit breakers for RPC calls
- Horizontal scaling with ECS

## License

[MIT](LICENSE)