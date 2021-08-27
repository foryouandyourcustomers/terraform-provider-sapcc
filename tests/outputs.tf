output "build" {
  value = data.sapcc_build.random_build_1
}
output "deploy" {
  value = data.sapcc_deployment.random_deployment
}
output "deploy2" {
  value = sapcc_deployment.random_deployment
}