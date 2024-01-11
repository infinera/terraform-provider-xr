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

/*provider "xrcm" {
  username = "dev"
  password = "xrSysArch3"
  host     = "https://10.100.204.48:7443"
}*/

// Note: id is computed field added as part of read/create
resource "xrcm_ethernet_diag" "ethernet_diag" {
  n            = "xr-regA_H1-L1"
  #n            = "XR LEAF 1"
  ethernetid   = 1
  termlb = "disabled"
  faclb = "disabled"
}

output "ethernet_diag" {
  value = xrcm_ethernet_diag.ethernet_diag
}

