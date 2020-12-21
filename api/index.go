package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var (
	// searcher 是协程安全的
	searcher   = riot.Engine{}
	bs         []blog
	contentMD5 [md5.Size]byte
)

const (
	indexJsonURL = "https://hugo-blog.qraffa.vercel.app/search/index.json"
)

type blog struct {
	Content   string
	Date      string
	Permalink string
	Title     string
	Tags      []string
}

type searchRes struct {
	types.BaseResp
	Docs []blog
}

func Handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	keyword := query.Get("s")
	log.Print("new search ", keyword)
	log.Print(fmt.Sprintf("index.json md5 ==> %x", contentMD5))

	bytes, err := json.Marshal(search(keyword))
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(bytes)
}

func init() {
	log.Print("making index")
	getIndex()
	log.Print("making index done")
}

func getIndex() {
	// 获取index.json文件
	res, err := http.Get(indexJsonURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	// index.json md5
	contentMD5 = md5.Sum(content)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &bs)
	if err != nil {
		log.Fatal(err)
	}
	makeIndex(bs)
}

// 初始化搜索器，为文章添加索引
func makeIndex(ba []blog) {
	// 初始化
	searcher.Init(types.EngineOpts{
		Using:   3,
		//GseDict: "zh",
		GseDict: "data/dictionary.txt", // for vercel includeFiles
	})
	for k, v := range ba {
		searcher.Index(strconv.Itoa(k+1), types.DocData{Content: v.Content})
	}
	searcher.Flush()
}

func search(keyword string) searchRes {
	res := searcher.Search(types.SearchReq{Text: keyword})
	baseRes := res.BaseResp
	docs := make([]blog, 0)
	if docsRes, ok := res.Docs.(types.ScoredDocs); ok {
		for k, v := range docsRes {
			id, _ := strconv.Atoi(v.DocId)
			docs = append(docs, bs[id-1])
			// 不返回content
			docs[k].Content = ""
		}
	}
	return searchRes{
		BaseResp: baseRes,
		Docs:     docs,
	}
}
