network:
  version: 2
  ethernets:
    ens3:
      dhcp4: false
      addresses:
        - ${ip}/24
      routes:
        - to: default
          via: 192.168.122.1
      nameservers:
        addresses: [2.189.44.44, 8.8.8.8]
