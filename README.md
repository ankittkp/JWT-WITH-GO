# JWT-WITH-GO

# What Is JWT?
JWT or JSON web token is a digitally signed string used to securely transmit information between parties. It’s an RFC7519 standard.
A JWT consists of three parts:
header.payload.signature
Below is a sample JWT.
```eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c```

# claims
The header is a Base64-encoded string and it contains the token type (JWT in this case) and the signing algorithm (HMAC SHA256 in this case, or HS256 for short).
```{
  "alg": "HS256",
  "typ": "JWT"
}
```

The payload is a Base64-encoded string that contains claims. Claims are a collection of data related to the user and the token itself. Example claims are: exp (expiration time), iat (issued at), name (user name), and sub (subject).
```{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022
}
```
The signature is a signed string. For HMAC signing algorithms, we use the Base64-encoded header, the Base64-encoded payload, and a signing secret to create it.
```HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
```
