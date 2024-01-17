
# default directive
all: synth pnr packing

synth:
	cd  $(SRC) \
	&& yosys -p 'synth_ice40 -top $(MODULE_NAME) -json $(MODULE_NAME).json' $(SYNTH_FILE)

pnr:
	cd $(SRC) \
	&& nextpnr-ice40 --package tq144 --hx1k --json $(MODULE_NAME).json --pcf $(CONSTRAINT_FILE) --asc $(MODULE_NAME).asc

packing:
	cd $(SRC) \
	&& icepack $(MODULE_NAME).asc $(OUTFILE)
