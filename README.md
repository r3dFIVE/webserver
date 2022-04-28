# webserver
Barebones Golang webserver for my personal "website"

## Port Forwarding

This webserver is meant to be run by non-root users thus will not bind to ports < 1000.

By default the webserver will use ports 8080 for HTTP and 8443 for HTTPS, if enabled.

You can simply forward traffic from 80/443 to their respective destinations.

Eg.

```
# HTTP
sudo iptables -t nat -I PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8080
sudo iptables -I INPUT -p tcp --dport 8080 -j ACCEPT

# HTTPS
sudo iptables -t nat -I PREROUTING -p tcp --dport 443 -j REDIRECT --to-port 8443
sudo iptables -I INPUT -p tcp --dport 8443 -j ACCEPT
```

