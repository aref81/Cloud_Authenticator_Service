package _rdPartyService

import (
	"Projeect/internal/model"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

var (
	imaggaSimilarityDetectionURL = os.Getenv("IMAGGA_SIMILARITY_DETECTION_URL")
)

func FaceSimilarity(faceID1, faceID2 string) float64 {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", imaggaSimilarityDetectionURL+"?face_id="+faceID1+"&second_face_id="+faceID2, nil)
	req.SetBasicAuth(imaggaApiKey, imaggaApiSecret)

	response, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error sending request to AIaaS: %v", err)
		return -1
	}

	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)

	var res model.ScoreRes
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		logrus.Errorf("JSON parse failed: %v", err)
		return -1
	}

	return res.Result.Score
}
