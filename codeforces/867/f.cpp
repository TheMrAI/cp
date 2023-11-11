#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

template<typename Output, typename Iter>
requires std::derived_from<Output, std::ostream> && std::weakly_incrementable<Iter> && std::indirectly_readable<Iter>
auto dump_to(Output& output, Iter begin, Iter end)
{
  while (begin != end) {
    output << *begin;
    if (std::next(begin) != end) { output << ' '; }
    ++begin;
  }
}

// 2,147,483,648
// 1,000,000,000

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

auto dfs_distance(vector<vector<int>> const& edge_list, int node, int from, vector<i6>& distances, int& farthest_node)
  -> void
{
  for (auto const& neighbour : edge_list[node]) {
    if (neighbour == from) { continue; }
    distances[neighbour] = distances[node] + 1;
    if (distances[farthest_node] < distances[neighbour]) { farthest_node = neighbour; }
    dfs_distance(edge_list, neighbour, node, distances, farthest_node);
  }
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto t = i6{ 0 };
  cin >> t;
  while (t--) {
    auto n = 0;
    auto k = i6{ 0 };
    auto c = i6{ 0 };
    cin >> n >> k >> c;
    auto edge_list = vector<vector<int>>(n, vector<int>());
    for (auto i = 0; i < n - 1; ++i) {
      auto from = 0;
      auto to = 0;
      cin >> from >> to;
      edge_list[from - 1].push_back(to - 1);
      edge_list[to - 1].push_back(from - 1);
    }

    auto distances_from_root = vector<i6>(n, 0);
    auto farthest_node_a = 0;
    dfs_distance(edge_list, 0, 0, distances_from_root, farthest_node_a);
    auto farthest_node_b = 0;
    auto distances_from_node_a = vector<i6>(n, 0);
    dfs_distance(edge_list, farthest_node_a, farthest_node_a, distances_from_node_a, farthest_node_b);
    auto distances_from_node_b = vector<i6>(n, 0);
    auto tmp = 0;
    dfs_distance(edge_list, farthest_node_b, farthest_node_b, distances_from_node_b, tmp);
    auto max_profit = i6{ 0 };
    for (auto i = 0; i < n; ++i) {
      max_profit =
        max(max_profit, (max(distances_from_node_a[i], distances_from_node_b[i]) * k) - (distances_from_root[i] * c));
    }
    cout << max_profit << endl;
  }

  return 0;
}
