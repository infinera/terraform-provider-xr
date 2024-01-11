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

// Note: id is computed field added as part of read/create
resource "xrcm_otu" "otu" {
  n = "xr-regA_H1-L1"
  otuid = "1"
  txtti = "testing123"
  expectedtti = "123testing"
}

output "otu" {
  value = xrcm_otu.otu
}
