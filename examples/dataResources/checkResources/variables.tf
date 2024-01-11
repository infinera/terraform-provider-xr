variable queries {
  type = list(object({n = string, resources = list(object({resourceid = string, attributevalues = list(object({attribute = optional(string), attributevalue = optional(string), controlattribute = optional(string), attributecontrolbyhost =  optional(bool), isvaluematch = optional(bool)}))}))}))
  default = [ {n="xr-regA_H1-Hub", 
      resources = [{resourceid="1", attributevalues = [{ attribute = "portSpeed", attributevalue = "200", controlattribute = "portSpeedControl", "attributecontrolbyhost" = false, "isvaluematch" = true}]}, 
                  {resourceid="2", attributevalues = [{ attribute = "portSpeed", attributevalue = "200", controlattribute = "portSpeedControl","controlbyhost" = false, "ismatch" = true}]}] } ]
}

variable "resourcetype" {
  type = string
  default = "Etherner" // Ethernet, DSC, DSCG, Carrier, Config, LC
}