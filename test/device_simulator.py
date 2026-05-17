#!/usr/bin/env python3
"""Simulate IoT devices publishing telemetry via MQTT.

Each device connects to the MQTT broker and periodically publishes random
property values to device/{device_id}/telemetry.

Usage:
    python device_simulator.py -n 10 -i 5.0
    python device_simulator.py -n 3 -i 2.0 --prefix sensor-
    python device_simulator.py -n 1 -u test1 -P 123456
    python device_simulator.py --properties my_props.json
"""

import argparse
import json
import logging
import random
import signal
import sys
import threading
import time
from pathlib import Path

import paho.mqtt.client as mqtt

DEFAULT_PROPERTIES = {
    "temperature": {"type": "float", "min": -10.0, "max": 50.0},
    "humidity": {"type": "float", "min": 0.0, "max": 100.0},
    "voltage": {"type": "float", "min": 3.0, "max": 4.2},
    "rssi": {"type": "int", "min": -90, "max": -30},
    "status": {"type": "choice", "values": ["online", "idle", "active"]},
}

logger = logging.getLogger("device-simulator")


def load_properties(path: str | None) -> dict:
    if path is None:
        return dict(DEFAULT_PROPERTIES)
    with open(path) as f:
        return json.load(f)


def generate_value(prop_def: dict) -> int | float | str | bool:
    kind = prop_def["type"]
    if kind == "float":
        v = random.uniform(prop_def["min"], prop_def["max"])
        return round(v, 2)
    if kind == "int":
        return random.randint(prop_def["min"], prop_def["max"])
    if kind == "choice":
        return random.choice(prop_def["values"])
    if kind == "bool":
        return random.choice([True, False])
    raise ValueError(f"Unknown property type: {kind}")


def generate_payload(properties: dict) -> dict:
    payload = {}
    keys = list(properties)
    n = random.randint(max(1, int(len(keys) * 0.6)), len(keys))
    for key in random.sample(keys, n):
        payload[key] = generate_value(properties[key])
    payload["ts"] = int(time.time() * 1000)
    return payload


def run_device(
    device_id: str,
    broker: str,
    port: int,
    interval: float,
    properties: dict,
    stop_event: threading.Event,
    mqtt_username: str = "",
    password: str = "",
):
    client = mqtt.Client(
        mqtt.CallbackAPIVersion.VERSION2,
        client_id=device_id,
    )
    if mqtt_username:
        client.username_pw_set(mqtt_username, password)

    conn_ok = threading.Event()
    conn_rc = [0]

    def on_connect(client, userdata, flags, reason_code, properties):
        conn_rc[0] = reason_code
        if reason_code == 0:
            logger.info("[%s] Connected to broker", device_id)
        else:
            logger.error("[%s] Connection rejected: %s", device_id, reason_code)
        conn_ok.set()

    def on_disconnect(client, userdata, flags, reason_code, properties):
        if reason_code != 0:
            logger.warning("[%s] Disconnected: %s", device_id, reason_code)

    client.on_connect = on_connect
    client.on_disconnect = on_disconnect

    try:
        client.connect(broker, port, keepalive=60)
    except OSError as e:
        logger.error("[%s] Cannot reach broker: %s", device_id, e)
        return

    client.loop_start()

    if not conn_ok.wait(timeout=10):
        logger.error("[%s] Connection timed out", device_id)
        client.loop_stop()
        return

    if conn_rc[0] != 0:
        client.loop_stop()
        return

    topic = f"device/{device_id}/telemetry"
    logger.info("[%s] Reporting to %s every %.1fs", device_id, topic, interval)

    try:
        while not stop_event.is_set():
            payload = generate_payload(properties)
            payload_bytes = json.dumps(payload)
            info = client.publish(topic, payload_bytes, qos=0)
            if info.rc == mqtt.MQTT_ERR_SUCCESS:
                logger.info("[%s] Published: %s", device_id, payload_bytes)
            else:
                logger.error("[%s] Publish failed: rc=%s", device_id, info.rc)
            stop_event.wait(interval)
    finally:
        client.loop_stop()
        client.disconnect()
        logger.info("[%s] Stopped", device_id)


def main():
    parser = argparse.ArgumentParser(
        description="Simulate IoT devices reporting telemetry via MQTT"
    )
    parser.add_argument(
        "-n", "--devices", type=int, default=5,
        help="Number of simulated devices (default: 5)",
    )
    parser.add_argument(
        "-i", "--interval", type=float, default=5.0,
        help="Seconds between each device's reports (default: 5.0)",
    )
    parser.add_argument(
        "-b", "--broker", default="localhost",
        help="MQTT broker host (default: localhost)",
    )
    parser.add_argument(
        "-p", "--port", type=int, default=1883,
        help="MQTT broker port (default: 1883)",
    )
    parser.add_argument(
        "--prefix", default="sim-",
        help="Device ID prefix when --username is not set (default: sim-)",
    )
    parser.add_argument(
        "-u", "--username", default=None,
        help="MQTT username. All devices share this username for auth. "
             "Device IDs are derived from it when set.",
    )
    parser.add_argument(
        "-P", "--password", default="",
        help="MQTT password. Used with --username.",
    )
    parser.add_argument(
        "--properties", default=None,
        help="Path to a JSON file defining device properties and their value ranges",
    )
    parser.add_argument(
        "--log-level", default="INFO",
        choices=["DEBUG", "INFO", "WARNING", "ERROR"],
        help="Logging level (default: INFO)",
    )
    args = parser.parse_args()

    logging.basicConfig(
        level=getattr(logging, args.log_level),
        format="%(asctime)s [%(levelname)s] %(message)s",
        datefmt="%H:%M:%S",
    )

    properties = load_properties(args.properties)
    logger.info("Loaded %d properties: %s", len(properties), ", ".join(properties))

    stop_event = threading.Event()
    threads: list[threading.Thread] = []

    for i in range(args.devices):
        if args.username:
            device_id = args.username if args.devices == 1 else f"{args.username}-{i + 1:03d}"
            mqtt_username = args.username
        else:
            device_id = f"{args.prefix}{i + 1:03d}"
            mqtt_username = device_id

        t = threading.Thread(
            target=run_device,
            args=(
                device_id, args.broker, args.port, args.interval,
                properties, stop_event, mqtt_username, args.password,
            ),
            daemon=True,
        )
        t.start()
        threads.append(t)

    logger.info(
        "%d device(s) started, reporting every %.1fs. Press Ctrl+C to stop.",
        args.devices,
        args.interval,
    )

    def shutdown(signum, frame):
        logger.info("Shutting down...")
        stop_event.set()

    signal.signal(signal.SIGINT, shutdown)
    signal.signal(signal.SIGTERM, shutdown)

    try:
        while any(t.is_alive() for t in threads):
            time.sleep(0.5)
    except KeyboardInterrupt:
        pass

    stop_event.set()
    for t in threads:
        t.join(timeout=5.0)

    logger.info("All devices stopped.")


if __name__ == "__main__":
    main()
