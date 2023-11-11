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

class WordDictionary
{
private:
  struct Entry
  {
    bool is_word;
    map<char, Entry> rest;
  };

  Entry root;

public:
  WordDictionary() : root{ false, {} } {}

  void addWord(string word)
  {
    auto* entry = &root;
    for (auto i = 0; i < word.size(); ++i) {
      if (entry->rest.count(word[i]) == 0) { entry->rest.emplace(word[i], Entry{ false, {} }); }
      if (i == word.size() - 1) {
        entry->rest.at(word[i]).is_word = true;
        break;
      }
      entry = &(entry->rest.at(word[i]));
    }
  }

  bool search(string word) { return _search(word, 0, root); }

private:
  bool _search(string const& word, int index, Entry const& entry)
  {
    if (index == word.size()) { return entry.is_word; }

    auto marker = word[index];
    if (marker == '.') {
      for (auto const& options : entry.rest) {
        if (_search(word, index + 1, options.second)) { return true; }
      }
      return false;
    }
    if (entry.rest.count(marker) == 1) { return _search(word, index + 1, entry.rest.at(marker)); }
    return false;
  }
};

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  return 0;
}
