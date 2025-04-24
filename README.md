# ChartSentinel 🛡️

**ChartSentinel** is a lightweight monitoring tool that watches ArgoCD applications and detects when deployed Helm charts are outdated.

## 🚀 Features
- Watches ArgoCD apps in real time
- Detects new versions of deployed Helm charts
- Notifies via logs, metrics, or external hooks (Slack, Webhook, etc.)
- Deployable via Helm chart

## 🛠️ Install

```bash
helm repo add chartsentinel https://your-domain.com/charts
helm install chartsentinel chartsentinel/chartsentinel

## ⚙️ Configuration

Via Helm values, environment, or flags. See values.yaml.

## 📡 Metrics

Exposes Prometheus metrics on /metrics.