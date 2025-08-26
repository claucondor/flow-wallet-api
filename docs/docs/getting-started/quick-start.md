# Quick Start Guide

Get up and running with Phoenix Wallet API in just a few minutes! This guide will walk you through creating your first account, checking balances, and making your first token transfer.

## 🚀 **Step 1: Start Phoenix Wallet API**

```bash
# Clone the repository
git clone https://github.com/flow-hydraulics/flow-wallet-api.git
cd flow-wallet-api

# Start in lightweight mode
make lightweight

# Wait for services to start (about 30 seconds)
# You'll see: "✅ Services started successfully!"
```

**Verify it's running:**
```bash
curl http://localhost:3000/v1/health/ready
# Expected response: {"status":"ready"}
```

## 📱 **Step 2: Create Your First Account**

Create a new Flow account that will be managed by Phoenix Wallet API:

```bash
curl -X POST http://localhost:3000/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Expected Response:**
```json
{
  "address": "0x01cf0e2f2f715450",
  "keys": [
    {
      "index": 0,
      "publicKey": "04a1b2c3d4e5f6789...",
      "signAlgo": "ECDSA_P256",
      "hashAlgo": "SHA3_256",
      "weight": 1000
    }
  ],
  "balance": "0.00000000"
}
```

**Save the address** - you'll need it for the next steps!

## 💰 **Step 3: Fund Your Account**

Since we're using the Flow emulator, we can fund accounts with test FLOW tokens:

```bash
# Replace YOUR_ADDRESS with the address from step 2
export ACCOUNT_ADDRESS="0x01cf0e2f2f715450"

# Fund the account with 100 FLOW tokens
curl -X POST "http://localhost:3000/v1/accounts/${ACCOUNT_ADDRESS}/fungible-tokens/FlowToken" \
  -H "Content-Type: application/json" \
  -d '{}'
```

The emulator automatically funds new accounts with FLOW tokens for testing.

## 🔍 **Step 4: Check Account Balance**

Verify your account has been funded:

```bash
curl "http://localhost:3000/v1/accounts/${ACCOUNT_ADDRESS}"
```

**Expected Response:**
```json
{
  "address": "0x01cf0e2f2f715450",
  "keys": [...],
  "balance": "1000.00000000"
}
```

You can also check specific token balances:

```bash
curl "http://localhost:3000/v1/accounts/${ACCOUNT_ADDRESS}/fungible-tokens/FlowToken"
```

## 💸 **Step 5: Create a Second Account**

Let's create another account to transfer tokens to:

```bash
curl -X POST http://localhost:3000/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Save the second address:**
```bash
export RECIPIENT_ADDRESS="0x179b6b1cb6755e31"
```

## 🔄 **Step 6: Transfer Tokens**

Now let's transfer 10 FLOW tokens from your first account to the second:

```bash
curl -X POST "http://localhost:3000/v1/accounts/${ACCOUNT_ADDRESS}/fungible-tokens/FlowToken/withdrawals" \
  -H "Content-Type: application/json" \
  -d "{
    \"recipient\": \"${RECIPIENT_ADDRESS}\",
    \"amount\": \"10.0\"
  }"
```

**Expected Response:**
```json
{
  "transactionId": "abc123def456...",
  "status": "pending",
  "amount": "10.0",
  "recipient": "0x179b6b1cb6755e31"
}
```

## ✅ **Step 7: Verify the Transfer**

Check that the transfer completed successfully:

```bash
# Check sender balance (should be reduced by 10)
curl "http://localhost:3000/v1/accounts/${ACCOUNT_ADDRESS}/fungible-tokens/FlowToken"

# Check recipient balance (should have 10 FLOW)
curl "http://localhost:3000/v1/accounts/${RECIPIENT_ADDRESS}/fungible-tokens/FlowToken"
```

## 🎉 **Congratulations!**

You've successfully:
- ✅ Started Phoenix Wallet API
- ✅ Created Flow accounts
- ✅ Checked account balances
- ✅ Transferred tokens between accounts

## 🔧 **Using the Interactive Documentation**

Phoenix Wallet API includes interactive documentation where you can test all endpoints:

1. **Open your browser** to http://localhost:8080
2. **Explore the API** using the interactive interface
3. **Test endpoints** directly from the documentation

## 📝 **Common Operations**

### **List All Accounts**
```bash
curl "http://localhost:3000/v1/accounts"
```

### **Get Transaction Details**
```bash
# Replace TRANSACTION_ID with actual transaction ID
curl "http://localhost:3000/v1/transactions/TRANSACTION_ID"
```

### **Check Available Tokens**
```bash
curl "http://localhost:3000/v1/fungible-tokens"
```

### **Get Account Transaction History**
```bash
curl "http://localhost:3000/v1/accounts/${ACCOUNT_ADDRESS}/transactions"
```

## 🔄 **Working with Idempotency**

For production applications, use idempotency keys to prevent duplicate operations:

```bash
# Generate a unique key (in real applications, use UUID)
IDEMPOTENCY_KEY="transfer-$(date +%s)-$(shuf -i 1000-9999 -n 1)"

curl -X POST "http://localhost:3000/v1/accounts/${ACCOUNT_ADDRESS}/fungible-tokens/FlowToken/withdrawals" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: ${IDEMPOTENCY_KEY}" \
  -d "{
    \"recipient\": \"${RECIPIENT_ADDRESS}\",
    \"amount\": \"5.0\"
  }"
```

If you run the same command again with the same `Idempotency-Key`, you'll get a `409 Conflict` response, preventing duplicate transfers.

## 🌐 **Testing with Different Networks**

### **Flow Testnet**
```bash
# Stop current instance
make lightweight-down

# Start with testnet (requires testnet account configuration)
make lightweight-testnet
```

### **Flow Mainnet**
```bash
# Stop current instance
make lightweight-down

# Start with mainnet (requires mainnet account and real FLOW)
make lightweight-mainnet-idempotent
```

## 🛠️ **Development Workflow**

### **Making Changes**
```bash
# Stop services
make lightweight-stop

# Make your changes to .env or configuration

# Restart services
make lightweight
```

### **Viewing Logs**
```bash
# View all logs
make lightweight-logs

# View specific service logs
docker compose -f docker-compose.lightweight.yml logs api
docker compose -f docker-compose.lightweight.yml logs emulator
```

### **Resetting Everything**
```bash
# Stop and remove all containers and data
make lightweight-down

# Start fresh
make lightweight
```

## 🔍 **Troubleshooting**

### **API Not Responding**
```bash
# Check if services are running
docker compose -f docker-compose.lightweight.yml ps

# Check API logs
docker compose -f docker-compose.lightweight.yml logs api
```

### **Account Creation Fails**
```bash
# Check emulator is running
curl http://localhost:3569/v1/blocks/sealed

# Check API can connect to emulator
docker compose -f docker-compose.lightweight.yml logs api | grep emulator
```

### **Transfer Fails**
```bash
# Verify account has sufficient balance
curl "http://localhost:3000/v1/accounts/${ACCOUNT_ADDRESS}/fungible-tokens/FlowToken"

# Check transaction status
curl "http://localhost:3000/v1/transactions/TRANSACTION_ID"
```

## 📚 **Next Steps**

Now that you've completed the quick start:

1. **[Learn Core Concepts](../concepts/architecture)** - Understand how Phoenix Wallet API works
2. **[Explore API Reference](../api-reference/overview)** - See all available endpoints
3. **[Integration Examples](../examples/basic-usage)** - Real-world usage patterns
4. **[Production Deployment](../deployment/production-setup)** - Deploy to production

## 💡 **Pro Tips**

- **Use environment variables** for addresses in scripts
- **Always check transaction status** before proceeding
- **Use idempotency keys** for critical operations
- **Monitor logs** when debugging issues
- **Test on emulator first** before using testnet/mainnet

Ready to build something amazing? Let's dive deeper into the [Core Concepts](../concepts/architecture)! 🚀