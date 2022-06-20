export TA_CONFIG_FILE=/app/config.yaml
if [[ -z "${TPA_HOST}" ]]; then
  TPA_HOST="host.docker.internal"
else
  TPA_HOST="${TPA_HOST}"
fi
sed -i "s/TPA_HOST/$TPA_HOST/g" /app/config.yaml
cd /app
/app/main
