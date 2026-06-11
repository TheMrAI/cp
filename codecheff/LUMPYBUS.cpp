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
    u64 n, p, q;
    cin >> n >> p >> q;

    auto nums = vector<u64>(n, 0);
    for (u64 i = 0; i < n; ++i)
    {
        i64 num;
        cin >> num;
        nums[i] = num;
    }

    sort(nums.begin(), nums.end());

    u64 changed = 0;
    u64 i = 0;
    while (i < nums.size() && nums[i] <= p + (2 * q))
    {
        u64 used_twos = min(q, nums[i] / 2);
        u64 used_ones = min(p, nums[i] - (2 * used_twos));
        if (used_ones + (2 * used_twos) == nums[i])
        {
            ++changed;
            p -= used_ones;
            q -= used_twos;
        }
        ++i;
    }
    cout << changed << brl;
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
