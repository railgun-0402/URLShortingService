# GET Shorten URL By ID API

## 概要

- 管理用としてIDを使用した短縮URLの取得API
- 管理画面/ログイン済みのUser等で実施する想定

### 仕様

- JSON（構造化データ）
- クライアントはレスポンスを「解釈」して画面に表示
- 管理画面等で実施する想定なので、スケール・パフォーマンスよりも以下を重視
  - 認可(誰のURL？)
  - 個人情報(作成者・設定・元URLの管理情報など)を管理
  - ログ(Reader・Change User)

---

## Endpoint

- /{id}
- GET Method

## GET /:id

---

### request

```bash
curl -i http://localhost:8080/aBcD12eF
```

### result
```bash
Status Code: 200
URL: https://example.com/foo/bar?param=123
```

## Not Exist ID

---

### request

```bash
curl -i http://localhost:8080/notfound123
```

### result
```bash
HTTP/1.1 404 Not Found
{"message":"short url not found"}
```
## empty URL

---

### request

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": ""}'
```

### result
```bash
HTTP/1.1 400 Bad Request
{"message":"url is required"}
```

## irregular URL

---

### request
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "hogehoge"}'
```

### result
```bash
HTTP/1.1 400 Bad Request
{"message":"invalid url format"}
```
