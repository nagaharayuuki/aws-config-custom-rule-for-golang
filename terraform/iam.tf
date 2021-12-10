data "aws_iam_policy_document" "lambda" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "s3_lifecycle" {
  name               = local.config_rule_name
  assume_role_policy = data.aws_iam_policy_document.lambda.json
}

data "aws_iam_policy_document" "s3_lifecycle" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:*:*:log-group:/aws/lambda/${local.config_rule_name}",
      "arn:aws:logs:*:*:log-group:/aws/lambda/${local.config_rule_name}:*"
    ]
  }
  statement {
    actions = [
      "config:PutEvaluations"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "s3_lifecycle" {
  name        = "s3_lifecycle"
  description = "IAM policy for logging and config from a lambda"
  policy      = data.aws_iam_policy_document.s3_lifecycle.json
}

resource "aws_iam_role_policy_attachment" "s3_lifecycle" {
  role       = aws_iam_role.s3_lifecycle.name
  policy_arn = aws_iam_policy.s3_lifecycle.arn
}
