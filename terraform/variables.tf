variable "project_id" {
  type        = string
  default     = "magic-8ball-416501"
  description = "Google Cloud Project ID"
}

variable "region" {
  type        = string
  default     = "us-central1"
  description = "Google Cloud Region"
}

variable "discord_bot_image_url" {
  type        = string
  default     = "gcr.io/magic-8ball-416501/discord-bot:latest"
  description = "URL of the Docker image for the Discord bot"
}
