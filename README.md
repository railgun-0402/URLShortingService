# URL Shorting Service
要件はシンプル
- 長いURLを “短い一意のID” に変換
- その短いURLにアクセスすると元のURLにリダイレクトする

### 例：
```bash
元URL：
https://example.com/articles/2025/11/27/why-url-shortener-is-powerful-for-system-design

短縮URL：
https://sho.rt/aBcD12
```

## 🟥 【機能要件（Functional Requirements）】
### ✔ ① 長いURLを短いURLに変換する
- ユーザーが長いURLを送る
- サービスが短い文字列（ID）を作る
- https://sho.rt/<短いID> を返す

### ✔ ② 短縮URLにアクセスすると元URLへリダイレクト
- ブラウザが sho.rt/abc123 にアクセス
- DB or キャッシュから元URLを検索
- 302 or 301 リダイレクトを返す

### ✔ ③ 有効期限を設定（検討中）
- 1時間
- 7日
- 永久保存
etc、用途に応じて TTL を決める

### ✔ ④ URLクリック数の記録（Analytics）
- 何回クリックされたか
- ユーザーの地域
- ブラウザ・デバイス
- 参照元（referrer）
etc.

Slack や Twitter が裏で使ってる理由がこの辺りにあるらしい

### ✔ ⑤ カスタム短縮URLをサポート

### ✔ ⑥ URLのバリデーション
- URLが正しい形式かチェック
- マルウェア / フィッシングURLは拒否
- 1秒で終わる必要あり（UX的に）

## 🟦 【非機能要件（Non-Functional Requirements）】

URL短縮は機能より 非機能要件 が鬼のように重要。

### ✔ ① 爆速である（レイテンシが超重要）
- URL短縮 → クリックした瞬間ページに飛びたい(遅いと使い物にならない)

### ✔ ② スケーラブル（高トラフィックに耐える）
Twitter や LINE と連携すると
1秒間に数十万リダイレクトが来る。

### ✔ ③ 高可用性（落ちたらサービスとして終わる）
- 短縮URLが死ぬとリンクが全部死ぬ

### ✔ ④ 一意性（短縮IDに衝突が起きてはならない）
- ID重複したら大事故

### ✔ ⑤ コスト効率（超大量アクセスに耐える必要）
- クリックイベントは激重（PVが爆発するなど仮定する）

  
