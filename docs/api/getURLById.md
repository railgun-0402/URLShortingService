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

- /admin/urls/{id}
- GET Method

## GET /admin/urls/:id

---

### request

```bash
curl -i http://localhost:8080/admin/urls/aBcD12eF
```

### result
```json
{
  "id": "aBcD12eF",
  "original_url": "https://example.com/foo/bar?param=123",
  "short_url": "https://sho.rt/r/aBcD12eF",
  "expires_at": "2026-02-01T00:00:00Z",
  "created_at": "2026-01-03T12:00:00Z",
  "created_by": "user_123"
}
```

## Not Exist ID

---

### request

```bash
curl -i http://localhost:8080/admin/urls/notfound123
```

### result
```json
{
  "error": {
    "code": 404,
    "message": "short url not found",
    "request_id": "..."
  }
}
```

## Expire URL

---

### request

```bash
curl -i http://localhost:8080/admin/urls/notfound123
```

### result
```json
{
  "error": {
    "code": 404,
    "message": "short url not found",
    "request_id": "..."
  }
}
```



## Auth Check

---

### 仕様

- 認証方式（例：Bearer JWT / Session）
- 認可ルール
  - url.tenant_id == user.tenant_id のみ取得可
  - それ以外は403にする
