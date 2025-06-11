
<div align=center>
<a href="https://dw.mangodemo.com" target="_blank">mgo-dw-demo</a>
</div>

## knowledge base

[mgo-go-sdk](https://github.com/MangoNet-Labs/mgo-go-sdk) : https://github.com/MangoNet-Labs/mgo-go-sdk

[solana-go](https://github.com/gagliardetto/solana-go) : https://github.com/gagliardetto/solana-go

[go-zero](https://github.com/zeromicro/go-zero) : https://github.com/zeromicro/go-zero

## 1. Overview

### 1.1 Project Introduction

> `mgo-dw-demo` is a deposit and withdrawal demo built on top of the `mgo-go-sdk`, featuring JWT-based authentication, address generation, transaction signing, event listeners, record queries, and integration with both the Mango and Solana blockchains. It includes a variety of example files so you can focus more of your time on business logic.

[Live Demo](https://dw.mangodemo.com): https://dw.mangodemo.com

### 1.2 Mango Blockchain Go SDK
Hi! Thank you for using mgo-go-sdk。

`mgo-go-sdk` is the official Go SDK for Mango Blockchain, providing capabilities to interact with the Mango chain, including account management, transaction building, token operations, and more.

- **Account Management**: Supports `ed25519` and `secp256k1` key pair generation and signing.
- **Transaction Operations**: Provides transaction creation, signing, and submission.
- **On-Chain Data Queries**: Supports querying token balances, events, objects, and more.
- **WebSocket Subscription**: Enables real-time event push notifications.
## 2. Usage Guide

### 2.1 backend project

```bash

# Clone the repository
git clone https://github.com/MangoNet-Labs/mgo-dw-demo.git
# Change into the backend folder
cd backend

# Generate code and install dependencies
go generate

# Run the server
go run . 

```
## 3. Project Architecture

### 3.1 Directory Structure

```
    backend
    ├── api                     (API definition layer; handlers, types, and clients are generated from .api files)
    ├── common                  (Common utilities and shared logic)
    │   └── response            (Unified response formats and helpers)
    ├── core                    (Core business logic)
    ├── etc                     (Configuration templates for different environments)
    │   └── user-api.yaml       (Service configuration file)
    ├── internal                (Service-private implementation; not importable by external packages)
    │   ├── config              (Config structs generated from etc/*.yaml)
    │   ├── handler             (HTTP/GRPC handlers; receive requests and invoke logic)
    │   ├── logic               (Business logic implementations)
    │   │   └── common          (Shared logic components for internal reuse)
    │   ├── svc                 (Service Context: bundles DB, Redis, RPC clients, etc.)
    │   └── types               (Request/Response structs generated from .api files)
    ├── middleware              (Middleware layer)
    │   └── jwtmiddleware.go    (JWT authentication and context injection)
    ├── model                   (Data models mapped to database tables with CRUD operations)
    ├── third                   (Third-party SDKs or client wrappers for external services)
    │   └── client.go           (Wrapper for external RPC/HTTP calls)
    ├── go.mod                  (Go module declaration and dependency versions)
    ├── go.sum                  (Dependency checksums)
    └── user.go                 (Service entry point; initializes and starts HTTP/GRPC server)

```

## 5.Quick Start

### 1. Initialize the Client

```go
package main

import (
    "fmt"
    "github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
)

func main() {
    cli := client.NewMgoClient(config.RpcMgoDevnetEndpoint)
    fmt.Println("Mango Client Initialized:", cli)
}
```

### 2. Generate a Key Pair

```go
package main

import (
    "fmt"
    "github.com/mangonet-labs/mgo-go-sdk/account/keypair"
)

func main() {
	kp, err := keypair.NewKeypair(config.Ed25519Flag)
    if err != nil {
		log.Fatalf("%v", err)
		return
	}
    fmt.Println("Public Key:", kp.PublicKeyHex())
}
```

### 3. Mango Blockchain PTB Transfer Example (Go SDK)

## ✅ Key Components

| Component         | Description                                                   |
|------------------|---------------------------------------------------------------|
| `walletKey`      | User's Mango wallet key used to sign the transaction          |
| `sysKey`         | Sponsor (system account) key used to pay gas fees             |
| `gasCoin`        | Coin object belonging to the sponsor account                  |
| `SplitCoins`     | Splits the merged coin to get the desired transfer amount     |
| `TransferObjects`| Transfers the coin to the recipient address                   |
| `Execute()`      | Sends the PTB transaction and waits for confirmation          |

---
## 💡 Flow Overview

1. **Validate Inputs**: Check recipient address format and user balance.
2. **Prepare Keys**: Load private keys for both user and sponsor.
3. **Get Sponsor Coin**: Fetch sponsor's gas coin object.
4. **Build Transaction**:
    - Setup signer and sponsor.
    - Set gas parameters.
    - Merge user MGO coins to get enough amount.
    - Split coins and transfer to the target address.
5. **Execute Transaction** and return the `Digest` (transaction hash).

---
```go

    //https://github.com/MangoNet-Labs/mgo-go-sdk/blob/main/test/ptb/ptb_test.go
    //https://github.com/MangoNet-Labs/mgo-dw-demo/blob/main/backend/internal/logic/common/withdrawallogic.go
    //This module demonstrates how to perform a **withdrawal transaction** using Mango's Go SDK.
    //It executes a **sponsored Programmable Transaction Block (PTB)**, where a system account pays 
    //the gas fee, and a user account initiates the MGO transfer.
	
    // Step 1: Validate recipient address length
    if len(req.ToAddress) != 66 {
        return nil, errors.New("receiving address error") // Invalid Mango address
    }
    // Step 2: Fetch current user's address record from DB
    var user model.Address
    txUser := l.svcCtx.DB.Where("id = ?", authUser.UserID).First(&user)
    if txUser.Error != nil {
        return nil, errors.New("user does not exist") // No matching record
    }
    // Step 3: Convert the total balance (string) to float64
    Balance, err := strconv.ParseFloat(TotalBalance, 10)
    if err != nil {
        return nil, err // Parsing error
    }
    // Step 4: Ensure user has enough balance for withdrawal
    if Balance < req.Amount {
        return nil, errors.New("insufficient cash withdrawal amount")
    }
    // Step 5: Create Mango SDK client and context
    cli := l.svcCtx.MgoCli
    var ctx = context.Background()
    // Step 6: Load the user's Mango private key into a keypair object
    walletKey, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, user.MgoPrivateKey)
    if err != nil {
        return nil, err // Invalid user private key
    }
    // Step 7: Load system sponsor's private key (used for paying gas)
    sysKey, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, l.svcCtx.Config.SysMgoPrivateKey)
    if err != nil {
        return nil, err // Invalid system private key
    }
    // Step 8: Fetch the sponsor’s gas coin object (to pay gas fees)
    gasCoinObj, err := cli.MgoGetObject(ctx, request.MgoGetObjectRequest{
        ObjectId: l.svcCtx.Config.SysGasObject,
    })
    if err != nil {
        return nil, err // Failed to fetch gas coin
    }
    // Step 9: Wrap the coin object into MgoObjectRef for PTB
    gasCoin, err := transaction.NewMgoObjectRef(
        mgoModel.MgoAddress(gasCoinObj.Data.ObjectId),
        gasCoinObj.Data.Version,
        mgoModel.ObjectDigest(gasCoinObj.Data.Digest),
    )
    if err != nil {
        return nil, err
    }
    // Step 10: Create and configure the PTB transaction
    tx := transaction.NewTransaction()
    tx.SetMgoClient(cli).                     // Assign Mango client
        SetSigner(walletKey).                // Signer is the user
        SetSponsoredSigner(sysKey).          // Sponsor pays the gas
        SetSender(mgoModel.MgoAddress(walletKey.MgoAddress())). // User address
        SetGasPrice(1000).                   // Set gas price
        SetGasBudget(50000000).              // Set max gas budget
        SetGasPayment([]transaction.MgoObjectRef{*gasCoin}).    // Sponsor’s gas coin
        SetGasOwner(mgoModel.MgoAddress(sysKey.MgoAddress()))   // Sponsor pays
    
    // Step 11: Calculate raw integer amount for transfer (in nano units)
    amountDecimal := decimal.NewFromFloat(req.Amount)
    amount := amountDecimal.Mul(decimal.New(1, 9)).BigInt().Uint64() // MGO has 9 decimals
    // Step 12: Collect a large enough coin to split from
    mergeCoin, err := l.GetEnoughMgo(ctx, walletKey.MgoAddress(), amountDecimal.Mul(decimal.New(1, 9)), tx)
    if err != nil {
        return nil, err
    }
    // Step 13: Split out the exact transfer amount
    splitCoin := tx.SplitCoins(mergeCoin, []transaction.Argument{
        tx.Pure(amount), // Specify amount to extract
    })
    // Step 14: Transfer the split coin to the recipient address
    tx.TransferObjects([]transaction.Argument{splitCoin}, tx.Pure(req.ToAddress))
    // Step 15: Execute the PTB transaction and wait for confirmation
    MgoTransactionBlockResponse, err := tx.Execute(
        ctx,
        request.MgoTransactionBlockOptions{
            ShowInput:    true,
            ShowRawInput: true,
            ShowEffects:  true,
            ShowEvents:   true,
        },
        "WaitForLocalExecution",
    )
    if err != nil {
        return nil, err // Execution failed
    }
    // Step 16: Return the transaction hash
    return &types.WithdrawalResp{
        Hash: MgoTransactionBlockResponse.Digest,
    }, nil
```




