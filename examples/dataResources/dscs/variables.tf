variable moduledscs {
  type = list(object({n = string, lineptpid = string, carrierid = string, dscids = optional(list(string))}))
  default = [ {n="xr-regA_H1-Hub", lineptpid="1", carrierid="1"},
  {n="xr-regA_H1-L1", lineptpid="1", carrierid="1"}]
}