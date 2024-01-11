variable moduleethernets {
  type = list(object({n = string, ethernetids = optional(list(string)),}))
  default = [ {n="xr-regA_H1-Hub", ethernetids=["1"]},
  {n="xr-regA_H1-L1", ethernetids=["1","2"]} ]
}