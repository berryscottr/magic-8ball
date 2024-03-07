resource "google_cloud_run_service" "discord_bot" {
  name     = "magic-8ball"
  location = var.region

  template {
    spec {
      containers {
        image = var.discord_bot_image_url
        ports {
          container_port = 8080
        }
        env {
          name  = "DISCORD_TOKEN"
          value = data.google_secret_manager_secret_version.discord_bot_token_latest.secret_data
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

