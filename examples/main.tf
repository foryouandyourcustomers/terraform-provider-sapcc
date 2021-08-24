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
  api_baseurl     = "http://localhost:8080"
  auth_token      = "xxxx"
  subscription_id = "myfakesubcription123"
}

data "sapcc_build" "random_build_1" {
  code = "20210819"
}

data "sapcc_deployment" "random_deployment" {
  code = "deploy"
}

resource "sapcc_deployment" "random_deployment" {
  // this is buildcode gives consistent results
  build_code           = "000000.0"
  environment_code     = "d0"
  strategy             = "ROLLING_UPDATE"
  database_update_mode = "NONE"
}


output "build" {
  value = data.sapcc_build.random_build_1
}

output "deploy" {
  value = data.sapcc_deployment.random_deployment
}
output "deploy2" {
  value = sapcc_deployment.random_deployment
}