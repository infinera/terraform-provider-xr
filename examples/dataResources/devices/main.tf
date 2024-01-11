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


data "xrcm_devices" "devices" {
  names = ["xr-regA_H1-L2", "xr-regA_H1-L1"] 
  state = "ONLINE"
}

output "xrcm_devices_devices" {
  value = data.xrcm_devices.devices
}