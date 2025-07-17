#!/bin/bash
set -e

echo "Installing Claude Code..."
npm install -g @anthropic-ai/claude-code

echo "Installing Taskfile..."
curl -sL https://taskfile.dev/install.sh | sh -s -- -d -b ~/.local/bin

echo "Setting up workspace..."
cd /workspace
echo "Dev environment ready! Clone your projects here."

echo "Setup complete!"