## logrotate: slightly better than `>>`

`logrotate` is a naïve log rotator which reads logs from stdin and writes them
to a file, gzipping and truncating when it grows too large. If you have daemons
that log to stdout, you can pipe them into this and get rotated logfiles.
