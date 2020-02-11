# auth-get
- A tool for getting compass hosts's authentication file.
- work by compass machine auth api.
- only worked in compass cluster(with machine auth api)
- usage example:
```bash
$ docker run --rm -it --network=host clayz95/auth-get sh
/opt/auth-get # ./auth-get -h
A tool for getting host authentication file.
Work by compass machine auth api.

Usage:
  auth-get [flags]

Flags:
  -h, --help              help for auth-get
  -m, --masterIp string   Master VIP or control cluster's master IP
  -o, --output string     output auth file by yaml, json or ansible inventory (yaml, json, inventory)  (default "yaml")
  -p, --password string   Master VIP or control cluster's master IP  (default "Pwd123456")
      --port string       Set SSH port in output file  (default "22")
  -u, --username string   Compass username (default "admin")

/opt/auth-get # ./auth-get -m 192.168.17.21
```