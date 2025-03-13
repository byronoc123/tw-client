# Application Load Balancer
resource "aws_lb" "app" {
  name               = "${var.app_name}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = [aws_subnet.public_a.id, aws_subnet.public_b.id]

  enable_deletion_protection = var.enable_deletion_protection

  # Enable access logs to track all requests
  access_logs {
    bucket  = aws_s3_bucket.alb_logs.id
    prefix  = "${var.app_name}-alb-logs"
    enabled = true
  }

  tags = {
    Name = "${var.app_name}-alb"
    Environment = var.environment
  }
}

# ALB Security Group
resource "aws_security_group" "alb" {
  name        = "${var.app_name}-alb-sg"
  description = "Controls access to the ALB"
  vpc_id      = aws_vpc.main.id

  # HTTPS ingress - restricted to allowed CIDR blocks
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = var.allowed_cidr_blocks
    description = "HTTPS from allowed CIDR blocks"
  }

  # HTTP ingress - only for redirect to HTTPS
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = var.allowed_cidr_blocks
    description = "HTTP redirect to HTTPS from allowed CIDR blocks"
  }

  # Health check from AWS health checkers
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["18.235.0.0/16", "35.172.0.0/14", "52.94.0.0/16"] # AWS health checker IPs
    description = "AWS health checker IPs"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }

  tags = {
    Name = "${var.app_name}-alb-sg"
    Environment = var.environment
  }
}

# S3 bucket for ALB access logs
resource "aws_s3_bucket" "alb_logs" {
  bucket = "${var.app_name}-alb-logs-${var.aws_region}-${var.environment}"
  force_destroy = true

  tags = {
    Name = "${var.app_name}-alb-logs"
    Environment = var.environment
  }
}

# S3 bucket policy to allow ALB to write logs
resource "aws_s3_bucket_policy" "alb_logs" {
  bucket = aws_s3_bucket.alb_logs.id
  policy = data.aws_iam_policy_document.alb_logs.json
}

# Policy document for S3 bucket
data "aws_iam_policy_document" "alb_logs" {
  statement {
    effect = "Allow"
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_elb_service_account.main.id}:root"]
    }
    actions = [
      "s3:PutObject",
    ]
    resources = [
      "${aws_s3_bucket.alb_logs.arn}/${var.app_name}-alb-logs/AWSLogs/${data.aws_caller_identity.current.account_id}/*",
    ]
  }
}

# Get AWS ELB service account
data "aws_elb_service_account" "main" {}

# Get current AWS account ID
data "aws_caller_identity" "current" {}

# Target Group
resource "aws_lb_target_group" "app" {
  name        = "${var.app_name}-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  health_check {
    path                = "/health"
    port                = "traffic-port"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
    matcher             = "200"
  }

  tags = {
    Name = "${var.app_name}-tg"
    Environment = var.environment
  }
}

# HTTPS Listener
resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.app.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-2017-01" # Modern TLS policy
  certificate_arn   = aws_acm_certificate.app_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app.arn
  }

  tags = {
    Name = "${var.app_name}-https-listener"
    Environment = var.environment
  }
}

# HTTP to HTTPS redirect
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.app.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }

  tags = {
    Name = "${var.app_name}-http-redirect"
    Environment = var.environment
  }
}

# Output the ALB DNS name
output "alb_dns_name" {
  value       = aws_lb.app.dns_name
  description = "The DNS name of the load balancer"
}

# Output the ALB HTTPS endpoint
output "alb_https_endpoint" {
  value       = "https://${var.domain_name}"
  description = "The HTTPS endpoint for the application"
}
