
variable "vnet_addr_prefix" {}

variable "region" {}

variable "client_id" {}

variable "client_secret" {}
variable "project_shortcut" {}

variable "tags" {}

variable "tfm_name" {}

variable "shared_rg_name" {}

variable environment {}

variable "newbit_size" {
  description = "Map the friendly name to our subnet bit mask"
  type        = "map"

  default = {
    //8
    size_8 = "5"
    //16
    size_16 = "4"
    //32
    size_32  = "3"
    //64
    size_64 = "2"
    //128
    size_128  = "1"
    //256
    size_256  = "0"
  }
}


variable "segments_size" {
  description = "Map the friendly name to our subnet bit mask"
  type        = "map"

  default = {
    //8
    size_8 = "3"
    //16
    size_16 = "2"
    //32
    size_32  = "1"
    //64
    size_64 = "0"
  }
}



