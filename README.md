# Kubernetes Native Development Platform (KNDP)

![KNDP Logo](./website/static/logo-transparent-small.png)

## Overview

The Kubernetes Native Development Platform (KNDP) is an innovative platform designed to streamline the development, deployment, and management of applications in a Kubernetes-native environment.

## Features

- **Kubernetes-Native**: Seamlessly integrate with Kubernetes, taking full advantage of its features and ecosystem.
- **Developer-Friendly**: Simplified workflows and tools tailored for developers.
- **Scalable**: Designed to scale with your needs, from small applications to large-scale deployments.
- **Extensible**: Easily extend and customize the platform with plugins and integrations.

## Getting Started

### Prerequisites

- Kubernetes cluster (v1.18+ recommended)
- `kubectl` command-line tool

### Installation

```bash
kubectl apply -f https://path-to-kndp-installation-file.yaml
```

### Quick Start

1. Deploy your first application:
    `kubectl apply -f ./examples/hello-world.yaml`

2. Access the application:
    `kubectl get svc hello-world`

## Documentation

For detailed documentation, tutorials, and guides, please visit our [official documentation](https://kndp.io/docs/introduction).

## Contributing

We welcome contributions from the community! Please read our [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on how to contribute.
