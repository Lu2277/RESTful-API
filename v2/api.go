package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
func GetBooks(c *gin.Context) {
	//获取请求的路径后面的查询参数xxx    ?f=xxx
	formatter, _ := c.GetQuery("f")
	for i := 0; i < len(books); i++ {
		books[i].Format(formatter, i+1)
	}
	//将结构体编码为JSON，加上Indent表示给生成的JSON加上缩进，方便阅读
	c.IndentedJSON(http.StatusOK, books)

}
func GetBooksID(c *gin.Context) {
	//从请求路径中将id作为一个值读取
	id := c.Param("id")
	formatter, _ := c.GetQuery("f")
	for i, b := range books {
		if b.ID == id {
			b.Format(formatter, i+1)
			c.IndentedJSON(http.StatusOK, b)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found!"})

}

// PostBooks 创建新的书籍，添加到books中
func PostBooks(c *gin.Context) {
	var newBook Book
	//Gin的Bind方法将请求的内容绑定到一个newBook结构体中
	err := c.Bind(&newBook)
	if err != nil {
		return
	}
	//将新的数据返回去
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	r := gin.Default()
	r.GET("/books", GetBooks)
	r.GET("/books/:id", GetBooksID)
	r.POST("/books", PostBooks)
	err := r.Run(":8000")
	if err != nil {
		return
	}
}
