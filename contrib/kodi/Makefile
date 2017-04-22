all: script.service.launchcontrol.zip

%.zip: %
	zip -FSr $@ $<

.PHONY: clean
clean:
	-rm *.zip
