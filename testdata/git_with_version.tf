## simple module call
module "something" {
  source    = "git::https://github.com/cloudposse/terraform-aws-vpc?ref=v1.2.0"
  namespace = "something"
}
