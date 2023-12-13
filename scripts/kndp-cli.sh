#!/bin/bash

function help() {
    echo "Usage: kndp [flags] [commands]"
    echo ""
    echo "Commands:"
    echo "   help             For more information about existing commands"
    echo "   install          Install KNDP"
    echo "   uninstall        Uninstall KNDP"
    echo "   upgrade          Upgrade KNDP"

    echo ""
    echo "Flags:"
    echo "   --cluster, -c    Set existing cluster or create new with given name"
    echo "   --config,  -f    Specify config in a YAML file"
}

function install_kndp() {
    echo "Installing required tools..."
    #############
    ### GIT #####
    #############
    if ! command -v git &>/dev/null; then
        echo " Installing Git..."
        sudo apt-get update
        sudo apt-get install -y git
    else
        echo " ✓ Git is already installed. Skipping..."
    fi

    ############
    ## NODEJS ##
    ############
    if ! command -v node &>/dev/null; then
        echo " Installing Node.js..."
        sudo apt-get update
        sudo apt-get install -y ca-certificates curl gnupg
        sudo mkdir -p /etc/apt/keyrings
        curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
        NODE_MAJOR=20
        echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
        sudo apt-get update
        sudo apt-get install nodejs -y
    else
        echo " ✓ Node.js is already installed. Skipping..."
    fi

    ############
    ## Docker ##
    ############
    if ! command -v docker &>/dev/null; then
        echo " Installing Docker..."
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
        echo " ✓ Docker installed. You may need to log out and log back in to use Docker without sudo."
    else
        echo " ✓ Docker is already installed. Skipping..."
    fi

    #######################################
    ### Install NX and create Workspace ###
    #######################################
    if ! command -v nx &>/dev/null; then
        echo " Installing NX..."
        npx create-nx-workspace --skipGit ./
        npm install -g nx
    else
        echo " ✓ NX is already installed."
    fi

    ###########################
    #  install and setup KIND #
    ###########################
    if ! command -v kind &>/dev/null; then
        echo " Installing kind..."
        curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
        chmod +x ./kind
        sudo mv ./kind /usr/local/bin/kind
    else
        echo " ✓ KIND is already installed."
    fi

    ##########
    ## HELM ##
    ##########
    if ! command -v helm &>/dev/null; then
        echo "Installing Helm..."
        curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
        chmod 700 get_helm.sh
        ./get_helm.sh
        rm get_helm.sh
    else
        echo " ✓ Helm is already installed."
    fi
    echo "KNDP basic tools are installed."
    echo ""

    ##########################################
    # create a Kubernetes cluster using kind #
    ##########################################
    echo "Creating Kubernetes cluster using kind..."
    echo " Found clusters:"
    if kind get clusters | grep $cluster_name ; then
        echo " ✓ Cluster $cluster_name already exists. Skipping cluster creation."
    else
        cat <<EOF | kind create cluster --name $cluster_name --config=-
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
    fi


    ###################
    # PULL KNDP CHART #
    ###################
    helm repo add kndp https://kndp.io &>/dev/null
    helm repo up kndp &>/dev/null

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

    if [ ! -z $config_file ]; then
        config_flag="-f $config_file"
    fi

    kndp_exists=$(helm list | grep -c "kndp")

    if [ $kndp_exists = "1" ]; then
        echo "Error: KNDP already installed, please use \"upgrade\" command for update it."
        exit 1
    fi

     # Start the loading animation in the background
    show_loading &
    loading_pid=$!

    echo " Installing KNDP into cluster '$cluster_name'..."

    helm_output=$(helm install kndp kndp/kndp $config_flag --kube-context "kind-$cluster_name" 2>&1)

    wait $loading_pid
    echo $helm_output
    exit 0
}

function upgrade_kndp() {
    echo "Upgrading KNDP..."

    helm repo up kndp &>/dev/null
    if [ ! -z $config_file ]; then
        config_flag="-f $config_file"
    fi
    helm upgrade kndp kndp/kndp $config_flag
    chart_version=$(helm ls -o yaml -f kndp | grep chart | awk -F "-" '{print $2}')
    echo "KNDP upgraded to version $chart_version."
}

function uninstall_kndp() {
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

cluster_name="kndp"
config_file=""
command=""

if [ -s "kndp.yaml" ]; then
    config_file="kndp.yaml"
fi

while [[ "${#}" -gt 0 ]]; do
    case "$1" in
        -c | --cluster)
            shift
            if [ -z "$1" ]; then
                echo "Error: Cluster name cannot be empty. Please provide a valid cluster name."
                exit 1
            fi
            cluster_name=$1
        ;;
        -f | --config)
            shift
            if [ -z "$1" ]; then
                echo "Error: Config name cannot be empty. Please provide a valid config name."
                exit 1
            fi
            config_file=$1
        ;;

        -h | --help)
            command="help"
        ;;

        install)
            command="install"
        ;;
    
        upgrade)
            command="upgrade"
        ;;
    
        uninstall)
            command="uninstall"
        ;;
    
        -*)
            printf "Error: Flag \"$1\" not found.\n\n" 
            help
            exit 1
        ;;
        
        *)
            printf "Error: Command \"$1\" not found.\n\n" 
            help
            exit 1
        ;;
    --)
        shift
        break
        ;;
    esac
    shift
done

case "${command}" in
    install)
        install_kndp
    ;;
    uninstall)
        uninstall_kndp
    ;;
    upgrade)
        upgrade_kndp
    ;;
    help)
        help
    ;;
    *)
        printf "Error: Command not provider.\n\n" 
        help
    ;;
esac
exit 0
