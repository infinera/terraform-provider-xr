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

resource "xrcm_lc" "lc" {
  n         = "xr-regA_H1-L1"
  clientaid  = 1
  dscgaid    = 1
  lineptpid = 1
  carrierid = 1

}

output "xrcm_lc" {
  value = xrcm_lc.lc
}
