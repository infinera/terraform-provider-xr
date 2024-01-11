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
data "xrcm_line_neighbors" "neighbors" {
  n             = "xr-regA_H2-Hub" #"XR-SFO_1-3"
  lineptpid      = 1
}

output "discovereddeighbors" {
  value =  data.xrcm_line_neighbors.neighbors.discoveredneighbors != null ? data.xrcm_line_neighbors.neighbors.discoveredneighbors[*] : []
}

output "controlplaneneighbors" {
  value =  data.xrcm_line_neighbors.neighbors.controlplaneneighbors != null ? data.xrcm_line_neighbors.neighbors.controlplaneneighbors[*] : []
}



