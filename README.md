## Face Detection Spike

Experiments with face detection and anti-spoof liveness checking.

### Setup

```
go mod download

aws configure
```

### Start Server

```
AWS_SDK_LOAD_CONFIG=1 go run server.go face.go
```

or without AWS:

```
FAKE_API=1 go run server.go face.go
```

Visit http://localhost:3000/

Sample output:

```
{
  FaceDetails: [{
      BoundingBox: {
        Height: 0.4014367461204529,
        Left: 0.3085518479347229,
        Top: 0.1626598834991455,
        Width: 0.2844625413417816
      },
      Confidence: 100,
      Landmarks: [
        {
          Type: "eyeLeft",
          X: 0.3641567528247833,
          Y: 0.3375476598739624
        },
        {
          Type: "eyeRight",
          X: 0.478274405002594,
          Y: 0.31672123074531555
        },
        {
          Type: "mouthLeft",
          X: 0.40021002292633057,
          Y: 0.47150078415870667
        },
        {
          Type: "mouthRight",
          X: 0.4938865303993225,
          Y: 0.4551337659358978
        },
        {
          Type: "nose",
          X: 0.4104361832141876,
          Y: 0.4025741219520569
        }
      ],
      Pose: {
        Pitch: -5.396634101867676,
        Roll: -11.281394958496094,
        Yaw: -28.795425415039062
      },
      Quality: {
        Brightness: 80.30354309082031,
        Sharpness: 89.85481262207031
      }
    }]
}
```
