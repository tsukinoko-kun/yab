require("env")

local yab = require("yab")

local bin_name = "yab-" .. yab.os_type() .. "-" .. yab.os_arch()

if yab.os_type() == "windows" then
    bin_name = bin_name .. ".exe"
end

yab.task(yab.find("**.go"), bin_name, function()
	os.execute('go build -ldflags="-s -w" -o ' .. bin_name .. " ./cmd/yab/")
end)

return bin_name
