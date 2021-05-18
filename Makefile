build:
	docker build -t standoff_bot .

run: build
	docker run -d --name standoff_bot -v /tmp/cache:/app/cache standoff_bot

clean:
	docker rm -vf standoff_bot

rmi: clean
	docker rmi standoff_bot

re: rmi run