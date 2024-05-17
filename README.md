# Website
This is everything that runs [bengholmes.com](https://bengholmes.com).

# Setup
To set up a new server, do the following:

## Raspberry Pi
Create a user `bengholmes` which will own all the website related stuff.

## Set up a firewall
```bash
sudo setup-firewall.sh
```

## Create a service
```bash
sudo cp systemd/website.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start website
sudo systemctl enable website
```

## Register a cron job for renewing certs
```bash
sudo crontab -e
```
Add the following
```bash
0 0 1 * * /home/bengholmes/website/utils/renew-certs.sh
```
