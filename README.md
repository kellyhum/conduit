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

### Architectural Pivots
- Started with Claude/ChatGPT-ing an architecture - my only constraints were that I knew I wanted AES and ed25519 encryption (googled the best digital signature algorithms -> ECDSA -> ed25519)
    - Realized that I didn't actually understand what was happening, so instead, started with user stories: what commands would a user input?
        - Worked out better :")

- Started with the Bubble Tea TUI package => realized that I was spending more time customizing the user-interface rather than building in the command functionality, and that I'd eventually need a way for the user to get commands - `if user_input == "conduit xyz"` didn't seem like a structurally best practice for code (very repetitive and fragile)
    - Searched for packages designed for building CLI commands, specifically noted `Cobra` and `urfav/cli`, and how they're built for standard CLI commands rather than my original idea of an "interactive shell"
        - Thought this was a better idea instead of an "interactive shell" b/c I only really wanted a welcome message, and the CLI packages are already fairly popular with a similar use-case
            - Chose urfav/cli b/c of its simpler usage (not building anything super complex so don't need Cobra) and can always pivot to Cobra if this grows (but don't anticipate it)