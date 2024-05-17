#!/bin/bash
apt update
apt install ufw

# Allow SSH
ufw allow ssh

# Allow connections on 443 from cloudflare
for ip in $(curl -s https://www.cloudflare.com/ips-v4); 
do 
  ufw allow from $ip to any port 443 proto tcp; 
done

# Deny all other incoming traffic
ufw default deny incoming

# Enable the firewall
ufw enable
