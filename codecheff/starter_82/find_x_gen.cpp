#include <bits/stdc++.h>

using namespace std;

auto main() -> int {
    ios::sync_with_stdio(0);
    cin.tie(0);
    
    auto cases = vector<vector<int>>{};

    auto limit = 20;

	for(auto a=1; a<limit; ++a) {
        for(auto a=1; a<limit; ++a) {
            for(auto b=1; b<limit; ++b) {
                for(auto c=1; c<limit; ++c) {
                    for(auto d=1; d<limit; ++d) {
                        if( a%b == c%d) {
                            cases.emplace_back(vector<int>{a,b,c,d});
                        }
                    }
                }
            }
        }
    }
    cout << cases.size() << '\n';
    for(auto i=0; i < cases.size(); ++i) {
        cout<<cases[i][0]<<" "<<cases[i][1]<<" "<<cases[i][2]<<" "<<cases[i][3]<<endl;

    }
	return 0;
}