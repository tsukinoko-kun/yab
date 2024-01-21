require("env")

local yab = require("yab")

local bin_name = yab.os_type() == "windows" and "yab.exe" or "yab"

yab.task(yab.find("**.go"), bin_name, function()
	os.execute('go build -ldflags="-s -w" -o ' .. bin_name .. " ./cmd/yab/")
end)

return bin_name
