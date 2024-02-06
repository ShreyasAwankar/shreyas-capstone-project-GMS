provider "google" {
  credentials = file("capstone-key.json")
  project     = "capstone-412502"
  region      = "us-central1"
}

locals {
  bucket_name = "shreyas-cloud-functions"
}

resource "google_storage_bucket_object" "function_zip" {
  for_each = var.functions
  name     = each.key
  bucket   = local.bucket_name
  source   = each.value.zip
}

resource "google_cloudfunctions2_function" "function" {
  for_each = var.functions
  name     = each.value.name
  location = var.region

  build_config {
    runtime     = each.value.runtime
    entry_point = each.value.entrypoint
    environment_variables = {
      JWT_KEY = "my_secret_key"
    }

    source {
      storage_source {
        bucket = local.bucket_name
        object = google_storage_bucket_object.function_zip[each.key].name
      }
    }
  }
  service_config {
    min_instance_count             = 1
    available_memory               = "128Mi"
    timeout_seconds                = 120
    all_traffic_on_latest_revision = true
    service_account_email          = "742350885173-compute@developer.gserviceaccount.com"
  }
}

resource "google_cloud_run_service_iam_member" "member" {
  for_each = var.functions

  location = google_cloudfunctions2_function.function[each.key].location
  service  = each.key
  role     = "roles/run.invoker"
  member   = "allUsers"
}
