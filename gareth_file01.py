#!/usr/bin/env python3.7

# Example private key.
private_key = 0x1E99423A4ED27608A15A2616A2B0E9E52CED330AC530EDCC32C8FFC6A526AEDD


# The characteristic of secp256k1; the order of the corresponding finite field.
p = 115792089237316195423570985008687907853269984665640564039457584007908834671663

# Check that the value of `p` given in the book matches the hexadecimal value
# given on en.bitcoin.it.
# source of hex number: https://en.bitcoin.it/wiki/Secp256k1
p_hex = 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f
assert p == p_hex

# Check that the value of `p` given in the book matches the mathematical
# expression given in the NIST spec.
p_math = 2**256 - 2**32 - 2**9 - 2**8 - 2**7 - 2**6 - 2**4 - 2**0
print(p_math)
assert p == p_math


# Example point on the curve 'secp256k1'.
P = (
    55066263022277343669578718895168534326250603453777594175500187360389116729240,
    32670510020758816978083085130507043184471273380659243275938904335757337482424,
    )

# Check that the above point actually lies on the elliptic curve
#     y^2 = x^3 + ax + b.
#		x = 0
# 		b = 7
x = P[0]
y = P[1]
left_side = (y**2) % p
right_side = (x**3 + 7) % p
print(f'y^2 = {left_side}')
print(f'x^3 + ax + b = {right_side}')
assert left_side == right_side
