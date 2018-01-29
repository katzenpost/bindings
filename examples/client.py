#!/usr/bin/env python
# client.py - python example mixnet client
# Copyright (C) 2017  Yawning Angel.
# Copyright (C) 2018  Ruben Pollan.
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

import katzenpost

linkKey = "4d488962dd5a7c2d2d2360a6bbe258bf75022eb39a05b8c877f3f92e99fd298c"
key = katzenpost.StringToKey(linkKey)

cfg = katzenpost.Config(
    PkiAddress="192.0.2.1:29483",
    PkiKey="900895721381C0756D28954524BB1D090F54C8DD9295F84B1D8A93F1E3C17AD8",
    User="alice",
    LinkKey=key,
    Provider="example.com",
    Log=katzenpost.LogConfig()
)

c = katzenpost.New(cfg)

mail = """From: alice@example.com
To: bob@panoramix.com
Subject: hello

Hello there.
"""
c.Send("bob@panoramix.com", mail)

while True:
    try:
        m = c.GetMessage(1)
    except RuntimeError:
        continue
    print("=================>" + m.Sender)
    print(m.Payload)
