find . -name '*.sql' | awk '{ print "source",$0 }' | mysql --batch -u root --password=123456 simking_dump
