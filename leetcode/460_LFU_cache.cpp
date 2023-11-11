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

class LFUCache
{
private:
  struct CacheEntry
  {
    int key;
    int value;
    int frequency;
  };

  size_t _capacity;
  list<CacheEntry> _cache;
  // technically the problem would require an unordered map, but in
  // practice for only this many numbers a map will be faster
  using cache_iterator = decltype(_cache)::iterator;
  using mapping = map<int, cache_iterator>;
  mapping _entries;
  // for the frequencies we need a map anyways
  mapping _frequencies;

public:
  LFUCache(size_t capacity) : _capacity{ capacity } {}

  int get(int key)
  {
    auto entry = _entries.find(key);
    if (entry == _entries.end()) { return -1; }
    auto data = *(entry->second);
    update_caches(entry, data);
    return data.value;
  }

  void put(int key, int value)
  {
    auto entry = _entries.find(key);
    auto data = CacheEntry{ key, value, 0 };
    if (entry != _entries.end()) {
      data = *(entry->second);
      data.value = value;
    }
    update_caches(entry, data);
  }

private:
  void update_caches(mapping::iterator& entry, CacheEntry& data)
  {
    if (entry != _entries.end()) { update_frequency_cache(entry->second, data.frequency); }
    data.frequency += 1;
    if (entry != _entries.end()) {
      _cache.erase(entry->second);
    } else if (_cache.size() == _capacity) {// delete last entry if we would insert a new element and overflow
      _entries.erase(_cache.front().key);
      update_frequency_cache(_cache.begin(), _cache.front().frequency);
      _cache.pop_front();
    }
    // insert
    auto insert_before =
      _frequencies.upper_bound(data.frequency);// yes, insert before the last entry with higher frequency
    auto insertion_iter = insert_before != _frequencies.end() ? insert_before->second : _cache.end();
    auto updated_entry = _cache.insert(insertion_iter, data);
    _entries[data.key] = updated_entry;
    if (_frequencies.count(data.frequency) == 0) { _frequencies.emplace(data.frequency, updated_entry); }
  }

  void update_frequency_cache(cache_iterator cache_entry, int frequency)
  {
    auto frequency_entry = _frequencies.find(frequency);
    if (frequency_entry != _frequencies.end() && frequency_entry->second == cache_entry) {
      auto next_entry = next(cache_entry);
      if (next_entry != _cache.end()) {
        if (next_entry->frequency == frequency) {
          frequency_entry->second = next_entry;
          return;
        }
      }
      _frequencies.erase(frequency_entry);
    }
  }
};

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  return 0;
}
