#include <bits/stdc++.h>
#include <cstdint>
#include <limits>

using namespace std;
#define bl '\n';

// Source: https://codeforces.com/blog/entry/78852
auto ceil2(int64_t a, int64_t b) -> int64_t { return (a + b - 1) / b; }

auto main() -> int
{
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);

    auto t = int64_t{ 0 };
    std::cin >> t;
    while (t--)
    {
        int64_t n;
        cin >> n;
        int64_t total = 0;
        int64_t level = numeric_limits<int64_t>::max();
        for (auto i = 1; i <= n; ++i)
        {
            int64_t num;
            cin >> num;
            total += num;
            level = min(level, total / i);
            cout << level << " ";
        }
        cout << bl;
    }

    return 0;
}
