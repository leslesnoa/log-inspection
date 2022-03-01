# プログラム実行方法

## 1.リポジトリをクローンする
```
$ https://github.com/leslesnoa/log-inspection.git
```

## 2.実行する設問のディレクトリに移動する
例：
```
$ cd question1
``` 

## 3.main.goを編集してパラメータを設定する
- trialCount: 故障とみなす試行回数
- confirmCount: 直近N回の過負荷確認回数
- averageResponseTime: 過負荷状態とみなす数値(miliSecond)

## 4.プログラムを実行する
```
$ go run main.go
```

# テスト実行方法

## テストする設問のディレクトリへ移動
```
$ cd qustion 1
```

## テスト実行
```
$ go test
```