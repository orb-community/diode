version: '3.7'

# Docker Services
services:

  service:
    image: orbcommunity/diode-service:IMAGE_TAG
    ports:
      - "4317:4317"
    environment:
      - DIODE_SERVICE_NETBOX_ENDPOINT=NETBOX_API_HOST
      - DIODE_SERVICE_NETBOX_TOKEN=NETBOX_API_TOKEN
