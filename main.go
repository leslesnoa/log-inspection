package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// var breakServer string
// var recoverTime string
// var breakStartTime string
// var breakServers []string

const (
	dateFormat = "20060102150405"
	timeoutStr = "-"
)

type FailedServer struct {
	FailedTime  string
	RecoverTime string
	IsBreak     bool
}

type Result struct {
	FailedHost string
	FailedSpan time.Duration
}

func main() {
	f, err := os.Open("log.txt")
	if err != nil {
		log.Fatal(err)
	}

	/* csvリーダーを生成 */
	r := csv.NewReader(f)

	failedServer := make(map[string]*FailedServer)
	var res []Result

	for {
		/* 監視ログを行ごとに読み込む */
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		/* 故障サーバを抽出・格納 */
		confirmTime := record[0]
		serverIP := record[1]
		serverResponse := record[2]

		if serverResponse == timeoutStr {
			if _, ok := failedServer[serverIP]; ok {
				// failedServer[server].BreakCount += 1
			} else {
				failedServer[serverIP] = &FailedServer{
					FailedTime: confirmTime,
					IsBreak:    true,
				}
			}
		} else {
			if _, ok := failedServer[serverIP]; ok {
				// fmt.Printf("key exists. The value is %#v", val)
				st := stringToTime(failedServer[serverIP].FailedTime)
				et := stringToTime(confirmTime)
				bt := et.Sub(st)
				res = append(res, Result{
					FailedHost: serverIP,
					FailedSpan: bt,
				})
			}
		}
	}

	/* 故障サーバ名、故障期間を出力 */
	for _, s := range res {
		fmt.Printf("故障サーバ: %s 故障期間: %s\n", s.FailedHost, s.FailedSpan)
	}
}

func stringToTime(str string) time.Time {
	t, _ := time.Parse(dateFormat, str)
	return t
}
