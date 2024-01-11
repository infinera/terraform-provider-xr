

**Sumilator Machine**
Simulator machine - sv-xrsim1-prd.infinera.com cred - sim/xrsim
Simulator directory

cd ~/marvel/xrsim

Device inventory file
~/inv/startup.cfg.src


Simulator specification
dockercompose.yaml
- serice name x replicas - each service name is the hub name and leaf is service name_1-1
replicas = 4 


Starting/Stoping simulator
start 
stop 

**plgd machine**
Kubenetes - sv-kube-prd.infinera.com - cred - dev/xrKube

Restarting Plgd
refresh 

status
kubectl get po -n am-1

Madhav's tutorial
https://infinera-my.sharepoint.com/:v:/p/mkothapalli/EWVIWrhYZX9FhLgDT26RNGEBVOAkRD8gskA5Yt4p_vBZIQ 