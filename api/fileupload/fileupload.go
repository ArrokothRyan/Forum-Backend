package fileupload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strings"
)

type fileUpload struct{}

type UploadFileInterface interface {
	UploadFile(c *gin.Context, file *multipart.FileHeader) (string, map[string]string)
}

//So what is exposed is Uploader
var FileUpload UploadFileInterface = &fileUpload{}

func (fu *fileUpload) UploadFile(c *gin.Context, file *multipart.FileHeader) (string, map[string]string) {
	errList := map[string]string{}

	f, err := file.Open()
	if err != nil {
		errList["Not_Image"] = "Please Upload a valid image"
		return "", errList
	}
	defer f.Close()

	size := file.Size
	//The image should not be more than 500KB
	fmt.Println("the size: ", size)
	if size > int64(512000) {
		errList["Too_large"] = "Sorry, Please upload an Image of 500KB or less"
		return "", errList

	}
	//only the first 512 bytes are used to sniff the content type of a file,
	//so, so no need to read the entire bytes of a file.
	buffer := make([]byte, size)
	f.Read(buffer)
	fileType := http.DetectContentType(buffer)
	//if the image is valid
	if !strings.HasPrefix(fileType, "image") {
		errList["Not_Image"] = "Please Upload a valid image"
		return "", errList
	}
	filePath := FormatFile(file.Filename)

	err = c.SaveUploadedFile(file, "upload/"+filePath)
	if err != nil {
		errList["Other_Err"] = "something went wrong"
		return "", errList
	}

	return filePath, nil
}
