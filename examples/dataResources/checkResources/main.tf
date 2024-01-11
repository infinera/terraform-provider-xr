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

locals {
  network = {
      configs = { portspeed = ""
                  trafficmode = "L2Mode"
                  modulation = "" 
                }
      setup = {
        xr-regA_H1-Hub = {
          moduleconfig = { configuredrole = "hub", trafficmode ="L1Mode"}
          moduleclients = [{ clientid = "1", portspeed="100"}, { clientid = "2",portspeed="100"}]
          modulecarriers = [{ lineptpid = "1", carrierid = "1", modulation ="16QAM"}]
        }
        xr-regA_H1-L1 = {
          moduleconfig = { configuredrole = "leaf", trafficmode ="L1Mode"}
          moduleclients = [{ clientid = "1", portspeed="100"}]
          modulecarriers = [{ lineptpid = "1", carrierid = "1", modulation ="16QAM"} ]
        }
        xr-regA_H1-L2 = {
          moduleconfig = { configuredrole = "leaf", trafficmode ="L1Mode"}
          moduleclients = [{ clientid = "1", portspeed="100"}]
          modulecarriers = [{ lineptpid = "1", carrierid = "1", modulation ="16QAM"} ]
        }
        xr-regA_H1-L3 = {
          moduleconfig = { configuredrole = "leaf", trafficmode ="L1Mode"}
          moduleclients = [{ clientid = "1", portspeed="100"}]
          modulecarriers = [{ lineptpid = "1", carrierid = "1", modulation ="16QAM"} ]
        }
        xr-regA_H1-L4 = {
          moduleconfig = { configuredrole = "leaf", trafficmode ="L1Mode"}
          moduleclients = [{ clientid = "1", portspeed="100"}]
          modulecarriers = [{ lineptpid = "1", carrierid = "1", modulation ="16QAM"} ]
        }
        xr-regA_H2-L1 = {
          moduleconfig = { configuredrole = "leaf", trafficmode ="L1Mode"}
          moduleclients = [{ clientid = "1", portspeed="100"}]
          modulecarriers = [{ lineptpid = "1", carrierid = "1", modulation ="16QAM"} ]
        }
        xr-regA_H2-L2 = {
          moduleconfig = { configuredrole = "leaf", trafficmode ="L1Mode"}
          moduleclients = [{ clientid = "1", portspeed="100"}]
          modulecarriers = [{ lineptpid = "1", carrierid = "1", modulation ="16QAM"} ]
        }
        xr-regA_H2-L3 = {
          moduleconfig = { configuredrole = "leaf", trafficmode ="L1Mode"}
          moduleclients = [{ clientid = "1", portspeed="100"}]
          modulecarriers = [{ lineptpid = "1", carrierid = "1", modulation ="16QAM"} ]
        }
        xr-regA_H2-L4 = {
          moduleconfig = { configuredrole = "leaf", trafficmode ="L1Mode"}
          moduleclients = [{ clientid = "1", portspeed="100"}]
          modulecarriers = [{ lineptpid = "1", carrierid = "1", modulation ="16QAM"} ]
        }
      }
    }

}

data "xrcm_check_resources" "check_ethernets" {
  queries = [ for k,v in local.network.setup: { n = k, resourcetype = "Ethernet", resources = [ for client in v["moduleclients"]: {resourceid = client.clientid, attributevalues = [{ attribute = "portSpeed", intentvalue = client.portspeed, controlattribute = "portSpeedControl"}]} ] } ]
}

data "xrcm_check_resources" "check_configs" {
  queries = [ for k,v in local.network.setup: { n = k, resourcetype = "Config", resources = [{resourceid = k, attributevalues = [{ attribute = "trafficMode", intentvalue = v["moduleconfig"].trafficmode, controlattribute = "trafficModeControl"}]} ] } ]
}

data "xrcm_check_resources" "check_carriers" {
  queries = [ for k,v in local.network.setup: { n = k, resourcetype = "Carrier", resources = [ for carrier in v["modulecarriers"]: {resourceid = carrier.carrierid, parentid = carrier.lineptpid,attributevalues = [{ attribute = "modulation", intentvalue = carrier.modulation, controlattribute = "modulationControl"}]} ] } ]
}
/*
data "xrcm_check_resources" "check_ethernets" {
  queries = [ {n="xr-regA_H1-Hub", resourcetype = "Ethernet", resources = [{resourceid="1", attributevalues = [{ attribute = "portSpeed", intentvalue = "200", controlattribute = "portSpeedControl"}]},  {resourceid="2", attributevalues = [{ attribute = "portSpeed", intentvalue = "100", controlattribute = "portSpeedControl"}]}] } ]
}

data "xrcm_check_resources" "check_config" {
  queries = [ {n="xr-regA_H1-Hub", resourcetype = "Config", resources = [{resourceid="1", attributevalues = [{ attribute = "trafficMode", intentvalue = "L1Mode", controlattribute = "trafficModeControl"}]} ] } ]
}

data "xrcm_check_resources" "check_carriers" {
  queries = [ {n="xr-regA_H1-Hub", resourcetype = "Carrier"
      resources = [{resourceid="1", parentid="1", attributevalues = [{ attribute = "modulation", intentvalue = "16QAM", controlattribute = "modulationControl"}]} ] } ]
}

data "xrcm_check_resources" "check_resources" {
  queries = [ {n="xr-regA_H1-Hub", resourcetype = "Carrier", resources = [{resourceid="1", parentid="1", attributevalues = [{ attribute = "modulation", intentvalue = "16QAM", controlattribute = "modulationControl"}]}]},
  {n="xr-regA_H1-Hub", resourcetype = "Config", resources = [{resourceid="1", attributevalues = [{ attribute = "trafficMode", intentvalue = "L1Mode", controlattribute = "trafficModeControl"}]} ] },
  {n="xr-regA_H1-Hub", resourcetype = "Ethernet", resources = [{resourceid="1", attributevalues = [{ attribute = "portSpeed", intentvalue = "200", controlattribute = "portSpeedControl"}]},  {resourceid="2", attributevalues = [{ attribute = "portSpeed", intentvalue = "100", controlattribute = "portSpeedControl"}]}] }
   ]
}*/

output "check_ethernets" {
  value = data.xrcm_check_resources.check_ethernets
}

output "check_configs" {
  value = data.xrcm_check_resources.check_configs
}

output "check_carriers" {
  value = data.xrcm_check_resources.check_carriers
}
