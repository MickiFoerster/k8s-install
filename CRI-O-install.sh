#!/bin/sh

modprobe overlay
modprobe br_netfilter

# Set up required sysctl params, these persist across reboots.
cat > /etc/sysctl.d/99-kubernetes-cri.conf <<EOF
net.bridge.bridge-nf-call-iptables  = 1
net.ipv4.ip_forward                 = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

sysctl --system

echo "Choose your OS:"
echo "1. Debian, Raspian"
echo "2. Ubuntu"
read a 
case $a in
    1)
        echo 'deb http://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/Raspbian_10/ /' > /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
        wget -nv https://download.opensuse.org/repositories/devel:kubic:libcontainers:stable/Raspbian_10/Release.key -O- | sudo apt-key add -
        ;;
    2)
        # Configure package repository
        . /etc/os-release
        sudo sh -c "echo 'deb http://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/x${NAME}_${VERSION_ID}/ /' > /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list"
        wget -nv https://download.opensuse.org/repositories/devel:kubic:libcontainers:stable/x${NAME}_${VERSION_ID}/Release.key -O- | sudo apt-key add -
        ;;
    *)
        echo "Wrong answer, so quit installation"
        exit 1 
        ;;
esac

sudo apt-get update
# Install CRI-O
sudo apt-get install -y cri-o-1.17

systemctl daemon-reload
systemctl start crio

