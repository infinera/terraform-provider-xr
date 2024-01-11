// Test Sample - to configure XR/Cfg resources

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

resource "xrcm_dsc_diag" "dsc_diag" {
  n = "xr-regA_H1-L1"
  lineptpid = 1
  carrierid = 1
  dscid = 1
  facprbsgen = true
  facprbsmon = true
}

output "dsc_diag" {
  value = xrcm_dsc_diag.dsc_diag
}
