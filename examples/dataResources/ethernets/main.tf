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

data "xrcm_ethernets" "ethernets" {
  moduleethernets = [ {n="xr-regA_H1-Hub", ethernetids=["1", "3"]},
  {n="xr-regA_H1-L1", ethernetids=["1","2"]} ]
}

output "ethernets" {
  value = data.xrcm_ethernets.ethernets
}