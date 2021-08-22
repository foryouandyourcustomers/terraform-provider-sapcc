terraform {
  required_providers {
    sapcc = {
      version = "~> 0.0.1"
      source  = "fyayc/sapcc"
    }

  }
  required_version = "~> 1.0.3"
}

provider "sapcc" {
  api_baseurl     = "http://localhost:3030"
  auth_token      = "xxxx"
  subscription_id = "fake-subscrription-123"
}

data "sapcc_build" "random_build_1" {
  code = "20210819.3"
}

data "sapcc_deployment" "random_deployment" {
  code = "200000"
}

output "build" {
  value = data.sapcc_build.random_build_1
}

output "deploy" {
  value = data.sapcc_deployment.random_deployment
}