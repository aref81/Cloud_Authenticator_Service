package internal

import (
	"Projeect/utils"
	_rdPartyService "Projeect/utils/3rdPartyService"
	"Projeect/utils/datasource"
	"github.com/sirupsen/logrus"
)

func processReq(encodedNationalCode string) {
	user, err := psql.FetchUser(encodedNationalCode)
	if err != nil {
		logrus.Errorf("User not found: %v", err)
	}

	uuid1 := utils.EncodeBase64(user.NationalCode) + "_IMAGE_1"
	uuid2 := utils.EncodeBase64(user.NationalCode) + "_IMAGE_2"

	pic1, err := datasource.DownloadPic(uuid1)
	if err != nil {
		logrus.Errorf("Failed to get image: %v", err)
	}
	pic2, err := datasource.DownloadPic(uuid2)
	if err != nil {
		logrus.Errorf("Failed to get image: %v", err)
	}

	imageID1, err := _rdPartyService.FaceDetection(pic1)
	if err != nil {
		logrus.Errorf("Failed to send image to AIaaS: %v", err)
	}
	imageID2, err := _rdPartyService.FaceDetection(pic2)
	if err != nil {
		logrus.Errorf("Failed to send image to AIaaS: %v", err)
	}

	score := _rdPartyService.FaceSimilarity(imageID1, imageID2)

	if score >= 80 {
		err := psql.UpdateStatus(user.NationalCode, "accepted")
		if err != nil {
			logrus.Errorf(err.Error())
		}
		_, err = _rdPartyService.SendMail("Congrats! \n Your authentication request is accepted", user.Email)
		if err != nil {
			logrus.Errorf(err.Error())
		}
		logrus.Infof("%s : Accepted", user.NationalCode)
	} else {
		err := psql.UpdateStatus(user.NationalCode, "rejected")
		if err != nil {
			logrus.Errorf(err.Error())
		}
		_, err = _rdPartyService.SendMail("sorry! \n Your authentication request is rejected", user.Email)
		if err != nil {
			logrus.Errorf(err.Error())
		}
		logrus.Infof("%s : Rejected", user.NationalCode)
	}
}
