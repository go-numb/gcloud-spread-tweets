# Spread to twitter/x with Google cloud services

├── cloud_functions/                 # Google Cloud Functions用のコード  
│   ├── scheduler_task/              # 定期実行タスク用のディレクトリ  
│   │   └── main.go                  # Cloud Functions用の関数  
│   └── ...                          # 他のCloud Functions関数がある場合  
├── cloud_run/                       # Cloud Run用のAPIコード  
│   ├── cmd/                         # エントリーポイント（mainパッケージ）  
│   │   └── api/                     # APIアプリケーション用  
│   │       └── main.go              # Cloud Runにデプロイされるアプリケーションのメイン  
│   └── pkg/                         # 再利用可能なパッケージやロジック  
│       ├── handlers/                # HTTPハンドラー  
│       ├── middleware/              # ミドルウェア  
│       └── ...                      # その他のパッケージ  
└── web_interface/                   # Webインターフェイス (Reactアプリ)用  
    ├── dist/                        # Reactビルドアウトプット (GCSにアップロードする)  
    └── src/                         # Reactソースファイル  
        ├── components/              # Reactコンポーネント  
        └── ...                      # その他のReactソース  


## Usage
- Google cloud service
  - Run         : API
  - Functions   : Post実行
  - Schedules   : 定期通知
  - Storage     : 静的ホスティング

## Author

[@_numbP](https://twitter.com/_numbP)

## License

[MIT](https://github.com/go-numb/gcloud-spread-tweets/blob/master/LICENSE)