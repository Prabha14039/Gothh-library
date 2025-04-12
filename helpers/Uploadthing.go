package helpers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"os"
	"github.com/labstack/echo"
)

func UploadFiles() {

	var key Key

	url := "https://api.uploadthing.com/v6/uploadFiles"

	payload := strings.NewReader("{\n  \"files\": [\n    {\n      \"name\": \"\",\n      \"size\": 1,\n      \"type\": \"\",\n      \"customId\": null\n    }\n  ],\n  \"acl\": \"public-read\",\n  \"metadata\": null,\n  \"contentDisposition\": \"inline\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("X-Uploadthing-Api-Key", key.UploadThing_Key)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}

func PrepareUpload(c echo.Context) error {

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("uploads/" + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully </p>", file.Filename))
}
