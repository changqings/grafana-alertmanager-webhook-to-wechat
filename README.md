# Copy from https://github.com/n0vad3v/g2ww and update

# G2WW (Grafana 2 Wechat Work)
> Proxy Grafana Webhook alert to WeChat Work.

Grafana doesn't support push alert to WeChat Work(企业微信) by it's design, this is a small adapter for supporting this.


## Build g2ww

```
go build *.go -o g2ww
```

Then g2ww will listen on `localhost:2408`, quite simple isn't it?

## Run g2ww

Run `g2ww` on server, it will listen on `http://0.0.0.0:2408` by default, keep it running in background (`systemd` or `screen`?).

## Run on k8s

kubectl apply -f ./k8s-deploy/deployment.yaml

## Let Nginx to proxy it

Like this:

```
server {
        listen 80;
        server_name some.domain.name;

        location / {
            proxy_pass http://127.0.0.1:2408;
        }
}
```

## Create a Wechat Work Bot

Create a Wechat Work Bot and get the webhook address.

![](./img/ww-bot.png)

For instance, the webhook address is `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=<wechat-secretkey>`.

## Configure Grafana

In the configuration above, we need to specify the address like this(the url second path should only be one of [grafana|alertmanager] ):

`https://g2ww.nova.moe/[grafana|alertmanager]/<wechat-secretkey>`

![](./img/grafana.png)

## Demo

![](./img/demo.png)

Quite simple, isn't it?