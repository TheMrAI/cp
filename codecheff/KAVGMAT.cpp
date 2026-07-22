#include <bits/stdc++.h>

using i64 = int64_t;
using u64 = uint64_t;
using i32 = int32_t;
using u32 = uint32_t;
using u8 = uint8_t;

constexpr char brl = '\n';

using namespace std;

auto matrix_average(vector<vector<u64>> const& prefix_sum, i64 s_i, i64 s_j, i64 e_i, i64 e_j) -> u64
{
    u64 total = prefix_sum[e_i][e_j];
    if (s_i > 0) { total -= prefix_sum[s_i - 1][e_j]; }
    if (s_j > 0) { total -= prefix_sum[e_i][s_j - 1]; }
    if (s_i > 0 && s_j > 0) { total += prefix_sum[s_i - 1][s_j - 1]; }
    u64 count = (1 + e_i - s_i) * (1 + e_j - s_j);

    return total / count;
}

auto solve() -> void
{
    u64 n, m, k;
    cin >> n >> m >> k;

    auto matrix = vector<vector<u64>>(n, vector<u64>(m, 0));
    auto prefix_sum = vector<vector<u64>>(n, vector<u64>(m, 0));
    for (u64 i = 0; i < n; ++i)
    {
        u64 row_sum = 0;
        for (u64 j = 0; j < m; ++j)
        {
            cin >> matrix[i][j];
            row_sum += matrix[i][j];
            prefix_sum[i][j] = row_sum;
            if (i > 0) { prefix_sum[i][j] += prefix_sum[i - 1][j]; }
        }
    }

    u64 worthy_count = 0;
    for (u64 i = 0; i < n; ++i)
    {
        for (u64 j = 0; j < m; ++j)
        {
            u64 max_length = min(n - i, m - j);
            u64 low = 0;
            u64 high = max_length;
            while (low < high)
            {
                u64 mid = low + ((high - low) / 2);

                if (matrix_average(prefix_sum, i, j, i + mid, j + mid) < k)
                {
                    low = mid + 1;
                } else
                {
                    high = mid;
                }
            }
            worthy_count += max_length - low;
        }
    }
    cout << worthy_count << brl;
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
