# outlook-monitor
Simple Go client to check unread Outlook emails every 30 seconds and send a desktop notification for Linux users.

## Setup

```bash
go build
```

## Usage

```bash
./outlook_monitor &
```

## Logs

Information from the email monitor will be logged to `log/outlook.log` every
30 seconds. To see what it's doing, just check the log:

```bash
cat log/outlook.log
```

## Stop Monitor

```bash
ps aux | grep outlook_monitor
kill -9 <PID_FROM_OUTPUT_ABOVE>
```
