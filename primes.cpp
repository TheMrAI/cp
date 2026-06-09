#include <bits/stdc++.h>
#include <vector>

using i64 = int64_t;
using u64 = uint64_t;
using i32 = int32_t;
using u32 = uint32_t;
using u8 = uint8_t;

constexpr char brl = '\n';

using namespace std;

// Sieve of Eratosthenes
//
// Calculate the minimum prime factors for each number
// [0, N].
// Complexity: O(n log n)
// Space: N
auto sieve_min_factors(u64 N) -> vector<u64> {
    auto min_factors = vector<u64>(N+1, 0);
    for (u64 i = 2; i <= N; ++i) {
        if (min_factors[i] > 0) {
            continue;
        }
        min_factors[i] = i;
        for (u64 j = i * i; j <= N; j += i) {
            if (min_factors[j]) {
                continue;
            }
            min_factors[j] = i;
        }
    }
    return min_factors;
}

// Sieve of Eratosthenes
//
// Calculate the minimum prime factors for each number
// [0, N].
// Return the individual primes as well.
// Complexity: O(n)
// Space: 2N
auto sieve_primes(u64 N) -> pair<vector<u64>, vector<u64>> {
    auto min_factors = vector<u64>(N+1, 0);
    auto primes = vector<u64>();
    for (u64 i = 2; i <= N; ++i) {
        if (min_factors[i] == 0) {
            min_factors[i] = i;
            primes.push_back(i);
        }
        for (u64 j = 0; i * primes[j] <= N; ++j) {
            min_factors[i * primes[j]] = primes[j];
            if (primes[j] == min_factors[i]) {
                break;
            }
        }
    }
    return {min_factors, primes};
}
