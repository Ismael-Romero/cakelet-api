# Authentication Flow
**Fecha de actualización**: 18 de Julio de 2026. <br/>
Actualizado por: Ismael Romero.
-----

## Introduction
Cakelet has been designed with security and access control as fundamental principles. 
For this reason, every interaction with the platform begins with the authentication 
process and the validation of user credentials.

Although the documentation in this repository is primarily focused on the API,
it is equally important for user interface components to understand the architecture
and workflows that underpin the platform. The business logic forms the core of the system
and defines the behavior of every interaction; consequently, the user interface acts as a consumer of this logic,
without which it would be impossible to ensure consistent, secure, and predictable system behavior.

## 1. Primary Authentication Flow

### 1.1 Access Request

The client submits the user's credentials (username and password) through the login form to the authentication endpoint.

### 1.2 Process Initiation

The authentication endpoint receives the request and initiates the identity validation process by querying the user information stored in the database.

### 1.3 User Existence Validation

The system verifies whether the provided username (or unique identifier) exists.

#### If the user does not exist

- Record a failed authentication attempt in the access log.
- Terminate the authentication flow.
- Return a generic authentication error to the client to prevent user enumeration.

### 1.4 Failed Attempts Lock Validation

The system verifies whether the account has been temporarily locked because the maximum number of failed login attempts has been exceeded.

### 1.5 Administrative Lock Validation

The system verifies whether the account has been explicitly locked by a system administrator.

#### If the account is locked

- Record a failed authentication attempt in the access log.
- Terminate the authentication flow.
- Return an account locked response to the client.

### 1.6 Password Verification

The submitted password is processed using the corresponding password hashing algorithm and compared against the securely stored password hash.

#### If the password is incorrect

- Record a failed authentication attempt.
- Increment the user's failed login attempts counter by one.
- If the counter reaches the configured threshold, update the user's status to **Locked**.
- Terminate the authentication flow.
- Return an authentication error to the client.

### 1.7 Two-Factor Authentication (2FA) Evaluation

The system verifies whether the user has Two-Factor Authentication enabled.

#### If 2FA is enabled

- Generate a cryptographically secure verification code.
- Send the verification code to the user's registered email address.
- Generate a signed Temporary Transition Token (JWT) with:
    - Maximum lifetime of **3 minutes**.
    - Embedded `user_id` claim.
- Return the temporary token to the client.
- Request completion of the 2FA challenge.
- Continue with the process described in **Section 2: 2FA Validation Flow**.

### 1.8 Successful Authentication (Without 2FA)

If the credentials are valid and 2FA is not enabled:

- Record a successful authentication event.
- Reset the failed login attempts counter to **0**.
- Retrieve the user's profile information:
    - `user_id`
    - First name
    - Last name
    - Nickname
    - Profile image URL
- Generate an **Access Token** with a lifetime of **1 hour**.
- Generate and persist a **Refresh Token** with a lifetime of **1 hour and 10 minutes**.
- Construct the JSON response containing:
    - Access Token
    - Refresh Token
    - User profile
- Return **HTTP 200 OK** to the client.

---

## 2. 2FA Validation Flow

### 2.1 Verification Code Reception

The user receives the verification code in their registered email inbox.

### 2.2 Challenge Submission

The user enters the verification code into the client application.

The frontend submits a request to the 2FA endpoint including:

- Verification code
- Temporary Transition Token (JWT)

### 2.3 Temporary Token Validation

The system:

- Decodes the Temporary Transition Token.
- Verifies its cryptographic signature.
- Confirms that the token has not expired (maximum validity: **3 minutes**).

### 2.4 Verification Code Validation

The submitted verification code is compared against the code generated for the corresponding authentication session.

### 2.5 Challenge Resolution

#### Successful Validation

If every validation succeeds:

- Record a successful authentication event.
- Reset the failed login attempts counter to **0**.
- Retrieve the user's profile:
    - `user_id`
    - First name
    - Last name
    - Nickname
    - Profile image URL
- Generate an **Access Token** (1 hour).
- Generate and persist a **Refresh Token** (1 hour and 10 minutes).
- Return a successful authentication response containing:
    - Access Token
    - Refresh Token
    - User profile

#### Failed Validation

If any validation fails (invalid code or expired token):

- Record a failed authentication attempt.
- Increment the failed login attempts counter by one.
- If the accumulated threshold is reached, lock the user account.
- Return the appropriate error response, for example:
    - `Invalid verification code`
    - `Verification session has expired`