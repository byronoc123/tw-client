# AWS WAF Web ACL for the Load Balancer
resource "aws_wafv2_web_acl" "app_waf" {
  name        = "${var.app_name}-waf-acl"
  description = "WAF Web ACL for ${var.app_name}"
  scope       = "REGIONAL"

  default_action {
    allow {}
  }

  # AWS Managed Rules - Core rule set
  rule {
    name     = "AWS-AWSManagedRulesCommonRuleSet"
    priority = 1

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"
        
        # Optionally exclude specific rules if they cause problems
        excluded_rule {
          name = "SizeRestrictions_BODY"
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "${var.app_name}-aws-common-rule"
      sampled_requests_enabled   = true
    }
  }

  # AWS Managed Rules - Known bad inputs
  rule {
    name     = "AWS-AWSManagedRulesKnownBadInputsRuleSet"
    priority = 2

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesKnownBadInputsRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "${var.app_name}-aws-bad-inputs-rule"
      sampled_requests_enabled   = true
    }
  }

  # Custom rate-based rule to prevent DDoS
  rule {
    name     = "RateLimitRule"
    priority = 3

    action {
      block {}
    }

    statement {
      rate_based_statement {
        limit              = var.waf_rate_limit
        aggregate_key_type = "IP"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "${var.app_name}-rate-limit-rule"
      sampled_requests_enabled   = true
    }
  }

  # Custom IP filtering rule based on blocklist
  dynamic "rule" {
    for_each = length(var.ip_blocklist) > 0 ? [1] : []
    content {
      name     = "IPBlocklistRule"
      priority = 4

      action {
        block {}
      }

      statement {
        ip_set_reference_statement {
          arn = aws_wafv2_ip_set.blocklist[0].arn
        }
      }

      visibility_config {
        cloudwatch_metrics_enabled = true
        metric_name                = "${var.app_name}-ip-blocklist-rule"
        sampled_requests_enabled   = true
      }
    }
  }

  # Custom rule to protect against SQL injection
  rule {
    name     = "SQLiProtectionRule"
    priority = 5

    action {
      block {}
    }

    statement {
      or_statement {
        statement {
          sqli_match_statement {
            field_to_match {
              body {}
            }
            text_transformation {
              priority = 1
              type     = "URL_DECODE"
            }
            text_transformation {
              priority = 2
              type     = "HTML_ENTITY_DECODE"
            }
          }
        }
        statement {
          sqli_match_statement {
            field_to_match {
              query_string {}
            }
            text_transformation {
              priority = 1
              type     = "URL_DECODE"
            }
            text_transformation {
              priority = 2
              type     = "HTML_ENTITY_DECODE"
            }
          }
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "${var.app_name}-sqli-rule"
      sampled_requests_enabled   = true
    }
  }

  # Logging configuration
  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name                = "${var.app_name}-web-acl"
    sampled_requests_enabled   = true
  }

  tags = {
    Name        = "${var.app_name}-waf-acl"
    Environment = var.environment
  }
}

# IP Set for blocklisted IPs
resource "aws_wafv2_ip_set" "blocklist" {
  count              = length(var.ip_blocklist) > 0 ? 1 : 0
  name               = "${var.app_name}-ip-blocklist"
  description        = "IP blocklist for ${var.app_name}"
  scope              = "REGIONAL"
  ip_address_version = "IPV4"
  addresses          = var.ip_blocklist

  tags = {
    Name        = "${var.app_name}-ip-blocklist"
    Environment = var.environment
  }
}

# WAF Web ACL Association with ALB
resource "aws_wafv2_web_acl_association" "app_waf_alb" {
  resource_arn = aws_lb.app.arn
  web_acl_arn  = aws_wafv2_web_acl.app_waf.arn
}

# Logging configuration for WAF
resource "aws_wafv2_web_acl_logging_configuration" "app_waf" {
  log_destination_configs = [aws_cloudwatch_log_group.waf_logs.arn]
  resource_arn            = aws_wafv2_web_acl.app_waf.arn
  redacted_fields {
    single_header {
      name = "authorization"
    }
    single_header {
      name = "cookie"
    }
  }
}

# CloudWatch Log Group for WAF logs
resource "aws_cloudwatch_log_group" "waf_logs" {
  name              = "/aws/waf/logs/${var.app_name}"
  retention_in_days = var.log_retention_days

  tags = {
    Name        = "${var.app_name}-waf-logs"
    Environment = var.environment
  }
}
