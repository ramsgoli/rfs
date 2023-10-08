locals {
  user_data_path = "${path.module}/files/userdata.sh"
}

module "master_ec2" {
  source               = "../../infra/terraform/ec2"
  key_name             = aws_key_pair.ec2_jumpbox_keypair.key_name
  user_data            = file(local.user_data_path)
  iam_instance_profile = aws_iam_instance_profile.master_ec2_instance_profile.name

  tags = {
    Name     = "master"
    NodeType = "master"
  }
}

# module "agent1_ec2" {
#   source    = "../../infra/terraform/ec2"
#   user_data = file(local.user_data_path)

#   tags = {
#     Name     = "agent"
#     NodeType = "agent"
#   }
# }

resource "aws_key_pair" "ec2_jumpbox_keypair" {
  key_name   = "ec2_jumpbox_keypair"
  public_key = file("~/.ssh/id_rsa.pub")
}
