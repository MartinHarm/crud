# Secrets Manager Secret for Database Password
resource "aws_secretsmanager_secret" "db_password" {
  name                    = "cruder/db/password"
  description             = "Database password for Cruder application"
  recovery_window_in_days = 7

  tags = {
    Name = "cruder-db-password"
  }
}

# Secrets Manager Secret Version for Database Password
resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = var.db_password
}

# Secrets Manager Secret for API Key
resource "aws_secretsmanager_secret" "api_key" {
  name                    = "cruder/api/key"
  description             = "API key for Cruder application authentication"
  recovery_window_in_days = 7

  tags = {
    Name = "cruder-api-key"
  }
}

# Secrets Manager Secret Version for API Key
resource "aws_secretsmanager_secret_version" "api_key" {
  secret_id     = aws_secretsmanager_secret.api_key.id
  secret_string = var.api_key
}