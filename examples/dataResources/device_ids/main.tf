terraform {
  required_providers {
    xrcm = {
      source = "infinera.com/poc/xrcm"
    }
  }
}

provider "xrcm" {
  username = "dev"
  password = "xrSysArch3"
  host     = "https://sv-kube-prd.infinera.com:443"
}


data "xrcm_devices_ids" "devices_ids" {
  deviceids = [{n="xr-regA_H1-L2"}, {n="xr-regA_H1-L1"}]
}

output "xrcm_devices_ids" {
  value = data.xrcm_devices_ids.devices_ids
}