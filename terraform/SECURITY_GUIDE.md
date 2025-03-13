# Simplified Security Guide for Blockchain Client

This guide explains the essential security features implemented in the Terraform configuration for the blockchain client application.

## Key Security Features

### 1. HTTPS with AWS Certificate Manager (ACM)

- **SSL/TLS Encryption**: All traffic to the application is encrypted using HTTPS
- **Certificate Management**: AWS Certificate Manager handles certificate creation and renewal
- **HTTP to HTTPS Redirection**: Automatic redirection ensures secure connections

### 2. IP Restrictions

- **Restricted Access**: Access is limited to specified CIDR blocks
- **AWS Health Checker Allowance**: Ensures health checks continue to work

### 3. AWS WAF Protection

- **Web Application Firewall**: Protects against common web exploits
- **Rate Limiting**: Prevents DDoS attacks
- **SQL Injection Protection**: Guards against SQL injection attempts

## Using the Configuration

### Required Variables

```hcl
# Required security variables
domain_name         = "your-domain.example.com"
route53_zone_id     = "Z1234567890ABCDEFGHI"
allowed_cidr_blocks = ["10.0.0.0/8", "192.168.1.0/24"]
```

### Optional Security Variables

```hcl
# Optional security variables
environment             = "prod"
enable_waf              = true
alarm_email             = "alerts@example.com"
```

## Security Best Practices

1. **Restrict Access**: Use the most restrictive CIDR blocks possible
2. **Rotate Credentials**: Regularly rotate AWS credentials and API keys
3. **Monitor Alarms**: Check CloudWatch alarm results regularly
4. **Update WAF Rules**: Keep WAF rules updated based on threat intelligence

## Implementation Checklist

- [ ] Domain name configured and DNS records validated
- [ ] HTTPS certificate validated and working
- [ ] IP restrictions properly configured
- [ ] WAF rules tested and confirmed working

For more information, refer to:
- [AWS ACM Documentation](https://docs.aws.amazon.com/acm/latest/userguide/acm-overview.html)
- [AWS WAF Documentation](https://docs.aws.amazon.com/waf/latest/developerguide/what-is-aws-waf.html)
