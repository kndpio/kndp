#!/bin/bash

# Deployment script to set up the necessary environment

#############
### GIT #####
#############
if ! command -v git &> /dev/null; then
    echo "Installing Git..."
    sudo apt-get update
    sudo apt-get install -y git
else
    echo "Git is already installed. Skipping..."
fi

############
## NODEJS ##
############
if ! command -v node &> /dev/null; then
    echo "Installing Node.js..."
    sudo apt install curl -y
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash - &&
    sudo apt-get install -y nodejs    
else
    echo "Node.js is already installed. Skipping..."
fi

############
## Docker ##
############
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    sudo apt-get update
    sudo apt-get install -y ca-certificates curl gnupg
    sudo install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg
    echo \
    "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
    "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    sudo usermod -aG docker $USER
    echo "Docker installed. You may need to log out and log back in to use Docker without sudo."
else
    echo "Docker is already installed. Skipping..."
fi

#######################################
### Install NX and create Workspace ###
#######################################

if ! command -v nx &> /dev/null; then
    echo "Installing NX..."    
    npx create-nx-workspace --skipGit ./
    npm install -g nx 
else
    echo "NX is already installed."
fi

###########################
#  install and setup KIND #
###########################
if ! command -v kind &> /dev/null; then
    echo "Installing kind..."
    curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
    chmod +x ./kind
    sudo mv ./kind /usr/local/bin/kind
else
    echo "KNDP basic tools are installed."
fi

##########################################
# create a Kubernetes cluster using kind # 
##########################################
echo "Creating Kubernetes cluster using kind..."
cat <<EOF > kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: worker
  extraMounts:
  - hostPath: ./
    containerPath: /storage

- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF

sudo kind create cluster --config kind-config.yaml

##########
## HELM ##
##########
if ! command -v helm &> /dev/null; then
    echo "Installing Helm..."
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
    chmod 700 get_helm.sh
    ./get_helm.sh
    rm get_helm.sh  # Clean up the installation script
else
    echo "Installed. Skipping..."
fi

###################
# pull KNDP CHART #
###################

helm repo add kndp https://kndp.io
helm repo up
helm install kndp kndp/kndp


echo "Installation succeded! Kind reminder to to log out and log back in to apply changes and run Docker without 'sudo' "