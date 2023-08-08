terraform {
  required_providers {
    resume = {
      source = "deutz.io/provider/resume"
    }
  }
}

provider "resume" {
  endpoint = "http://localhost:3000"
  token = "test"
}

resource "resume_resume" "this" {
  name = "Michael G Scott"
  address = "1 Paper St"
  phone_number = "555-555-5555"
  website = "https://michaelthesco.tt"
}

output "id" {
  value = resume_resume.this.id
}