resource "aws_config_config_rule" "s3_lifecycle" {
  description = "Checks if S3 Lifecycle configuration is set for s3 buckets"
  name        = local.config_rule_name
  input_parameters = jsonencode(
    {
      excludeBuckets = var.exclude_buckets
    }
  )
  scope {
    compliance_resource_types = [
      "AWS::S3::Bucket",
    ]
  }
  source {
    owner             = "CUSTOM_LAMBDA"
    source_identifier = aws_lambda_function.s3_lifecycle.arn
    source_detail {
      event_source = "aws.config"
      message_type = "ConfigurationItemChangeNotification"
    }
    source_detail {
      event_source = "aws.config"
      message_type = "OversizedConfigurationItemChangeNotification"
    }
  }
}
