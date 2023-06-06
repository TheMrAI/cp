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
        auto lamps = vector<pair<int, int>>{};
        while(n--) {
            auto a = 0;
            auto b = 0;
            cin>>a>>b;
            lamps.emplace_back(a, b);
        }
        ranges::sort(lamps, [](auto const& lhs, auto const& rhs) {
            if(lhs.first == rhs.first) {
                return lhs.second > rhs.second;
            }
            return lhs.first < rhs.first;
        });
        auto on_lamps = priority_queue<int, std::vector<int>, std::greater<int>>{};
        auto checking = 0;
        auto total_score = i6{0};
        while(checking < lamps.size()) {
            total_score += lamps[checking].second;
            on_lamps.push(lamps[checking].first);
            auto lamps_breaking_until = on_lamps.size();
            while(!on_lamps.empty() && on_lamps.top() <= lamps_breaking_until) {
                on_lamps.pop();
            }
            if(lamps[checking].first <= lamps_breaking_until) {
                while(checking < lamps.size() && lamps[checking].first <= lamps_breaking_until) {
                    ++checking;
                }
            } else {
                ++checking;
            }
        }
        cout << total_score << endl;
    }

    return 0;
}
