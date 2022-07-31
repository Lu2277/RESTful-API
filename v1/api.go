package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Book struct {
	ID     string
	Title  string
	Author string
	Type   string

	Formatted string //格式化后的文本
}

//Format 用于格式化书籍为GB/T格式
func (b *Book) Format(formatter string, idx int) {
	switch formatter {
	case "gbt":
		b.Formatted = fmt.Sprintf("[%v] %v. %v[%v]", idx, b.Author, b.Title, b.Type)
	default:
		b.Formatted = ""
	}
}

var books = []Book{
	{
		ID:     "sicp",
		Title:  "计算机程序的构造及解释",
		Author: "Jal Sussman",
		Type:   "M",
	},
	{
		ID:     "go",
		Title:  "Go语言设计及实现",
		Author: "左书祺",
		Type:   "N",
	},
	{
		ID:     "sql",
		Title:  "SQL数据库介绍",
		Author: "Jerry",
		Type:   "M",
	},
}

// GetBooks 获取所有books,若指定f=gbt ，则格式化为GB/T格式
func GetBooks(w http.ResponseWriter, r *http.Request) {
	//日志记录
	log.Println(r.Method, r.URL.String())
	//获取请求的路径GET后面的查询参数
	formatter := r.URL.Query().Get("f")
	for i := 0; i < len(books); i++ {
		books[i].Format(formatter, i+1)
	}
	//json.Marshal将结构体编码为JSON，加上Indent表示给生成的JSON加上缩进，方便阅读
	j, _ := json.MarshalIndent(books, "", "   ")
	w.Write(j)

}

// PostBooks 创建新的书籍，添加到books中
func PostBooks(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.String())
	var newBook Book
	//从r.Body获取请求中传来的JSON数据
	body, _ := ioutil.ReadAll(r.Body)
	//将JSON数据读取中一个结构体中
	err := json.Unmarshal(body, &newBook)
	if err != nil {
		return
	}
	//将新的数据返回去
	books = append(books, newBook)
	j, _ := json.MarshalIndent(newBook, "", "   ")
	w.Write(j)
}
func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetBooks(w, r)
	case http.MethodPost:
		PostBooks(w, r)

	}
}
func main() {
	http.HandleFunc("/books", handleBooks)
	err := http.ListenAndServe("127.0.0.1:8000", nil)
	if err != nil {
		return
	}
}
