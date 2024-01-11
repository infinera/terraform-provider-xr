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
data "xrcm_lcs" "lcsdata" {
  n             = "xr-regA_H2-Hub" #"XR-SFO_1-3"
}

output "xrcm_lcs" {
  value =  data.xrcm_lcs.lcsdata.lcs != null ? data.xrcm_lcs.lcsdata.lcs[*] : []
}



