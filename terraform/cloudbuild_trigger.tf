# google_cloudbuild_trigger.magic-8ball:
resource "google_cloudbuild_trigger" "magic-8ball" {
  count          = var.project_id == "magic-8ball-416501" ? 1 : 0
  description    = "Build and deploy to Cloud Run service magic-8ball on push to \"^main$\""
  disabled       = false
  ignored_files  = []
  included_files = []
  name           = "magic-8ball"
  project        = var.project_id
  substitutions = {
    "_DEPLOY_REGION" = var.region
    "_GCR_HOSTNAME"  = "us.gcr.io"
    "_PLATFORM"      = "managed"
    "_SERVICE_NAME"  = "magic-8ball"
  }
  tags = [
    "gcp-cloud-build-deploy-cloud-run",
    "gcp-cloud-build-deploy-cloud-run-managed",
    "magic-8ball",
  ]

  build {
    images = [
      "$_GCR_HOSTNAME/$PROJECT_ID/$REPO_NAME/$_SERVICE_NAME:$COMMIT_SHA",
    ]
    substitutions = {
      "_DEPLOY_REGION" = var.region
      "_GCR_HOSTNAME"  = "us.gcr.io"
      "_PLATFORM"      = "managed"
      "_SERVICE_NAME"  = "magic-8ball"
    }
    tags = [
      "gcp-cloud-build-deploy-cloud-run",
      "gcp-cloud-build-deploy-cloud-run-managed",
      "magic-8ball",
    ]

    options {
      disk_size_gb           = 0
      dynamic_substitutions  = false
      env                    = []
      secret_env             = []
      source_provenance_hash = []
      substitution_option    = "ALLOW_LOOSE"
    }

    step {
      args = [
        "build",
        "--no-cache",
        "-t",
        "$_GCR_HOSTNAME/$PROJECT_ID/$REPO_NAME/$_SERVICE_NAME:$COMMIT_SHA",
        ".",
        "-f",
        "Dockerfile",
      ]
      env        = []
      id         = "Build"
      name       = "gcr.io/cloud-builders/docker"
      secret_env = []
      wait_for   = []
      timeout    = "1200s"
    }
    step {
      args = [
        "tag",
        "$_GCR_HOSTNAME/$PROJECT_ID/$REPO_NAME/$_SERVICE_NAME:$COMMIT_SHA",
        "$_GCR_HOSTNAME/$PROJECT_ID/$REPO_NAME/$_SERVICE_NAME:latest",
      ]
      env        = []
      id         = "Tag"
      name       = "gcr.io/cloud-builders/docker"
      secret_env = []
      wait_for   = []
    }
    step {
      args = [
        "push",
        "$_GCR_HOSTNAME/$PROJECT_ID/$REPO_NAME/$_SERVICE_NAME:$COMMIT_SHA",
      ]
      env        = []
      id         = "Push COMMIT_SHA"
      name       = "gcr.io/cloud-builders/docker"
      secret_env = []
      wait_for   = ["Tag"]
    }
    step {
      args = [
        "push",
        "$_GCR_HOSTNAME/$PROJECT_ID/$REPO_NAME/$_SERVICE_NAME:latest",
      ]
      env        = []
      id         = "Push latest tag"
      name       = "gcr.io/cloud-builders/docker"
      secret_env = []
      wait_for   = ["Tag"]
    }
    step {
      args = [
        "run",
        "deploy",
        "$_SERVICE_NAME",
        "--port=8080",
        "--cpu=1000m",
        "--memory=512Mi",
        "--platform=managed",
        "--set-secrets=BOT_TOKEN=magic-8ball:latest",
        "--image=$_GCR_HOSTNAME/$PROJECT_ID/$REPO_NAME/$_SERVICE_NAME:$COMMIT_SHA",
        "--labels=managed-by=gcp-cloud-build-deploy-cloud-run,commit-sha=$COMMIT_SHA",
        "--region=$_DEPLOY_REGION",
        "--quiet",
        "--min-instances=0",
      ]
      entrypoint = "gcloud"
      env        = []
      id         = "Deploy"
      name       = "gcr.io/google.com/cloudsdktool/cloud-sdk:slim"
      secret_env = []
      wait_for   = ["Push COMMIT_SHA"]
    }
    timeout = "1800s"
  }

  github {
    name  = "magic-8ball"
    owner = "berryscottr"

    push {
      branch       = "^main$"
      invert_regex = false
    }
  }

  timeouts {}
}