#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto dfs_largest(vector<vector<int>>& matrix, int i, int j, int& depth, int& largest_lake) -> void {
    if(i < 0 || j < 0 || i >= matrix.size() || j >= matrix[0].size()) {
        return;
    }
    if(matrix[i][j] == -1 || matrix[i][j] == 0) {
        return;
    }
    depth += matrix[i][j];
    largest_lake = max(largest_lake, depth);
    matrix[i][j] = -1;
    constexpr auto dirs = array<pair<int, int>, 4>{make_pair(1,0), make_pair(-1,0), make_pair(0,1), make_pair(0,-1)};
    for(auto const& dir : dirs) {
        dfs_largest(matrix, i + dir.first, j + dir.second, depth, largest_lake);
    }
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
        auto matrix = vector<vector<int>>(n, vector(m, 0));
        for(auto i = 0; i < n; ++i) {
            for(auto j = 0; j < m; ++j) {
                cin>>matrix[i][j];
            }
        }
        auto largest_lake = 0;
        for(auto i = 0; i < n; ++i) {
            for(auto j = 0; j < m; ++j) {
                if(matrix[i][j] == -1) {
                    continue;
                }
                auto depth = 0;
                dfs_largest(matrix, i, j, depth, largest_lake);
            }
        }
        cout << largest_lake << endl;
    }

    return 0;
}
