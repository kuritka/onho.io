locals {
  suffix = "${var.project_shortcut}-${var.environment_short}"

  shared_rg_name = "rg-${local.suffix}"

  //fs storage is automatically replicated to paired region
  fs_name = "st${var.project_shortcut}${var.environment_short}${var.primary_region}"

}
