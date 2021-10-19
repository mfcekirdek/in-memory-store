helm repo add gitlab https://charts.gitlab.io
helm search repo -l gitlab/gitlab-runner
kubectl create ns gitlab-runner
helm install --namespace gitlab-runner gitlab-runner --dry-run -f values.yaml gitlab/gitlab-runner