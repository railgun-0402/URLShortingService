# Admin: Get Shorten URL by ID API

管理用（認可あり）の短縮URL取得API。
管理画面・ログイン済みユーザー向けに、短縮URLのメタ情報をJSONで返す。

---

## Purpose / Responsibilities

- 管理用途として、短縮URLを `id` 指定で取得する
- 管理画面で表示・編集・監査を行う前提
- リダイレクト用途のAPI（Public）とは責務を分離する

---

## Characteristics

- Response: JSON（構造化データ）
- 認証・認可が必須
- パフォーマンスよりも以下を優先
    - 認可（誰のURLか）
    - 監査（誰が閲覧/変更したかを追跡）
    - PII取り扱い（ログに機微情報を漏らさない）

---

## Endpoint

- `GET /admin/urls/{id}`

---

## Authentication / Authorization

### Authentication
- 例：Bearer JWT / Session（実装に合わせる）

### Authorization Rule
- `url.tenant_id == actor.tenant_id` のみ取得可能

### Authorization Failure Handling (Security Policy)
- 認可NGの場合も **404** を返す（存在を推測されないようにする）
    - `NOT_FOUND` は「存在しない」または「権限がない」の両方を含む
    - 管理画面側は 404 を「閲覧不可 or 存在しない」として扱う

---

## Response

### Success (200)
- 期限切れでも管理用途では取得可能（`is_expired: true` を返す）
- リダイレクトAPI側で 410（Gone）にするのが責務

```json
{
  "id": "aBcD12eF",
  "original_url": "https://example.com/foo/bar?param=123",
  "short_url": "https://sho.rt/r/aBcD12eF",
  "expires_at": "2026-02-01T00:00:00Z",
  "is_expired": false,
  "created_at": "2026-01-03T12:00:00Z",
  "created_by": "user_123",
  "tenant_id": "tenant_abc"
}
```

## Error Model

```json
{
  "error": {
    "code": "SHORT_URL_NOT_FOUND",
    "message": "short url not found",
    "request_id": "req-01HRXXXX"
  }
}
```
## Error Codes

### SHORT_URL_NOT_FOUND (404)
- idが存在しない、または認可NG

### UNAUTHORIZED (401)
- 認証情報がない/無効

### INTERNAL (500)

- 想定外のエラー

## Example
### Get URL (Success)

---

```bash
curl -i http://localhost:8080/admin/urls/aBcD12eF
```

```json
{
"id": "aBcD12eF",
"original_url": "https://example.com/foo/bar?param=123",
"short_url": "https://sho.rt/r/aBcD12eF",
"expires_at": "2026-02-01T00:00:00Z",
"is_expired": false,
"created_at": "2026-01-03T12:00:00Z",
"created_by": "user_123",
"tenant_id": "tenant_abc"
}
```

### Not Found / Forbidden (Both -> 404)

---

```bash
curl -i http://localhost:8080/admin/urls/notfound123
```

```json
{
    "error": {
    "code": "SHORT_URL_NOT_FOUND",
    "message": "short url not found",
    "request_id": "req-01HRXXXX"
    }
}
```

### Expired URL (Still 200 in Admin API)

---

```bash
curl -i http://localhost:8080/admin/urls/expired123
```

```json
{
"id": "expired123",
"original_url": "https://example.com/foo/bar?param=123",
"short_url": "https://sho.rt/r/expired123",
"expires_at": "2026-01-01T00:00:00Z",
"is_expired": true,
"created_at": "2025-12-01T12:00:00Z",
"created_by": "user_123",
"tenant_id": "tenant_abc"
}
```

## Audit Log
### Purpose

- 誰が・いつ・何を・どこから・どうした を追跡
- 不正調査・問い合わせ対応・運用に利用する

### Logged Actions

- SHORT_URL_READ : GET /admin/urls/{id}

- SHORT_URL_UPDATE : PATCH /admin/urls/{id}

- SHORT_URL_DELETE : DELETE /admin/urls/{id}

### Recommended Fields

- timestamp

- action

- actor_user_id

- actor_tenant_id

- target_short_url_id

- result : SUCCESS / DENIED / NOT_FOUND

- request_id

- ip_hash（生IPは保存しない）

- user_agent（保持とサイズに注意）

- diff（UPDATE時のみ。PIIマスキング必須）

### Example (Audit Log JSON)

```json

{
"timestamp": "2026-01-04T01:23:45Z",
"action": "SHORT_URL_READ",
"actor_user_id": "user_123",
"actor_tenant_id": "tenant_abc",
"target_short_url_id": "aBcD12eF",
"result": "SUCCESS",
"request_id": "req-01HRXXXX",
"ip_hash": "sha256:....",
"user_agent": "Mozilla/5.0 ..."
}
```

## PII / Security

---

### Policy

- original_url はクエリ文字列に機微情報が含まれる可能性があるため、
アプリログ/監査ログに生値を出力しない

- 必要がある場合は正規化して出す（クエリを除去）

### Recommendations

- original_url のログ出力が必要なら normalized_url（query除去）を使用

- ip は生値ではなく ip_hash を使用

- request_id をレスポンス・ログ双方に含めて相関できるようにする

### Normalization Example

- https://example.com/foo?token=abc&email=a@b.com
-> https://example.com/foo
