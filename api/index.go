package api

import (
	"encoding/json"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var (
	// searcher 是协程安全的
	searcher = riot.Engine{}
	bs       []blog
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

func Handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	keyword := query.Get("s")
	log.Print("new search ", keyword)
	resBlog := search(keyword)
	bytes, err := json.Marshal(resBlog)
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
		GseDict: "zh",
	})

	for k, v := range ba {
		searcher.Index(strconv.Itoa(k+1), types.DocData{Content: v.Content})
	}
	searcher.Flush()
}

// 查找文章
func search(keyword string) []blog {
	res := make([]blog, 0)
	if docs, ok := searcher.Search(types.SearchReq{Text: keyword}).Docs.(types.ScoredDocs); ok {
		for _, v := range docs {
			id, _ := strconv.Atoi(v.DocId)
			res = append(res, bs[id-1])
		}
	}
	return res
}
