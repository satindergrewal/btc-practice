#!/usr/bin/env python3.7

import hashlib
try:
    import base58
except ImportError:
    pass


class Point:
    def __init__(self, x=0, y=0):
        self.x = x
        self.y = y


# Secp256k1 parameters.  See:
#     https://en.bitcoin.it/wiki/Secp256k1
#     https://www.secg.org/sec2-v2.pdf - Section 2.4.1.
p = 2**256 - 2**32 - 2**9 - 2**8 - 2**7 - 2**6 - 2**4 - 1
G = Point(
    0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798,
    0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8,
    )

# Elliptic curve point addition.  Unneeded side-cases are omitted for
# simplicity.  See:
#     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Point_addition
#     https://stackoverflow.com/a/31089415
#     https://en.wikipedia.org/wiki/Modular_multiplicative_inverse#Using_Euler's_theorem
#     https://crypto.stanford.edu/pbc/notes/elliptic/explicit.html
def ec_point_add(P, Q):
    if not P:
        return Q

    if P == Q:
        slope = 3*P.x**2 * pow(2*P.y, p-2, p)
    else:
        slope = (Q.y - P.y) * pow(Q.x - P.x, p-2, p)
    R = Point()
    R.x = (slope**2 - P.x - Q.x) % p
    R.y = (slope*(P.x - R.x) - P.y) % p
    return R

# Elliptic curve point multiplication.  This is an implimentation of the
# Double-and-add algorithm with increasing index described here:
#     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
def ec_point_multiply(d, P):
    N = P
    Q = None
    while d:
        if d & 1:
            Q = ec_point_add(Q, N)
        N = ec_point_add(N, N)
        d >>= 1
    return Q

# The compressed serialisation of the public key.  See:
#     https://www.ntirawen.com/2019/03/bitcoin-compressed-and-uncompressed.html
#     https://www.secg.org/sec2-v2.pdf - Section 2.4.1.
#     https://www.secg.org/sec1-v2.pdf - Section 2.3.3.
def serialize(K):
    return bytes([2 + K.y % 2]) + K.x.to_bytes(32, byteorder='big')


print('Example from Mastering Bitcoin, pages 69-70.')

# A private key must be a whole number from 1 to
#     0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140,
# one less than the order of the base point G.  See:
#     https://en.bitcoin.it/wiki/Private_key
private_key = 0x038109007313a5807b2eccc082c8c3fbb988a973cacf1a7df9ce725c31b14776
print(f'\tprivate_key:\n{private_key}')

public_key = ec_point_multiply(private_key, G)
print(f'\tpublic_key:\nx={public_key.x}\ny={public_key.y}')

serialized_public_key = serialize(public_key)
print(f'\tserialized_public_key:\n{serialized_public_key.hex()}')

hash1 = hashlib.sha256(serialized_public_key).digest()
hash2 = hashlib.new('ripemd160', hash1).digest()
public_key_hash = hash2
print(f'\tpublic_key_hash:\n{public_key_hash.hex()}')

try:
    address = base58.b58encode_check(b'\x00' + public_key_hash).decode()
except NameError:
    print()
    print('...Install the python package "base58" to see the address.')
else:
    print(f'\taddress:\n{address}')
