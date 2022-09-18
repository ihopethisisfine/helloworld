locals {
  name   = "helloworld"
  region = "eu-west-1"

  vpc_cidr    = "10.0.0.0/16"
  azs         = slice(data.aws_availability_zones.available.names, 0, 3)
  domain_name = "helloworld.brunoferreira.me"
}
