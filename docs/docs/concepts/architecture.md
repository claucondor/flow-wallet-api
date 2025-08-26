# Architecture Overview

Phoenix Wallet API is built with a modular, service-oriented architecture that provides scalability, maintainability, and flexibility for custodial wallet management on Flow blockchain.

## 🏗️ **High-Level Architecture**

```mermaid
graph TB
    subgraph "Client Layer"
        WebApp[Web Application]
        Mobile[Mobile App]
        CLI[CLI Tools]
        API_Client[API Clients]
    end
    
    subgraph "Phoenix Wallet API"
        Gateway[HTTP Gateway]
        Auth[Authentication]
        Middleware[Middleware Layer]
        
        subgraph "Core Services"
            AccountSvc[Account Service]
            TokenSvc[Token Service]
            TxSvc[Transaction Service]
            JobSvc[Job Service]
            SystemSvc[System Service]
        end
        
        subgraph "Infrastructure"
            KeyMgr[Key Manager]
            WorkerPool[Worker Pool]
            ChainListener[Chain Event Listener]
        end
    end
    
    subgraph "External Dependencies"
        FlowNet[Flow Blockchain]
        Database[(Database)]
        KMS[Key Management Service]
        Redis[(Redis Cache)]
    end
    
    WebApp --> Gateway
    Mobile --> Gateway
    CLI --> Gateway
    API_Client --> Gateway
    
    Gateway --> Auth
    Auth --> Middleware
    Middleware --> AccountSvc
    Middleware --> TokenSvc
    Middleware --> TxSvc
    Middleware --> JobSvc
    Middleware --> SystemSvc
    
    AccountSvc --> KeyMgr
    TokenSvc --> WorkerPool
    TxSvc --> WorkerPool
    JobSvc --> WorkerPool
    
    KeyMgr --> KMS
    WorkerPool --> FlowNet
    ChainListener --> FlowNet
    
    AccountSvc --> Database
    TokenSvc --> Database
    TxSvc --> Database
    JobSvc --> Database
    SystemSvc --> Database
    
    Middleware --> Redis
```

## 🔧 **Core Components**

### **HTTP Gateway Layer**
The entry point for all API requests, handling:
- **Request Routing**: Directs requests to appropriate services
- **Authentication**: Validates API keys and permissions
- **Rate Limiting**: Prevents abuse and ensures fair usage
- **CORS Handling**: Enables cross-origin requests for web applications

### **Middleware Layer**
Provides cross-cutting concerns:
- **Logging**: Comprehensive request/response logging
- **Compression**: Reduces bandwidth usage
- **Idempotency**: Prevents duplicate operations
- **Error Handling**: Standardized error responses

### **Service Layer**
Business logic organized into focused services:

#### **Account Service**
```mermaid
graph LR
    AccountAPI[Account API] --> AccountService[Account Service]
    AccountService --> AccountStore[Account Store]
    AccountService --> KeyManager[Key Manager]
    AccountService --> FlowClient[Flow Client]
    AccountStore --> Database[(Database)]
    KeyManager --> KMS[Key Management Service]
    FlowClient --> FlowBlockchain[Flow Blockchain]
```

**Responsibilities:**
- Create and manage Flow accounts
- Handle account key operations
- Manage account metadata and state

#### **Transaction Service**
```mermaid
graph LR
    TxAPI[Transaction API] --> TxService[Transaction Service]
    TxService --> TxStore[Transaction Store]
    TxService --> WorkerPool[Worker Pool]
    TxService --> RateLimiter[Rate Limiter]
    TxStore --> Database[(Database)]
    WorkerPool --> FlowClient[Flow Client]
    FlowClient --> FlowBlockchain[Flow Blockchain]
```

**Responsibilities:**
- Execute transactions on Flow blockchain
- Manage transaction lifecycle and status
- Handle transaction signing and submission
- Provide transaction history and details

#### **Token Service**
```mermaid
graph LR
    TokenAPI[Token API] --> TokenService[Token Service]
    TokenService --> TokenStore[Token Store]
    TokenService --> TemplateService[Template Service]
    TokenService --> AccountService[Account Service]
    TokenStore --> Database[(Database)]
    TemplateService --> CadenceTemplates[Cadence Templates]
```

**Responsibilities:**
- Manage fungible and non-fungible tokens
- Handle token transfers and deposits
- Track token balances and metadata
- Support custom token implementations

#### **Job Service**
```mermaid
graph LR
    JobAPI[Job API] --> JobService[Job Service]
    JobService --> JobStore[Job Store]
    JobService --> WorkerPool[Worker Pool]
    JobStore --> Database[(Database)]
    WorkerPool --> JobProcessor[Job Processor]
    JobProcessor --> FlowClient[Flow Client]
```

**Responsibilities:**
- Manage asynchronous operations
- Queue and process background jobs
- Provide job status and progress tracking
- Handle job retries and error recovery

## 🔄 **Data Flow Architecture**

### **Synchronous Operations**
```mermaid
sequenceDiagram
    participant Client
    participant API
    participant Service
    participant Database
    participant Flow
    
    Client->>API: HTTP Request
    API->>Service: Process Request
    Service->>Database: Read/Write Data
    Service->>Flow: Blockchain Query
    Flow-->>Service: Response
    Database-->>Service: Data
    Service-->>API: Result
    API-->>Client: HTTP Response
```

### **Asynchronous Operations**
```mermaid
sequenceDiagram
    participant Client
    participant API
    participant JobService
    participant WorkerPool
    participant Flow
    participant Database
    
    Client->>API: Async Request
    API->>JobService: Create Job
    JobService->>Database: Store Job
    JobService-->>API: Job ID
    API-->>Client: Job Created Response
    
    WorkerPool->>Database: Poll for Jobs
    Database-->>WorkerPool: Job Details
    WorkerPool->>Flow: Execute Operation
    Flow-->>WorkerPool: Result
    WorkerPool->>Database: Update Job Status
    
    Note over Client: Client can poll job status
    Client->>API: Check Job Status
    API->>Database: Query Job
    Database-->>API: Job Status
    API-->>Client: Status Response
```

## 🗄️ **Data Architecture**

### **Database Schema Overview**
```mermaid
erDiagram
    ACCOUNTS ||--o{ ACCOUNT_KEYS : has
    ACCOUNTS ||--o{ TRANSACTIONS : owns
    ACCOUNTS ||--o{ ACCOUNT_TOKENS : holds
    
    TRANSACTIONS ||--o{ JOBS : creates
    TOKENS ||--o{ ACCOUNT_TOKENS : referenced_by
    
    ACCOUNTS {
        string address PK
        string type
        jsonb keys
        timestamp created_at
        timestamp updated_at
    }
    
    ACCOUNT_KEYS {
        uuid id PK
        string account_address FK
        int key_index
        string public_key
        string private_key_encrypted
        string key_type
    }
    
    TRANSACTIONS {
        uuid id PK
        string account_address FK
        string transaction_id
        string status
        jsonb result
        timestamp created_at
    }
    
    TOKENS {
        uuid id PK
        string name
        string address
        string type
        jsonb configuration
    }
    
    ACCOUNT_TOKENS {
        uuid id PK
        string account_address FK
        uuid token_id FK
        decimal balance
        timestamp last_updated
    }
    
    JOBS {
        uuid id PK
        string type
        string status
        jsonb attributes
        int error_count
        timestamp created_at
        timestamp updated_at
    }
```

## 🔐 **Security Architecture**

### **Key Management Flow**
```mermaid
graph TB
    subgraph "Key Generation"
        Generate[Generate Key Pair] --> Encrypt[Encrypt Private Key]
        Encrypt --> Store[Store in Database]
    end
    
    subgraph "Key Usage"
        Retrieve[Retrieve Encrypted Key] --> Decrypt[Decrypt Private Key]
        Decrypt --> Sign[Sign Transaction]
        Sign --> Clear[Clear from Memory]
    end
    
    subgraph "Key Management Services"
        Local[Local Encryption]
        GoogleKMS[Google Cloud KMS]
        AWSKMS[AWS KMS]
    end
    
    Encrypt --> Local
    Encrypt --> GoogleKMS
    Encrypt --> AWSKMS
    
    Decrypt --> Local
    Decrypt --> GoogleKMS
    Decrypt --> AWSKMS
```

### **Authentication & Authorization**
```mermaid
graph LR
    Request[API Request] --> Auth{Authentication}
    Auth -->|Valid| RateLimit[Rate Limiting]
    Auth -->|Invalid| Reject[Reject Request]
    RateLimit -->|Within Limits| Process[Process Request]
    RateLimit -->|Exceeded| Throttle[Throttle Request]
    Process --> Response[API Response]
```

## 🚀 **Deployment Architectures**

### **Lightweight Mode**
```mermaid
graph TB
    subgraph "Single Container"
        API[Phoenix Wallet API]
        SQLite[(SQLite Database)]
        InMemory[(In-Memory Queue)]
    end
    
    API --> SQLite
    API --> InMemory
    API --> FlowNetwork[Flow Network]
    
    Client[Client Applications] --> API
```

**Characteristics:**
- Single container deployment
- SQLite for data persistence
- In-memory job processing
- Perfect for development and small-scale production

### **Standard Mode**
```mermaid
graph TB
    subgraph "Application Tier"
        API1[API Instance 1]
        API2[API Instance 2]
        LoadBalancer[Load Balancer]
    end
    
    subgraph "Data Tier"
        PostgreSQL[(PostgreSQL)]
        Redis[(Redis)]
    end
    
    subgraph "External Services"
        FlowNetwork[Flow Network]
        KMS[Key Management Service]
    end
    
    Client[Client Applications] --> LoadBalancer
    LoadBalancer --> API1
    LoadBalancer --> API2
    
    API1 --> PostgreSQL
    API1 --> Redis
    API1 --> FlowNetwork
    API1 --> KMS
    
    API2 --> PostgreSQL
    API2 --> Redis
    API2 --> FlowNetwork
    API2 --> KMS
```

**Characteristics:**
- Horizontally scalable
- Shared PostgreSQL database
- Redis for caching and job queues
- Production-ready with high availability

## 📊 **Performance Considerations**

### **Scalability Patterns**
- **Horizontal Scaling**: Multiple API instances behind load balancer
- **Database Optimization**: Connection pooling and query optimization
- **Caching Strategy**: Redis for frequently accessed data
- **Rate Limiting**: Prevent abuse and ensure fair resource usage

### **Monitoring & Observability**
- **Health Checks**: Liveness and readiness probes
- **Metrics Collection**: Performance and usage metrics
- **Logging**: Structured logging for debugging and audit
- **Error Tracking**: Comprehensive error reporting and alerting

This architecture provides a solid foundation for building scalable, secure, and maintainable custodial wallet solutions on Flow blockchain.