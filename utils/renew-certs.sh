#!/bin/bash
certbot renew --quiet

cp /etc/letsencrypt/live/bengholmes.com/fullchain.pem /home/bengholmes/website/certs
cp /etc/letsencrypt/live/bengholmes.com/privkey.pem /home/bengholmes/website/certs

chown bengholmes /home/bengholmes/website/certs/*

systemctl restart website
