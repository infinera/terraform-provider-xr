variable moduleacs {
  type = list(object({n = string, ethernetid = string, acids = optional(list(string)),}))
  default = [ {n="xr-regA_H1-Hub", ethernetid="1"},
  {n="xr-regA_H1-L1", ethernetid="1"}]
}