#!/bin/bash
exec docker run -p 6831:6831/udp \
		-p 6832:6832/udp \
		-p 16686:16686 \
		-p 14268:14268 \
		-p 9411:9411 \
                jaegertracing/all-in-one:1.11
