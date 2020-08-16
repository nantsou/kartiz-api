# kartiz-api

---

### Quick Guide

**build**

```
# build the image
$ docker build -t kartiz .

# clear unnecesary images
$ docker image prune
```

**run docker**

```
# put the configure file in $(pwd)/config/{file_name} for docker as bellow
$ docker run --rm -p 8080:8080 -v $(pwd)/config/docker.json:/config.json kartiz
```

**note**

```
# build image for docker hub
$ docker build -t ddrsmile/kartiz-api:<version> .

# push the image to docker hub
$ docker push ddrsmile/kartiz-api:<version>
```