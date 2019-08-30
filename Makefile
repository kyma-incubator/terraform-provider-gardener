.PHONY: build
build: 
	./before-commit.sh ci

.PHONY: clean
clean:
	rm -rf bin