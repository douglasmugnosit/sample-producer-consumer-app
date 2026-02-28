terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# ---------------------------------------------------------------------------
# SQS Queue — messages flow from producer to consumer through this queue
# ---------------------------------------------------------------------------
resource "aws_sqs_queue" "orders" {
  name                       = "${var.project_name}-orders"
  message_retention_seconds  = 86400  # 1 day
  visibility_timeout_seconds = 30
  receive_wait_time_seconds  = 5 # long polling

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# ---------------------------------------------------------------------------
# DynamoDB Table — consumer stores processed messages here
# ---------------------------------------------------------------------------
resource "aws_dynamodb_table" "orders" {
  name         = "${var.project_name}-orders"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# ---------------------------------------------------------------------------
# IAM Role — shared role for both producer and consumer (ECS / EC2 / local)
# ---------------------------------------------------------------------------
resource "aws_iam_role" "app_role" {
  name = "${var.project_name}-app-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# Producer permissions — send messages to SQS
resource "aws_iam_role_policy" "producer_policy" {
  name = "${var.project_name}-producer"
  role = aws_iam_role.app_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["sqs:SendMessage"]
        Effect   = "Allow"
        Resource = aws_sqs_queue.orders.arn
      }
    ]
  })
}

# Consumer permissions — read/delete from SQS + write to DynamoDB
resource "aws_iam_role_policy" "consumer_policy" {
  name = "${var.project_name}-consumer"
  role = aws_iam_role.app_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ]
        Effect   = "Allow"
        Resource = aws_sqs_queue.orders.arn
      },
      {
        Action = [
          "dynamodb:PutItem"
        ]
        Effect   = "Allow"
        Resource = aws_dynamodb_table.orders.arn
      }
    ]
  })
}
