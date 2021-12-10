data "archive_file" "s3_lifecycle" {
  type        = "zip"
  source_file = "../out/s3-lifecycle"
  output_path = "../out/s3-lifecycle.zip"
}

resource "aws_lambda_function" "s3_lifecycle" {
  filename         = data.archive_file.s3_lifecycle.output_path
  source_code_hash = data.archive_file.s3_lifecycle.output_base64sha256
  function_name    = local.config_rule_name
  role             = aws_iam_role.s3_lifecycle.arn
  handler          = "s3-lifecycle"
  runtime          = "go1.x"
  architectures    = ["x86_64"]
  depends_on       = [aws_cloudwatch_log_group.s3_lifecycle]
}

resource "aws_lambda_permission" "s3_lifecycle" {
  statement_id   = "AllowConfigToInvoke"
  action         = "lambda:InvokeFunction"
  function_name  = aws_lambda_function.s3_lifecycle.function_name
  principal      = "config.amazonaws.com"
  source_account = data.aws_caller_identity.current.account_id
}

resource "aws_cloudwatch_log_group" "s3_lifecycle" {
  name              = "/aws/lambda/${local.config_rule_name}"
  retention_in_days = 7
}
