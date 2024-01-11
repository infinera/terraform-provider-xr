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

data "xrcm_acs" "acs" {
  moduleacs = [ {n="xr-regA_H1-Hub", ethernetid="1"},
  {n="xr-regA_H1-L1", ethernetid="1"}]
}

output "acs" {
  value = data.xrcm_acs.acs
}