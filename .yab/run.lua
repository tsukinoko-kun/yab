require("env")
local yab = require("yab")

local args = yab.args()
local join = function(t)
	local s = ""
	for _, v in ipairs(t) do
		s = s .. " " .. v
	end
	return s
end

local bin_name = require("build")

os.execute("./" .. bin_name .. join(args))
