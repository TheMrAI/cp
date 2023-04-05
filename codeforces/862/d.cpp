#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto dfs(int index, int from, int distance, vector<int>& distances, vector<vector<int>> const& edges) -> void {
    distances[index] = distance;
    for (auto const edge : edges[index]) {
        if (edge != from) {
            dfs(edge, index, distance+1, distances, edges);
        }
    }
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto n=0;
    cin>>n;
    
    auto edges = vector<vector<int>>(n);
    for (auto i=0; i < n - 1; ++i) {
        auto u = 0;
        auto v = 0;
        cin>>u>>v;
        edges[u-1].push_back(v-1);
        edges[v-1].push_back(u-1);
    }

    auto distances_1 = vector<int>(n, 0);
    dfs(0, -1, 0, distances_1, edges);
    auto farthest_node = distance(distances_1.begin(), max_element(distances_1.begin(), distances_1.end()));
    dfs(farthest_node, -1, 0, distances_1, edges); // now we will know the longest path length
    auto distances_2 = vector<int>(n, 0);
    farthest_node = distance(distances_1.begin(), max_element(distances_1.begin(), distances_1.end()));
    dfs(farthest_node, -1, 0, distances_2, edges); // now we will know the path lengths from the other nodes

    // we have found the diameter for the graph, we know all distances from the two furthest nodes
    // we merge these distances together by the maximums
    for (auto i = 0; i < n; ++i) {
        distances_2[i] = max(distances_2[i], distances_1[i]);
    }
    sort(distances_2.begin(), distances_2.end()); // sort the distances in order

    auto component_count = 1;
    auto checking = distances_2.begin();
    for (auto k = 1; k <= n; ++k) {
        while(component_count < n && checking != distances_2.end() && *checking < k) {
            ++component_count;
            ++checking;
        }
        cout << component_count << " ";
    }

    return 0;
}
