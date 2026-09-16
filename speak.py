import sys

import soundfile as sf
from kokoro_onnx import Kokoro

model, voices, voice, speed, wav, text = sys.argv[1:7]
k = Kokoro(model, voices)
samples, sr = k.create(text, voice=voice, speed=float(speed), lang="en-us")
sf.write(wav, samples, sr)
