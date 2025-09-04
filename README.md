# NATS Connector Dev Container

A ready-to-use development container for working with NATS messaging and JetStream.

## 🚀 Quick Start

### 1. Open in Dev Container

This project is designed to run in a VS Code Dev Container:

1. **Install Prerequisites:**
   - [VS Code](https://code.visualstudio.com/)
   - [Docker Desktop](https://www.docker.com/products/docker-desktop)
   - VS Code "Dev Containers" extension

2. **Open the Container:**
   - Clone this repository
   - Open the folder in VS Code
   - When prompted, click "Reopen in Container"
   - Or use Command Palette (F1) → "Dev Containers: Reopen in Container"

The container will automatically set up your complete development environment with Go, Node.js, Git, and all necessary tools.

### 2. Add NATS Credentials

**Important:** You need NATS credentials to connect to NGS (NATS Global Service).

1. Get your credentials file from [Synadia NGS](https://app.ngs.global/) or your NATS administrator
2. Place the `.creds` file in the `Credentials/` folder:
   ```
   Credentials/NGS-Default-CLI.creds
   ```
3. The credentials file should look something like:
   ```
   -----BEGIN NATS USER JWT-----
   eyJ0eXAiOiJKV1QiLCJhbGc...
   -----END NATS USER JWT-----

   ************************* IMPORTANT *************************
   NKEY Seed printed below can be used to sign and prove identity.
   NKEYs are sensitive and should be treated as secrets.

   -----BEGIN USER NKEY SEED-----
   SUABC123...
   -----END USER NKEY SEED-----
   ```

> ⚠️ **Security Note:** The `Credentials/` folder is gitignored. Never commit credentials to version control.

### 3. Run the Publisher

Once your credentials are in place, start the temperature data publisher:

```bash
task publisher
```

Or run directly:
```bash
cd Publisher
go run main.go
```

## 📁 Project Structure

```
connector-intro-series/
├── Credentials/           # Place your .creds file here (gitignored)
│   └── NGS-Default-CLI.creds
├── Publisher/             # Sample data publisher
│   └── main.go
├── Taskfile.yaml          # Task automation
└── .devcontainer/         # Dev container configuration
```

## 🛠️ What's Included

The dev container comes pre-configured with:
- **Go** - For NATS client development
- **Node.js LTS** - For web-based tools and scripts
- **Git & GitHub CLI** - Version control
- **Task Runner** - Build automation
- **NATS Ports** - 4222 (client), 8222 (monitoring)

## 📝 Common Tasks

| Command | Description |
|---------|-------------|
| `task deps` | Install Go dependencies |
| `task publisher` | Run the sample publisher |

## 🔧 Troubleshooting

### Credentials Not Found
If you see `Credentials file not found`, ensure:
1. Your `.creds` file is in the `Credentials/` folder
2. The file is named exactly `NGS-Default-CLI.creds`
3. The file has proper permissions

### Connection Failed
If connection fails:
1. Check your credentials are valid and not expired
2. Ensure you have internet connectivity
3. Verify NGS service status at [status.ngs.global](https://status.ngs.global)

## 🔗 Resources

- [Get NGS Credentials](https://app.ngs.global/)
- [NATS Documentation](https://docs.nats.io/)
- [VS Code Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)