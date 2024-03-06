require("env")

local bin_name = require("installer.build")

os.execute("./" .. bin_name)
