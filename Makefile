.PHONY: tidy install uninstall run apply delete status resources

tidy:
	go mod tidy

install:
	kubectl apply -f config/crd/simpleapps.yaml

uninstall:
	kubectl delete -f config/crd/simpleapps.yaml --ignore-not-found

run:
	go run .

apply:
	kubectl apply -f config/samples/apps_v1_simpleapp.yaml

delete:
	kubectl delete -f config/samples/apps_v1_simpleapp.yaml --ignore-not-found

status:
	kubectl get simpleapps

resources:
	kubectl get simpleapps,deployments,pods
