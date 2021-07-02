#/bin/sh
MAX_POD_NUMBER=10
#POD_NUMBER=$(echo $(MINOMBREPOD) | cut -d '-' -f 2)
POD_NUMBER=1
for i in 'seq 0 $MAX_POD_NUMBER'
do  
     $(echo 'RAFT LEADER' | 
     nc "kvdbrep-$i.kvdbservice.default.svc.cluster.local" 11001)
done
