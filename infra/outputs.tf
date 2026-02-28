output "sqs_queue_url" {
  description = "URL of the SQS queue — set as SQS_QUEUE_URL env var for both apps"
  value       = aws_sqs_queue.orders.url
}

output "sqs_queue_arn" {
  description = "ARN of the SQS queue"
  value       = aws_sqs_queue.orders.arn
}

output "dynamodb_table_name" {
  description = "DynamoDB table name — set as DYNAMODB_TABLE env var for the consumer"
  value       = aws_dynamodb_table.orders.name
}

output "dynamodb_table_arn" {
  description = "ARN of the DynamoDB table"
  value       = aws_dynamodb_table.orders.arn
}

output "app_role_arn" {
  description = "IAM role ARN for the applications"
  value       = aws_iam_role.app_role.arn
}
