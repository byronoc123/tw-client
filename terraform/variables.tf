variable "aws_region" {
  description = "The AWS region to deploy resources"
  type        = string
  default     = "us-west-2"
}

variable "app_name" {
  description = "Name of the application"
  type        = string
  default     = "blockchain-client"
}

variable "app_count" {
  description = "Number of Docker containers to run"
  type        = number
  default     = 2
}

variable "app_cpu" {
  description = "Fargate instance CPU units to provision (1 vCPU = 1024 CPU units)"
  type        = string
  default     = "256"
}

variable "app_memory" {
  description = "Fargate instance memory to provision (in MiB)"
  type        = string
  default     = "512"
}

variable "rpc_url" {
  description = "Blockchain RPC URL"
  type        = string
  default     = "https://polygon-rpc.com/"
}

# Security-related variables

variable "domain_name" {
  description = "Domain name for the application (e.g., app.example.com)"
  type        = string
  default     = "blockchain-client.example.com"
}

variable "route53_zone_id" {
  description = "Route 53 hosted zone ID for the domain"
  type        = string
  default     = ""
}

variable "environment" {
  description = "Environment name (e.g., dev, staging, prod)"
  type        = string
  default     = "prod"
}

variable "allowed_cidr_blocks" {
  description = "List of CIDR blocks allowed to access the application"
  type        = list(string)
  default     = ["0.0.0.0/0"] # Default is open, but should be restricted in production
}

variable "enable_deletion_protection" {
  description = "Enable deletion protection for the ALB"
  type        = bool
  default     = true
}

variable "waf_rate_limit" {
  description = "Maximum requests per 5 minutes per IP address"
  type        = number
  default     = 1000
}

variable "ip_blocklist" {
  description = "List of IP addresses to block"
  type        = list(string)
  default     = []
}

variable "log_retention_days" {
  description = "Number of days to retain CloudWatch logs"
  type        = number
  default     = 90
}

variable "enable_waf" {
  description = "Whether to enable AWS WAF protection"
  type        = bool
  default     = true
}

variable "additional_security_headers" {
  description = "Whether to add security headers to HTTP responses"
  type        = bool
  default     = true
}

variable "tls_policy" {
  description = "TLS policy for HTTPS listeners"
  type        = string
  default     = "ELBSecurityPolicy-TLS-1-2-2017-01"
}
