# Read data and Post data

投稿データを読み込み、選別、投稿

Google cloud runに展開、Schedulesで定期発火を行う

1. 発火
2. Google cloud firestoreからUserAccountsを取得
3. 条件選別
4. UserAccountsのPostを取得
5. 指定する各条件で投稿を選別
6. Twitter/Xに投稿