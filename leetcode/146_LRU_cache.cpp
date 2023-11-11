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

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

class LRUCache
{
private:
  int capacity;
  list<pair<int, int>> data_list;
  unordered_map<int, list<pair<int, int>>::iterator> mapping;

public:
  LRUCache(int capacity_) : capacity{ capacity_ } {}

  int get(int key)
  {
    auto match = mapping.find(key);
    if (match == mapping.end()) { return -1; }
    auto value = match->second->second;
    data_list.erase(match->second);
    data_list.emplace_front(key, value);
    mapping[key] = data_list.begin();
    return value;
  }

  void put(int key, int value)
  {
    auto match = mapping.find(key);
    if (match != mapping.end()) { data_list.erase(match->second); }
    data_list.emplace_front(key, value);
    mapping[key] = data_list.begin();
    if (data_list.size() > capacity) {
      mapping.erase(data_list.back().first);
      data_list.pop_back();
    }
  }
};

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  return 0;
}
