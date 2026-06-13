#include <bits/stdc++.h>

using i64 = int64_t;
using u64 = uint64_t;
using i32 = int32_t;
using u32 = uint32_t;
using u8 = uint8_t;

constexpr char brl = '\n';

using namespace std;

auto solve() -> void
{
    i64 n, k;
    cin >> n >> k;
    auto cost = vector<i64>(n, 0);
    char prev;
    cin >> prev;
    for (i64 i = 1; i < n; ++i)
    {
        char cur;
        cin >> cur;
        if (prev != cur) { cost[i] = 1; }
        prev = cur;
    }

    i64 window = 0;
    for (i64 i = 1; i < 1 + k; ++i) { window += cost[i]; }
    i64 total = window;
    for (i64 i = 1 + k; i < n; ++i)
    {
        window += cost[i];
        window -= cost[i - k];
        total += window;
    }
    cout << total << brl;
}

auto main() -> int
{
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);

    auto t = i64{ 0 };
    std::cin >> t;
    while (t--) { solve(); }

    return 0;
}
