# Idempotency

Idempotency is a critical concept in Phoenix Wallet API that ensures the same operation can be performed multiple times without causing unintended side effects. This is especially important for financial operations where duplicate transactions could result in loss of funds.

## 🔄 **What is Idempotency?**

**Idempotency** means that performing the same operation multiple times produces the same result as performing it once. In the context of Phoenix Wallet API, this prevents:

- **Duplicate transactions** from network retries
- **Double spending** from user interface issues  
- **Race conditions** in concurrent systems
- **Accidental repeated operations** from client errors

## 🎯 **Why Idempotency Matters**

### **The Problem Without Idempotency**

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant API
    participant Blockchain
    
    User->>Frontend: Click "Send 100 FLOW"
    Frontend->>API: POST /transfer (100 FLOW)
    Note over API,Blockchain: Network timeout
    Frontend->>API: POST /transfer (100 FLOW) [Retry]
    API->>Blockchain: Transfer 100 FLOW
    API->>Blockchain: Transfer 100 FLOW [Duplicate!]
    
    Note over User: User loses 200 FLOW instead of 100!
```

### **The Solution With Idempotency**

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant API
    participant Database
    participant Blockchain
    
    User->>Frontend: Click "Send 100 FLOW"
    Frontend->>API: POST /transfer<br/>Idempotency-Key: uuid-123
    API->>Database: Check if uuid-123 exists
    Database-->>API: Not found
    API->>Database: Store uuid-123
    API->>Blockchain: Transfer 100 FLOW
    
    Note over API,Blockchain: Network timeout, retry
    Frontend->>API: POST /transfer<br/>Idempotency-Key: uuid-123
    API->>Database: Check if uuid-123 exists
    Database-->>API: Found! (409 Conflict)
    API-->>Frontend: 409 Conflict - Already processed
    
    Note over User: Only 100 FLOW transferred ✓
```

## 🛠️ **How Idempotency Works**

### **Idempotency Key Flow**

```mermaid
graph TB
    Request[HTTP Request] --> Header{Has Idempotency-Key?}
    Header -->|No| RequiredCheck{POST Request?}
    RequiredCheck -->|Yes| Reject[400 Bad Request]
    RequiredCheck -->|No| Process[Process Request]
    
    Header -->|Yes| CheckStore[Check Key Store]
    CheckStore --> Exists{Key Exists?}
    Exists -->|Yes| Conflict[409 Conflict]
    Exists -->|No| StoreKey[Store Key]
    StoreKey --> Process
    Process --> Response[Return Response]
```

### **Storage Mechanisms**

Phoenix Wallet API supports multiple storage backends for idempotency keys:

#### **Redis Storage (Production)**
```mermaid
graph LR
    API[API Request] --> Redis[(Redis)]
    Redis --> TTL[Time-To-Live: 1 hour]
    TTL --> Expire[Auto Expire]
```

**Characteristics:**
- High performance in-memory storage
- Automatic key expiration
- Supports distributed deployments
- Recommended for production

#### **Database Storage (Shared)**
```mermaid
graph LR
    API[API Request] --> Database[(SQL Database)]
    Database --> Table[idempotency_keys table]
    Table --> Cleanup[Manual Cleanup Required]
```

**Characteristics:**
- Uses same database as application
- Persistent storage
- Requires manual cleanup
- Good for single-instance deployments

#### **In-Memory Storage (Development)**
```mermaid
graph LR
    API[API Request] --> Memory[(In-Memory Map)]
    Memory --> Process[Single Process Only]
    Process --> Restart[Lost on Restart]
```

**Characteristics:**
- Fastest access time
- No external dependencies
- Not suitable for production
- Perfect for development and testing

## ⚙️ **Configuration Options**

### **Standard Mode Configuration**

```bash
# Enable idempotency (default: false)
FLOW_WALLET_DISABLE_IDEMPOTENCY_MIDDLEWARE=false

# Storage type: redis, shared, local
FLOW_WALLET_IDEMPOTENCY_MIDDLEWARE_DATABASE_TYPE=redis

# Redis connection (if using redis)
FLOW_WALLET_IDEMPOTENCY_MIDDLEWARE_REDIS_URL=redis://localhost:6379/
```

### **Lightweight Mode Configuration**

```bash
# Enable lightweight mode
FLOW_WALLET_LIGHTWEIGHT_MODE=true

# Optional: Enable idempotency in lightweight mode
FLOW_WALLET_LIGHTWEIGHT_IDEMPOTENCY=true
```

When `FLOW_WALLET_LIGHTWEIGHT_IDEMPOTENCY=true`:
- Automatically uses SQLite database for key storage
- No Redis dependency required
- Keys persist across container restarts
- Perfect for production deployments without Redis

## 🔧 **Implementation Details**

### **Idempotency Key Requirements**

```http
POST /v1/accounts/0x123.../fungible-tokens/FlowToken/withdrawals
Content-Type: application/json
Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000

{
  "recipient": "0x456...",
  "amount": "100.0"
}
```

**Key Requirements:**
- **Unique**: Each operation must have a unique key
- **UUID Format**: Recommended to use UUID v4
- **Client Generated**: Client is responsible for generating keys
- **Consistent**: Same operation should use same key for retries

### **Response Behavior**

| Scenario | HTTP Status | Response |
|----------|-------------|----------|
| First request with key | `200/201` | Normal response |
| Duplicate key | `409 Conflict` | "Idempotency-Key conflict" |
| Missing key (POST) | `400 Bad Request` | "Idempotency-Key header not found" |
| GET/DELETE requests | `200` | Key not required |

### **Key Expiration**

```mermaid
timeline
    title Idempotency Key Lifecycle
    
    section Creation
        Request Received : Key stored with TTL
    
    section Active Period
        1 Hour Window : Key prevents duplicates
        
    section Expiration
        After 1 Hour : Key automatically expires
        New Request : Same key can be reused
```

**Default Expiration**: 1 hour
- Balances security with usability
- Prevents indefinite key accumulation
- Allows eventual key reuse

## 🌐 **Network-Specific Behavior**

### **Deployment Modes Comparison**

| Mode | Idempotency | Storage | Persistence | Use Case |
|------|-------------|---------|-------------|----------|
| **Standard** | ✅ Redis | Redis | ✅ | Production |
| **Lightweight** | ❌ Default | - | - | Development |
| **Lightweight + Idempotent** | ✅ SQLite | SQLite | ✅ | Simple Production |

### **Network Commands**

```bash
# Local development (no idempotency)
make lightweight

# Local development (with idempotency)
make lightweight-idempotent

# Testnet (no idempotency)
make lightweight-testnet

# Testnet (with idempotency)
make lightweight-testnet-idempotent

# Mainnet (with idempotency - recommended!)
make lightweight-mainnet-idempotent
```

## 💡 **Best Practices**

### **Client Implementation**

```javascript
// Good: Generate unique key per operation
const idempotencyKey = crypto.randomUUID();

const response = await fetch('/v1/accounts/0x123.../withdrawals', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Idempotency-Key': idempotencyKey
  },
  body: JSON.stringify({
    recipient: '0x456...',
    amount: '100.0'
  })
});

// Handle idempotency conflicts
if (response.status === 409) {
  console.log('Operation already processed');
  // Don't retry, operation was successful
}
```

### **Key Generation Strategies**

#### **UUID v4 (Recommended)**
```javascript
const key = crypto.randomUUID();
// Example: 550e8400-e29b-41d4-a716-446655440000
```

#### **Timestamp + Random**
```javascript
const key = `${Date.now()}-${Math.random().toString(36)}`;
// Example: 1640995200000-k2j3h4g5f6
```

#### **Operation-Based**
```javascript
const key = `transfer-${fromAccount}-${toAccount}-${amount}-${timestamp}`;
// Example: transfer-0x123-0x456-100.0-1640995200000
```

### **Error Handling**

```javascript
async function safeTransfer(transferData) {
  const idempotencyKey = crypto.randomUUID();
  
  try {
    const response = await apiCall(transferData, idempotencyKey);
    return response;
  } catch (error) {
    if (error.status === 409) {
      // Idempotency conflict - operation already processed
      console.log('Transfer already completed');
      return { success: true, duplicate: true };
    }
    
    if (error.status >= 500) {
      // Server error - safe to retry with same key
      return await apiCall(transferData, idempotencyKey);
    }
    
    // Client error - don't retry
    throw error;
  }
}
```

## 🔍 **Monitoring and Debugging**

### **Idempotency Metrics**

```mermaid
graph TB
    subgraph "Metrics Collection"
        Requests[Total Requests] --> Unique[Unique Operations]
        Requests --> Duplicates[Duplicate Attempts]
        Duplicates --> Conflicts[409 Conflicts]
    end
    
    subgraph "Analysis"
        Conflicts --> Rate[Conflict Rate %]
        Rate --> Alert{High Rate?}
        Alert -->|Yes| Investigation[Investigate Client Issues]
        Alert -->|No| Normal[Normal Operation]
    end
```

### **Common Issues and Solutions**

| Issue | Symptom | Solution |
|-------|---------|----------|
| **High conflict rate** | Many 409 responses | Check client retry logic |
| **Missing keys** | 400 Bad Request | Ensure client sends keys |
| **Key reuse** | Unexpected conflicts | Generate truly unique keys |
| **Storage full** | Performance issues | Configure key expiration |

## 🚀 **Production Recommendations**

### **For High-Volume Applications**

1. **Use Redis Storage**
   ```bash
   FLOW_WALLET_IDEMPOTENCY_MIDDLEWARE_DATABASE_TYPE=redis
   FLOW_WALLET_IDEMPOTENCY_MIDDLEWARE_REDIS_URL=redis://redis-cluster:6379/
   ```

2. **Monitor Key Usage**
   ```bash
   # Enable detailed logging
   FLOW_WALLET_LOG_LEVEL=info
   ```

3. **Configure Appropriate TTL**
   ```go
   // Adjust expiration based on your use case
   IdempotencyHandlerOptions{
     Expiry: 2 * time.Hour, // Longer for critical operations
   }
   ```

### **For Simple Deployments**

1. **Use Lightweight Mode with Idempotency**
   ```bash
   make lightweight-idempotent
   ```

2. **Monitor SQLite Performance**
   ```bash
   # Check database size and performance
   ls -la data/wallet.db
   ```

Idempotency in Phoenix Wallet API provides essential protection against duplicate operations, ensuring the reliability and safety of your custodial wallet operations.