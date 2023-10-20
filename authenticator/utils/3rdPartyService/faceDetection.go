package _rdPartyService

import (
	"Projeect/internal/model"
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	imaggaApiKey           = os.Getenv("IMAGGA_API_KEY")
	imaggaApiSecret        = os.Getenv("IMAGGA_API_SECRET")
	imaggaFaceDetectionURL = os.Getenv("IMAGGA_FACE_DETECTION_URL")
)

func FaceDetection(imageBinary *bytes.Buffer) (faceID string, err error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	defer writer.Close()

	formValue, err := writer.CreateFormFile("image", "user_pic")
	if err != nil {
		logrus.Errorf("Error creating form file: %v", err)
		return "", err
	}

	_, err = io.Copy(formValue, imageBinary)
	if err != nil {
		logrus.Errorf("Error injecting binary data: %v", err)
		return "", err
	}

	url := imaggaFaceDetectionURL + "?return_face_id=1"
	request, err := http.NewRequest("POST", url, bytes.NewReader(requestBody.Bytes()))
	if err != nil {
		logrus.Errorf("Error making http request: %v", err)
		return "", err
	}
	request.SetBasicAuth(imaggaApiKey, imaggaApiSecret)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logrus.Errorf("Error sending request to AIaaS: %v", err)
		return "", err
	}

	defer response.Body.Close()
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, response.Body)
	if err != nil {
		logrus.Errorf("Error reading AIaaS response: %v", err)
		return "", err
	}

	var res model.FaceDetectionRes
	err = json.Unmarshal([]byte(buffer.String()), &res)
	if err != nil {
		logrus.Errorf("JSON parse failed: %v", err)
		return "", err
	}

	if len(res.Result.Faces) > 0 {
		return res.Result.Faces[0].FaceID, nil
	}
	return "", nil
}
