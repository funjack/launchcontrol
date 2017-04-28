--[[
Launchcontrol Extension for VLC

Copyright 2017 Funjack

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
may be used to endorse or promote products derived from this software without
specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
--]]

--[[ VLC extension hooks ]]--

function descriptor()
  return {
    title = "Launchcontrol 0.0.1",
    version = "0.0.1",
    author = "Funjack",
    url = "https://github.com/funjack/launchcontrol/",
    shortdesc = "Launchcontrol",
    description = [[
      Extension that will scan for haptic scripts on playback. The haptics will
      be loaded into a Launchcontrol server and played in sync. Actions like
      pausing playback or seeking to a new postion are relayed to the
      Launchcontrol server.
    ]],
    capabilities = {"menu", "input-listener", "playing-listener"},
    icon = icon_string,
  }
end

function activate()
  vlc.msg.dbg("[Launchcontrol] Get ready for the Launch")

  -- all input changed when a file al already open
  if vlc.input.item() then
    input_changed()
  end
 
  vlc.msg.dbg("[Launchcontrol] Activated")
end

function close()
  vlc.msg.dbg("[Launchcontrol] Close")
  vlc.deactivate()
end

function deactivate()
  vlc.msg.dbg("[Launchcontrol] Stopped")
end

function menu()
  return {
    "Test connection",
  }
end

function meta_changed()
  return false
end

function trigger_menu(id)
  vlc.msg.dbg("[Launchcontrol] Trigger Menu")
  -- Test connection
  if id == 1 then
    http_post("http://localhost:6969/v1/play",
              "x-text/kiiroo",
              "{0.50:4,1.00:0,2.50:4,3.00:0}")
  end
end

function input_changed()
  vlc.msg.dbg("[Launchcontrol] Input changed")
  local item = vlc.input.item()
  if item then
    local uri = item:uri()
    if uri then
      vlc.msg.dbg("[Launchcontrol] Searching script for "..uri)
      local data, mediaType = read_script(uri)
      if data then
        vlc.msg.dbg("[Launchcontrol] Found "..mediaType.." script!")
        launch_play(data, mediaType)
        launch_skip_to_current_time()
      end
    end
  else
    launch_stop()
  end
end

function playing_changed(status)
  vlc.msg.dbg("[Launchcontrol] Playing changed "..status)
  if status == 2 then launch_resume() end
  if status == 3 then launch_pause() end
end


--[[ Actions ]]--

--- Skip to the current time code.
function launch_skip_to_current_time()
  local input = vlc.object.input()
  if input then
    -- XXX will change in future versions form sec to microsec.
    time = math.floor(vlc.var.get(input, "time")*1000)
    launch_skip(time)
  end
end

--- read_script uses VLC stream() to detect and read a scrip file.
function read_script(file)
  local dotSplit = {}
  for p in string.gmatch(file, "[^\\.]+") do
    table.insert(dotSplit, p)
  end
  if #dotSplit > 1 then
    table.remove(dotSplit)
  end

  local baseFilename = table.concat(dotSplit, ".")

  for _, scriptType in ipairs(scriptTypes) do
    for _, extension in ipairs(scriptType["extensions"]) do
      s, err = vlc.stream(baseFilename.."."..extension)
      if s then
        local data = ""
        local line = s:readline()
        while line do
          data = data..line
          line = s:readline()
        end
        if data ~= "" then
          return data, scriptType["mediaType"]
        end
      else
        vlc.msg.dbg("[Launchcontrol] could not open file: " .. err)
      end
    end
  end
  return nil
end

--[[ Launch client ]]--

scriptTypes = {
  {
    name = "kiiroo",
    extensions = {"kiiroo"},
    mediaType = "x-text/kiiroo",
  },
  {
    name = "realtouch",
    extensions = {"realtouch", "ott"},
    mediaType = "x-text/realtouch",
  },
  {
    name = "vorze",
    extensions = {"vorze", "csv"},
    mediaType = "x-text/vorze",
  },
}

clientConfig = {
  url = "http://127.0.0.1:6969",
  latency = 0,
  positionMin = 0,
  positionMax = 100,
  speedMin = 20,
  speedMax = 100,
}

--- Play by sending data as specified mediatype.
-- @param data      Raw script data.
-- @param mediaType Mimetype of the script in data
function launch_play(data, mediaType)
  -- TODO personalization
  http_post(clientConfig["url"].."/v1/play", mediaType, data)
end

--- Stop playback.
function launch_stop()
  http_get(clientConfig["url"].."/v1/stop")
end

--- Pause playback.
function launch_pause()
  http_get(clientConfig["url"].."/v1/pause")
end

--- Resume playback.
function launch_resume()
  http_get(clientConfig["url"].."/v1/resume")
  launch_skip_to_current_time()
end

--- Skip jumps to a timecode
-- @param time Time position in script to jump to in ms.
function launch_skip(time)
  http_get(clientConfig["url"].."/v1/skip?p="..time.."ms")
end


--[[ HTTP Client ]]--

--- HTTP GET
function http_get(url)
  return http_request("GET", url)
end

--- HTTP POST
-- @param url URL
-- @param contenttype Type in data
-- @param data        Data to send
function http_post(url, contenttype, data)
  local cHdr = "Content-Type: "..contenttype
  return http_request("POST", url, cHdr, data)
end

function http_request(method, url, headers, body)
  local u
  if vlc.strings.url_parse then
    u = vlc.strings.url_parse(url)
  else
    u = vlc.net.url_parse(url)
  end

  if u["protocol"] ~= "http" then return false end

  local host, path, port = u["host"], u["path"], u["port"]
  local header = {
    string.upper(method).." "..path.." HTTP/1.0",
    "Host: "..host,
  }
  if body then table.insert(header, "Content-Length: "..#body) end
  if headers then
    if type(headers) == "table" then
      for v in headers do
        table.insert(header, v)
      end
    else
      table.insert(header, headers)
    end
  end
  -- header break
  table.insert(header, "")
  table.insert(header, "")
  local request = table.concat(header, "\r\n")

  if body then
    request = request..body
  end

  --return status, response
  return http_execute(host, port, request)
end

function http_execute(host, port, request)
  local fd = vlc.net.connect_tcp(host, tonumber(port))
  if not fd then return false end
  local pollfds = {}
 
  pollfds[fd] = vlc.net.POLLIN
  vlc.net.send(fd, request)
  vlc.net.poll(pollfds)
 
  local chunk = vlc.net.recv(fd, 2048)
  local response = ""
  local headerStr, header, body
  local status
 
  while chunk do
    response = response..chunk
    if not header then
      headerStr, body = response:match("(.-\r?\n)\r?\n(.*)")
      if headerStr then
        response = body
        header = http_parse_header(headerStr)
        status = tonumber(header["statuscode"])
      end
    end
    vlc.net.poll(pollfds)
    chunk = vlc.net.recv(fd, 1024)
  end

  vlc.net.close(fd)
  return status, response
end

function http_parse_header(data)
  local header = {}
  for name, s, val in string.gmatch(data, "([^%s:]+)(:?)%s([^\n]+)\r?\n") do
    if s == "" then
      header['statuscode'] = tonumber(string.sub(val, 1 , 3))
    else
      header[name] = val
    end
  end
  return header
end

icon_string = "\137\80\78\71\13\10\26\10\0\0\0\13\73\72\68\82\0\0\0\32\0\0\0\32\8\6\0\0\0\115\122\122\244\0\0\0\6\98\75\71\68\0\228\0\150\0\148\42\235\73\127\0\0\0\9\112\72\89\115\0\0\11\19\0\0\11\19\1\0\154\156\24\0\0\0\7\116\73\77\69\7\225\4\27\13\36\6\115\7\195\25\0\0\0\25\116\69\88\116\67\111\109\109\101\110\116\0\67\114\101\97\116\101\100\32\119\105\116\104\32\71\73\77\80\87\129\14\23\0\0\1\141\73\68\65\84\88\195\189\151\81\110\194\48\12\134\29\175\90\197\99\159\144\246\100\229\6\168\23\224\10\200\231\172\118\3\196\5\184\66\84\237\2\60\79\20\135\7\136\88\105\25\113\211\244\151\42\129\19\236\207\142\235\82\3\211\229\159\190\155\41\78\204\164\192\8\192\59\238\25\155\239\6\64\244\62\117\0\8\158\119\12\206\185\209\101\107\109\0\49\57\0\60\243\35\248\233\116\234\45\86\85\245\128\104\154\104\223\133\34\123\112\206\13\2\7\5\187\115\14\0\33\28\199\108\242\204\236\137\200\223\155\239\229\69\68\158\153\253\72\147\166\1\108\54\155\183\193\195\245\103\111\76\97\227\118\89\107\163\105\173\181\177\158\99\183\229\19\198\222\122\90\241\142\1\240\253\49\168\42\192\204\64\68\47\215\137\8\152\121\230\10\60\169\174\235\81\8\34\130\186\174\251\70\153\99\14\200\237\30\15\131\38\64\12\130\101\235\129\145\169\23\163\87\227\122\94\0\196\100\104\76\202\72\36\41\123\21\128\230\24\52\123\11\205\185\30\143\199\183\19\81\147\189\26\32\4\176\214\194\225\112\232\217\183\219\173\58\248\228\81\188\223\239\239\45\112\1\145\75\207\182\8\64\255\255\140\73\242\80\36\13\17\196\101\6\209\40\121\89\130\136\128\136\64\81\150\89\1\76\251\211\14\140\221\249\60\250\57\232\254\27\147\165\2\68\4\210\117\128\136\128\136\32\93\247\239\83\50\91\15\200\69\18\94\73\18\1\62\87\43\232\62\126\111\78\86\229\242\0\95\235\245\114\143\227\156\138\61\61\31\154\47\70\109\219\70\251\87\1\228\72\80\245\110\152\163\186\87\49\209\167\2\90\159\6\188\0\0\0\0\73\69\78\68\174\66\96\130"
