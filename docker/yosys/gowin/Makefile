
# default directive
all: synth pnr packing

synth:
	cd  $(SRC) \
	&& yosys -p 'synth_$(BOARD_MODEL) -top $(MODULE_NAME) -json $(MODULE_NAME).json' $(SYNTH_FILE)

pnr:
	cd $(SRC) \
	&& nextpnr-$(BOARD_MODEL) --package tq144 --hx1k --json $(MODULE_NAME).json --pcf $(MODULE_NAME).pcf --asc $(MODULE_NAME).asc

packing:
	cd $(SRC) \
	&& icepack $(MODULE_NAME).asc $(OUTFILE)
