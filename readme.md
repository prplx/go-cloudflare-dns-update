## Cloudflare DNS update

#### Run

Create .env file and add the following variables to it

```
CF_API_URL=https://api.cloudflare.com/client/v4
CF_API_TOKEN=Cloudflare API token with updating DNS permission
CF_ZONE_ID=Cloudflare zone id
CF_DNS_RECORD_ID-Cloudflare record id. Multiple values are supported, separate them by a comma
```

run `make run`
