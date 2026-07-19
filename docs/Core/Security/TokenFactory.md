# TokenFactory Security Component Specification

**Last Updated:** July 18, 2026  
**Author:** Ismael Romero

---

# 1. Introduction & Responsibility

The `TokenFactory` is a core component within the security subsystem responsible for
the generation, signing, and validation of JSON Web Tokens (JWT).

Its primary responsibility is to provide a secure abstraction over JWT cryptographic
operations while keeping authentication workflows independent from token implementation
details.

The component is responsible for:

- Generating signed JWT tokens.
- Applying standard JWT registered claims.
- Validating JWT signatures.
- Validating token expiration and registered claims.
- Extracting trusted token payload information.

The `TokenFactory` does not manage:

- User authentication.
- User accounts.
- Authorization decisions.
- Persistence operations.
- Session storage.
- Token revocation.
- Refresh token rotation.
- MFA verification workflows.

Token purpose and lifetime are defined by higher-level authentication services. The
`TokenFactory` only provides the cryptographic mechanism required to create and validate
tokens.

The architectural boundary can be summarized as:

```text
Authentication Service
          |
          |
          v
     TokenFactory
          |
          |
          v
 Signed JWT Generation / Validation
```

The responsibility of the component ends once the token has been generated or validated.

---

# 2. Design & Architecture

The `TokenFactory` follows the security subsystem architectural principles:

- Stateless operation.
- Cryptographic isolation.
- Interface-based abstraction.
- Constructor dependency injection.
- Independent testability.

The component does not store generated tokens or authentication state.

---

# 2.1 Interface-Based Abstraction

The component exposes the `TokenFactory` interface as its public contract.

Example:

```go
type TokenFactory interface {
    GenerateToken(
        userID string,
        duration time.Duration,
    ) (string, error)

    ValidateToken(
        tokenString string,
    ) (*TokenClaims, error)
}
```

Higher-level services depend on the interface rather than the concrete implementation.

This provides:

- Easier unit testing.
- Reduced coupling.
- Ability to replace the underlying JWT implementation.

---

# 2.2 Dependency Injection

The component integrates with Uber Fx through constructor-based dependency injection.

The constructor receives the application configuration dependency:

```go
func NewTokenFactory(cfg *config.Config) TokenFactory
```

Uber Fx resolves the required dependencies automatically through the application
dependency graph.

The security module registers the component using:

```go
var Module = fx.Options(
    fx.Provide(
        factory.NewPasswordFactory,
        factory.NewTokenFactory,
        factory.NewMFAFactory,
    ),
)
```

Uber Fx is responsible for:

- Resolving constructor dependencies.
- Creating component instances.
- Managing the dependency graph lifecycle.

The `TokenFactory` does not receive:

- User repositories.
- Authentication services.
- Persistence dependencies.
- MFA services.

Its only external dependency is the application configuration required to obtain the
JWT signing secret.

---

# 3. Token Structure & Claims

Generated tokens contain user identity information through the `TokenClaims`
structure.

The structure extends the standard JWT claims provided by the JWT library.

---

# 3.1 TokenClaims

The custom claims include:

## Custom Claims

### `user_id`

Stores the unique identifier of the authenticated user.

Example:

```json
{
  "user_id": "12345"
}
```

---

## Registered JWT Claims

The component automatically manages standard JWT security claims:

| Claim | Description |
|---|---|
| `exp` | Token expiration timestamp |
| `iat` | Token issued-at timestamp |
| `nbf` | Time before which the token must not be accepted |

These claims ensure that generated tokens have a controlled lifetime and cannot be used
outside their intended validity period.

---

# 4. Public API

The `TokenFactory` exposes two primary operations:

- JWT generation.
- JWT validation.

---

# 4.1 GenerateToken

## Signature

```go
GenerateToken(
    userID string,
    duration time.Duration,
) (string, error)
```

## Description

Creates a signed JWT containing the authenticated user's identity.

The method performs:

1. Validates that the provided user identifier is not empty.
2. Creates the JWT claims structure.
3. Calculates issuance and expiration timestamps.
4. Signs the token using HS256.
5. Returns the encoded JWT string.

---

## Input

### userID

The unique identifier of the authenticated user.

The operation rejects empty identifiers because a token without an identity reference
cannot represent an authenticated principal.

---

### duration

The requested token lifetime.

The lifetime is defined by the authentication workflow.

Examples:

```text
Access Token:
1 hour

Temporary Authentication Token:
3 minutes
```

---

# 4.2 ValidateToken

## Signature

```go
ValidateToken(
    tokenString string,
) (*TokenClaims, error)
```

## Description

Validates a JWT received from an external caller and confirms its authenticity.

The method performs:

1. JWT parsing.
2. Signing algorithm validation.
3. Cryptographic signature verification.
4. Registered claims validation.
5. Claims extraction.

If every validation succeeds, the method returns the trusted `TokenClaims`.

---

## Cryptographic Validation

The expected signing algorithm is:

```text
HS256 (HMAC-SHA256)
```

The validation process rejects tokens using unexpected signing methods.

This prevents algorithm confusion attacks where an attacker attempts to use an
unintended cryptographic algorithm.

---

# 5. Error Handling

Errors generated by `TokenFactory` follow the security module convention:

```text
security/factory:
```

| Error | Description |
|---|---|
| `ErrInvalidToken` | Returned when the token is invalid or claims cannot be trusted. |
| `ErrEmptyUserID` | Returned when token generation receives an empty user identifier. |
| `ErrTokenGeneration` | Returned when JWT signing fails internally. |
| `ErrUnexpectedSigningMethod` | Returned when the token uses an unsupported signing algorithm. |
| `ErrTokenParsing` | Returned when the JWT structure cannot be parsed. |

---

# 6. Architectural Integration Flow

The `TokenFactory` is consumed by authentication services after successful identity
verification.

The component remains isolated from authentication decisions and only performs JWT
operations.

```mermaid
sequenceDiagram
    autonumber

    participant Auth as Authentication Service
    participant Factory as TokenFactory
    participant Client

    Note over Auth,Factory: Token Generation Flow

    Auth->>Factory: GenerateToken(userID, duration)

    Factory->>Factory: Build JWT Claims
    Factory->>Factory: Sign JWT using HS256

    Factory-->>Auth: Signed JWT Token

    Auth-->>Client: Return Authentication Token


    Note over Auth,Factory: Token Validation Flow

    Client->>Auth: Send JWT Token

    Auth->>Factory: ValidateToken(token)

    Factory->>Factory: Validate Signature
    Factory->>Factory: Validate Signing Algorithm
    Factory->>Factory: Validate Registered Claims

    Factory-->>Auth: Token Claims

    Auth->>Auth: Continue Protected Operation
```

---

# 7. Security Requirements

## 7.1 Secret Protection

The JWT signing secret must:

- Never be hardcoded.
- Never be logged.
- Never be exposed through API responses.
- Be managed through secure configuration.

---

## 7.2 Token Lifetime Control

Every generated token must include an explicit expiration timestamp.

Token lifetime decisions belong to authentication workflows.

Example:

| Token Type | Lifetime |
|---|---|
| Access Token | Defined by authentication policy |
| Refresh Token | Defined by authentication policy |
| Temporary Authentication Token | Defined by authentication policy |

---

## 7.3 Signature Verification

Every received token must be validated before trusting its claims.

The system must never rely solely on decoded JWT payload information without verifying
the cryptographic signature.

---

## 7.4 Stateless Operation

The `TokenFactory`:

- Does not persist tokens.
- Does not store sessions.
- Does not track token usage.
- Does not maintain revocation state.

Token revocation, rotation, and lifecycle policies belong to higher-level services.

---

# 8. Design Principles Summary

The `TokenFactory` follows these architectural principles:

- **Single Responsibility:** Handles JWT generation and validation only.
- **Interface Driven Design:** Provides abstraction for testing and maintenance.
- **Constructor Dependency Injection:** Integrates with Uber Fx through registered constructors.
- **Cryptographic Isolation:** Encapsulates signing and verification details.
- **Secure Defaults:** Enforces expiration handling and signing algorithm validation.
- **Stateless Operation:** Does not persist authentication state.