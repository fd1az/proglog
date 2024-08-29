#!/bin/sh

# Simular el valor de HOSTNAME como si estuviera en Kubernetes
HOSTNAME="proglog-1"  # Cambia el número aquí para probar diferentes escenarios

# Simular los valores de Helm (reemplaza estos valores con los valores reales que deseas probar)
RPC_PORT=8400
RELEASE_NAMESPACE="default"
SERF_PORT=8401

# Obtener el ID del hostname
ID=$(echo $HOSTNAME | rev | cut -d- -f1 | rev)

# Comando de generación del archivo config.yaml
cat <<EOD > config.yaml
data-dir: /var/run/proglog/data
rpc-port: ${RPC_PORT}
bind-addr: "$HOSTNAME.proglog.${RELEASE_NAMESPACE}.svc.cluster.local:${SERF_PORT}"
$([ $ID != 0 ] && echo "start-join-addrs: \"proglog-0.proglog.${RELEASE_NAMESPACE}.svc.cluster.local:${SERF_PORT}\"")
bootstrap: $([ $ID = 0 ] && echo true || echo false)
EOD

# Mostrar el contenido del archivo generado
cat config.yaml