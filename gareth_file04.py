#!/usr/bin/env python3.7
#
# Please note that the following code is a demo.  Edge cases and error
# checking are intentionally omitted where they might otherwise distract
# us from the core ideas.
#     "Before we can understand how something can go wrong, we must
#      learn how it can go right."

print('Example from Mastering Bitcoin, pages 69-70.')


# A private key must be a whole number from 1 to
#     0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140,
# one less than the order of the base point (or "generator point") G.
# See:
#     https://en.bitcoin.it/wiki/Private_key#Range_of_valid_ECDSA_private_keys
private_key = 0x038109007313a5807b2eccc082c8c3fbb988a973cacf1a7df9ce725c31b14776
print(f'\tprivate_key:\n{private_key}')


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
        slope = 3*P.x**2 * pow(2*P.y, p-2, p)  # 3Px^2 / 2Py
    else:
        slope = (Q.y - P.y) * pow(Q.x - P.x, p-2, p)  # (Qy - Py) / (Qx - Px)
    R = Point()
    R.x = (slope**2 - P.x - Q.x) % p  # (slope^2 - Px - Qx)
    R.y = (slope*(P.x - R.x) - P.y) % p  # slope*(Px - Rx) - Py
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

# Generate the public key.  See:
#     Mastering Bitcoin, page 63.
public_key = ec_point_multiply(private_key, G)
print(f'\tpublic_key:\nx={public_key.x}\ny={public_key.y}')


# The compressed serialisation of the public key.  See:
#     Mastering Bitcoin, pages 73-75.
#     https://www.ntirawen.com/2019/03/bitcoin-compressed-and-uncompressed.html
#     https://www.secg.org/sec2-v2.pdf - Section 2.4.1.
#     https://www.secg.org/sec1-v2.pdf - Section 2.3.3.
def serialize(K):
    if K.y % 2:
        leading_byte = b'\x03'
    else:
        leading_byte = b'\x02'
    return leading_byte + K.x.to_bytes(32, byteorder='big')

serialized_public_key = serialize(public_key)
print(f'\tserialized_public_key:\n{serialized_public_key.hex()}')


import hashlib  # For sha256 and ripemd160.

def sha256(data):
    return hashlib.new('sha256', data).digest()

def ripemd160(data):
    return hashlib.new('ripemd160', data).digest()

public_key_hash = ripemd160(sha256(serialized_public_key))
print(f'\tpublic_key_hash:\n{public_key_hash.hex()}')


# Calculate the checksum needed for Bitcoin's Base58Check format.  See:
#     Mastering Bitcoin, page 58
#     https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address - Steps 5-7.
version = b'\x00'
checksum = sha256(sha256(version + public_key_hash))[:4]
print(f'\tchecksum:\n{checksum.hex()}')


# Encode the data in Bitcoin's Base58 format.  See:
#     https://en.bitcoin.it/wiki/Base58Check_encoding#Base58_symbol_chart
def base58(data):
    alphabet = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz'
    x = int.from_bytes(data, byteorder='big')
    result = []
    while x:
        x, i = divmod(x, 58)
        result.append(alphabet[i])
    for byte in data:
        if byte:
            break
        result.append(alphabet[0])
    result.reverse()
    return ''.join(result)

# A Bitcoin address is just the public key hash encoded in Bitcoin's
# Base58Check format.  See:
#     Mastering Bitcoin, page 66.
address = base58(version + public_key_hash + checksum)
print(f'\taddress:\n{address}')
