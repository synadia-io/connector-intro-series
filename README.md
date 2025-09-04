# Development Environment

General devcontainer for Node.js, Go, Git, and GitHub CLI development.

## Setup

### 1. Clone and Open
```bash
# From WSL2
cd ~/Projects
git clone https://github.com/phumulock/DevEnvironment.git
cd DevEnvironment
code .
```

### 2. Environment Variables

Set these in your WSL2 environment **before** opening the devcontainer:

```bash
# Add to ~/.bashrc
echo 'export PRIVATE_STORYBLOK_KEY="your-key-here"' >> ~/.bashrc
source ~/.bashrc

# Verify
echo $PRIVATE_STORYBLOK_KEY
```

### 3. Open in Container
VS Code will prompt to "Reopen in Container" - click it and wait for setup.

## Structure
```
DemoEnvironment/
├── .devcontainer/devcontainer.json
├── setup.sh
└── [your-project-files]    # Your projects go here
```

## Usage
```bash
# Inside devcontainer, your repo is at:
cd /workspaces/DemoEnvironment
# Add your project files here
npm run dev  # Port 4321 auto-forwards
```

## Troubleshooting

**Environment variables not working:**
```bash
# Check if set in WSL2 host
echo $PRIVATE_STORYBLOK_KEY

# If empty, add to ~/.bashrc and rebuild container
```

**Slow performance:**
- Make sure project is in WSL2 filesystem (`/home/rphum/Projects/`)
- Avoid Windows filesystem (`/mnt/c/`)

## Features
- Node.js LTS + npm
- Go development
- Git + GitHub CLI
- VS Code extensions pre-configured
- Port forwarding (4321 for Astro)