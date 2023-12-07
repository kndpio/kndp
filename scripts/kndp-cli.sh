#!/bin/bash

function help() {
    echo "Usage: kndp [options]"
    echo ""
    echo "Options:"
    echo "   help, -h             For more information about existing commands"
    echo "   install, -i          Install KNDP"
    echo "   uninstall, -u        Uninstall KNDP"
    echo ""
}

function install_kndp() {
    #############
    ### GIT #####
    #############
    if ! command -v git &>/dev/null; then
        echo "Installing Git..."
        sudo apt-get update
        sudo apt-get install -y git
    else
        echo "Git is already installed. Skipping..."
    fi

    ############
    ## NODEJS ##
    ############
    if ! command -v node &>/dev/null; then
        echo "Installing Node.js..."
        sudo apt-get update
        sudo apt-get install -y ca-certificates curl gnupg
        sudo mkdir -p /etc/apt/keyrings
        curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
        NODE_MAJOR=20
        echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
        sudo apt-get update
        sudo apt-get install nodejs -y
    else
        echo "Node.js is already installed. Skipping..."
    fi

    ############
    ## Docker ##
    ############
    if ! command -v docker &>/dev/null; then
        echo "Installing Docker..."
        sudo apt-get update
        sudo apt-get install -y ca-certificates curl gnupg
        sudo install -m 0755 -d /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        sudo chmod a+r /etc/apt/keyrings/docker.gpg
        echo \
            "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
        "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" |
            sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
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
    if ! command -v nx &>/dev/null; then
        echo "Installing NX..."
        npx create-nx-workspace --skipGit ./
        npm install -g nx
    else
        echo "NX is already installed."
    fi

    ###########################
    #  install and setup KIND #
    ###########################
    if ! command -v kind &>/dev/null; then
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
    cat <<EOF >kind-config.yaml
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

    if kind get clusters | grep "kndp"; then
        echo "Cluster 'kndp' already exists. Skipping cluster creation."
    else
        kind create cluster --name kndp --config kind-config.yaml
    fi

    ##########
    ## HELM ##
    ##########
    echo "Installing Helm..."
    if ! command -v helm &>/dev/null; then
        curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
        chmod 700 get_helm.sh
        ./get_helm.sh
        rm get_helm.sh
    fi

    ###################
    # PULL KNDP CHART #
    ###################
    echo "Adding kndp repository..."
    helm repo add kndp https://kndp.io
    echo "Updating kndp repository..."
    helm repo up kndp

    count=0
    total=41
    pstr="[=================================================================================================]"
    loading_pid=""

    function show_loading() {
        while [ $count -lt $total ]; do
            sleep 0.5
            count=$(($count + 1))
            pd=$(($count * 100 / $total))
            printf "\r%3d.%1d%% %.${pd}s" $(($count * 100 / $total)) $((($count * 1000 / $total) % 10)) $pstr
        done
        echo
    }

    # Start the loading animation in the background
    show_loading &
    loading_pid=$!

    echo "Installing kndp chart..."
    helm_output=$(helm install kndp kndp/kndp 2>&1)

    wait $loading_pid
    echo $helm_output
    exit 0

}

function uninstall_kndp() {
    # Uninstall KNDP
    if kind get clusters | grep -q "kndp"; then
        kind delete cluster --name kndp
        echo "KNDP cluster removed."
    else
        if helm list | grep -q "kndp"; then
            helm uninstall kndp
            echo "KNDP helm chart removed."
        else
            echo "KNDP cluster or helm chart not found."
        fi
    fi
    exit 0

}

# Check for the command
if [ "$1" == "install" ]; then
    install_kndp
    exit 0
elif [ "$1" == "uninstall" ]; then
    uninstall_kndp
    exit 0
elif [ "$1" == "help" ]; then
    help
    exit 0
fi

# Parse command line arguments
LONGOPTS="help,install,uninstall"
ARGS=$(getopt -o "hiu" --long "$LONGOPTS" -n "$(basename "$0")" -- "$@")
eval set -- "$ARGS"

while true; do
    case "$1" in
    -h | help)
        help
        exit 0
        ;;
    -i | install)
        echo "Installing KNDP..."
        install_kndp
        ;;

    -u | uninstall)
        uninstall_kndp
        ;;
    --)
        shift
        break
        ;;
    \?)
        echo "Invalid option: $OPTARG" >&2
        exit 1
        ;;
    esac
done

# No options provided, display help
if [ -z "$OPTARG" ]; then
    help
    exit 0
fi
