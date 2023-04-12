#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b) {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

auto make_set(vector<i6>& parent, i6 value) -> void {
    parent[value] = value;
}

auto find_set(vector<i6>& parent, i6 value) -> i6 {
    if(parent[value] == value) {
        return value;
    }
    auto _parent = find_set(parent, parent[value]);
    parent[value] = _parent;
    return _parent;
}

auto union_sets(vector<i6>& parents, vector<i6>& ranks, i6 lhs, i6 rhs) -> void {
    auto parent_lhs = find_set(parents, lhs);
    auto parent_rhs = find_set(parents, rhs);
    if(parent_lhs != parent_rhs) {
        if (ranks[parent_lhs] < ranks[parent_rhs]) {
            swap(parent_lhs, parent_rhs);
        }
        parents[parent_rhs] = parent_lhs;
        if (ranks[parent_lhs] == ranks[parent_rhs]) {
            ++ranks[parent_lhs];
        }
    }
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto t=i6{0};
    cin>>t;
    while(t--) {
        auto n=0;
        cin>>n;
        auto k=0;
        cin>>k;
        auto spell = string{};
        cin>>spell;
        auto to_spell = string{};
        cin>>to_spell;

        auto spell_chars = vector<int>(26, 0);
        for(auto c : spell) {
            ++spell_chars[c-'a'];
        }
        auto to_spell_chars = vector<int>(26, 0);
        for(auto c : to_spell) {
            ++to_spell_chars[c-'a'];
        }

        auto match = true;
        for(auto i=0; i < 26; ++i) {
            if(spell_chars[i] != to_spell_chars[i]) {
                match = false;
                break;
            }
        }

        if(!match) {
            cout << "NO" << endl;
            continue;
        }

        auto parents = vector<i6>(spell.size());
        auto ranks = vector<i6>(spell.size(), 0);
        for(auto i = 0; i < spell.size(); ++i) {
            make_set(parents, i);
        }
        for(auto i = 0; i < spell.size(); ++i) {
            if(0 <= i - k) {
                union_sets(parents, ranks, i, i-k);
            }
            if(0 <= i - k - 1) {
                union_sets(parents, ranks, i, i-k-1);
            }
            if(i + k <= spell.size() - 1) {
                union_sets(parents, ranks, i, i+k);
            }
            if(i + k + 1 <= spell.size() - 1) {
                union_sets(parents, ranks, i, i+k+1);
            }
        }

        auto available_chars_in_set = map<int, vector<int>>{};
        for(auto i = 0; i < spell.size(); ++i) {
            auto group = find_set(parents, i);
            auto& counter = available_chars_in_set[group];
            if(counter.size() == 0) {
                counter.assign(26, 0);
            }
            ++counter[spell[i]-'a'];
        }

        auto ok = true;
        for(auto i = 0; i < spell.size(); ++i) {
            auto group = find_set(parents, i);
            auto& counter = available_chars_in_set[group];
            auto available = counter[to_spell[i]-'a'];
            if (available == 0) {
                ok = false;
                break;
            }
            --counter[to_spell[i]-'a'];
        }

        if(ok) {
            cout << "YES" << endl;
        } else {
            cout << "NO" << endl;
        }
    }

    return 0;
}
