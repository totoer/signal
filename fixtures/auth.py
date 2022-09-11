
import hashlib

# HA1=MD5(username:realm:password)
# HA2=MD5(method:digestURI)
# response=MD5(HA1:nonce:HA2)

# 5a0854b2573eb60928de254c2d4dc8e0

def response(method, ruri, username, nonce, realm, password):
    ha1s = f"{username}:{realm}:{password}"
    ha2s = f"{method}:{ruri}"
    ha1 = hashlib.md5(ha1s.encode()).hexdigest()
    ha2 = hashlib.md5(ha2s.encode()).hexdigest()
    response = hashlib.md5(f"{ha1}:{nonce}:{ha2}".encode()).hexdigest()
    print("ha1:", ha1s)
    print("ha1:", ha2s)
    print(response)


response(method="REGISTER", ruri="sip:127.0.0.1:5080", username="test", nonce="test", realm="127.0.0.1:5080", password="test")