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


resource "xrcm_carrier" "carrier" {
  #n            = "xr-regA_H1-L1"
  n            = "xr-regA_H2-Hub"
  lineptpid              = 1
  carrierid              = 1
  clientportmode         = "ethernet"
  modulation             = "QPSK"
  //modulationcontrol      = "IPM"
  constellationfrequency = 0
  //modulation = "16QAM"
  //capacity = 100
  //txpowertargetperdsc = -5
}

output "xrcm_carrier" {
  value = xrcm_carrier.carrier
}
