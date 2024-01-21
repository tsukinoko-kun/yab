require("env")
local yab = require("yab")

local bin_name = "install-yab-" .. yab.os_type() .. "-" .. yab.os_arch()

if yab.os_type() == "windows" then
	bin_name = bin_name .. ".exe"
end

if yab.os_type() == "windows" then
	yab.task(yab.find("./cmd/installer/**.go"), bin_name, function()
		os.execute('go build -ldflags="-s -w -H=windowsgui" -o ' .. bin_name .. " ./cmd/installer/")
	end)
else
	yab.task(yab.find("./cmd/installer/**.go"), bin_name, function()
		os.execute('go build -ldflags="-s -w" -o ' .. bin_name .. " ./cmd/installer/")
	end)
end

return bin_name
