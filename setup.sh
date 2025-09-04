#!/bin/bash
set -e

echo "Installing basic tools..."
sudo apt update
sudo apt install -y telnet curl wget unzip

echo "Installing ngrok..."
curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null
echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | sudo tee /etc/apt/sources.list.d/ngrok.list
sudo apt update
sudo apt install -y ngrok

echo "Installing Claude Code..."
npm install -g @anthropic-ai/claude-code

echo "Installing Taskfile..."
curl -sL https://taskfile.dev/install.sh | sh -s -- -d -b ~/.local/bin

echo "Installing NATS Server..."
wget https://github.com/nats-io/nats-server/releases/download/v2.11.6/nats-server-v2.11.6-linux-amd64.tar.gz
tar -xzf nats-server-v2.11.6-linux-amd64.tar.gz
sudo mv nats-server-v2.11.6-linux-amd64/nats-server /usr/local/bin/
rm -rf nats-server-v2.11.6-linux-amd64*

echo "Installing NATS CLI..."
wget https://github.com/nats-io/natscli/releases/download/v0.1.4/nats-0.1.4-linux-amd64.zip
unzip nats-0.1.4-linux-amd64.zip
sudo mv nats-0.1.4-linux-amd64/nats /usr/local/bin/
rm -rf nats-0.1.4-linux-amd64*

echo "Installing Synadia Connect..."
# Clone, build, and install connect binary
git clone https://github.com/synadia-io/connect.git /tmp/connect
cd /tmp/connect
task build
task install  # Installs to ~/.local/bin
cd /workspaces/connector-intro-series
rm -rf /tmp/connect
echo "Connect installed to ~/.local/bin/connect"

echo "Creating NATS config directory..."
mkdir -p /workspaces/connector-intro-series/.nats

echo "Creating basic NATS config..."
cat > /workspaces/connector-intro-series/.nats/nats.conf << 'EOF'
# NATS Server Configuration

# Client port
port: 4222

# HTTP monitoring port
http_port: 8222

# Logging
log_file: "/workspaces/connector-intro-series/.nats/nats.log"
logtime: true
debug: false
trace: false

# Limits
max_connections: 64K
max_control_line: 4KB
max_payload: 1MB
max_pending: 64MB

# Authentication (uncomment to enable)
# authorization {
#   user: "admin"
#   password: "password"
# }

# JetStream (uncomment to enable)
# jetstream {
#   store_dir: "/workspaces/connector-intro-series/.nats/jetstream"
#   max_memory_store: 1GB
#   max_file_store: 10GB
# }
EOF

echo "Setting up environment variables..."
echo 'export EDITOR="code --wait"' >> ~/.bashrc
echo 'export VISUAL="code --wait"' >> ~/.bashrc

echo "Setting up workspace..."
cd /workspaces/connector-intro-series
echo "Dev environment ready! Clone your projects here."

echo "NATS Server installed! Usage:"
echo "  Start NATS: nats-server -c /workspaces/connector-intro-series/.nats/nats.conf"
echo "  Check status: nats server check"
echo "  Monitor: http://localhost:8222"

echo "Setup complete!"