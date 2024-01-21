require("env")

local yab = require("yab")

yab.task(yab.find("**.go"), "./DOCS.md", function()
	os.execute("go run ./cmd/docs")
end)
