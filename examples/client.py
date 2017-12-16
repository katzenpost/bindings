#!/usr/bin/env python
# client.py - python example mixnet client
# Copyright (C) 2017  Yawning Angel.
# Copyright (C) 2017  Ruben Pollan.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

import minclient

name = "alice"
provider = "example.com"
keyStr = "4d488962dd5a7c2d2d2360a6bbe258bf75022eb39a05b8c877f3f92e99fd298c"
pkiAddr = "192.0.2.1:29483"
pkiKey = "900895721381C0756D28954524BB1D090F54C8DD9295F84B1D8A93F1E3C17AD8"

client = minclient.NewClient(pkiAddr, pkiKey, minclient.LogConfig())
key = minclient.StringToKey(keyStr)
session = client.NewSession(name, provider, key)
session.WaitToConnect()

session.SendMessage("bob", "panoramix.org", "hello bob!!!")
print("Message sent")

while 1:
    try:
        print(session.GetMessage(1))
    except RuntimeError:
        pass
