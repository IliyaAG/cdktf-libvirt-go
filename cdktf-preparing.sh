#!/usr/bin/env bash
set -e

RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
RESET='\033[0m'


echo -e "${YELLOW} -> Starting CDKTF setup for Linux ${RESET}"

sudo apt update -y

# 1. Git
if ! command -v git &> /dev/null; then
    echo -e "${YELLOW} -> Installing git...${RESET}"
  sudo apt install -y git
else
  echo -e "${GREEN} -> git already installed.${RESET}"
fi

# 2. Curl
if ! command -v curl &> /dev/null; then
  echo -e "${YELLOW} -> Installing curl...${RESET}"
  sudo apt install -y curl
else
  echo -e "${GREEN} -> curl already installed.${RESET}"
fi

# 3. Build tools
if ! dpkg -s build-essential &> /dev/null; then
  echo -e "${YELLOW} -> Installing build-essential...${RESET}"
  sudo apt install -y build-essential
else
  echo -e "${GREEN} -> build-essential already installed.${RESET}"
fi

# 4. Node.js (>= 18) and npm
NODE_VERSION=$(node -v 2>/dev/null || echo "v0.0.0")
if [[ "$NODE_VERSION" < "v18" ]]; then
  echo -e "${YELLOW} -> Installing Node.js 18.x...${RESET}"
  curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
  sudo apt install -y nodejs
else
  echo -e "${GREEN} -> Node.js $NODE_VERSION already installed.${RESET}"
fi

# 5. Go (>= 1.20)
if ! command -v go &> /dev/null; then
  echo -e "${YELLOW} -> Installing Go 1.21...${RESET}"
  GO_VERSION="1.21.0"
  wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
  rm go${GO_VERSION}.linux-amd64.tar.gz
  echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
  source ~/.bashrc
else
  echo -e "${GREEN} Go $(go version) already installed.${RESET}"
fi

# --- Install CDKTF CLI ---
if ! command -v cdktf &> /dev/null; then
  echo -e "${YELLOW} -> Installing CDKTF CLI...${RESET}"
  sudo npm install -g cdktf-cli
else
  echo -e "${GREEN} -> CDKTF CLI already installed ($(cdktf --version))${RESET}"
fi

echo -e "${GREEN} -> Setup complete! You can now use 'cdktf' command.${RESET}"
