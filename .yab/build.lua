require("env")

local yab = require("yab")

local bin_name = "yab-" .. yab.os_type() .. "-" .. yab.os_arch()

yab.task(yab.find("**.go"), bin_name, function()
	os.execute('go build -ldflags="-s -w" -o ' .. bin_name .. " .")
end)

return bin_name
