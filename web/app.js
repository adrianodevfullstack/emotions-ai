import {
  DrawingUtils,
  FaceLandmarker,
  FilesetResolver,
} from "https://cdn.jsdelivr.net/npm/@mediapipe/tasks-vision@0.10.14";

const EMOJI = {
  joy: "😀",
  sadness: "😢",
  anger: "😠",
  surprise: "😮",
  fear: "😨",
  disgust: "🤢",
  neutral: "😐",
};

const startBtn = document.getElementById("start-btn");
const video = document.getElementById("video");
const overlay = document.getElementById("overlay");
const ctxOverlay = overlay.getContext("2d");
const emotionLabel = document.getElementById("emotion-label");
const emoji = document.getElementById("emoji");
const confidenceEl = document.getElementById("confidence");
const confidenceBar = document.getElementById("confidence-bar");
const statusEl = document.getElementById("status");
const interpretationText = document.getElementById("interpretation-text");

let faceLandmarker = null;
let ws = null;
let running = false;

function score(blendshapes) {
  const s = Object.fromEntries(
    blendshapes.map((b) => [b.categoryName, b.score]),
  );
  const g = (name) => s[name] ?? 0;

  const scores = {
    joy:
      g("mouthSmileLeft") +
      g("mouthSmileRight") +
      0.5 * (g("cheekSquintLeft") + g("cheekSquintRight")),
    sadness:
      g("mouthFrownLeft") +
      g("mouthFrownRight") +
      0.6 * g("browInnerUp") +
      0.3 * (g("browDownLeft") + g("browDownRight")),
    anger:
      g("browDownLeft") +
      g("browDownRight") +
      0.7 * (g("noseSneerLeft") + g("noseSneerRight")) +
      0.5 * (g("mouthPressLeft") + g("mouthPressRight")),
    surprise:
      g("browOuterUpLeft") +
      g("browOuterUpRight") +
      0.6 * g("browInnerUp") +
      (g("eyeWideLeft") + g("eyeWideRight")) +
      0.5 * g("jawOpen"),
    fear:
      0.6 * g("browInnerUp") +
      (g("eyeWideLeft") + g("eyeWideRight")) +
      0.3 * g("jawOpen"),
    disgust:
      g("noseSneerLeft") +
      g("noseSneerRight") +
      0.6 * (g("mouthUpperUpLeft") + g("mouthUpperUpRight")),
  };

  let best = "neutral";
  let bestScore = 0.35;
  for (const [emotionName, value] of Object.entries(scores)) {
    if (value > bestScore) {
      best = emotionName;
      bestScore = value;
    }
  }

  const confidence = Math.max(0, Math.min(1, bestScore / 1.5));
  return { emotion: best, confidence };
}

async function start() {
  startBtn.disabled = true;
  statusEl.textContent = "loading model…";

  const fileset = await FilesetResolver.forVisionTasks(
    "https://cdn.jsdelivr.net/npm/@mediapipe/tasks-vision@0.10.14/wasm",
  );
  faceLandmarker = await FaceLandmarker.createFromOptions(fileset, {
    baseOptions: {
      modelAssetPath:
        "https://storage.googleapis.com/mediapipe-models/face_landmarker/face_landmarker/float16/1/face_landmarker.task",
      delegate: "GPU",
    },
    outputFaceBlendshapes: true,
    runningMode: "VIDEO",
    numFaces: 1,
  });

  statusEl.textContent = "requesting webcam access…";
  const stream = await navigator.mediaDevices.getUserMedia({
    video: { width: 480, height: 360 },
    audio: false,
  });
  video.srcObject = stream;
  await video.play();

  overlay.width = video.videoWidth;
  overlay.height = video.videoHeight;

  const protocol = location.protocol === "https:" ? "wss" : "ws";
  ws = new WebSocket(`${protocol}://${location.host}/ws`);
  ws.onopen = () => {
    statusEl.textContent = "connected — detecting";
    statusEl.classList.add("connected");
  };
  ws.onclose = () => {
    statusEl.textContent = "disconnected";
    statusEl.classList.remove("connected");
  };
  ws.onmessage = (ev) => {
    const summary = JSON.parse(ev.data);
    interpretationText.textContent =
      summary.interpretation || "(no interpretation for this window)";
  };

  running = true;
  startBtn.textContent = "Running…";
  loop();
}

function loop() {
  if (!running) return;
  const now = performance.now();
  const result = faceLandmarker.detectForVideo(video, now);

  ctxOverlay.clearRect(0, 0, overlay.width, overlay.height);

  if (result.faceBlendshapes && result.faceBlendshapes.length > 0) {
    const blendshapes = result.faceBlendshapes[0].categories;
    const { emotion, confidence } = score(blendshapes);

    emotionLabel.textContent = emotion;
    emoji.textContent = EMOJI[emotion] ?? "🙂";
    confidenceEl.textContent = `confidence: ${(confidence * 100).toFixed(0)}%`;
    confidenceBar.style.width = `${(confidence * 100).toFixed(0)}%`;

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          emotion,
          trust: confidence,
          timestampMs: Date.now(),
        }),
      );
    }

    if (result.faceLandmarks && result.faceLandmarks.length > 0) {
      const drawing = new DrawingUtils(ctxOverlay);
      drawing.drawConnectors(
        result.faceLandmarks[0],
        FaceLandmarker.FACE_LANDMARKS_TESSELATION,
        { color: "#6ea8fe33", lineWidth: 1 },
      );
    }
  } else {
    emotionLabel.textContent = "no face detected";
    emoji.textContent = "🙈";
    confidenceBar.style.width = "0%";
  }

  requestAnimationFrame(loop);
}

startBtn.addEventListener("click", () => {
  start().catch((err) => {
    console.error(err);
    statusEl.textContent = "error: " + err.message;
    startBtn.disabled = false;
  });
});
