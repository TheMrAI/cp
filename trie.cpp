#include <bits/stdc++.h>

using namespace std;

class Trie {
private:
    struct TrieNode {
        bool is_word;
        map<char, TrieNode> next;
    };

    map<char, TrieNode> root;
public:
    Trie() {}
    
    void insert(string word) {
        if(word.size() == 0) {
            return;
        }
        map<char, TrieNode>* node = &root;

        for(auto i = 0; i < word.size(); ++i) {
            if(node->count(word[i]) == 0) {
                node->insert(make_pair(word[i], TrieNode{false, {}} ));
            }
            if (i == word.size() - 1) {
                (*node)[word[i]].is_word = true;
            }
            node = &((*node)[word[i]].next);
        }
    }
    
    bool search(string word) {
        return _search(word, true);
    }
    
    bool startsWith(string prefix) {
        return _search(prefix, false);
    }

private:
    bool _search(string const& word, bool has_to_be_word) {
        if(word.size() == 0) {
            return false;
        }
        map<char, TrieNode>* node = &root;
        bool was_word = false;

        for(auto i = 0; i < word.size(); ++i) {
            if(node->count(word[i]) == 0) {
                return false;
            }
            was_word = (*node)[word[i]].is_word;
            node = &((*node)[word[i]].next);
        }
        return has_to_be_word ? was_word : true;
    }
};

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto trie = Trie{};
    trie.insert("apple");
    cout << trie.search("apple") << endl;
    cout << trie.search("apple") << endl;

    return 0;
}
