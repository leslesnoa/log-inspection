package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	dateFormat          = "20060102150405"
	timeoutStr          = "-"
	trialCount          = 3
	confirmCount        = 3
	averageResponseTime = 1000
)

var filepath = "testlog1.txt"

var failedServer = make(map[string]*FailedServer)

type FailedServer struct {
	ServerIP    string
	FailedTime  time.Time
	RecoverTime time.Time
	FailedCount int32
}

type Result struct {
	FailedHost string
	FailedSpan time.Duration
}

type OverLoadServer struct {
	ServerIP string
	Span     time.Duration
}

var res []Result
var overLoadServers = make(map[string]*OverLoadServer)

type ResponseResult struct {
	ServerIP         string
	AverageResponses []Response
}

type Response struct {
	RecordTime time.Time
	Result     string
}

func main() {

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	/* csvリーダーを生成 */
	r := csv.NewReader(f)

	var rr ResponseResult

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
			// 故障リストに該当IPが無ければ追加する
			if _, ok := failedServer[serverIP]; !ok {
				failedServer[serverIP] = &FailedServer{
					ServerIP:    serverIP,
					FailedTime:  confirmTime,
					FailedCount: 1,
				}
			} else {
				// 故障リストに該当IPが存在する場合、試行回数を1追加する
				failedServer[serverIP].FailedCount += 1
			}
		} else {
			if _, ok := failedServer[serverIP]; ok {
				failedServer[serverIP].CheckFailedServer(confirmTime, trialCount)
			} else {
				// サーバから応答がある場合は、各サーバの過負荷状態をチェックする
				rr.CheckAverageResponse(serverIP, serverResponse, confirmTime, confirmCount, averageResponseTime)
			}
		}
	}

	/* 故障サーバ名、故障期間を出力する */
	for _, s := range res {
		fmt.Printf("故障サーバー: %s 故障期間: %s\n", s.FailedHost, s.FailedSpan)
	}
	/* 過負荷状態となっているサーバ名、期間を出力する */
	for _, r := range overLoadServers {
		fmt.Printf("過負荷サーバー: %s 過負荷期間: %s\n", r.ServerIP, r.Span)
	}
}

/* 監視ログのフォーマットで日時をTime型に変換する */
func stringToTime(str string) time.Time {
	t, _ := time.Parse(dateFormat, str)
	return t
}

func (f *FailedServer) CheckFailedServer(c time.Time, count int32) {
	/* N回以内にレスポンスがあった場合、故障リストから除外する */
	if f.FailedCount < trialCount {
		delete(failedServer, f.ServerIP)
	}
	if f.FailedCount >= count {
		bt := c.Sub(f.FailedTime)
		res = append(res, Result{
			FailedHost: f.ServerIP,
			FailedSpan: bt,
		})
	}
}

func (r *ResponseResult) CheckAverageResponse(ip string, res string, ct time.Time, c int, t int) {

	r.ServerIP = ip
	r.AverageResponses = append(r.AverageResponses, Response{
		RecordTime: ct,
		Result:     res,
	})

	if len(r.AverageResponses) < c {
		return
	}

	sl := len(r.AverageResponses) - c
	el := len(r.AverageResponses)
	inspectSlice := r.AverageResponses[sl:el]

	average := CalculateAverageResponse(inspectSlice)

	if average > t {
		st := inspectSlice[0].RecordTime
		et := inspectSlice[c-1].RecordTime
		if _, ok := overLoadServers[ip]; !ok {
			overLoadServers[ip] = &OverLoadServer{
				ServerIP: r.ServerIP,
				Span:     et.Sub(st),
			}
		} else {
			overLoadServers[ip].Span = et.Sub(st)
		}
	}
}

func CalculateAverageResponse(s []Response) int {
	var num int
	for _, t := range s {
		pt, _ := strconv.Atoi(t.Result)
		num += pt
	}
	return num
}
