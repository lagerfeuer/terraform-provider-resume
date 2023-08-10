terraform {
  required_providers {
    resume = {
      source = "deutz.io/provider/resume"
    }
  }
}

provider "resume" {
  endpoint = "http://localhost:3000"
  token    = "test"
}

data "resume_info" "this" {}

output "name" {
  value = data.resume_info.this.name
}

output "version" {
  value = data.resume_info.this.version
}

output "environment" {
  value = data.resume_info.this.environment
}
