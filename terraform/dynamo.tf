resource "aws_dynamodb_table" "users_table" {
  name         = "users"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "username"

  attribute {
    name = "username"
    type = "S"
  }
}
