
Readme.md: analytics.go
	godocdown > $@

clean:
	rm -f Readme.md

.PHONY: clean