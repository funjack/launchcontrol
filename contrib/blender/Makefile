all: funscripting.zip

%.zip: %.py
	zip -FSr $@ $<

.PHONY: clean
clean:
	-rm *.zip
