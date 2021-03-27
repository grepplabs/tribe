## Implementation ideas

* Phase 1
- [ ] Storage
    - [ ] Object store (Minio API)
    - [ ] SQL (Postgres) - realm is a new database
- [ ] Automatic JWT rotation 
- [ ] Server TLS Listener cert rotation 
- [ ] OIDC client flows
- [ ] Master Password , HashiCorp Vault, AWS KMS, GCP Cloud KMS, Azure Key Vault
- [ ] Optional token / refresh token / payload encryption e.g. JWE (RSA_OAEP_256 / A128CBC_HS256)
    - [ ] Client configures own certificates 
    - [ ] Client send certificate in each request when data should be encrypted
- [ ] -o yaml or -o json options 

Others
- Legacy: Password Grant / rotation of user credentials in S3
