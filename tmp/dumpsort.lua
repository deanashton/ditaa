--fn = 'out'

function sep(imgs, img)
	if img then
		imgs[#imgs+1] = img
	end
	img = nil
	
	table.sort(imgs)
	for _, img in ipairs(imgs) do
		print(img)
	end
	imgs = {}
	
	return imgs, img
end

function main()
--	local fh = assert(io.open(fn, 'r'))
	local fh = io.stdin

	local imgs = {}
	local img, prefix

	for line in fh:lines() do
		if line == "Closed boundaries:" or line == "Open boundaries:" then
			prefix = line
		elseif line:match '^    0123456789' then
			--print("START")
			-- store previous img
			if img then
				imgs[#imgs+1] = img
			end
			if not prefix then
				prefix = ""
			end
			img = prefix .. '\n' .. line
			prefix = nil
			--print(line)
		elseif line:match '^ *[0-9]+ %(' or line:match '^%-%-%-' then
			--print("CONT")
			img = img .. '\n' .. line
			--print(line)
		else
			--print("SEP")
			imgs, img = sep(imgs, img)
			print(line)
		end
	end
	
	sep(imgs, img)

--	fh:close()
end

main()
