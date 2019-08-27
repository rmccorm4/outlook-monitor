# outlook-monitor
Simple Go client to check unread Outlook emails every 30 seconds and send a desktop notification for Linux users.

## Setup

```bash
git clone https://github.com/rmccorm4/outlook-monitor.git
cd outlook-monitor

go build
```

## Usage

You'll need to set the `OUTLOOK_EMAIL` and `OUTLOOK_PASSWORD` environment variables
for the client to login and check your emails.

```bash
# Set email/password as environment variables
export OUTLOOK_EMAIL=username@email.com
export OUTLOOK_PASSWORD=password

# Start monitor
./outlook-monitor &
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
