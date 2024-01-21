require("env")
local yab = require("yab")

if yab.os_type() == "windows" then
	local bin_name = "install-yab.exe"
	yab.task(yab.find("./cmd/installer/**.go"), bin_name, function()
		os.execute('go build -ldflags="-s -w -H=windowsgui" -o ' .. bin_name .. " ./cmd/installer/")
	end)
	return bin_name
else
	local bin_name = "install-yab"
	yab.task(yab.find("./cmd/installer/**.go"), bin_name, function()
		os.execute('go build -ldflags="-s -w" -o ' .. bin_name .. " ./cmd/installer/")
	end)
	return bin_name
end
