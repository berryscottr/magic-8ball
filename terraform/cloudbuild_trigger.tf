resource "google_cloudbuild_trigger" "discord_bot_trigger" {
  name     = "discord-bot-trigger"
  description = "Trigger to build and deploy Discord bot"
  project  = var.project_id
  filename = "cloudbuild.yaml"

  substitutions = {
    _CLOUD_RUN_SERVICE_NAME = google_cloud_run_service.discord_bot.name
    _IMAGE_URL = var.discord_bot_image_url
  }

  github {
    owner      = "berryscottr"
    name       = "magic-8ball"
    push {
      branch = "main"
    }
  }
}
