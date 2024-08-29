#!/bin/bash

# Nombre del release de Helm
HELM_RELEASE_NAME="proglog"
NAMESPACE="default"

# 1. Eliminar el release de Helm
echo "Eliminando release de Helm: $HELM_RELEASE_NAME"
helm uninstall $HELM_RELEASE_NAME -n $NAMESPACE

# 2. Eliminar todos los pods en el namespace especificado
echo "Eliminando todos los pods en el namespace: $NAMESPACE"
kubectl delete pods --all -n $NAMESPACE

# 3. Eliminar todos los Persistent Volume Claims (PVC) en el namespace especificado
echo "Eliminando todos los PVCs en el namespace: $NAMESPACE"
kubectl delete pvc --all -n $NAMESPACE

# 4. Eliminar todos los Persistent Volumes (PV) no reclamados
echo "Eliminando todos los PVs no reclamados"
kubectl get pv --no-headers | awk '/Released/ {print $1}' | xargs kubectl delete pv

# 5. Eliminar todos los StatefulSets en el namespace especificado
echo "Eliminando todos los StatefulSets en el namespace: $NAMESPACE"
kubectl delete statefulsets --all -n $NAMESPACE

echo "Recursos eliminados con éxito."
