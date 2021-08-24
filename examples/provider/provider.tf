provider "sapcc" {
  api_baseurl     = var.base_url        # optional
  auth_token      = var.auth_token      # or export SAPCC_AUTH_TOKEN=..
  subscription_id = var.subscription_id # or export SAPCC_SUBSCRIPTION_ID=..
}