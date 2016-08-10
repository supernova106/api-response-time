run:
	docker run -d --name=golang-base -v $(PWD)/app/:/go/src/app/ --net=host supernova106/golang-base:latest
jump:
	docker exec -it golang-base bash	
clean:
	docker rm -f golang-base
