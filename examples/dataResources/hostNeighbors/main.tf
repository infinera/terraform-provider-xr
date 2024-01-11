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
data "xrcm_host_neighbors" "neighbors" {
  n             = "xr-regA_H2-Hub" #"XR-SFO_1-3"
  ethernetid      = 1
}

output "xrcm_neighbors" {
  value =  data.xrcm_host_neighbors.neighbors.neighbors != null ? data.xrcm_host_neighbors.neighbors.neighbors[*] : []
}



