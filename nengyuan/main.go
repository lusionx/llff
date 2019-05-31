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

var argSOURCE = "source"
var argDICT = "dict"

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

func appAction(ctx *cli.Context) {
	var dicts, err = readDict(ctx.String(argDICT))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dicts)
	var sdata, err1 = readSource(ctx.String(argSOURCE))
	if err1 != nil {
		log.Fatal(err)
	}
	var sdata1 = []mSource{}
	for _, ss := range sdata {
		for _, d := range dicts {
			if ss.Code == d.Name {
				ss.Code = d.Code
				sdata1 = append(sdata1, ss)
			}
		}
	}
	log.Println(sdata1)
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
			Name:  "target",
			Usage: "辅助余额表",
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
