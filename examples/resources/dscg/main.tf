// Test Sample - to configure XR/Cfg resources

terraform {
  required_providers {
    xrcm = {
      source = "infinera.com/poc/xrcm"
    }
  }
  required_version = "~> 1.3.4"
}

provider "xrcm" {
  username = "dev"
  password = "xrSysArch3"
  host     = "https://sv-kube-prd.infinera.com:443"
}

resource "xrcm_dscg" "dscg" {
  n         = "xr-regA_H1-L1"
  lineptpid    = 1
  carrierid = 1
  dscgid    = 20
  txcdscs = [1,2,5]
  rxcdscs = [2, 3, 5, 6]

}

output "dscg" {
  value = xrcm_dscg.dscg
}
