# Pet with peers

Mini simple [pet](http://cloudscaling.com/blog/cloud-computing/the-history-of-pets-vs-cattle/#pets) application written in golang for container environement. 
Just testing purpose only. Can use curl with `POST` method to sent data and can request data with `GET` method from application.

### Build image 
```
$ make TAG_NAME=pet-with-peer
```


### Run on ` k8s`

```
$ kubectl create -f pvs-hostpath.yaml
$ kubectl create -f hola-statefulset-peers.yaml
$ kubectl create -f hola-service-headless.yaml
$ kubectl create -f hola-service-public.yaml
```

#### Reference 
Kubernetes in Action
