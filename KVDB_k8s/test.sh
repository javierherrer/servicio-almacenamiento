#/bin/sh
#POD=$(echo 'adsafsdf-2' | grep -oE '[^0-9]+([0-9])+$')
#echo $POD
POD=$(echo 'adsafsdf-2' | cut -d '-' -f 2)
echo $POD
if [ $? -eq 0 ]; then
    echo "funsiona"
else
    echo "no funca"
fi


#!/bin/sh                                                                                                    
MAX_POD_NUMBER=10                                                                        
MINOMBREPOD="holads-1"                                                                   
POD_NUMBER=$(echo $MINOMBREPOD | cut -d '-' -f 2)                                        
POD_NUMBER=1                                                                             
                                                                                         
for i in $(seq 0 $MAX_POD_NUMBER); do                                                    
        IP_LEADER=$(echo 'RAFT LEADER' |  nc kvdbrep-$i.kvdbservice.default.svc.cluster.l
        echo $IP_LEADER | grep -qE '[0-9]+.*'                                            
        if [ $? -eq 0  ] ; then                                                          
                                                                                         
        fi                                                                               
done                                 






#!/bin/sh
MAX_POD_NUMBER=10
for i in $(seq 0 $MAX_POD_NUMBER); do            
        IP_LEADER=$(echo 'RAFT LEADER' |  nc kvdbrep-$i.kvdbservice.default.svc.cluster.l
        echo $IP_LEADER | grep -qE '[0-9]+.*'                                            
        if [ $? -eq 0  ] ; then                                                          
                echo "es ip"                                                             
                /kvdb -n $(MINOMBREPOD) -a $(MINOMBREPOD).$(MISUBDOMINIO):11001 -j kvdbre
                break                                                                    
        fi                                                                               
done  
~          