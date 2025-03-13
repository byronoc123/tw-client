# Simplified Observability Guide

This document describes the simplified observability approach implemented for the blockchain client application.

## Overview

The monitoring infrastructure consists of:

1. **Basic CloudWatch Alarms** - Alerts for critical error conditions
2. **SNS Notifications** - Alert delivery via email

## Essential CloudWatch Alarms

The following critical alarms are maintained:

- **ALB 5XX Errors** - Alerts when the number of 5XX errors exceeds the threshold
- **ALB Response Time** - Alerts when response time exceeds the threshold
- **Target Health** - Monitors the health of target instances
- **ECS CPU/Memory** - Monitors resource utilization
- **WAF Blocked Requests** - Tracks potential security threats

## Alert Notifications

Alerts are delivered via SNS to the configured email address:

- Alarm state changes (both ALARM and OK transitions)
- Includes alarm details and timestamp

## Configuration Variables

The monitoring infrastructure can be customized using the following variables:

```hcl
# Alarm thresholds
alb_5xx_error_threshold = 5     # Count per minute
alb_response_time_threshold = 1.0  # Seconds
ecs_cpu_threshold = 80          # Percentage
ecs_memory_threshold = 80       # Percentage
waf_blocked_requests_threshold = 100  # Count per 5 minutes
healthy_host_threshold = 1      # Minimum healthy targets

# Notification settings
alarm_email = "alerts@example.com"  # Email to receive alerts
```

## Best Practices

1. **Set appropriate thresholds**: Adjust alarm thresholds based on your application's traffic patterns
2. **Implement log retention policies**: Balance data retention with cost considerations
3. **Regularly review alarms**: Check alarm history to spot trends and potential issues

---

For more advanced monitoring needs, consider implementing:
- Custom application metrics using CloudWatch agent
- Detailed dashboards for visualizing application performance
- Synthetic canary tests for endpoint availability monitoring
