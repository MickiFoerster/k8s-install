GOFILES=$(shell ls -1 *.go)

k8s-install : $(GOFILES)
	go build -o $@ $^


clean:
	rm -f k8s-install

.PHONY: clean
