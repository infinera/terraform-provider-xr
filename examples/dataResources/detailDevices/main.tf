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

data "xrcm_detaildevices" "onlinedevices" {
  names = ["xr-regA_H1-L2", "xr-regA_H1-L3", "xr-regA_H1-L4"]
  //state = "OFFLINE"
}

output "devices" {
  value = data.xrcm_detaildevices.onlinedevices
}