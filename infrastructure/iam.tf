resource "aws_iam_role" "ec2_iam_role" {
  name = "master-iam-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy" "master_iam_role_policy" {
  name = "master-iam-role-policy"
  role = aws_iam_role.ec2_iam_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "ec2:Describe*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_instance_profile" "master_ec2_instance_profile" {
  name = "master-instance-profile"
  role = aws_iam_role.ec2_iam_role.name
}
