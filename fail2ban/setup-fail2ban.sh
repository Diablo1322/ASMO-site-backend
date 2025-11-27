#!/bin/bash

echo "🔒 Setting up Fail2Ban..."

# Install Fail2Ban
apt update
apt install -y fail2ban

# Copy configuration files
cp fail2ban/jail.local /etc/fail2ban/
cp fail2ban/nginx-auth.conf /etc/fail2ban/filter.d/
cp fail2ban/nginx-botsearch.conf /etc/fail2ban/filter.d/

# Set proper permissions
chmod 644 /etc/fail2ban/jail.local
chmod 644 /etc/fail2ban/filter.d/nginx-*.conf

# Create log directory
mkdir -p /var/log/fail2ban

# Start and enable Fail2Ban
systemctl enable fail2ban
systemctl start fail2ban

# Check status
echo "📊 Fail2Ban status:"
fail2ban-client status

echo "✅ Fail2Ban setup completed!"
echo "🔧 Commands:"
echo "   fail2ban-client status"
echo "   fail2ban-client status nginx-http-auth"
echo "   tail -f /var/log/fail2ban.log"