module "eks_blueprints" {
  source = "github.com/aws-ia/terraform-aws-eks-blueprints?ref=v4.10.0"

  cluster_name    = local.name
  cluster_version = "1.23"

  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnets

  managed_node_groups = {
    spot_2vcpu_8mem = {
      node_group_name        = "mng-spot-2vcpu-8mem"
      create_launch_template = true
      launch_template_os     = "bottlerocket"
      public_ip              = false
      capacity_type          = "SPOT"
      ami_type               = "BOTTLEROCKET_ARM_64"
      instance_types         = ["t4g.large", "m6g.large"] // Instances with same specs for memory and CPU so Cluster Autoscaler scales efficiently
      subnet_ids             = module.vpc.private_subnets
      desired_size           = 2
      max_size               = 10
      min_size               = 2
      max_unavailable        = 1
    }
  }

}

module "eks_blueprints_kubernetes_addons" {
  source = "github.com/aws-ia/terraform-aws-eks-blueprints//modules/kubernetes-addons?ref=v4.10.0"

  eks_cluster_id                 = module.eks_blueprints.eks_cluster_id
  eks_cluster_endpoint           = module.eks_blueprints.eks_cluster_endpoint
  eks_oidc_provider              = module.eks_blueprints.oidc_provider
  eks_cluster_version            = module.eks_blueprints.eks_cluster_version
  eks_cluster_domain             = local.domain_name
  external_dns_route53_zone_arns = [aws_route53_zone.helloworld.arn]

  enable_ingress_nginx = true
  ingress_nginx_helm_config = {
    values = [templatefile("${path.module}/nginx-values.yaml", {
      hostname = local.domain_name
    })]
  }

  enable_aws_load_balancer_controller = true
  enable_external_dns                 = true
  # this is to avoid getting an error because of a datasource used on the external-dns module
  depends_on = [
    aws_route53_zone.helloworld,
    aws_route53_record.helloworld_NS
  ]
}
