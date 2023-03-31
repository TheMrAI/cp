    // #pragma GCC optimize("O3")
    // #pragma GCC optimization ("unroll-loops")
    // #pragma GCC target("avx,avx2,fma")
    // #pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,tune=native")
    #include <bits/stdc++.h>
    // #include <ext/pb_ds/assoc_container.hpp>
    // #include <ext/pb_ds/tree_policy.hpp>
    // using namespace __gnu_pbds;
    using namespace std;
     
    #define pb push_back
    #define ff first
    #define ss second
     
    typedef long long ll;
    typedef long double ld;
    typedef pair<int, int> pii;
    typedef pair<ll, ll> pll;
    typedef pair<ld, ld> pld;
     
    const int INF = 1e9;
    const ll LLINF = 1e18;
    const int MOD = 1e9 + 7;
     
    // template<class K> using sset =  tree<K, null_type, less<K>, rb_tree_tag, tree_order_statistics_node_update>;
     
    inline ll ceil0(ll a, ll b) {
        return a / b + ((a ^ b) > 0 && a % b);
    }
     
    void setIO() {
        ios_base::sync_with_stdio(0); cin.tie(0);
    }
     
    int query(int l, int r){
        cout << "? ";
        cout << r - l + 1 << " ";
        for(int i = l; i <= r; i++) cout << i << " ";
        cout << endl;
        int ret;
        cin >> ret;
        return ret;
    }
     
    int main(){
        setIO();
        int t;
        cin >> t;
        for(int tt = 1; tt <= t; tt++){
            int n;
            cin >> n;
            int arr[n + 1];
            for(int i = 1; i <= n; i++) cin >> arr[i];
            for(int i = 1; i <= n; i++) arr[i] += arr[i - 1];
            int l = 1, r = n;
            while(l < r){
                int mid = (l + r)/2;
                int tot = query(l, mid);
                if(tot == arr[mid] - arr[l - 1]){
                    l = mid + 1;
                } else {
                    r = mid;
                }
            }
            cout << "! " << l << endl;
        }
    }