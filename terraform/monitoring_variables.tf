# Simplified monitoring variables

variable "alb_5xx_error_threshold" {
  description = "Threshold for ALB 5XX error alarm (count per minute)"
  type        = number
  default     = 5
}

variable "alarm_email" {
  description = "Email address for alarm notifications"
  type        = string
  default     = ""
}
