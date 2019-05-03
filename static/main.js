// @ts-check

class Socket {
  constructor(url) {
    this.socket = new WebSocket(url);
    this.socket.binaryType = "blob";
    this.socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      this.data = msg;
      // console.log(msg);
    };
  }

  send(data) {
    if (this.socket.readyState !== 1) {
      console.error(`socket not open: ${this.socket.readyState}`);
      return;
    }
    this.socket.send(data);
  }

  lastMessage() {
    return this.data;
  }

  close() {
    this.socket.close();
  }
}

class VideoFeed {
  constructor(videoID) {
    this.video = document.getElementById(videoID);
  }

  async aquireCamera() {
    const constraints = {
      video: {
        facingMode: "user",
        width: {min: 1280},
        height: {min: 720},
      },
      audio: false,
    };

    this.stream = await navigator.mediaDevices.getUserMedia(constraints);
    this.video.setAttribute("autoplay", "true");
    this.video.srcObject = this.stream;
  }

  startCapture(server, interval) {
    this.server = server;
    this.interval = interval;
    this.snapshot = document.createElement("canvas");
    this.snapshotCtx = this.snapshot.getContext("2d");

    this._captureLoop();
  }

  stop() {
    window.clearTimeout(this.timeoutHandler);
    this.stream.getVideoTracks()[0].stop();
  }

  _captureLoop() {
    this.timeoutHandler = window.setTimeout(this._takeSnapshot.bind(this), this.interval);
  }

  _takeSnapshot() {
    this.snapshot.width = this.video.videoWidth;
    this.snapshot.height = this.video.videoHeight;
    this.snapshotCtx.drawImage(this.video, 0, 0);

    const mimeType = "image/jpeg";
    const quality = 0.88;
    this.snapshot.toBlob(this._send.bind(this), mimeType, quality);

    this._captureLoop();
  }

  _send(blob) {
    this.server.send(blob);
  }
}

const boxColor = "rgb(255, 255, 255)";
const labelColor = "rgb(255, 255, 255)";
const shadowColor = "rgb(0, 0, 0)";
const labelFont = "18px sans-serif";
const commandFont = "48px sans-serif";

class Overlay {
  constructor(canvasID) {
    const canvas = document.getElementById(canvasID);
    this.width = canvas.width;
    this.height = canvas.height;
    this.ctx = canvas.getContext("2d");
    this.canvas = canvas;
  }

  clear() {
    this.ctx.clearRect(0, 0, this.width, this.height);
  }

  drawFace(f) {
    const bb = f.BoundingBox;
    const x = Math.floor(bb.Left * this.width);
    const y = Math.floor(bb.Top * this.height);
    const w = Math.floor(bb.Width * this.width);
    const h = Math.floor(bb.Height * this.height);
    this.ctx.strokeStyle = boxColor;
    this.ctx.strokeRect(x, y, w, h);

    const pitch = f.Pose.Pitch.toFixed(2);
    const roll = f.Pose.Roll.toFixed(2);
    const yaw = f.Pose.Yaw.toFixed(2);
    this.ctx.font = labelFont;
    this.drawText(`Pitch ${pitch}`, x + 4, y + 20);
    this.drawText(`Roll ${roll}`, x + 4, y + 40);
    this.drawText(`Yaw ${yaw}`, x + 4, y + 60);
  }

  drawCommand(cmd) {
    this.ctx.font = commandFont;
    this.drawText(cmd, 4, 50);
  }

  drawText(text, x, y) {
    this.ctx.fillStyle = shadowColor;
    this.ctx.fillText(text, x + 1, y + 1);
    this.ctx.fillStyle = labelColor;
    this.ctx.fillText(text, x, y);
  }
}

var animationHandle = null;

function showOverlay(server) {
  animationHandle = requestAnimationFrame(() => showOverlay(server));
  const msg = server.lastMessage();
  if (!msg) {
    return;
  }
  overlay.clear();
  overlay.drawCommand(msg.Command);
  for (let f of msg.FaceDetails) {
    overlay.drawFace(f);
  }
}

const videoFeed = new VideoFeed("viewfinder");
const overlay = new Overlay("overlay");

const startButton = document.getElementById("go");
const stopButton = document.getElementById("stop");
const videoInterface = document.getElementById("interface");

var server = null;

startButton.addEventListener("click", async () => {
  startButton.disabled = true;
  videoInterface.style.display = "block";
  await videoFeed.aquireCamera();
  stopButton.disabled = false;

  server = new Socket("ws://localhost:3000/socket");
  videoFeed.startCapture(server, 1000);
  showOverlay(server);
});

stopButton.addEventListener("click", () => {
  cancelAnimationFrame(animationHandle);
  overlay.clear();
  stopButton.disabled = true;
  videoFeed.stop();
  server.close();
  videoInterface.style.display = "none";
  startButton.disabled = false;
});
