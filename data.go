timeout 30 mosquitto_sub -h 127.0.0.1 -p 1883 -t '#' -F '%t' 2>/dev/null |
awk '!seen[$0]++ {print; fflush()}'