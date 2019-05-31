package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

const (
	argSOURCE = "source"
	argDICT   = "dict"
	argTARGET = "target"
)

type mDict struct {
	Name string
	Code string
}

func readDict(fname string) ([]mDict, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	var lines = []mDict{}
	for _, ln := range bytes.Split(data, []byte("\n")) {
		if len(ln) == 0 {
			continue
		}
		var data = strings.Split(string(ln), " ")
		var o = mDict{
			Name: data[0],
			Code: data[1],
		}
		lines = append(lines, o)
	}
	return lines, nil
}

type mSource struct {
	Name string
	Code string
	Cash float32
}

func readSource(fn string) ([]mSource, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	var lines = []mSource{}
	for _, ln := range bytes.Split(data, []byte("\n")) {
		if len(ln) == 0 {
			continue
		}
		var vs = strings.Split(string(ln), "\t")
		if len(vs) == 3 {
			var o = mSource{
				Name: vs[0],
				Code: vs[2],
			}
			v, err := strconv.ParseFloat(vs[1], 32)
			if err == nil {
				o.Cash = float32(v)
			}
			lines = append(lines, o)
		}
	}
	return lines, nil
}

type mTarget struct {
	Code string
	Name string
	Cash float32
}

func readTarget(fn string) ([]mTarget, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	var lines = []mTarget{}
	for _, ln := range bytes.Split(data, []byte("\n")) {
		var vs = strings.Split(string(ln), "\t")
		if len(vs) == 3 {
			var o = mTarget{
				Code: vs[0],
				Name: vs[1],
			}
			v, err := strconv.ParseFloat(vs[2], 32)
			if err == nil {
				o.Cash = float32(v)
			}
			lines = append(lines, o)
		}
	}
	return lines, nil
}

type mKeyFloat struct {
	Key   string
	Value float32
}

func appAction(ctx *cli.Context) {
	var dicts, err = readDict(ctx.String(argDICT))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dicts)
	var sdata, err1 = readSource(ctx.String(argSOURCE))
	if err1 != nil {
		log.Fatal(err1)
	}
	var sdata1 = []mKeyFloat{}
	for _, ss := range sdata {
		for _, d := range dicts {
			if ss.Code == d.Name {
				var m = mKeyFloat{
					Key:   d.Code + ss.Name,
					Value: ss.Cash,
				}
				var find = false
				for _, me := range sdata1 {
					if me.Key == m.Key {
						me.Value += m.Value
						find = true
						break
					}
				}
				if !find {
					sdata1 = append(sdata1, m)
				}
			}
		}
	}
	log.Println(sdata1, len(sdata1))
	var tdata, err2 = readTarget(ctx.String(argTARGET))
	if err1 != nil {
		log.Fatal(err2)
	}
	var tdata1 = []mKeyFloat{}
	for _, tt := range tdata {
		var m = mKeyFloat{
			Key:   tt.Code + tt.Name,
			Value: tt.Cash,
		}
		var find = false
		for _, me := range tdata1 {
			if me.Key == m.Key {
				me.Value += m.Value
			}
		}
		if !find {
			tdata1 = append(tdata1, m)
		}
	}
	for _, ee := range sdata1 {
		find := false
		for _, ee1 := range tdata1 {
			if ee.Key == ee1.Key && ee.Value == ee1.Value {
				find = true
				break
			}
		}
		if !find {
			log.Println(ee)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "nsq2kafka"
	app.Usage = "copy nsq topic to kafka"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: argSOURCE,
			Usage: "开票明细: 购方名称	合计金额	主要商品名称",
		},
		cli.StringFlag{
			Name: argTARGET,
			Usage: "辅助余额表: 科目编码	客户信息(说明)	本期贷方(本位币)",
		},
		cli.StringFlag{
			Name:  argDICT,
			Usage: "科目对应",
		},
	}
	app.Action = appAction
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
