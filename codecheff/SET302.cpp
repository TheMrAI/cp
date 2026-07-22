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
    u64 n, h;
    cin >> n >> h;
    auto ays = vector<u64>(n, 0);

    for (auto i = 0; i < n; ++i)
    {
        u64 a;
        cin >> a;
        ays[i] = a;
    }
    sort(ays.begin(), ays.end());
    // for (auto const a : ays) { cout << a << " "; }
    // cout << brl;

    u64 low = 1;
    u64 high = ays.back();
    while (low < high)
    {
        u64 mid = low + ((high - low) / 2);

        u64 chomp_count = 0;
        for (auto const& a : ays) { chomp_count += (a + mid - 1) / mid; }
        if (chomp_count <= h)
        {
            high = mid;
        } else
        {
            low = mid + 1;
        }
    }

    cout << low << brl;
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
