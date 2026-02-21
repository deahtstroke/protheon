record.sum = record["sum"].A + record["sum"].B
record["hello world"] = record["hello world"].foo .. " " .. record["hello world"].bar

local total = 0
for _, value in ipairs(record.sums) do
	total = total + value.A + value.B
end
record.total = total
