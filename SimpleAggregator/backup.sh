cd $TELEPORT_HOME/SimpleAggregator/logs
tar cvzf teleport-nslog-`date +%F`.tgz *.log  
rm -f *.log
