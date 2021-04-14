package big_file

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func saveChunk(name string, buf []byte) {
	save, err := os.OpenFile("./files/"+ name, os.O_CREATE |os.O_TRUNC| os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer save.Close()
	_, err = save.Write(buf)
	if err != nil {
		panic(err)
	}
}

func NewChunkFile() {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		defer func(){
			if e := recover(); e != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": e})
			}
		}()
		c.Next()
	})



	r.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.Header().Set("Content-type", "audio/mp4")
		for i := 0; i <= 100; i++ {
			f, _ := os.Open(fmt.Sprintf("./files/img_%d.mp4", i))
			time.Sleep(time.Millisecond*1)
			b, _ := ioutil.ReadAll(f)
			c.Writer.Write(b)
			c.Writer.(http.Flusher).Flush()
		}
	})

	r.POST("/file", func(c *gin.Context) {
		file, head, _ := c.Request.FormFile("file")
		fmt.Printf("receive request, file:%v\n", head.Filename)
		// save, _ := os.OpenFile("./files/" + head.Filename, os.O_CREATE | os.O_RDWR, 0666)
		block := head.Size/100
		i := 0
		for {
			buf := make([]byte, block)
			n, err := file.Read(buf)
			if err != nil && err != io.EOF{
				panic(err.Error())
			}

			if n == 0 {
				break
			}

			time.Sleep(time.Millisecond*100)
			saveChunk(fmt.Sprintf("img_%d.mp4", i), buf)
			fmt.Printf("save chunk %d ok.\n", i)
			i++
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	r.Run(":8080")
}