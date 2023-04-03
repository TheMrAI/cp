// DSU, Source: https://cp-algorithms.com/data_structures/disjoint_set_union.html#union-by-size-rank

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