# Changelog
All notable changes to this project will be documented in this file.

## [unreleased]

### Bug Fixes

#### API

- Handling empty responses ([e3e01ce](e3e01ceba8a00e7d371ad1f27b20a081671c92e2))
- Adapting to edge cases ([bd7ef00](bd7ef0073e3a548744d783ebd948948afcd90121))

#### CI

- Testing quick fix for CI? ([a1990aa](a1990aa08e778bd047fd9d6def65ba9b8af86511))
- Adding priorities to mock responses for failed builds/deployments ([4083e2c](4083e2c36d720b7368fcea3d837d4f0419063ec2))

#### TEST

- Fixing the acceptance test helper to accomodate CI ([30b57c3](30b57c307e3e2879e6c5e07c09907a6bdeb22977))

### Features

#### API

- Adding dedicated http client ([4afc0f1](4afc0f13068fd47f746b72b84d22f979405626d7))
- Adding dedicated http client for deployment resource ([3510e78](3510e78a7015f9cbedb17503500036392e9474b9))
- Adding support for tracking deployment progress ([2064a07](2064a07cf478c0522519e0012609c99524a04664))
- Adding support for updating deployments on changes ([05ddf7d](05ddf7df80699ae4c4464d6996772ae5c84f4ef2))

#### TEST

- Sending out errors & warn list ([1c95906](1c959069a473aef4b63d22bc1a04b9367ef82fe8))

### Miscellaneous Tasks

#### BUILD

- Removing test rule ([8e7b42b](8e7b42b97922e6aa227dffef8fe5a00c4814bd03))

#### DEPS

- Bump github.com/hashicorp/go-hclog from 0.15.0 to 0.16.2 (#2) ([0cfd0ad](0cfd0ad5fe17c6ebbdae4a0e1cdb5cab11c36276))

#### DOCS

- Updated Readme documentaton for acceptance tests ([b10c352](b10c3523e9978a4e560c4867c38d1d55ba932fff))

#### TEST

- Adding missing test signature ([3a06fb3](3a06fb322327bd60f28097ff9b2009aab7b6519d))
- Cleaning up the error messages ([e83f609](e83f609034ff4364772bca59af01b53afb97e233))

### Refactor

#### DEPS

- Bumping to 0.3.0 of plugin framework ([367862d](367862da6c5a87acf4e547f13ff2dd3a852afefa))

#### MODEL

- Moving models separate package to avoid cycle dependency ([2826d52](2826d5251190c5606108fb800a7fb548ab2e4051))

### Testing

