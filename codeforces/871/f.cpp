#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto dfs_last(vector<vector<int>>& edge_list, int node, int prev) -> int {
    for(auto const& neighbour : edge_list[node]) {
        if(neighbour == prev) {
            continue;
        }
        return dfs_last(edge_list, neighbour, node);
    }
    return node;
}

auto dfs_through_center(vector<vector<int>>& edge_list, int node, vector<int>& path, int prev) -> bool {
    path.push_back(node);
    if(path.size() == 5) {
        return true;
    }
    for(auto const& neighbour : edge_list[node]) {
        if(neighbour == prev) {
            continue;
        }
        if(dfs_through_center(edge_list, neighbour, path, node)) {
            return true;
        }
    }
    path.pop_back();
    return false;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto nodes = 0;
        cin>>nodes;
        auto edges = 0;
        cin >> edges;
        auto edge_list = vector<vector<int>>(nodes, vector<int>{});
        for(auto i = 0; i < edges; ++i) {
            auto n1 = 0;
            auto n2 = 0;
            cin>>n1>>n2;
            edge_list[n1-1].push_back(n2-1);
            edge_list[n2-1].push_back(n1-1);
        }
        auto from = dfs_last(edge_list, 0, -1);
        auto path = vector<int>{};
        dfs_through_center(edge_list, from, path, -1);
        auto x = edge_list[path[2]].size();
        auto y = edge_list[path[1]].size() - 1;
        cout << x << " " << y << endl;
    }

    return 0;
}