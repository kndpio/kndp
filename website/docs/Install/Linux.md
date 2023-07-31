
##




 # <div align="center"> 🚀  Install KNDP 🚀 </div>


###  Prerequisites

- 🖥️🐧 Linux-based operating system. This script is designed for Linux-based operating systems, specifically tested on Ubuntu.
  
- 🌐  Internet access for package downloads.
  
- 🔒🔑 Administrative privileges. Ensure you have administrative privileges to execute certain commands (e.g., sudo).

# Run the scrypt:


    sudo apt update && sudo apt install curl && curl -sSf https://raw.githubusercontent.com/web-seven/kndp/release/0.1/scripts/install.sh | bash
       
 
#

## Setup

### Git

The script checks if Git is installed and installs it if not. Git is a version control system used for tracking changes in source code during software development.

### Node.js

If Node.js is not installed, the script installs it using `nvm` (Node Version Manager). NVM allows you to manage multiple Node.js versions on the same system.

### Docker

If Docker is not installed, the script installs it and configures the necessary repositories. Docker is a platform that enables developers to develop, ship, and run applications in containers.

### NX workspace

If NX (Nrwl Extensions) is not installed, the script installs it using npm. NX is a set of extensible dev tools for monorepos, which allows you to manage and scale large codebases efficiently.

### KIND 

If KIND  is not installed, the script installs it. KIND allows you to install KNDP foundamental tools using Docker containers, ideal for local development and testing.


### KNDP base utilities
Installs locally our Internal Development Platform

## !! DISCLAIMER !!
#### If Docker was installed by the script and was not present before, it is recommended to log out and log back in to apply changes and run Docker without ``sudo``.