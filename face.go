package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type FaceDetector struct {
	rekog     *rekognition.Rekognition
	fakeAPI   bool
	fakeIndex int
}

func NewFaceDetector(fakeAPI bool) *FaceDetector {
	if fakeAPI {
		return &FaceDetector{fakeAPI: true}
	}
	return &FaceDetector{
		rekog: rekognition.New(session.New()),
	}
}

func (fd *FaceDetector) DetectFaces(image []byte) (*rekognition.DetectFacesOutput, error) {
	if fd.fakeAPI {
		var dfo rekognition.DetectFacesOutput
		err := json.Unmarshal([]byte(fake[fd.fakeIndex]), &dfo)
		if err != nil {
			return nil, err
		}
		fd.fakeIndex++
		if fd.fakeIndex >= len(fake) {
			fd.fakeIndex = 0
		}
		return &dfo, nil
	}
	return fd.rekog.DetectFaces(&rekognition.DetectFacesInput{
		Image: &rekognition.Image{Bytes: image},
	})
}

var fake = []string{
	`{"FaceDetails":[]}`,
	`{"FaceDetails":[{"AgeRange":null,"Beard":null,"BoundingBox":{"Height":0.5120102763175964,"Left":0.0444658026099205,"Top":-0.2300969511270523,"Width":0.22225633263587952},"Confidence":99.99996185302734,"Emotions":null,"Eyeglasses":null,"EyesOpen":null,"Gender":null,"Landmarks":[{"Type":"eyeLeft","X":0.13331173360347748,"Y":-0.06670842319726944},{"Type":"eyeRight","X":0.2305816411972046,"Y":-0.02748301438987255},{"Type":"mouthLeft","X":0.11967450380325317,"Y":0.10469245165586472},{"Type":"mouthRight","X":0.2008291780948639,"Y":0.13726593554019928},{"Type":"nose","X":0.1809072047472,"Y":0.038752343505620956}],"MouthOpen":null,"Mustache":null,"Pose":{"Pitch":6.568717956542969,"Roll":14.095389366149902,"Yaw":6.412472724914551},"Quality":{"Brightness":50.58533477783203,"Sharpness":60.49041748046875},"Smile":null,"Sunglasses":null}]}`,
	`{"FaceDetails":[{"AgeRange":null,"Beard":null,"BoundingBox":{"Height":0.836955726146698,"Left":0.41928306221961975,"Top":-0.04353904351592064,"Width":0.3449418544769287},"Confidence":99.99998474121094,"Emotions":null,"Eyeglasses":null,"EyesOpen":null,"Gender":null,"Landmarks":[{"Type":"eyeLeft","X":0.5351470112800598,"Y":0.2619765102863312},{"Type":"eyeRight","X":0.6973856687545776,"Y":0.2662754952907562},{"Type":"mouthLeft","X":0.5362831950187683,"Y":0.5976361036300659},{"Type":"mouthRight","X":0.6706182956695557,"Y":0.5999993085861206},{"Type":"nose","X":0.6247119903564453,"Y":0.48171505331993103}],"MouthOpen":null,"Mustache":null,"Pose":{"Pitch":-17.631322860717773,"Roll":1.0789055824279785,"Yaw":8.349883079528809},"Quality":{"Brightness":76.4839859008789,"Sharpness":86.86019134521484},"Smile":null,"Sunglasses":null}]}`,
	`{"FaceDetails":[{"AgeRange":null,"Beard":null,"BoundingBox":{"Height":0.6982589960098267,"Left":0.2486524134874344,"Top":-0.022139623761177063,"Width":0.28551608324050903},"Confidence":100,"Emotions":null,"Eyeglasses":null,"EyesOpen":null,"Gender":null,"Landmarks":[{"Type":"eyeLeft","X":0.29664716124534607,"Y":0.2613499164581299},{"Type":"eyeRight","X":0.3432553708553314,"Y":0.24820420145988464},{"Type":"mouthLeft","X":0.3270845413208008,"Y":0.5215925574302673},{"Type":"mouthRight","X":0.36074304580688477,"Y":0.5160477161407471},{"Type":"nose","X":0.2614690065383911,"Y":0.39391008019447327}],"MouthOpen":null,"Mustache":null,"Pose":{"Pitch":9.773420333862305,"Roll":-21.155019760131836,"Yaw":-65.08802032470703},"Quality":{"Brightness":79.98331451416016,"Sharpness":89.85481262207031},"Smile":null,"Sunglasses":null}]}`,
}

type FaceState struct {
	state int
}

func (fs *FaceState) AnalyzeFaces(faces []*rekognition.FaceDetail) string {
	if len(faces) == 0 {
		fs.state = 0
		return "Look at the camera"
	}
	if len(faces) > 1 {
		fs.state = 0
		return "Focus on your face"
	}

	face := faces[0]
	if face.BoundingBox == nil || face.Pose == nil {
		fs.state = 0
		return "Focus on your face"
	}

	bb := *face.BoundingBox
	if *bb.Width < 0.25 || *bb.Height < 0.25 {
		fs.state = 0
		return "Move the camera closer"
	}

	pose := *face.Pose
	dir := facing(*pose.Yaw)
	switch fs.state {
	case 0:
		if dir == Facing && *pose.Pitch > -20 && *pose.Pitch < 10 {
			fs.state = 1
			return "Turn to the right"
		}
		return "Face the camera"
	case 1:
		if dir == FacingRight {
			fs.state = 2
			return "Thank you"
		}
		return "Turn to the right"
	}
	return "Thank you"
}

const (
	FacingLeft            = -2
	FacingLeftTransition  = -1
	Facing                = 0
	FacingRightTransition = 1
	FacingRight           = 2
)

func facing(yaw float64) int {
	if yaw > -20 && yaw < 20 {
		return Facing
	}
	if yaw > -45 && yaw <= -20 {
		return FacingRightTransition
	}
	if yaw <= -45 {
		return FacingRight
	}
	if yaw >= 20 && yaw < 45 {
		return FacingLeftTransition
	}
	if yaw >= 45 {
		return FacingLeft
	}
	return Facing
}
