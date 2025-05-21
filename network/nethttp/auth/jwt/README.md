# JWT的三元结构

```go
header.payload.signature
┌────────┐ ┌───────────┐ ┌───────────┐
| base64 | | base64    | | base64    |
| encoded| | encoded   | | encoded   |
| JSON   | | JSON      | | signature |
└────────┘ └───────────┘ └───────────┘
```

```
Header: {"alg":"HS256", "typ":"JWT"}
Payload: {"sub":"123","name":"bob","exp":1747832913}
Signature: HMACSHA256(base64(header)+"."+base64(payload), secret)
```
