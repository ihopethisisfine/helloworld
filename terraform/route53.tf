resource "aws_route53_zone" "helloworld" {
  name = local.domain_name
}

resource "aws_route53_record" "helloworld_NS" {
  allow_overwrite = true
  name            = local.domain_name
  ttl             = 60
  type            = "NS"
  zone_id         = aws_route53_zone.helloworld.zone_id

  records = [
    aws_route53_zone.helloworld.name_servers[0],
    aws_route53_zone.helloworld.name_servers[1],
    aws_route53_zone.helloworld.name_servers[2],
    aws_route53_zone.helloworld.name_servers[3],
  ]
}
