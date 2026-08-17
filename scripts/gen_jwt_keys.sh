#!/bin/sh
set -e

ENV_FILE=".env"

if grep -q "^JWT_PRIVATE_KEY=" "$ENV_FILE" 2>/dev/null; then
  exit 0
fi

echo "🔑 JWT_PRIVATE_KEY não encontrada em $ENV_FILE. Gerando par de chaves RSA..."

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out "$TMP_DIR/private.pem" 2>/dev/null
openssl rsa -in "$TMP_DIR/private.pem" -pubout -out "$TMP_DIR/public.pem" 2>/dev/null

PRIVATE_KEY=$(awk '{printf "%s\\n", $0}' "$TMP_DIR/private.pem")
PUBLIC_KEY=$(awk '{printf "%s\\n", $0}' "$TMP_DIR/public.pem")

{
  echo ""
  echo "JWT_PRIVATE_KEY=\"$PRIVATE_KEY\""
  echo "JWT_PUBLIC_KEY=\"$PUBLIC_KEY\""
} >> "$ENV_FILE"

echo "✅ Chaves JWT geradas e adicionadas ao $ENV_FILE"
