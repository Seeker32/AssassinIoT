#!/usr/bin/env python3
"""Simulate IoT devices publishing telemetry data via MQTT.

Usage:
    uv run test/stress_test.py
    uv run test/stress_test.py --devices 100 --rate 500
    uv run test/stress_test.py --properties test/stress_properties.json
"""

from __future__ import annotations

import argparse
import asyncio
import json
import logging
import random
import signal
import sys
import time
from pathlib import Path

import aiomqtt

logger = logging.getLogger("device-sim")

DEFAULT_PROPERTIES: dict = {
    "temperature": {"type": "float", "min": -10.0, "max": 50.0},
    "humidity":    {"type": "float", "min": 0.0,   "max": 100.0},
    "voltage":     {"type": "float", "min": 3.0,   "max": 4.2},
    "rssi":        {"type": "int",   "min": -90,   "max": -30},
    "status":      {"type": "choice","values": ["online", "idle", "active"]},
}


def load_properties(path: str | None) -> dict:
    if path is None:
        return dict(DEFAULT_PROPERTIES)
    return json.loads(Path(path).read_text())


def generate_value(prop_def: dict) -> int | float | str | bool:
    kind = prop_def["type"]
    if kind == "float":
        return round(random.uniform(prop_def["min"], prop_def["max"]), 2)
    if kind == "int":
        return random.randint(prop_def["min"], prop_def["max"])
    if kind == "choice":
        return random.choice(prop_def["values"])
    if kind == "bool":
        return random.choice([True, False])
    raise ValueError(f"Unknown property type: {kind}")


def generate_payload(properties: dict) -> dict:
    keys = list(properties)
    n = random.randint(max(1, int(len(keys) * 0.6)), len(keys))
    payload = {}
    for key in random.sample(keys, n):
        payload[key] = generate_value(properties[key])
    payload["ts"] = int(time.time() * 1000) + random.random()
    return payload


async def run_device(
    broker: str,
    port: int,
    device_id: str,
    username: str,
    password: str,
    properties: dict,
    interval: float,
    stop: asyncio.Event,
) -> None:
    kwargs: dict = dict(hostname=broker, port=port, identifier=device_id)
    if username:
        kwargs["username"] = username
    if password:
        kwargs["password"] = password
    try:
        async with aiomqtt.Client(**kwargs) as client:
            topic = f"device/{device_id}/telemetry"
            while not stop.is_set():
                payload = generate_payload(properties)
                try:
                    await client.publish(topic, json.dumps(payload), qos=0)
                except Exception:
                    logger.exception("[%s] Publish failed", device_id)
                await asyncio.sleep(interval)
    except Exception:
        logger.exception("[%s] Connection failed", device_id)


async def main_async() -> None:
    args = parse_args()
    logging.basicConfig(
        level=getattr(logging, args.log_level),
        format="%(asctime)s [%(levelname)s] %(message)s",
        datefmt="%H:%M:%S",
    )

    properties = load_properties(args.properties)
    logger.info("Loaded %d properties: %s", len(properties), ", ".join(properties))

    # Device IDs
    if args.username:
        device_ids = [f"{args.username}-{i + 1:04d}" for i in range(args.devices)]
        mqtt_user = args.username
    else:
        device_ids = [f"sim-{i + 1:04d}" for i in range(args.devices)]
        mqtt_user = ""

    # Per-device interval: total rate = devices / interval
    n = len(device_ids)
    rate = args.rate if args.rate > 0 else n  # default: 1 msg/s per device
    interval = n / rate

    logger.info(
        "Starting %d devices, %.0f msg/s total (%.3fs per device)",
        n, rate, interval,
    )

    stop = asyncio.Event()

    # Handle shutdown signals
    loop = asyncio.get_running_loop()
    for sig in (signal.SIGTERM, signal.SIGINT):
        loop.add_signal_handler(sig, stop.set)

    tasks = []
    for i, device_id in enumerate(device_ids):
        stagger = (i / n) * min(interval, 0.5)
        tasks.append(
            asyncio.create_task(
                run_device(args.broker, args.port, device_id, mqtt_user,
                           args.password, properties, interval, stop)
            )
        )

    await stop.wait()
    logger.info("Shutting down...")
    await asyncio.gather(*tasks, return_exceptions=True)


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Simulate IoT devices publishing telemetry")
    p.add_argument("--broker", default="localhost", help="MQTT broker host")
    p.add_argument("--port", type=int, default=1883, help="MQTT broker port")
    p.add_argument("-u", "--username", default=None, help="MQTT username")
    p.add_argument("-P", "--password", default="", help="MQTT password")
    p.add_argument("-n", "--devices", type=int, default=50, help="Number of simulated devices")
    p.add_argument("--rate", type=int, default=0, help="Total messages per second (default: devices)")
    p.add_argument("--properties", default=None, help="Path to JSON file defining properties")
    p.add_argument("--log-level", default="INFO", choices=["DEBUG", "INFO", "WARNING", "ERROR"])
    return p.parse_args()


def main() -> None:
    try:
        asyncio.run(main_async())
    except KeyboardInterrupt:
        print("\nStopped.")


if __name__ == "__main__":
    main()
