terraform {
  required_providers {
    resume = {
      source = "deutz.io/provider/resume"
    }
  }
}

provider "resume" {
  endpoint = "https://localhost:3000"
}
