import {
    to=aws_iam_role.lambda_role
    id="payment-email-sender-lambda-role"
}

import {
    to=aws_s3_bucket.lambda_deployments
    id="arlequines-lambda-deployments"
}

import {
    to=aws_lambda_function.payment_email
    id="payment-email-sender"
}

import {
    to=aws_cloudwatch_log_group.lambda_logs
    id="/aws/lambda/payment-email-sender"
}

import {
  to = aws_iam_role_policy.lambda_policy
  id = "payment-email-sender-lambda-role:payment-email-sender-lambda-policy"
}

import {
    to = aws_lambda_event_source_mapping.sqs_trigger
    id = "93e3f6ec-7806-4d85-a73f-6ca32db641d0"
}