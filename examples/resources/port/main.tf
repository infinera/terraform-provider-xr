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

#provider "xrcm" {
#  username = "dev"
#  password = "xrSysArch3"
#  host     = "https://10.100.204.48:7443"
#}

resource "xrcm_port" "port1" {
  n                   = "xr-regA_H1-L1"
  portid = 1
}
output "port1" {
  value = xrcm_port.port1
}
