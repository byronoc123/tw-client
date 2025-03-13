# CloudWatch Dashboard for the application
resource "aws_cloudwatch_dashboard" "app_dashboard" {
  dashboard_name = "${var.app_name}-dashboard"
  
  dashboard_body = jsonencode({
    widgets = [
      # ALB metrics
      {
        type   = "metric"
        x      = 0
        y      = 0
        width  = 12
        height = 6
        properties = {
          metrics = [
            ["AWS/ApplicationELB", "RequestCount", "LoadBalancer", aws_lb.app.arn_suffix, { "stat" = "Sum", "period" = 60 }],
            ["AWS/ApplicationELB", "TargetResponseTime", "LoadBalancer", aws_lb.app.arn_suffix, { "stat" = "Average", "period" = 60 }],
            ["AWS/ApplicationELB", "HTTPCode_ELB_5XX_Count", "LoadBalancer", aws_lb.app.arn_suffix, { "stat" = "Sum", "period" = 60 }],
            ["AWS/ApplicationELB", "HTTPCode_ELB_4XX_Count", "LoadBalancer", aws_lb.app.arn_suffix, { "stat" = "Sum", "period" = 60 }]
          ]
          view    = "timeSeries"
          stacked = false
          region  = var.aws_region
          title   = "ALB Metrics"
        }
      },
      # ECS Service metrics
      {
        type   = "metric"
        x      = 12
        y      = 0
        width  = 12
        height = 6
        properties = {
          metrics = [
            ["AWS/ECS", "CPUUtilization", "ServiceName", aws_ecs_service.app.name, "ClusterName", aws_ecs_cluster.main.name, { "stat" = "Average", "period" = 60 }],
            ["AWS/ECS", "MemoryUtilization", "ServiceName", aws_ecs_service.app.name, "ClusterName", aws_ecs_cluster.main.name, { "stat" = "Average", "period" = 60 }]
          ]
          view    = "timeSeries"
          stacked = false
          region  = var.aws_region
          title   = "ECS Service Metrics"
        }
      },
      # WAF metrics
      {
        type   = "metric"
        x      = 0
        y      = 6
        width  = 12
        height = 6
        properties = {
          metrics = [
            ["AWS/WAFV2", "BlockedRequests", "WebACL", aws_wafv2_web_acl.app_waf.name, "Region", var.aws_region, { "stat" = "Sum", "period" = 60 }],
            ["AWS/WAFV2", "AllowedRequests", "WebACL", aws_wafv2_web_acl.app_waf.name, "Region", var.aws_region, { "stat" = "Sum", "period" = 60 }]
          ]
          view    = "timeSeries"
          stacked = false
          region  = var.aws_region
          title   = "WAF Metrics"
        }
      },
      # Custom application metrics (if configured with CloudWatch agent)
      {
        type   = "metric"
        x      = 12
        y      = 6
        width  = 12
        height = 6
        properties = {
          metrics = [
            ["${var.app_name}", "ApiLatency", { "stat" = "Average", "period" = 60 }],
            ["${var.app_name}", "RpcLatency", { "stat" = "Average", "period" = 60 }],
            ["${var.app_name}", "BlockchainRequests", { "stat" = "Sum", "period" = 60 }]
          ]
          view    = "timeSeries"
          stacked = false
          region  = var.aws_region
          title   = "Application Metrics"
        }
      }
    ]
  })
}

# CloudWatch Alarms - ALB 5XX errors
resource "aws_cloudwatch_metric_alarm" "alb_5xx_errors" {
  alarm_name          = "${var.app_name}-alb-5xx-errors"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "HTTPCode_ELB_5XX_Count"
  namespace           = "AWS/ApplicationELB"
  period              = 60
  statistic           = "Sum"
  threshold           = var.alb_5xx_error_threshold
  alarm_description   = "This alarm monitors for ALB 5XX errors"
  
  dimensions = {
    LoadBalancer = aws_lb.app.arn_suffix
  }
  
  alarm_actions = [aws_sns_topic.app_alarms.arn]
  ok_actions    = [aws_sns_topic.app_alarms.arn]
  
  tags = {
    Name        = "${var.app_name}-alb-5xx-errors"
    Environment = var.environment
  }
}

# CloudWatch Alarms - ALB target response time
resource "aws_cloudwatch_metric_alarm" "alb_response_time" {
  alarm_name          = "${var.app_name}-alb-response-time"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 3
  metric_name         = "TargetResponseTime"
  namespace           = "AWS/ApplicationELB"
  period              = 60
  statistic           = "Average"
  threshold           = var.alb_response_time_threshold
  alarm_description   = "This alarm monitors ALB target response time"
  
  dimensions = {
    LoadBalancer = aws_lb.app.arn_suffix
  }
  
  alarm_actions = [aws_sns_topic.app_alarms.arn]
  ok_actions    = [aws_sns_topic.app_alarms.arn]
  
  tags = {
    Name        = "${var.app_name}-alb-response-time"
    Environment = var.environment
  }
}

# CloudWatch Alarms - ECS Service CPU utilization
resource "aws_cloudwatch_metric_alarm" "ecs_cpu_utilization" {
  alarm_name          = "${var.app_name}-ecs-cpu-utilization"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 3
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = 60
  statistic           = "Average"
  threshold           = var.ecs_cpu_threshold
  alarm_description   = "This alarm monitors ECS CPU utilization"
  
  dimensions = {
    ClusterName = aws_ecs_cluster.main.name
    ServiceName = aws_ecs_service.app.name
  }
  
  alarm_actions = [aws_sns_topic.app_alarms.arn]
  ok_actions    = [aws_sns_topic.app_alarms.arn]
  
  tags = {
    Name        = "${var.app_name}-ecs-cpu-utilization"
    Environment = var.environment
  }
}

# CloudWatch Alarms - ECS Service memory utilization
resource "aws_cloudwatch_metric_alarm" "ecs_memory_utilization" {
  alarm_name          = "${var.app_name}-ecs-memory-utilization"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 3
  metric_name         = "MemoryUtilization"
  namespace           = "AWS/ECS"
  period              = 60
  statistic           = "Average"
  threshold           = var.ecs_memory_threshold
  alarm_description   = "This alarm monitors ECS memory utilization"
  
  dimensions = {
    ClusterName = aws_ecs_cluster.main.name
    ServiceName = aws_ecs_service.app.name
  }
  
  alarm_actions = [aws_sns_topic.app_alarms.arn]
  ok_actions    = [aws_sns_topic.app_alarms.arn]
  
  tags = {
    Name        = "${var.app_name}-ecs-memory-utilization"
    Environment = var.environment
  }
}

# CloudWatch Alarms - WAF blocked requests
resource "aws_cloudwatch_metric_alarm" "waf_blocked_requests" {
  alarm_name          = "${var.app_name}-waf-blocked-requests"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "BlockedRequests"
  namespace           = "AWS/WAFV2"
  period              = 300
  statistic           = "Sum"
  threshold           = var.waf_blocked_requests_threshold
  alarm_description   = "This alarm monitors elevated WAF blocked requests (potential attack)"
  
  dimensions = {
    WebACL = aws_wafv2_web_acl.app_waf.name
    Region = var.aws_region
  }
  
  alarm_actions = [aws_sns_topic.app_alarms.arn]
  
  tags = {
    Name        = "${var.app_name}-waf-blocked-requests"
    Environment = var.environment
  }
}

# CloudWatch Alarms - Target health
resource "aws_cloudwatch_metric_alarm" "target_health" {
  alarm_name          = "${var.app_name}-target-health"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 2
  metric_name         = "HealthyHostCount"
  namespace           = "AWS/ApplicationELB"
  period              = 60
  statistic           = "Average"
  threshold           = var.healthy_host_threshold
  alarm_description   = "This alarm monitors for unhealthy targets"
  
  dimensions = {
    TargetGroup  = aws_lb_target_group.app.arn_suffix
    LoadBalancer = aws_lb.app.arn_suffix
  }
  
  alarm_actions = [aws_sns_topic.app_alarms.arn]
  ok_actions    = [aws_sns_topic.app_alarms.arn]
  
  tags = {
    Name        = "${var.app_name}-target-health"
    Environment = var.environment
  }
}

# SNS Topic for alarms
resource "aws_sns_topic" "app_alarms" {
  name = "${var.app_name}-alarms"
  
  tags = {
    Name        = "${var.app_name}-alarms"
    Environment = var.environment
  }
}

# SNS Topic subscription for email notifications
resource "aws_sns_topic_subscription" "app_alarms_email" {
  count     = var.alarm_email != "" ? 1 : 0
  topic_arn = aws_sns_topic.app_alarms.arn
  protocol  = "email"
  endpoint  = var.alarm_email
}

# Create a CloudWatch Log Group for the application
resource "aws_cloudwatch_log_group" "app_logs" {
  name              = "/ecs/${var.app_name}"
  retention_in_days = 30
  
  tags = {
    Name        = "/ecs/${var.app_name}"
    Environment = var.environment
  }
}

# CloudWatch Composite Alarm for critical system state
resource "aws_cloudwatch_composite_alarm" "critical_system_state" {
  alarm_name          = "${var.app_name}-critical-system-state"
  alarm_description   = "Composite alarm that triggers when multiple system components are in alarm state"
  
  alarm_rule = join(" OR ", [
    "ALARM(${aws_cloudwatch_metric_alarm.alb_5xx_errors.alarm_name})",
    "ALARM(${aws_cloudwatch_metric_alarm.target_health.alarm_name})"
  ])
  
  alarm_actions = [aws_sns_topic.app_alarms.arn]
  ok_actions    = [aws_sns_topic.app_alarms.arn]
  
  tags = {
    Name        = "${var.app_name}-critical-system-state"
    Environment = var.environment
  }
}

# Outputs for monitoring endpoints
output "cloudwatch_dashboard_url" {
  description = "URL to the CloudWatch Dashboard"
  value       = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#dashboards:name=${aws_cloudwatch_dashboard.app_dashboard.dashboard_name}"
}

output "alarm_topic_arn" {
  description = "ARN of the SNS topic for alarms"
  value       = aws_sns_topic.app_alarms.arn
}
