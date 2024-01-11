trafficmode = "L1Mode"

hub_names = ["xr-regA_H1-Hub"]
leaf_names = [
  "xr-regA_H1-L1",
  "xr-regA_H1-L2",
  "xr-regA_H1-L3",
  "xr-regA_H1-L4"
]

hub_bandwidth = {
  xr-regA_H1-Hub = {
    xr-regA_H1-Hub-BW1 = { hubdscgid = "1", leafdscgid = "1", hubdscidlist = ["5", "1", "7", "3"], leafdscidlist = ["1", "2", "3", "4"], direction = "ds" },
    xr-regA_H1-Hub-BW2 = { hubdscgid = "2", leafdscgid = "1", hubdscidlist = ["9", "11", "13", "15"], leafdscidlist = ["1", "2", "3", "4"], direction = "ds" },
    xr-regA_H1-Hub-BW3 = { hubdscgid = "3", leafdscgid = "1", hubdscidlist = ["2", "4", "6", "8"], leafdscidlist = ["1", "2", "3", "4"], direction = "ds" },
    xr-regA_H1-Hub-BW4 = { hubdscgid = "4", leafdscgid = "1", hubdscidlist = ["14", "10", "12", "16"], leafdscidlist = ["1", "2", "3", "4"], direction = "ds" }
  }
}

leaf_bandwidth = {
  xr-regA_H1-L1 = {
    xr-regA_H1-Hub-BW1 = { hubdscgid = "1", leafdscgid = "1", hubdscidlist = ["5", "1", "7", "3"], leafdscidlist = ["1", "2", "3", "4"], direction = "us" }
  }
  xr-regA_H1-L2 = {
    xr-regA_H1-Hub-BW2 = { hubdscgid = "2", leafdscgid = "1", hubdscidlist = ["9", "11", "13", "15"], leafdscidlist = ["1", "2", "3", "4"], direction = "us" }
  }
  xr-regA_H1-L3 = {
    xr-regA_H1-Hub-BW3 = { hubdscgid = "3", leafdscgid = "1", hubdscidlist = ["2", "4", "6", "8"], leafdscidlist = ["1", "2", "3", "4"], direction = "us" }
  }
  xr-regA_H1-L4 = {
    xr-regA_H1-Hub-BW4 = { hubdscgid = "4", leafdscgid = "1", hubdscidlist = ["14", "10", "12", "16"], leafdscidlist = ["1", "2", "3", "4"], direction = "us" }
  }
}

client-2-dscg = {
  xr-regA_H1-Hub = {
    lc-xr-regA_H1-L1 = {
      clientid = "1"
      dscgid   = "1"
    },
    lc-xr-regA_H1-L2 = {
      clientid = "2"
      dscgid   = "2"
    },
    lc-xr-regA_H1-L3 = {
      clientid = "3"
      dscgid   = "3"
    },
    lc-xr-regA_H1-L4 = {
      clientid = "4"
      dscgid   = "4"
  } },
  xr-regA_H1-L1 = {
    lc-xr-regA_H1-Hub-1 = {
      clientid = "1"
      dscgid   = "1"
  } },
  xr-regA_H1-L2 = {
    lc-xr-regA_H1-Hub-2 = {
      clientid = "1"
      dscgid   = "1"
  } },
  xr-regA_H1-L3 = {
    lc-xr-regA_H1-Hub-3 = {
      clientid = "1"
      dscgid   = "1"
  } },
  xr-regA_H1-L4 = {
    lc-XR-SFO_4-Hub-3 = {
      clientid = "1"
      dscgid   = "1"
  } }
}

