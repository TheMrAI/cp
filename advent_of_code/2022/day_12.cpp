#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b)
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

auto get_node_height(char c) -> int
{
  auto node_height = c - 'a';
  if (c == 'S') {
    node_height = 0;
  } else if (c == 'E') {
    node_height = 'z' - 'a';
  }
  return node_height;
}

auto bfs_shortest_path(vector<vector<int>> const& edge_list, int start_id, int destination_id)
{
  auto next_nodes = queue<pair<int, int>>{};
  auto visited = set<int>{};

  next_nodes.emplace(start_id, 0);

  while (!next_nodes.empty()) {
    auto current_node = next_nodes.front();
    next_nodes.pop();
    if (visited.count(current_node.first) == 1) { continue; }
    for (auto neighbour_id : edge_list[current_node.first]) {
      auto distance = current_node.second + 1;
      if (neighbour_id == destination_id) { return distance; }
      next_nodes.emplace(neighbour_id, distance);
    }
    visited.insert(current_node.first);
  }
  return numeric_limits<int>::min();
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto geography = vector<vector<char>>{};
  auto line = string{};
  while (getline(cin, line)) { geography.emplace_back(line.begin(), line.end()); }

  auto edge_list = vector<vector<int>>(geography.size() * geography[0].size());
  auto start = 0;
  auto destination = 0;
  for (auto i = 0; i < geography.size(); ++i) {
    for (auto j = 0; j < geography[0].size(); ++j) {
      auto node_id = (i * geography[0].size()) + j;
      auto node_height = get_node_height(geography[i][j]);
      if (geography[i][j] == 'S') {
        start = node_id;
      } else if (geography[i][j] == 'E') {
        destination = node_id;
      }

      // right
      if (j < geography[0].size() - 1) {
        auto right_node_id = (i * geography[0].size()) + j + 1;
        auto right_node_height = get_node_height(geography[i][j + 1]);

        if (right_node_height - node_height <= 1) { edge_list[node_id].push_back(right_node_id); }
        if (node_height - right_node_height <= 1) { edge_list[right_node_id].push_back(node_id); }
      }
      // down
      if (i < geography.size() - 1) {
        auto down_node_id = ((i + 1) * geography[0].size()) + j;
        auto down_node_height = get_node_height(geography[i + 1][j]);

        if (down_node_height - node_height <= 1) { edge_list[node_id].push_back(down_node_id); }
        if (node_height - down_node_height <= 1) { edge_list[down_node_id].push_back(node_id); }
      }
    }
  }

  auto shortest_scenic_route_min_length = bfs_shortest_path(edge_list, start, destination);
  cout << "Shortest path has length: " << shortest_scenic_route_min_length << endl;
  for (auto i = 0; i < geography.size(); ++i) {
    for (auto j = 0; j < geography[0].size(); ++j) {
      if (geography[i][j] == 'a') {
        auto path_length = bfs_shortest_path(edge_list, i * geography[0].size() + j, destination);
        if (path_length == numeric_limits<int>::min()) { continue; }
        shortest_scenic_route_min_length = min(shortest_scenic_route_min_length, path_length);
      }
    }
  }
  cout << "Shortest scenic path has length: " << shortest_scenic_route_min_length << endl;
  return 0;
}
