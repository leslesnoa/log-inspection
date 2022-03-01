package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	dateFormat = "20060102150405"
	timeoutStr = "-"
)

var filepath = "testlog1.txt"

type FailedServer struct {
	ServerIP    string
	FailedTime  time.Time
	RecoverTime time.Time
}

type Result struct {
	FailedHost string
	FailedSpan time.Duration
}

var res []Result

func main() {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	/* csvリーダーを生成 */
	r := csv.NewReader(f)

	failedServer := make(map[string]*FailedServer)

	for {
		/* 監視ログを行ごとに読み込む */
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(record) != 3 {
			log.Println("Error log file format is invalid.")
			os.Exit(1)
		}

		/* 故障サーバを抽出・格納 */
		confirmTime := stringToTime(record[0])
		serverIP := record[1]
		serverResponse := record[2]

		if serverResponse == timeoutStr {
			// FailedServerに該当IPが無ければ追加する
			if _, ok := failedServer[serverIP]; !ok {
				failedServer[serverIP] = &FailedServer{
					ServerIP:   serverIP,
					FailedTime: confirmTime,
				}
			}
		} else {
			if _, ok := failedServer[serverIP]; ok {
				failedServer[serverIP].SetFailedServerToResult(confirmTime)
			}
		}
	}

	/* 故障サーバ名、故障期間を出力 */
	for _, s := range res {
		fmt.Printf("故障サーバ: %s 故障期間: %s\n", s.FailedHost, s.FailedSpan)
	}
}

/* 監視ログのフォーマットで日時をTime型に変換する */
func stringToTime(str string) time.Time {
	t, _ := time.Parse(dateFormat, str)
	return t
}

func (f *FailedServer) SetFailedServerToResult(c time.Time) {
	bt := c.Sub(f.FailedTime)
	res = append(res, Result{
		FailedHost: f.ServerIP,
		FailedSpan: bt,
	})
}
