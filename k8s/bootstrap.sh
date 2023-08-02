kubectl apply -f ./argocd/ns.yaml
kubectl apply -k ./argocd/kustomize
kubectl apply -f ./argocd-apps.yaml
