# ArgoSourceTracker

Petit serveur web en Go qui expose une API REST pour surveiller les Applications (ArgoCD) dans un cluster Kubernetes

Généres un rapport listant les tags disponibles pour chaque charts supérieur à la version actuellement installée

## ⚙️ Installation
*Incoming*


## Développement

### 1. Pré-requis

- Go 1.20+
- Un accès à un cluster Kubernetes (via `~/.kube/config`)
- ArgoCD installé dans le cluster (ou au moins des applications)

### 2. Clone du projet

```bash
git clone https://github.com/ArnaudTA/ArgoSourceTracker.git
cd ArgoSourceTracker
go mod tidy
