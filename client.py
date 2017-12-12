import minclient
import time

name = "alice"
provider = "example.com"
keyStr = "4d488962dd5a7c2d2d2360a6bbe258bf75022eb39a05b8c877f3f92e99fd298c"
pkiAddr = "192.0.2.1:29483"
pkiKey = "900895721381C0756D28954524BB1D090F54C8DD9295F84B1D8A93F1E3C17AD8"

client = minclient.NewClient(pkiAddr, pkiKey, minclient.LogConfig())
key = minclient.StringToKey(keyStr)
session = client.NewSession(name, provider, key)

sent = False
while not sent:
    try:
        session.SendMessage("bob", "panoramix.org", "hello bob!!!")
        time.sleep(1)
        sent = True
    except RuntimeError:
        pass
print("Message sent")

while 1:
    print(session.GetMessage())
