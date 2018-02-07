FROM scratch
LABEL maintainer="Tyler Davis <tydavis@gmail.com>"

EXPOSE 80
COPY showreadiness .
CMD ["./showreadiness"]
