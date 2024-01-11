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

data "xrcm_carriers" "carriers" {
  modulecarriers = [ {n="xr-regA_H1-Hub", lineptpid="1", carrierids=["1"]},
  {n="xr-regA_H1-L1", lineptpid="1", carrierids=["1"]}]
}
output "carriers" {
  value = data.xrcm_carriers.carriers
}