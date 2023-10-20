package model

type FaceDetectionRes struct {
	Result struct {
		Faces []struct {
			FaceID string `json:"face_id"`
		} `json:"faces"`
	} `json:"result"`
}

type ScoreRes struct {
	Result struct {
		Score float64 `json:"score"`
	} `json:"result"`
}
