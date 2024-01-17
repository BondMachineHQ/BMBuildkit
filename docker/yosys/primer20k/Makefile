
# default directive
all: synth pnr packing

synth:
	cd  $(SRC) \
	&& yosys $(YOSYS_OPTIONS) -p "read_verilog $(SYNTH_FILE); synth_gowin -json $(MODULE_NAME)-synth.json"

pnr:
	cd $(SRC) \
	&& nextpnr-himbaechel --device GW2A-LV18PG256C8/I7 --vopt family=GW2A-18 --vopt cst=$(CONSTRAINT_FILE) --json $(MODULE_NAME)-synth.json --write $(MODULE_NAME).json

packing:
	cd $(SRC) \
	&& gowin_pack -d GW2A-18 -o $(OUTFILE) $(MODULE_NAME).json

