FROM golang:alpine

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab

#RUN go get github.com/samuelandersson97/Kademlia-Lab.git

WORKDIR /go/src/main
COPY Kademlia-Lab-master/main/main.go /go/src/main/
COPY Kademlia-Lab-master/d7024e/* /go/src/d7024e/
RUN go build main.go
CMD ./main
