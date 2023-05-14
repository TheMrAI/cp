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

auto dir_to_index(pair<int, int> const& direction) -> int {
    if(direction.first == 1) {
        if(direction.second == 1) {
            return 0;
        } else {
            return 1;
        }
    } else {
        if(direction.second == 1) {
            return 2;
        } else {
            return 3;
        }
    }
}

auto simulate(int start_i, int start_j, pair<int, int> direction, int n, int m, int target_i, int target_j) -> int {
    auto bounces = 0;
    auto i = start_i;
    auto j = start_j;
    auto seen = vector<vector<vector<bool>>>(4, vector<vector<bool>>(n, vector<bool>(m, false)));
    while(true) {
        if(seen[dir_to_index(direction)][i][j]) {
            break;
        }
        seen[dir_to_index(direction)][i][j] = true;
        if(i == target_i && j == target_j) {
            return bounces;
        }
        i += direction.first;
        j += direction.second;
        if(i == -1 || j == -1 || i == n || j == m) {
            i -= direction.first;
            j -= direction.second;
            if(i == 0 || i == n-1) {
                direction.first *= -1;
            }
            if(j == 0 || j == m-1) {
                direction.second *= -1;
            }
            // edge case
            if(i == 0) {
                if(j == 0) {
                    direction = make_pair(1,1);
                } else if(j == m - 1) {
                    direction = make_pair(1, -1);
                }
            } else if(i == n - 1) {
                if(j == 0) {
                    direction = make_pair(-1, 1);
                } else if(j == m - 1) {
                    direction = make_pair(-1, -1);
                }
            }
            ++bounces;
        }
    }
    return -1;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto n = 0;
        auto m = 0;
        cin>>n>>m;
        auto start_i = 0;
        auto start_j = 0;
        auto end_i = 0;
        auto end_j = 0;
        cin>>start_i>>start_j>>end_i>>end_j;
        start_i -= 1;
        start_j -= 1;
        end_i -= 1;
        end_j -= 1;
        auto start_dir = string{};
        cin>>start_dir;
        auto direction = make_pair(1,1); // DR
        if(start_dir == "DL") {
            direction = make_pair(1, -1);
        } else if(start_dir == "UR") {
            direction = make_pair(-1, 1);
        } else if(start_dir == "UL") {
            direction = make_pair(-1, -1);
        }

        cout << simulate(start_i, start_j, direction, n, m, end_i, end_j) << endl;
    }

    return 0;
}
