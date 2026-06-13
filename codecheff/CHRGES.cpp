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
    i64 n;
    cin >> n;
    string data;
    cin >> data;

    auto pos_closest = vector<i64>(n, 1e5 + 5);
    auto neg_closest = vector<i64>(n, 1e5 + 5);
    i64 left_pos_dis = 1e5 + 5;
    i64 left_neg_dis = 1e5 + 5;
    for (i64 i = 0; i < n; ++i)
    {
        if (data[i] == '+')
        {
            left_pos_dis = 0;
        } else
        {
            ++left_pos_dis;
        }
        if (data[i] == '-')
        {
            left_neg_dis = 0;
        } else
        {
            ++left_neg_dis;
        }
        pos_closest[i] = left_pos_dis;
        neg_closest[i] = left_neg_dis;
    }

    i64 right_pos_dis = 1e5 + 5;
    i64 right_neg_dis = 1e5 + 5;
    for (i64 i = n - 1; i >= 0; --i)
    {
        if (data[i] == '+')
        {
            right_pos_dis = 0;
        } else
        {
            ++right_pos_dis;
        }
        if (data[i] == '-')
        {
            right_neg_dis = 0;
        } else
        {
            ++right_neg_dis;
        }
        pos_closest[i] = min(right_pos_dis, pos_closest[i]);
        neg_closest[i] = min(right_neg_dis, neg_closest[i]);
    }

    u64 count = 0;
    for (i64 i = 0; i < n; ++i)
    {
        if (pos_closest[i] == neg_closest[i]) { ++count; }
    }
    cout << count << brl;
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
