root.key = this.timestamp.string()
root.data = {
    "data": {
        "sensor_id": this.sensor_id,
        "temperature": this.temperature,
        "unit": this.unit,
        "location": this.location,
        "original_timestamp": this.timestamp
    }
}