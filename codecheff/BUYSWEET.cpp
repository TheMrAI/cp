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
    i64 n, r;
    cin >> n >> r;

    auto delta_costs = vector<pair<i64, i64>>();
    {
        auto sweets = vector<i64>(n, 0);
        auto rebates = vector<i64>(n, 0);
        for (i64 i = 0; i < n; ++i)
        {
            i64 val;
            cin >> val;
            sweets[i] = val;
        }

        for (i64 i = 0; i < n; ++i)
        {
            i64 val;
            cin >> val;
            rebates[i] = val;
        }

        for (i64 i = 0; i < n; ++i) { delta_costs.push_back({ sweets[i] - rebates[i], sweets[i] }); }
    }
    sort(delta_costs.begin(), delta_costs.end());

    i64 sweet_count = 0;
    while (true)
    {
        i64 selection = -1;
        for (i64 i = 0; i < n; ++i)
        {
            if (delta_costs[i].second <= r)
            {
                selection = i;
                break;
            }
        }
        if (selection == -1) { break; }
        i64 times = 1 + ((r - delta_costs[selection].second) / delta_costs[selection].first);
        r -= delta_costs[selection].first * times;
        sweet_count += times;
    }
    cout << sweet_count << brl;
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
