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
    sudo apt update && sudo apt install curl -y
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.4/install.sh | bash
    nvm install node
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

###########
### NX  ###
###########
if ! command -v nx &> /dev/null; then
    echo "Installing NX..."
    npm install -g nx
    npx create-nx-workspace my-org --preset=ts --no-nx-cloud --no-nx-cache --nx-cloud=false -y
else
    echo "NX is already installed. Skipping..."
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
    echo "Kind is already installed. Skipping..."
fi

################################
### install and setup KUBECTL ##
################################
if ! command -v kubectl &> /dev/null; then
    echo "Installing kubectl..."
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x kubectl
    sudo mv kubectl /usr/local/bin/
else
    echo "kubectl is already installed. Skipping..."
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

kind create cluster --config kind-config.yaml

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
    echo "Helm is already installed. Skipping..."
fi

#################
# pull kndp chart
#################
# pulling 0.1.6 for testing 
helm pull https://github.com/web-seven/kndp/releases/download/kndp-0.1.6/kndp-0.1.6.tgz

if [ -s "kndp-0.1.6.tgz" ]; then
    echo "Helm chart downloaded successfully, unpacking"
    tar -xzf kndp-0.1.6.tgz -C . && rm kndp-0.1.6.tgz
else
    echo "Helm chart is empty or download failed."
fi

helm install kndp kndp