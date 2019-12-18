package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

const IMGUR_TOKEN = "xxxxxx"

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST(`/img`, ImgUpload)
	if err := r.Run(`:8083`); err != nil {
		log.Print(err)
	}

}
func ImgUpload(c *gin.Context) {
	file, err := c.FormFile(`IMG`)
	if err != nil {
		log.Print(`file read failure`)
	} else {
		fileOpen, err := file.Open()
		if err != nil {
			log.Print(err)
		}
		result := upload(fileOpen, IMGUR_TOKEN)
		c.String(200, result)
	}
}

func upload(image io.Reader, token string) string {
	APIURL := "https://api.imgur.com/3/image"
	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	part, _ := writer.CreateFormFile("image", "dont care about name")

	_, err := io.Copy(part, image)
	if err != nil {
		log.Print(err)
	}
	if err := writer.Close(); err != nil {
		log.Print(err)
	}
	req, _ := http.NewRequest("POST", APIURL, buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	res, _ := client.Do(req)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Print(err)
		}
	}()
	b, _ := ioutil.ReadAll(res.Body)
	return string(b)
}
