resource "google_secret_manager_secret" "discord_bot_token" {
  secret_id = "magic-8ball-token"
  replication {
    auto {}
  }
}

data "google_secret_manager_secret_version" "discord_bot_token_latest" {
  secret = google_secret_manager_secret.discord_bot_token.id
  version   = "latest"
}

