resource "sapcc_deployment" "deployment" {
  # the build code to be used for deployment
  build_code = "000000.0"
  # the environment code to be used for deployment
  environment_code = "d0"
  # the strategy used for deployment
  strategy = "ROLLING_UPDATE"
  # the database update mode for deployment
  database_update_mode = "NONE"
}
