FROM busybox:latest
MAINTAINER Ty Davis <tydavis@gmail.com>

EXPOSE 80
COPY showReadiness / 
CMD /showReadiness
