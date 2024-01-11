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

resource "xrcm_cfg" "cfg" {
  n             = "xr-regA_H2-L1"
  trafficmode   = "L1Mode"
  tcmode        = true
  //trafficmodecontrol = "IPM"

}
output "cfg" {
  value = xrcm_cfg.cfg
}
