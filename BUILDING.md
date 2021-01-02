# Building

## TLS Certificates for Testing purposes

We use some test / dummy TLS certificates when running test servers and clients.
Below is how we generate these certificates

`server-csr.conf`

```bash
# server-csr.conf
cat > server-csr.conf <<EOF
[ req ]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C = "  "
ST = " "
L = " "
O = " "
OU = " "
CN = HELMSERVER

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
IP.1 = 127.0.0.1

[ v3_ext ]
subjectAltName = @alt_names

EOF
```

`client-csr.conf`

```bash
# client-csr.conf
cat > client-csr.conf <<EOF
[ req ]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C = "  "
ST = " "
L = " "
O = " "
OU = " "
CN = HELMCLIENT

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
IP.1 = 127.0.0.1

[ v3_ext ]
subjectAltName = @alt_names

EOF
```

Run the below commands to create the certificates

```bash
{
    EXPIRATION_IN_DAYS=3650

    openssl req -x509 -newkey rsa:4096 -keyout server_ca.key -out server_ca.crt -days $EXPIRATION_IN_DAYS -nodes -subj "/C=  /ST= /L= /O= /CN=HELMSERVERCERTCA"

    openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.csr -subj "/C=  /ST= /L= /O= /CN=HELMSERVER" -config server-csr.conf

    openssl x509 -req -in server.csr -CA server_ca.crt -CAkey server_ca.key -CAcreateserial -out server.crt -days $EXPIRATION_IN_DAYS -extensions v3_ext -extfile server-csr.conf

    openssl req -x509 -newkey rsa:4096 -keyout client_ca.key -out client_ca.crt -days $EXPIRATION_IN_DAYS -nodes -subj "/C=  /ST= /L= /O= /CN=HELMCLIENTCERTCA"

    openssl req -newkey rsa:4096 -nodes -keyout client.key -out client.csr -subj "/C=  /ST= /L= /O= /CN=HELMCLIENT" -config client-csr.conf

    openssl x509 -req -in client.csr -CA client_ca.crt -CAkey client_ca.key -CAcreateserial -out client.crt -days $EXPIRATION_IN_DAYS -extfile client-csr.conf
}
```

Now move the necessary certificates and keys to the right places

```bash
mv ca.crt server_ca.crt test_cert.crt test_key.key test_server.crt test_server.key testdata/tls/
```

You can cleanup the other unnecessary files

```bash
$ rm -v *.crt *.key *.srl server-csr.conf client-csr.conf
```
