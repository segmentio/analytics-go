
Readme.md: analytics.go
	godocdown --heading analytics-go > $@

clean:
	rm -f Readme.md

.PHONY: clean