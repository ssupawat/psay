import json
import os
import queue
import shutil
import socket
import subprocess
import sys
import tempfile
import threading
import time

import numpy as np
import soundfile as sf
from kokoro_onnx import Kokoro

IDLE_EXIT_SECONDS = 900
CACHE_LIMIT_BYTES = 200 * 1024 * 1024


def pick_model(model_dir):
    int8 = os.path.join(model_dir, "kokoro-v1.0.int8.onnx")
    if os.path.exists(int8):
        return int8
    return os.path.join(model_dir, "kokoro-v1.0.onnx")


def prune_cache(cache_dir):
    total = 0
    for name in os.listdir(cache_dir):
        try:
            total += os.path.getsize(os.path.join(cache_dir, name))
        except OSError:
            pass
    if total > CACHE_LIMIT_BYTES:
        shutil.rmtree(cache_dir)


def serve_request(kokoro, req):
    chunks = req["chunks"]
    voice = req["voice"]
    speed = req["speed"]
    cache_dir = os.path.dirname(req["cache"])
    os.makedirs(cache_dir, exist_ok=True)
    prune_cache(cache_dir)

    made = queue.Queue()
    failures = []

    def produce():
        try:
            for chunk in chunks:
                samples, sr = kokoro.create(
                    chunk, voice=voice, speed=speed, lang="en-us"
                )
                made.put((samples, sr))
        except Exception as exc:
            failures.append(exc)
        finally:
            made.put(None)

    threading.Thread(target=produce, daemon=True).start()

    pieces = []
    sr = 24000
    while True:
        item = made.get()
        if item is None:
            break
        samples, sr = item
        pieces.append(samples)
        fd, wav = tempfile.mkstemp(suffix=".wav", dir=req["dir"])
        os.close(fd)
        sf.write(wav, samples, sr)
        subprocess.run(["afplay", wav], check=False)
        os.unlink(wav)

    if failures:
        raise failures[0]
    if pieces:
        sf.write(req["cache"], np.concatenate(pieces), sr)


def handle_conn(kokoro, conn):
    payload = b""
    while not payload.endswith(b"\n"):
        part = conn.recv(65536)
        if not part:
            return
        payload += part
    try:
        serve_request(kokoro, json.loads(payload))
        reply = {"event": "done"}
    except Exception as exc:
        reply = {"event": "error", "message": str(exc)}
    conn.sendall((json.dumps(reply) + "\n").encode())


def daemon(sock_path, model_dir):
    kokoro = Kokoro(pick_model(model_dir), os.path.join(model_dir, "voices-v1.0.bin"))
    server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    server.bind(sock_path)
    server.listen(4)
    server.settimeout(5)
    last_activity = time.time()
    try:
        while True:
            try:
                conn, _ = server.accept()
            except socket.timeout:
                if time.time() - last_activity > IDLE_EXIT_SECONDS:
                    return
                continue
            last_activity = time.time()
            with conn:
                handle_conn(kokoro, conn)
    finally:
        try:
            os.unlink(sock_path)
        except FileNotFoundError:
            pass


def oneshot(model_dir):
    payload = b""
    while not payload.endswith(b"\n"):
        part = sys.stdin.buffer.read(65536)
        if not part:
            break
        payload += part
    kokoro = Kokoro(pick_model(model_dir), os.path.join(model_dir, "voices-v1.0.bin"))
    serve_request(kokoro, json.loads(payload))


def main():
    mode = sys.argv[1]
    if mode == "daemon":
        daemon(sys.argv[2], sys.argv[3])
    elif mode == "oneshot":
        oneshot(sys.argv[2])
    else:
        raise SystemExit(f"unknown mode: {mode}")


if __name__ == "__main__":
    main()
