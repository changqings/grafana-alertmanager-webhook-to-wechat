FROM golang:stretch

COPY g2ww /go/bin/g2ww

EXPOSE 2408

ENTRYPOINT ["g2ww"]
