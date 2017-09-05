#!/bin/bash
cd /etc/coredns/
echo '$ORIGIN x.benjojo.co.uk.' > /tmp/a
echo '@       3600 IN SOA sns.dns.icann.org. noc.dns.icann.org. ('  >> /tmp/a
TS=$(date +%s)
echo "                                $TS ; serial" >> /tmp/a
echo "                                7200       ; refresh (2 hours)" >> /tmp/a
echo "                                3600       ; retry (1 hour)" >> /tmp/a
echo "                                1209600    ; expire (2 weeks)" >> /tmp/a
echo "                                3600       ; minimum (1 hour)" >> /tmp/a
echo "                                )" >> /tmp/a

echo "    3600 IN NS rdns1.benjojo.co.uk." >> /tmp/a
echo "       3600 IN NS rdns2.benjojo.co.uk." >> /tmp/a

shuf x.benjojo.co.uk.sources | head -n 6 >> /tmp/a
cp /tmp/a /etc/coredns/x.benjojo.co.uk
