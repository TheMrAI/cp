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

auto valid(vector<int> const& matrix) -> bool {
    if(matrix[0] < matrix[1] && matrix[2] < matrix[3] && matrix[0] < matrix[2] && matrix[1] < matrix[3]) {
        return true;
    }
    return false;
}

auto check(vector<int>& matrix) -> bool {
    for(auto i = 0; i < 4; ++i) {
        if(valid(matrix)) {
            return true;
        }
        auto tmp = matrix[0];
        matrix[0] = matrix[2];
        matrix[2] = matrix[3];
        matrix[3] = matrix[1];
        matrix[1] = tmp;
    }
    return false;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto matrix = vector<int>(4, 0);
        cin >> matrix[0];
        cin >> matrix[1];
        cin >> matrix[2];
        cin >> matrix[3];

        if(check(matrix)) {
            cout << "YES" << endl;
            continue;
        }
        cout << "NO" << endl;
    }

    return 0;
}
