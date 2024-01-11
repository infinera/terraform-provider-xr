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


resource "xrcm_carrier_diag" "carrier_diag" {
  n                      = "xr-regA_H1-L1"
  lineptpid                 = 1
  carrierid              = 1
  termlb               = "loopback"
}

output "carrier_diag" {
  value = xrcm_carrier_diag.carrier_diag
}
