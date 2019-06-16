package main

import (
	"bytes"
	"fmt"
	"github.com/huichen/sego"
	"github.com/tealeg/xlsx"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

//C000008	财经新闻
//C000010	IT新闻
//C000013	健康新闻
//C000014	体育新闻
//C000016	旅游新闻
//C000020	教育新闻
//C000022	招聘新闻
//C000023	文化新闻
//C000024	军事新闻

func main() {
	xlsxfile,_ := xlsx.OpenFile("result.xlsx")
	sheet2,_ := xlsxfile.AddSheet("详细词频统计")
	row := sheet2.AddRow()
	ceil := row.AddCell()
	sheet1 := xlsxfile.Sheet["总览"]
	checklist := make([]string,0)
	for i,sheet1_row := range sheet1.Rows{
		if i == 0{
			continue
		}else if i == 331{
			break
		}
		ceil = row.AddCell()
		ceil.Value = sheet1_row.Cells[0].Value
		checklist = append(checklist,ceil.Value)
	}
	words := make(map[string]int)
	//载入词典
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("E:/GOPATH/src/github.com/huichen/sego/data/dictionary.txt")

	//载入停用词表
	stopwordsfile, _ := ioutil.ReadFile("stopwords.txt")
	stopwords := strings.Split(string(stopwordsfile), "\n")
	for p := range stopwords {
		stopwords[p] = strings.Replace(stopwords[p], string(13), "", 1)
	}
	stopwords = append(stopwords, "\n")

	//对文件进行分词与去除停用词
	corpus, _ := ioutil.ReadDir("SogouC.reduced/Reduced")
	i := 1
	for _, f := range corpus {
		k := 1
		files, _ := ioutil.ReadDir("SogouC.reduced/Reduced/" + f.Name())
		for _, file := range files {
			fmt.Printf("正在处理第%d个文件夹,第%d个文件……\n", i, k)
			d, _ := ioutil.ReadFile("SogouC.reduced/Reduced/" + f.Name() + "/" + file.Name())
			reader := transform.NewReader(bytes.NewReader(d), simplifiedchinese.GBK.NewDecoder())
			doc, _ := ioutil.ReadAll(reader)
			segments := segmenter.Segment(doc)
			re := make([]int,330)
			for _,seg := range segments{
				words[seg.Token().Text()] ++
			}
			for m,checkword := range checklist{
				for _,seg := range segments{
					if seg.Token().Text() == checkword{
						re[m]++
					}
				}
			}
			row = sheet2.AddRow()
			ceil = row.AddCell()
			ceil.Value = f.Name()+"-"+file.Name()
			for _,num :=range re{
				ceil = row.AddCell()
				ceil.SetInt(num)
			}
			ceil = row.AddCell()
			ceil.SetInt(len(segments))
			k ++
		}
		i ++
	}
	xlsxfile.Save("result.xlsx")
}