# JWT-WITH-GO

# What Is JWT?
JWT or JSON web token is a digitally signed string used to securely transmit information between parties. It's an RFC7519 standard. A JWT consists of three parts: header. payload.signature

Below is a sample JWT.
```eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c```


# flow
![Flow](flow.png)

# claims
The header is a Base64-encoded string and it contains the token type (JWT in this case) and the signing algorithm (HMAC SHA256 in this case, or HS256 for short).
```{
  "alg": "HS256",
  "typ": "JWT"
}
```

The payload is a Base64-encoded string that contains claims. Claims are a collection of data related to the user and the token itself. Example claims are exp (expiration time), it (issued at), name (user name), and sub (subject).
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

1. main(): ```Contain all the routes```
2. init():```Each time we run main.go redis will automatically connect```
3. login() and newlogin() : ``` It takes the credentials and checks from databases and logs them if the credentials are valid.```
4. logout() : ```Extract the JWT metadata. If true then delete the metadata, and so JWT invalid immediately.```
5. middleware() : ```two routes that require authentication: /login and /logout as anybody can access that, Middleware will secure these routes```



# Redis
It has a feature that automatically deletes data whose expiration time has reached. Redis can also handle a lot of writes and can scale horizontally.
Since Redis is key-value storage, its keys need to be unique, to achieve this, use UUID as the key and use the user id as the value. Also, the UUID is unique each time it is created, a user can create more than one token. This happens when a user is logged in on different devices. The user can also log out from any of the devices without them being logged out from all devices.
If the expiration time is reached for either the refresh token or the access token, the JWT is automatically deleted from Redis.

# Note: Install Redis if is not installed

# Access Token
Generating JWT Token which is valid for 15 minutes Some Implementation Loopholes :
1. user can log in, then decide to log out immediately, but the user's JWT remains valid until the expiration time is reached.
2. The JWT might be hijacked and used by a hacker without the user doing anything about it until the token expires.
3. The user will need to re-login after the token expires, thereby leading to a poor user experience.

# Refresh Token
A refresh token has a longer lifespan, usually 7 days. This token is used to generate new access and refresh tokens. In the event that the access token expires, new sets of access and refresh tokens are created when the refresh token route is hit (from our application), handles the above problem well, and also uses persistence storage(i.e. Redis) layer to store the JWT metadata. This will enable us to invalidate a JWT the very second the user logs out, thereby improving security.

1. Get refresh token from the headers. Verify the signing method of the token. Check it is still valid or not.
3. The refresh_uuid and the user_id are then extracted, which are metadata used as claims when creating the refresh token.
4. Search for the metadata in redis store and delete it using the refresh_uuid as key.
5. Create a new pair of access and refresh tokens that will now be used for future requests.
6. The metadata of the access and refresh tokens are saved in redis. The created tokens are returned to the caller.
7. Else, if the refresh token is not valid, the user will not be allowed to create a new pair of tokens. We will need to relogin to get new tokens.

