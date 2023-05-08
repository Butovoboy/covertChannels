#!/bin/bash

echo 1 > /proc/sys/net/ipv4/ip_forward
iptables -A OUTPUT -p icmp --icmp-type any -m limit --limit 3/minute --limit-burst 3 -j ACCEPT
iptables -A OUTPUT -p icmp -j DROP
iptables -S
