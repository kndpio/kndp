
# <div align="center">   KNDP Chart Deployment Script 🚀


##  Prerequisites

Before running the script, make sure you have the following prerequisites:

    🖥️ Linux-based operating system (tested on Ubuntu).
    🌐 Internet access for package downloads.
    🔒 Administrative privileges (for installing packages using sudo).

## Script Execution Steps 

The script performs the following steps:

    📦 Installs Necessary Packages:
        Installs Git if it is not already installed.
        Installs Node.js using NVM (Node Version Manager) if it is not already installed.
        Installs NX (NestJS CLI) globally if it is not already installed.

    🛠️ Installs and Set Up Kubernetes (Using Kind):
        Installs Docker if it is not already installed.
        Installs kind (Kubernetes in Docker) if it is not already installed.
        Creates a Kubernetes cluster using the configuration specified in kind-config.yaml.

    ⚙️ Installs Helm:
        Installs Helm (Helm 3) if it is not already installed.

    🚀 Deploys KNDP Helm Chart:
        Pulls the KNDP Helm chart.
        Extracts the Helm chart and installs it into the Kubernetes cluster using the helm install command.

 
## Important Considerations

    🐧 This script is designed for Linux-based operating systems, specifically tested on Ubuntu. Some parts may require adjustments for other distributions.
    🔑 Ensure you have administrative privileges to execute certain commands (e.g., sudo).
    🧪 It is always advisable to thoroughly test the script in a controlled environment before running it in production.

# Disclaimer 
 The script is provided as-is, and the author holds no responsibility for any adverse consequences or damages arising from its usage. Always review and understand the script's content before executing it on your system.