#include <bits/stdc++.h>

using namespace std;

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto n = 0;
    cin>>n;
    auto odd_length_words = vector<pair<uint32_t, uint32_t>>{};
    auto even_length_words = vector<pair<uint32_t, uint32_t>>{};
    while(n--) {
        auto word = string{};
        cin>>word;

        auto char_present = uint32_t{0};
        auto char_oddness = uint32_t{0};
        for(auto c : word) {
            char_present |= 1 << (c-'a');
            char_oddness ^= 1 << (c-'a');
        }

        if(word.size() % 2 == 0) {
            even_length_words.emplace_back(char_present, char_oddness);
        } else {
            odd_length_words.emplace_back(char_present, char_oddness);
        }
    }

    sort(odd_length_words.begin(), odd_length_words.end(), [](auto const& lhs, auto const& rhs) {
        return lhs.second < rhs.second;
    });
    sort(even_length_words.begin(), even_length_words.end(), [](auto const& lhs, auto const& rhs) {
        return lhs.second < rhs.second;
    });

    auto nightmares = int64_t{0};
    for (auto i = 0; i < 26; ++i) {
        auto filtered_odd = vector<pair<bitset<26>, uint32_t>>{};
        copy_if(odd_length_words.begin(), odd_length_words.end(), back_inserter(filtered_odd),
        [i](auto const& entry) {
            return (entry.first & (1<<i)) == 0;
        });
        auto filtered_even = vector<pair<bitset<26>, uint32_t>>{};
        copy_if(even_length_words.begin(), even_length_words.end(), back_inserter(filtered_even),
        [i](auto const& entry) {
            return (entry.first & (1<<i)) == 0;
        });

        for(auto const& odd_word : filtered_odd) {
            auto target = odd_word.second ^ ((1 << 26) - 1) ^ (1 << i);
            auto match_start = lower_bound(filtered_even.cbegin(), filtered_even.cend(), target,
                [](auto const& entry, auto const& target) {
                    return entry.second < target;
                });
            auto match_end = upper_bound(filtered_even.cbegin(), filtered_even.cend(), target,
                [](auto const& target, auto const& entry) {
                    return target < entry.second;
                });
            nightmares += distance(match_start, match_end);
            // while(match != filtered_even.cend()) {
            //     if(match->second != target) {
            //         break;
            //     }
            //     // this is already satisfied as long as target matches
            //     // if((odd_word.second ^ even_word.second).count() != 25) {
            //     //     continue;
            //     // } 
            //     if((odd_word.first | match->first).count() == 25) {
            //         ++nightmares;
            //     }
            //     ++match;
            // }
        }
    }
    cout << nightmares << endl;

    return 0;
}
