#!/usr/bin/env python3.7

# Example private key.
private_key = 0x1E99423A4ED27608A15A2616A2B0E9E52CED330AC530EDCC32C8FFC6A526AEDD


# The characteristic of secp256k1; the order of the corresponding finite field.
p = 115792089237316195423570985008687907853269984665640564039457584007908834671663

# Check that the value of `p` given in the book matches the hexadecimal value
# given on en.bitcoin.it.
p_hex = 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f
assert p == p_hex

# Check that the value of `p` given in the book matches the mathematical
# expression given in the NIST spec.
p_math = 2**256 - 2**32 - 2**9 - 2**8 - 2**7 - 2**6 - 2**4 - 2**0
assert p == p_math


# Example point on the curve 'secp256k1'.
P = (
    55066263022277343669578718895168534326250603453777594175500187360389116729240,
    32670510020758816978083085130507043184471273380659243275938904335757337482424,
    )

# Check that the above point actually lies on the elliptic curve
#     y^2 = x^3 + ax + b.
x = P[0]
y = P[1]
left_side = (y**2) % p
right_side = (x**3 + 7) % p
print(f'y^2 = {left_side}')
print(f'x^3 + ax + b = {right_side}')
assert left_side == right_side

# For secp256k1, the elliptic curve used is:
#     y^2 = x^3 + 7
# I.e. the curve parameters `a` and `b` are:
#     a = 0
#     b = 7
# See:
#     https://en.bitcoin.it/wiki/Secp256k1
a = 0
b = 7

# Elliptic curve arithmetic involves a special "point at infinity", usually
# represented by O, to serve as the additive identity.  I.e., if P is some
# point on the elliptic curve then
#     P + O = O + P = P
#     P - P = O
#     O + O = O
# We set O to be the string 'identity' and catch it as a special case for each
# operation.
O = 'identity'

def inv(P):
    if P == O:
        return P
    else:
        return (P[0], (-P[1])%p)

# Elliptic curve point addition.  See:
#     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Point_addition
def add(P, Q):
    if P == O:
        return Q
    elif Q == O:
        return P
    elif Q == inv(P):
        return O
    else:
        if P == Q:
            l = 3*P[0]**2 * pow(2*P[1], p-2, p)
        else:
            l = (Q[1] - P[1]) * pow(Q[0] - P[0], p-2, p)
        x = (l**2 - P[0] - Q[0]) % p
        y = (l*(P[0] - x) - P[1]) % p
        return (x, y)

print(P)
print(inv(P))
print(add(P, inv(P)))
print(add(P, add(P, add(P, P))))
print(add(add(P, P), add(P, P)))
