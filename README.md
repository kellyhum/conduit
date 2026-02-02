# Conduit

A file-transfer CLI tool that lets users upload and download files with E2E encryption

### Cryptography
- File is encrypted with AES-GCM
- Public/private keys are generated with Ed25519, used for confidentiality + integrity
- Transport will be using HTTPS/TLS
- SHA256 is used for displaying the user id

### Tech Stack
- Golang (Bubble Tea TUI)
- Python/Flask
- SQL