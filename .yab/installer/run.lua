require("env")
local yab = require("yab")

local bin_name = require("installer.build")

if yab.os_type() == "windows" then
    os.execute(bin_name)
else
    os.execute("./" .. bin_name)
end
