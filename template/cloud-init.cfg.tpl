#cloud-config
ssh_pwauth: true
disable_root: false
chpasswd:
  list: |
    root:password
  expire: false

ssh_authorized_keys:
%{ for key in sshKeys ~}
  - ${key}
%{ endfor ~}

runcmd:
  - [ sed, -i, -e, '/^#PermitRootLogin .*\$/a PermitRootLogin yes', /etc/ssh/sshd_config ]
