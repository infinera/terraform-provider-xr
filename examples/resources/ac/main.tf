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


// Note: id is computed field added as part of read/create
resource "xrcm_ac" "ac1" {
  n             = "xr-regA_H1-L1" #"XR-SFO_1-3"
  ethernetid      = 1
  acid          = 3
  capacity      = 1
  imc           = "MatchAll"
  imc_outer_vid = ""
  emc           = "MatchAll"
  emc_outer_vid = ""
}

output "xrcm_ac1" {
  value = xrcm_ac.ac1
}



