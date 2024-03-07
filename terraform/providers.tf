provider "google" {
  project = var.project_id
  region  = var.region
}

terraform {
  required_version = ">=1.0.5"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.50.0"
    }
  }
}
