### [1.1.2](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/compare/v1.1.1...v1.1.2) (2021-10-05)


### Bug Fixes

* **api:** Improving validators to fix issues reading inputs from variables ([#40](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/issues/40)) ([7c6793b](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/commit/7c6793b1da150138d33364373318a97d95e87b8c)), closes [#39](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/issues/39)

### [1.1.1](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/compare/v1.1.0...v1.1.1) (2021-09-27)


### Bug Fixes

* **api:** Removing 'Unknown' value validation ([#38](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/issues/38)) ([9dd83e9](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/commit/9dd83e9323f670bf895734e2c5c67a9a31953e43)), closes [#37](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/issues/37)

## [1.1.0](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/compare/v1.0.1...v1.1.0) (2021-09-26)


### Features

* **bin:** Releasing binaries for freebsd and windows ([#33](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/issues/33)) ([543a0e6](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/commit/543a0e68733bc2676f614cd10632a0bbc7458dd7)), closes [#26](https://github.com/foryouandyourcustomers/terraform-provider-sapcc/issues/26)

### [1.0.1](https://github.com/tckb/terraform-provider-sapcc/compare/v1.0.0...v1.0.1) (2021-09-13)


### Bug Fixes

* **api:** Re-enforcing the required variables ([#19](https://github.com/tckb/terraform-provider-sapcc/issues/19)) ([5234fb5](https://github.com/tckb/terraform-provider-sapcc/commit/5234fb5ce0aa44c44bb54176bdd01e00ba01c42e)), closes [#18](https://github.com/tckb/terraform-provider-sapcc/issues/18)

## 1.0.0 (2021-09-12)


### Features

* **api:** Adding dedicated http client ([4afc0f1](https://github.com/tckb/terraform-provider-sapcc/commit/4afc0f13068fd47f746b72b84d22f979405626d7))
* **api:** Adding dedicated http client for deployment resource ([3510e78](https://github.com/tckb/terraform-provider-sapcc/commit/3510e78a7015f9cbedb17503500036392e9474b9))
* **api:** Adding support for tracking deployment progress ([2064a07](https://github.com/tckb/terraform-provider-sapcc/commit/2064a07cf478c0522519e0012609c99524a04664))
* **api:** Adding support for updating deployments on changes ([05ddf7d](https://github.com/tckb/terraform-provider-sapcc/commit/05ddf7df80699ae4c4464d6996772ae5c84f4ef2)), closes [#4](https://github.com/tckb/terraform-provider-sapcc/issues/4)
* **api:** Validating attributes for resource deployment ([#17](https://github.com/tckb/terraform-provider-sapcc/issues/17)) ([0e714ad](https://github.com/tckb/terraform-provider-sapcc/commit/0e714ada63b4e51b528d6ca7f335718a7186ed74)), closes [#11](https://github.com/tckb/terraform-provider-sapcc/issues/11)
* **test:** Sending out errors & warn list ([1c95906](https://github.com/tckb/terraform-provider-sapcc/commit/1c959069a473aef4b63d22bc1a04b9367ef82fe8))


### Bug Fixes

* **api:** Handling edge cases ([35d1c8e](https://github.com/tckb/terraform-provider-sapcc/commit/35d1c8e5cbdbb42e3823bab95541820d6e04e838))
* **api:** Handling empty responses ([e3e01ce](https://github.com/tckb/terraform-provider-sapcc/commit/e3e01ceba8a00e7d371ad1f27b20a081671c92e2))
* **ci:** Adding missing dependencies ([08c58e1](https://github.com/tckb/terraform-provider-sapcc/commit/08c58e1a42e836f83a4d9f05264836a8ffaaa4bc)), closes [#16](https://github.com/tckb/terraform-provider-sapcc/issues/16)
* **ci:** Adding priorities to mock responses for failed builds/deployments ([4083e2c](https://github.com/tckb/terraform-provider-sapcc/commit/4083e2c36d720b7368fcea3d837d4f0419063ec2))
* **ci:** Installing goreleaser ([752e45f](https://github.com/tckb/terraform-provider-sapcc/commit/752e45fc9fb3bb0e668ece221674edac1dc0f32f))
* **ci:** merging goreleaser to sem-version releaser ([594b022](https://github.com/tckb/terraform-provider-sapcc/commit/594b022c1fd72df6eac3166c0a029ce29b6f4308))
* **ci:** testing quick fix for CI? ([a1990aa](https://github.com/tckb/terraform-provider-sapcc/commit/a1990aa08e778bd047fd9d6def65ba9b8af86511))
* **test:** Fixing the acceptance test helper to accomodate CI ([30b57c3](https://github.com/tckb/terraform-provider-sapcc/commit/30b57c307e3e2879e6c5e07c09907a6bdeb22977))

## [1.0.0-beta.2](https://github.com/tckb/terraform-provider-sapcc/compare/v1.0.0-beta.1...v1.0.0-beta.2) (2021-09-12)


### Bug Fixes

* **ci:** Adding missing dependencies ([08c58e1](https://github.com/tckb/terraform-provider-sapcc/commit/08c58e1a42e836f83a4d9f05264836a8ffaaa4bc)), closes [#16](https://github.com/tckb/terraform-provider-sapcc/issues/16)
* **ci:** Installing goreleaser ([752e45f](https://github.com/tckb/terraform-provider-sapcc/commit/752e45fc9fb3bb0e668ece221674edac1dc0f32f))
* **ci:** merging goreleaser to sem-version releaser ([594b022](https://github.com/tckb/terraform-provider-sapcc/commit/594b022c1fd72df6eac3166c0a029ce29b6f4308))

## 1.0.0-beta.1 (2021-09-12)


### Features

* **api:** Adding dedicated http client ([4afc0f1](https://github.com/tckb/terraform-provider-sapcc/commit/4afc0f13068fd47f746b72b84d22f979405626d7))
* **api:** Adding dedicated http client for deployment resource ([3510e78](https://github.com/tckb/terraform-provider-sapcc/commit/3510e78a7015f9cbedb17503500036392e9474b9))
* **api:** Adding support for tracking deployment progress ([2064a07](https://github.com/tckb/terraform-provider-sapcc/commit/2064a07cf478c0522519e0012609c99524a04664))
* **api:** Adding support for updating deployments on changes ([05ddf7d](https://github.com/tckb/terraform-provider-sapcc/commit/05ddf7df80699ae4c4464d6996772ae5c84f4ef2)), closes [#4](https://github.com/tckb/terraform-provider-sapcc/issues/4)
* **test:** Sending out errors & warn list ([1c95906](https://github.com/tckb/terraform-provider-sapcc/commit/1c959069a473aef4b63d22bc1a04b9367ef82fe8))


### Bug Fixes

* **api:** Handling edge cases ([35d1c8e](https://github.com/tckb/terraform-provider-sapcc/commit/35d1c8e5cbdbb42e3823bab95541820d6e04e838))
* **api:** Handling empty responses ([e3e01ce](https://github.com/tckb/terraform-provider-sapcc/commit/e3e01ceba8a00e7d371ad1f27b20a081671c92e2))
* **ci:** Adding priorities to mock responses for failed builds/deployments ([4083e2c](https://github.com/tckb/terraform-provider-sapcc/commit/4083e2c36d720b7368fcea3d837d4f0419063ec2))
* **ci:** testing quick fix for CI? ([a1990aa](https://github.com/tckb/terraform-provider-sapcc/commit/a1990aa08e778bd047fd9d6def65ba9b8af86511))
* **test:** Fixing the acceptance test helper to accomodate CI ([30b57c3](https://github.com/tckb/terraform-provider-sapcc/commit/30b57c307e3e2879e6c5e07c09907a6bdeb22977))
