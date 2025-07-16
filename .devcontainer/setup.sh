#!/bin/bash
set -e

echo "Installing Claude Code..."
npm install -g @anthropic-ai/claude-code

echo "Installing Taskfile..."
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

echo "Setting up workspace..."
cd /workspace
echo "Dev environment ready! Clone your projects here."

echo "Setup complete!"