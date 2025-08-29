resource "google_project_service" "apis" {
  for_each = toset([
    "cloudfunctions.googleapis.com",
    "cloudbuild.googleapis.com",
    "run.googleapis.com",
    "pubsub.googleapis.com",
    "bigquery.googleapis.com",
    "artifactregistry.googleapis.com"
  ])
  service                    = each.key
  disable_dependent_services = true
}

resource "google_pubsub_topic" "webhooks_topic" {
  name    = "webhooks-topic"
  project = var.gcp_project_id
  depends_on = [google_project_service.apis]
}

resource "google_bigquery_dataset" "webhooks_dataset" {
  dataset_id = "webhooks_dataset"
  project    = var.gcp_project_id
  location   = var.gcp_region
  depends_on = [google_project_service.apis]
}

resource "google_bigquery_table" "webhooks_table" {
  dataset_id = google_bigquery_dataset.webhooks_dataset.dataset_id
  table_id   = "events"
  project    = var.gcp_project_id
  schema = jsonencode([
    {
      "name" : "event_type",
      "type" : "STRING",
      "mode" : "NULLABLE"
    },
    {
      "name" : "event_timestamp",
      "type" : "TIMESTAMP",
      "mode" : "NULLABLE"
    },
    {
      "name" : "raw_data",
      "type" : "JSON",
      "mode" : "NULLABLE"
    }
  ])
}

data "archive_file" "source" {
  type        = "zip"
  source_dir  = var.function_source_code_path
  output_path = "/tmp/function-source.zip"
  excludes = [
    ".git",
    ".gitignore",
    "terraform"
  ]
}

resource "google_storage_bucket" "source_bucket" {
  name     = "${var.gcp_project_id}-functions-source"
  location = var.gcp_region
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "source_object" {
  name   = "source.zip"
  bucket = google_storage_bucket.source_bucket.name
  source = data.archive_file.source.output_path
}

resource "google_cloudfunctions2_function" "receiver_http" {
  name     = "go-receiver-http"
  location = var.gcp_region
  project  = var.gcp_project_id

  build_config {
    runtime     = "go122"
    entry_point = "ReceiverFunction"
    source {
      storage_source {
        bucket = google_storage_bucket.source_bucket.name
        object = google_storage_bucket_object.source_object.name
      }
    }
  }

  service_config {
    max_instance_count = 1
    min_instance_count = 0
    available_memory   = "256Mi"
    timeout_seconds    = 60
    environment_variables = {
      GCP_PROJECT_ID  = var.gcp_project_id
      PUBSUB_TOPIC_ID = google_pubsub_topic.webhooks_topic.name
    }
    ingress_settings               = "ALLOW_ALL"
    all_traffic_on_latest_revision = true
  }

  depends_on = [google_project_service.apis]
}

resource "google_cloud_run_service_iam_member" "invoker" {
  location = google_cloudfunctions2_function.receiver_http.location
  service  = google_cloudfunctions2_function.receiver_http.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloudfunctions2_function" "processor_pubsub" {
  name     = "go-processor-pubsub"
  location = var.gcp_region
  project  = var.gcp_project_id

  build_config {
    runtime     = "go122"
    entry_point = "ProcessorFunction"
    source {
      storage_source {
        bucket = google_storage_bucket.source_bucket.name
        object = google_storage_bucket_object.source_object.name
      }
    }
  }

  service_config {
    max_instance_count = 5
    min_instance_count = 0
    available_memory   = "256Mi"
    timeout_seconds    = 300
    environment_variables = {
      GCP_PROJECT_ID      = var.gcp_project_id
      BIGQUERY_DATASET_ID = google_bigquery_dataset.webhooks_dataset.dataset_id
      BIGQUERY_TABLE_ID   = google_bigquery_table.webhooks_table.table_id
    }
  }

  event_trigger {
    trigger_region = var.gcp_region
    event_type     = "google.cloud.pubsub.topic.v1.messagePublished"
    pubsub_topic   = google_pubsub_topic.webhooks_topic.id
    retry_policy   = "RETRY_POLICY_RETRY"
  }

  depends_on = [google_project_service.apis]
}