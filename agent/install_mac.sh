#!/bin/bash

# Ensure running as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (sudo ./install_mac.sh)"
  exit 1
fi

echo "Installing SentinelDesk Agent for Mac..."

AGENT_DIR="/Library/Application Support/SentinelDesk"
PLIST_FILE="/Library/LaunchDaemons/com.sentineldesk.agent.plist"

# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "arm64" ]; then
    BIN_SOURCE="Sentdesk_mac_apple_silicon"
else
    BIN_SOURCE="Sentdesk_mac_intel"
fi

# 1. Create installation directory
mkdir -p "$AGENT_DIR"

# 2. Copy files
if [ ! -f "$BIN_SOURCE" ]; then
    echo "Error: $BIN_SOURCE not found in current directory!"
    exit 1
fi

if [ ! -f "config.yaml" ]; then
    echo "Error: config.yaml not found in current directory!"
    exit 1
fi

cp "$BIN_SOURCE" "$AGENT_DIR/Sentdesk_mac"
cp "config.yaml" "$AGENT_DIR/config.yaml"
chmod +x "$AGENT_DIR/Sentdesk_mac"

# 3. Create LaunchDaemon plist
cat <<EOF > "$PLIST_FILE"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.sentineldesk.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>$AGENT_DIR/Sentdesk_mac</string>
    </array>
    <key>WorkingDirectory</key>
    <string>$AGENT_DIR</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/sentineldesk-agent.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/sentineldesk-agent.log</string>
</dict>
</plist>
EOF

# 4. Load the LaunchDaemon
launchctl unload "$PLIST_FILE" 2>/dev/null
launchctl load -w "$PLIST_FILE"

echo "Installation complete! SentinelDesk agent is now running in the background."
