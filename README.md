YARIMASUNE

PUHAAAAAAAAAAAA!

## ビルドと実行 (v4.2)

```bash
go build -ldflags="-s -w" -trimpath
```

Windows (PowerShell / cmd):

```powershell
.\yajuws.exe
```

macOS / Linux:

```bash
./yajuws
```

## バージョン履歴

### v4.2
- ターミナル枠の描画処理を共通化し、全角対応の幅計算を導入
- ファイルマネージャーの列揃えを実表示幅ベースに調整
- gopsutil v4系へ移行

### v4.1
- 起動シーケンス (Enterでスキップ可)
- 便利ツール (時計/タイマー/稼働時間/じゃんけん)
- 語録履歴 + お気に入り管理
- テーマ/カラー/高速起動の設定
- タスクマネージャー (CPU使用率/プロセス数/メモリ)
### v4.0
go言語で書き直した
