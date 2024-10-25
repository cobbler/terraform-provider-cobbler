terraform {
  required_providers {
    cobbler = {
      source = "cobbler/cobbler"
      version = "4.1.0"
    }
  }
}

variable "username" {
  type = string
}

variable "password" {
  type = string
}

variable "url" {
  type = string
}

variable "insecure" {
  type = bool
}

provider "cobbler" {
  username = var.username # optionally use COBBLER_USERNAME env var
  password = var.password # optionally use COBBLER_PASSWORD env var
  url  = var.url          # optionally use COBBLER_URL env var

  # You may need to allow insecure TLS communications unless you
  # have configured certificates
  insecure = var.insecure # optionally use COBBLER_INSECURE env var
}
