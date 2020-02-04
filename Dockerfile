FROM scratch

WORKDIR $GOPATH/gin-blog
COPY    . $GOPATH/gin-blog

EXPOSE 8000
CMD ["./gin-blog"]