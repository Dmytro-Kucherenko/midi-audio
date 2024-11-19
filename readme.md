# Midi Audio Control

## Setup

- Fill local environment: `.env`. Use reference: `.env.example`

- Inspect ports, apps: `make inspect`

- Fill binds using inspected connections: `binds.json`. Use  reference: `binds.example.json`

- Run and test: `make run`

## Run on Linux startup

- Build app: `make build`

- Create service with filled prod environment: `sudo nano /etc/systemd/system/midi-audio.service`
  ```
  [Unit]
  Description=Midi Audio Control
  After=pulseaudio.socket

  [Service]
  ExecStart=/home/dmytro/Projects/midi-audio/bin/main
  Environment="PULSE_SERVER=unix:/run/user/1000/pulse/native"
  Environment="DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus"
  Environment="DISPLAY=:0"
  Environment="XAUTHORITY=/home/dmytro/.Xauthority"
  Environment="BINDS_PATH=/path/to/midi-audio/binds.json"
  Environment="HISTORY_PATH=/path/to/midi-audio/history.json"
  User=dmytro
  Restart=always
  RestartSec=2
  StandardOutput=journal
  StandardError=journal

  [Install]
  WantedBy=multi-user.target
- Reload daemon: `sudo systemctl daemon-reload`

- Enable service startup: `sudo systemctl enable midi-audio`

- Start service: `sudo systemctl start midi-audio`

## Control startup

- Reload daemon: `sudo systemctl daemon-reload`

- Enable service startup: `sudo systemctl enable midi-audio`

- Start service: `sudo systemctl start midi-audio`

- Restart service: `sudo systemctl restart midi-audio`

- Stop service: `sudo systemctl stop midi-audio`

- Disable service startup: `sudo systemctl disable midi-audio`

- Check service status: `sudo systemctl status midi-audio`