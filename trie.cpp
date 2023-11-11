#include <bits/stdc++.h>

using namespace std;

class Trie
{
private:
  struct Entry
  {
    bool is_word;
    map<char, Entry> rest;
  };

  Entry root;

public:
  Trie() : root{ false, {} } {}

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
    if (entry.rest.count(marker) == 1) { return _search(word, index + 1, entry.rest.at(marker)); }
    return false;
  }
};

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto trie = Trie{};
  trie.addWord("apple");
  cout << trie.search("apple") << endl;
  cout << trie.search("apple") << endl;

  return 0;
}
