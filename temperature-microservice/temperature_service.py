#!/usr/bin/env python3
import asyncio
import json
import logging
import os
import signal
from dataclasses import dataclass
from typing import Optional, Tuple

import nats
from nats.aio.client import Client as NATS
from nats.micro import add_service, Request  # type: ignore

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


@dataclass
class TemperatureRequest:
    """Represents the incoming temperature data"""
    temperature: float  # Temperature in Celsius


@dataclass
class MongoChangeStreamEvent:
    """Represents MongoDB change stream data"""
    operation_type: str
    full_document: dict


@dataclass
class TemperatureResponse:
    """Represents the temperature classification"""
    status: str  # "cold", "warm", or "hot"
    temperature: float  # Temperature in Celsius
    description: str
    sensor_id: Optional[str] = None
    location: Optional[str] = None

    def to_dict(self):
        result = {
            "status": self.status,
            "temperature": self.temperature,
            "description": self.description
        }
        if self.sensor_id:
            result["sensor_id"] = self.sensor_id
        if self.location:
            result["location"] = self.location
        return result


def classify_temperature(temp_celsius: float) -> Tuple[str, str]:
    """Classify temperature based on Celsius value"""
    if temp_celsius < 10:
        status = "cold"
        description = f"{temp_celsius:.1f}C is cold"
    elif 10 <= temp_celsius <= 25:
        status = "warm"
        description = f"{temp_celsius:.1f}C is warm"
    else:
        status = "hot"
        description = f"{temp_celsius:.1f}C is hot"

    return status, description


async def temperature_handler(req: Request):
    """Handle temperature analysis requests from MongoDB change streams"""
    try:
        data = json.loads(req.data.decode())

        # Expect MongoDB change stream format
        if "operationType" not in data or "fullDocument" not in data:
            # type: ignore
            # pyright: ignore[reportArgumentType]
            await req.respond_error(code=400, description="Expected MongoDB change stream event")
            return

        # Parse MongoDB change stream event
        operation_type = data["operationType"]
        full_doc = data["fullDocument"]

        # Extract temperature data
        temperature = full_doc.get("temperature")
        sensor_id = full_doc.get("sensor_id")
        location = full_doc.get("location")

        if temperature is None:
            # type: ignore
            # pyright: ignore[reportArgumentType]
            await req.respond_error(code=400, description="Missing temperature field in document")
            return

        # Classify the temperature
        status, description = classify_temperature(temperature)

        # Create response
        response = TemperatureResponse(
            status=status,
            temperature=temperature,
            description=description,
            sensor_id=sensor_id,
            location=location
        )

        # Send response
        await req.respond(json.dumps(response.to_dict()).encode())
        logger.info("Processed change stream event: operation=%s temp=%.1fC status=%s sensor=%s",
                    operation_type, temperature, status, sensor_id or "N/A")

    except json.JSONDecodeError as e:
        logger.error(f"Error parsing JSON: {e}")
        # type: ignore
        await req.respond_error(code=400, description="Invalid JSON")
    except Exception as e:
        logger.error(f"Error processing request: {e}")
        # type: ignore
        await req.respond_error(code=500, description="Internal error")


async def main():
    """Main entry point"""
    # Connect to NATS
    nats_url = os.getenv("NATS_URL", "nats://localhost:4222")

    # Build connection options
    connect_opts = {}

    # Check for NGS credentials
    creds_path = os.getenv("NATS_CREDS")
    if creds_path:
        connect_opts["user_credentials"] = creds_path
        logger.info(f"Using NATS credentials from: {creds_path}")

    # Add performance tuning options
    connect_opts["max_outstanding_pings"] = 5
    connect_opts["ping_interval"] = 120

    nc = await nats.connect(nats_url, **connect_opts)
    logger.info(f"Connected to NATS at {nats_url}")

    # Create microservice
    service = await add_service(
        nc,
        name="temperature-analyzer",
        version="1.0.0",
        description="Analyzes temperature data and classifies as cold, warm, or hot"
    )

    # Add endpoint
    await service.add_endpoint(
        name="analyze",
        handler=temperature_handler,
        subject="temperature.analyze"
    )

    logger.info("Temperature Analysis Service started - subject: temperature.analyze | service: temperature-analyzer v1.0.0 | Ready to process requests")

    # Setup signal handlers for graceful shutdown
    def signal_handler():
        logger.info("Shutting down service...")
        asyncio.create_task(shutdown(nc))

    loop = asyncio.get_event_loop()
    for sig in (signal.SIGTERM, signal.SIGINT):
        loop.add_signal_handler(sig, signal_handler)

    # Keep running
    try:
        await asyncio.Future()  # Run forever
    except asyncio.CancelledError:
        pass


async def shutdown(nc: NATS):
    """Gracefully shutdown the service"""
    await nc.drain()


if __name__ == "__main__":
    asyncio.run(main())
