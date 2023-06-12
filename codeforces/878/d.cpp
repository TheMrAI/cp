#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

template<typename Output, typename Iter>
requires 
    std::derived_from<Output, std::ostream> &&
    std::weakly_incrementable<Iter> && 
    std::indirectly_readable<Iter>
auto dump_to(Output& output, Iter begin, Iter end) {
    while(begin != end) {
        output << *begin;
        if(std::next(begin) != end) {
            output << ' ';
        }
        ++begin;
    }
}

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto n = 0;
        cin>>n;
        auto a = vector<int>(n, 0);
        for(auto i = 0; i < n; ++i) {
            cin>>a[i];
        }
        ranges::sort(a);

        auto l = -1;
        auto h = static_cast<int>(1e9);
        while(h - l > 1) {
            auto half_segment_width = l + (h - l) / 2;

            auto i = 0;
            while(i + 1 < n && a[i + 1] - a[0] <= 2*half_segment_width) {
                ++i;
            }

            auto j = n-1;
            while(j - 1 >= 0 && a.back() - a[j-1] <= 2*half_segment_width) {
                --j;
            }

            ++i;
            --j;
            if(i > j || a[j] - a[i] <= 2*half_segment_width) {
                h = half_segment_width;
            } else {
                l = half_segment_width;
            }
        }
        cout << h << endl;
    }

    return 0;
}
