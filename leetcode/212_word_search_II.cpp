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

class Solution
{
public:
  vector<string> findWords(vector<vector<char>>& board, vector<string>& words)
  {
    auto available_chars = array<int, 26>{};
    available_chars.fill(0);
    for (auto i = 0; i < board.size(); ++i) {
      for (auto j = 0; j < board[0].size(); ++j) { ++available_chars[board[i][j] - 'a']; }
    }

    for (auto const& word : words) {
      auto needed_chars = array<int, 26>{};
      auto could_be_found = true;
      for (auto c : word) {
        ++needed_chars[c - 'a'];
        if (needed_chars[c - 'a'] > available_chars[c - 'a']) {
          could_be_found = false;
          break;
        }
      }
      if (could_be_found) { addWord(word); }
    }

    auto found = vector<string>{};
    auto visited = set<pair<int, int>>{};
    auto empty = string{};
    for (auto i = 0; i < board.size(); ++i) {
      for (auto j = 0; j < board[0].size(); ++j) {
        _find_words_from(board, make_pair(i, j), visited, found, empty, root);
      }
    }

    return found;
  }

  Solution() : root{ false, {} } {}

private:
  struct Entry
  {
    bool is_word;
    map<char, Entry> rest;
  };
  Entry root;

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

  auto _find_words_from(vector<vector<char>> const& board,
    pair<int, int> location,
    set<pair<int, int>>& visited,
    vector<string>& found,
    string& word,
    Entry& trie_node) -> void
  {
    if (trie_node.rest.count(board[location.first][location.second]) == 0) { return; }
    auto& next_node = trie_node.rest[board[location.first][location.second]];
    auto cur_word = word + board[location.first][location.second];
    if (next_node.is_word) {
      next_node.is_word = false;
      found.push_back(cur_word);
    }
    auto offsets = array<pair<int, int>, 4>{ make_pair(1, 0), make_pair(-1, 0), make_pair(0, 1), make_pair(0, -1) };
    visited.insert(location);
    for (auto const& offset : offsets) {
      auto next_location = make_pair(location.first + offset.first, location.second + offset.second);
      if (next_location.first < 0 || next_location.first >= board.size() || next_location.second < 0
          || next_location.second >= board[0].size()) {
        continue;
      }
      if (visited.count(next_location) == 1) { continue; }
      _find_words_from(board, next_location, visited, found, cur_word, next_node);
    }
    if (next_node.is_word == false && next_node.rest.empty()) {
      trie_node.rest.erase(board[location.first][location.second]);
    }
    visited.erase(location);
    return;
  }
};

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  return 0;
}
