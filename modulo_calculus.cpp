#include <bits/stdc++.h>

using namespace std;

using i64 = int64_t;
using u64 = uint64_t;
using i32 = int32_t;
using u32 = uint32_t;
using u8 = uint8_t;

constexpr char brl = '\n';

constexpr u64 PRIME = 998244353;

// Binary exponentiation % prime
// Complexity: O(log(b))
auto bin_pow_mod(u64 a, u64 b, u64 prime) -> u64 {
    u64 result = 1;
    a = a % prime;
    while (b > 0) {
        if (b & 1) {
            result = result * a % prime;
        }
        a = a * a % prime;
        b >>= 1;
    }
    return result;
}

// Inverse of n modulo prime: n^-1 % prime
// Uses binary exponentiation.
// Complexity O(log(prime))
auto inv_mod(u64 n, u64 prime) {
    return bin_pow_mod(n, prime - 2, prime);
}

constexpr u64 MAXN = 2*1e5;
// Precalculate factorial(i) in range [0, MAXN]
// All factorials are stored as % prime
u64 factorial[MAXN + 1];
// We can precalculate the inverse factorials as well
// if necessary.
u64 inverse_factorial[MAXN + 1];

auto precalcualte_factorials(u64 prime) {
    factorial[0] = 1;
    for (i64 i = 1; i <= MAXN; ++i) {
        factorial[i] = factorial[i-1] * i % prime;
    }

    inverse_factorial[MAXN] = inv_mod(factorial[MAXN], prime);
    for (i64 i = MAXN-1; i >= 0; --i) {
        inverse_factorial[i] = inverse_factorial[i + 1] * (i + 1) % prime;
    }
}

// Binomial coefficient (n chose k) % prime
// Prerequisite: precalculated factorials
// Complexity O(log(prime))
auto binomial_coefficient_mod(u64 n, u64 k, u64 prime) -> u64 {
    if (k > n) {
        return 0;
    }
    return (factorial[n] * inv_mod(factorial[k] * factorial[n-k], prime)) % prime;
}

// Binomial coefficient (n chose k) % prime
// Prerequisite: precalculated factorials, inverse factorials
// Complexity O(log(prime))
auto binomial_coefficient_mod_precalc(u64 n, u64 k, u64 prime) -> u64 {
    if (k > n) {
        return 0;
    }
    return (factorial[n] * inverse_factorial[k] % prime * inverse_factorial[n-k]) % prime;
}
