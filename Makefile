all: push

TAG = latest
PREFIX = quay.io/pulsecode/kubernetes-haproxy
HAPROXY_IMAGE = contrib-haproxy
SECRET_CONFIG = "secrets.json"
SECRET = haproxy-stats-secrets.json

server: service_loadbalancer.go
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w' -o service_loadbalancer ./service_loadbalancer.go ./loadbalancer_log.go

container: server haproxy
	docker build -t $(PREFIX):$(TAG) .

push: container
	gcloud docker push $(PREFIX):$(TAG)

haproxy:
	docker build -t $(HAPROXY_IMAGE):$(TAG) build
	docker create --name $(HAPROXY_IMAGE) $(HAPROXY_IMAGE):$(TAG) true
	# docker cp semantics changed between 1.7 and 1.8, so we cp the file to cwd and rename it.
	docker cp $(HAPROXY_IMAGE):/work/x86_64/haproxy-1.6-r0.apk .
	docker rm -f $(HAPROXY_IMAGE)
	mv haproxy-1.6-r0.apk haproxy.apk

secret:
	godep go run make_secret.go -config $(SECRET_CONFIG) > $(SECRET)

clean:
	rm -f service_loadbalancer haproxy.apk
	# remove servicelb and contrib-haproxy images
	docker rmi -f $(HAPROXY_IMAGE):$(TAG) || true
	docker rmi -f $(PREFIX):$(TAG) || true
