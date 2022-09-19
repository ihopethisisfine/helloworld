locals {
  oidc_fully_qualified_subjects = format("system:serviceaccount:%s:%s", "default", local.name)
}

#Role to be assumed by helloworld's service account
resource "aws_iam_role" "helloworld" {
  name = local.name
  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRoleWithWebIdentity"
      Effect = "Allow"
      Principal = {
        Federated = module.eks_blueprints.eks_oidc_issuer_url
      }
      Condition = {
        StringEquals = {
          format("%s:sub", module.eks_blueprints.eks_oidc_provider_arn) = local.oidc_fully_qualified_subjects
        }
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy" "dynamodb_access_policy" {
  name = "dynamodb_access"
  role = aws_iam_role.helloworld.id

  policy = data.aws_iam_policy_document.dynamodb_access.json
}

data "aws_iam_policy_document" "dynamodb_access" {
  statement {
    sid       = "ListAndDescribe"
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "dynamodb:List*",
      "dynamodb:DescribeReservedCapacity*",
      "dynamodb:DescribeLimits",
      "dynamodb:DescribeTimeToLive",
    ]
  }

  statement {
    sid       = "UsersTable"
    effect    = "Allow"
    resources = [aws_dynamodb_table.users_table.arn]

    actions = [
      "dynamodb:BatchGet*",
      "dynamodb:DescribeStream",
      "dynamodb:DescribeTable",
      "dynamodb:Get*",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchWrite*",
      "dynamodb:CreateTable",
      "dynamodb:Delete*",
      "dynamodb:Update*",
      "dynamodb:PutItem",
    ]
  }
}
