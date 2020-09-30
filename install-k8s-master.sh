#!/bin/bash

function docker_install() {
    # (Install Docker CE)
    ## Set up the repository:
    ### Install packages to allow apt to use a repository over HTTPS
    apt-get update && apt-get install -y \
      apt-transport-https ca-certificates curl software-properties-common gnupg2


    # Add Docker’s official GPG key:
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -


    # Add the Docker apt repository:
    add-apt-repository \
      "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) \
      stable"

    # Install Docker CE
    apt-get update && apt-get install -y \
      containerd.io=1.2.13-2 \
      docker-ce=5:19.03.11~3-0~ubuntu-$(lsb_release -cs) \
      docker-ce-cli=5:19.03.11~3-0~ubuntu-$(lsb_release -cs)

    # Set up the Docker daemon
    cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

    mkdir -p /etc/systemd/system/docker.service.d

    # Restart Docker
    systemctl daemon-reload
    systemctl restart docker

    sudo systemctl enable docker


    }


function let_iptables_see_bridged_traffic() {
    # let iptables see bridged traffic

    cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
    sudo sysctl --system
}

function install_kubexxx() {
    sudo apt-get update && sudo apt-get install -y apt-transport-https curl
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
    cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
    sudo apt-get update
    sudo apt-get install -y kubelet kubeadm kubectl
    sudo apt-mark hold kubelet kubeadm kubectl
}

docker_install
let_iptables_see_bridged_traffic
install_kubexxx

echo "kubeXXX tools installed, now create cluster ..."
export KUBECONFIG=/etc/kubernetes/admin.conf
curl https://docs.projectcalico.org/manifests/calico.yaml -O
kubectl apply -f calico.yaml
mkdir -p $HOME/.kube
sudo chown -R $(id -u):$(id -g) $HOME/.kube
