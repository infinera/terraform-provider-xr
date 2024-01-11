variable modulecarriers {
  type = list(object({n = string, lineptpid = string, carrierids = optional(list(string))}))
  default = [ {n="xr-regA_H1-Hub", lineptpid="1", carrierids=["1"]},
  {n="xr-regA_H1-L1", lineptpid="1", carrierids=["1"]}]
}