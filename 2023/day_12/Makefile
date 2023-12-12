INPUT_FILE := seed/12.txt
.PHONY: clean prepare execute
prepare: main
execute: main
	bun index.ts $(INPUT_FILE) o
watch: main
	bun --watch index.ts $(INPUT_FILE)
clean:
	-rm output.json
main:
	bun