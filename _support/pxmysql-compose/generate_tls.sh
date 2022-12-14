#!/usr/bin/env sh
#
# Copyright (c) 2022, Geert JM Vanderkelen
#

set -e

#
# This script generates the MySQL Server Certificate Authority (CA).
# It also generates server and client certificates.
#
# The server certificate has X.509 subjectAltName extension with value
# of 'IP:127.0.0.1'.
#
#

OUT_DIR=conf.d
CA_KEY="${OUT_DIR}/ca-key.pem"
CA_CERT="${OUT_DIR}/ca.pem"
SERVER_KEY="${OUT_DIR}/server-key.pem"
SERVER_REQ="${OUT_DIR}/server-req.pem"
SERVER_CERT="${OUT_DIR}/server-cert.pem"
CLIENT_KEY="${OUT_DIR}/client-key.pem"
CLIENT_REQ="${OUT_DIR}/client-req.pem"
CLIENT_CERT="${OUT_DIR}/client-cert.pem"
DAYS=3600

OPENSSL=$(command -v openssl)
if [ "${OPENSSL}" = "" ]; then
  echo "Error: openssl command not available in path"
  exit 1
fi

v=$($OPENSSL version)
case "${v}" in
"OpenSSL 1.1"*) ;;
"LibreSSL 3.3.6"*) ;;
*)
    echo "Error: expecting OpenSSL v1.1 or greater, or LibreSSL v3.3 or greater (got ${v})"
    exit 1
esac

# Create CA certificate
${OPENSSL} genrsa -out ${CA_KEY} 2048
${OPENSSL} req -new \
 -subj "/C=DE/CN=pxmysql-test-CA" \
 -x509 -days ${DAYS} -nodes -key ${CA_KEY} -out ${CA_CERT}

# Create server certificate (removing passphrase)
${OPENSSL} genrsa -out ${SERVER_KEY} 2048
${OPENSSL} req -new -nodes \
  -subj "/C=DE/CN=pxmysql-test-server" \
  -key ${SERVER_KEY} -out ${SERVER_REQ}
${OPENSSL} x509 -req -days ${DAYS} -in ${SERVER_REQ} \
  -CA ${CA_CERT} -CAkey ${CA_KEY} -CAcreateserial \
  -extfile server_x509_ext.conf \
  -out ${SERVER_CERT}

# Create client certificate (removing passphrase)
${OPENSSL} genrsa -out ${CLIENT_KEY} 2048
${OPENSSL} req -new -nodes \
  -subj "/C=DE/CN=pxmysql-test-client" \
  -key ${CLIENT_KEY} -out ${CLIENT_REQ}
${OPENSSL} x509 -req -days ${DAYS} -in ${CLIENT_REQ}  \
  -CA ${CA_CERT} -CAkey ${CA_KEY} -CAcreateserial \
  -out ${CLIENT_CERT}

rm ${OUT_DIR}/*-req.pem
