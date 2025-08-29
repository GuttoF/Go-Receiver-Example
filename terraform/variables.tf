variable "gcp_project_id" {
  description = "O ID do seu projeto no Google Cloud."
  type        = string
}

variable "gcp_region" {
  description = "A região onde os recursos serão criados."
  type        = string
  default     = "southamerica-east1"
}

variable "function_source_code_path" {
  description = "Caminho para o código-fonte das funções."
  type        = string
  default     = ".."
}