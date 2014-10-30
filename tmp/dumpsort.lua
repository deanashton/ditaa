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

function rewind(line)
	local t = {}
	local min, imin = {x=999, y=999}, 1
	line:gsub('%(([0-9]+), ([0-9]+)%)', function(x, y)
		local x, y = 0+x, 0+y
		local p = {x=x, y=y}
		t[#t+1] = p
		if x<min.x or (x==min.x and y<min.y) then
			min, imin = p, #t
		end
	end)
	
	local function wrapped(i)
		return (i-1) % #t +1
	end
	
	local handed = 0
	for i=1, #t-2 do
		local j = imin+i-1
		local p0, p1, p2 = t[wrapped(j)], t[wrapped(j+1)], t[wrapped(i+2)]
		local u = {x=p1.x-p0.x, y=p1.y-p0.y}
		local v = {x=p2.x-p1.x, y=p2.y-p1.y}
		handed = 2*handed + u.x*v.y - u.y*v.x
	end
	
	local out = {}
	local n = #t
	if handed>=0 then
		for i=1, #t do
			local p = t[wrapped(imin+i-1)]
			out[#out+1] = '(' .. p.x .. ', ' .. p.y .. ')'
		end
	else
		for i=1, #t do
			local p = t[wrapped(imin-i+1+#t)]
			out[#out+1] = '(' .. p.x .. ', ' .. p.y .. ')'
		end
	end
	return table.concat(out, '/')
end

function main()
--	local fh = assert(io.open(fn, 'r'))
	local fh = io.stdin

	local imgs = {}
	local img
	local prefix = ''

	for line in fh:lines() do
		if line == "Closed boundaries:" or line == "Open boundaries:" or line:match '^ *$' then
			prefix = line
		elseif line:match '^    0123456789' then
			--print("START")
			-- store previous img
			if img then
				imgs[#imgs+1] = img
			end
			img = prefix .. '\n' .. line
			prefix = ''
			--print(line)
		elseif line:match '^ *[0-9]+ %(' or line:match '^%-%-%-' then
			--print("CONT")
			img = img .. '\n' .. line
			--print(line)
		else
			--print("SEP")
			imgs, img = sep(imgs, img)
			
			if line:match '^%([0-9(, )/]+%)$' then
				--print("REWIND!")
				line = rewind(line)
			end
			
			print(line)
		end
	end
	
	sep(imgs, img)

--	fh:close()
end

main()
