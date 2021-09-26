name: Provider Bug
description: Report a bug in the provider
title: "[Provider Bug]: ..."
labels: ["bug"]
assignees: ''
body:
- type: markdown
  attributes:
  value: |
  Thank you for reporting the bug

- type: textarea
  id: tf-version
  attributes:
  label: Terraform version
  description: Run `terraform -v` to show the version. If you are not running the latest version of Terraform.
  placeholder: "Terraform v1.0.3 on darwin_amd64"
  validations:
  required: true

- type: textarea
  id: provider-version
  attributes:
  label: Provider version
  description: The version of provider you are on
  placeholder: Current version of the provider
  value: "1.0.0"
  validations:
  required: true

- type: dropdown
  id: OS
  attributes:
  label: Which operating system are you on?
  multiple: false
  options:
  - linux
  - darwin

- type: dropdown
  id: Arch
  attributes:
  label: Which architecture are you on?
  multiple: false
  options:
  - amd64
  - 386
  - arm
  - arm64
  
- type: textarea
  id: what-happened
  attributes:
  label: Describe the bug
  description: A clear and concise description of what the bug is, include any `crash.log, make sure you include all affected resource(s).
  placeholder: Tell us what you see!
  value: "Problem is..."
  validations:
  required: true

  
- type: textarea
  id: logs
  attributes:
  label: Debug logs for `Apply` or `Plan`
  description: 
  render: shell
