module "local" {
  source = "../../local"
}

module "local" {
  source = "../../modules/vpc-endpoints"
}

module "remote" {
  source = "git::https://github.com/something/module?ref=1.0.0"
}
