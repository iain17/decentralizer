# CA key and certificate
openssl genrsa -out ca.key 4096
openssl req -new -x509 \
-subj "/emailAddress=ca@example.com/CN=ca.example.com/OU=ops/O=example.com services/L=City/ST=state/C=Country" \
-days 365 -key ca.key -out ca.crt

# server key
openssl genrsa -out server.key 1024
# CSR (certificate sign request) to obtain certificate
openssl req -new \
-subj "/emailAddress=svc@example.com/CN=svc.example.com/OU=ops/O=example.com services/L=City/ST=state/C=Country" \
-key server.key -out server.csr

# sign server CSR with CA certificate and key
openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key \
-extensions v3_server -extfile ./ssl-extensions-x509-example \
-set_serial 01 -out server.crt

# client key
openssl genrsa -out client.key 1024
# CSR to obtain certificate
openssl req -new \
-subj "/emailAddress=u@example.com/CN=user.example.com/OU=ops/O=example.com services/L=City/ST=state/C=Country" \
-key client.key -out client.csr

# sign client CSR with CA certificate and key
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key \
-extensions v3_client -extfile ./ssl-extensions-x509-example \
-set_serial 01 -out client.crt

# server key out to temp.key
openssl rsa -in server.key -out temp.key
# remove server.key
rm server.key
# make temp.key as server key
mv temp.key server.key

# client key out to temp.key
openssl rsa -in client.key -out temp.key
# remove client.key
rm client.key
# make temp.key as client key
mv temp.key client.key

openssl pkcs12 -export -clcerts -in client.crt -inkey client.key -out client.p12
openssl pkcs12 -in client.p12 -out client.pem -clcerts

cat ca.crt ca.key  > ca.pem

echo Run
echo defaults write com.google.Chrome AutoSelectCertificateForUrls -array
echo defaults write com.google.Chrome AutoSelectCertificateForUrls -array-add -string '{"pattern":"https://[*.]example.com","filter":{"ISSUER":{"CN":"ca.example.com"}}}'