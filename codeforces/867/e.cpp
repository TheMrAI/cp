#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b)
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto t = i6{ 0 };
  cin >> t;
  while (t--) {
    auto n = i6{ 0 };
    cin >> n;
    auto data = string{};
    cin >> data;
    if (n % 2 != 0) {
      cout << -1 << endl;
      continue;
    }

    auto start = 0;
    auto end = data.size() - 1;
    auto invalid_chars = map<char, vector<int>>{};
    while (start < end) {
      if (data[start] == data[end]) { invalid_chars[data[start]].push_back(start); }
      ++start;
      --end;
    }
    auto swap_count = 0;
    auto by_length = [](auto const& lhs, auto const& rhs) { return lhs.size() < rhs.size(); };
    auto char_indexes = priority_queue<vector<int>, vector<vector<int>>, decltype(by_length)>{ by_length };
    for (auto& [_, v] : invalid_chars) { char_indexes.push(move(v)); }
    while (!char_indexes.empty() && char_indexes.size() >= 2) {
      auto most_common = char_indexes.top();
      char_indexes.pop();
      auto second_most_common = char_indexes.top();
      char_indexes.pop();
      if (!most_common.empty() && !second_most_common.empty()) {
        swap(data[most_common.back()], data[second_most_common.back()]);
        most_common.pop_back();
        second_most_common.pop_back();
        ++swap_count;
      }
      if (!most_common.empty()) { char_indexes.push(most_common); }
      if (!second_most_common.empty()) { char_indexes.push(second_most_common); }
    }

    start = 0;
    end = data.size() - 1;
    while (!char_indexes.empty() && start < end) {
      auto top = char_indexes.top();
      char_indexes.pop();
      while (!top.empty() && start < end) {
        if (data[start] != data[end]) {
          if (data[top.back()] != data[start] && data[top.back()] != data[end]) {
            swap(data[top.back()], data[start]);
            ++swap_count;
            top.pop_back();
          }
        }
        ++start;
        --end;
      }
      if (!top.empty()) { char_indexes.push(top); }
    }
    if (char_indexes.empty()) {
      cout << swap_count << endl;
    } else {
      cout << -1 << endl;
    }
  }

  return 0;
}
