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

data "xrcm_lcs" "lcs" {
  n = "xr-regA_H1-Hub"
  #lcids = ["1", "2"]
}

output "lcs" {
  value = data.xrcm_lcs.lcs
}