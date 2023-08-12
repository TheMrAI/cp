#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

template<typename Output, typename Iter>
requires 
    std::derived_from<Output, std::ostream> &&
    std::weakly_incrementable<Iter> && 
    std::indirectly_readable<Iter>
auto dump_to(Output& output, Iter begin, Iter end) {
    while(begin != end) {
        output << *begin;
        if(std::next(begin) != end) {
            output << ' ';
        }
        ++begin;
    }
}

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
requires std::integral<T>
auto ceil2(T a, T b) -> T {
    if (a == 0) return 0;
    return ((a - 1)/b) + 1;
}

struct Valve {
    int flowRate;
    std::vector<std::string> tunnels;
};

auto parseGraph() -> std::unordered_map<std::string, Valve> {
    auto graph = std::unordered_map<std::string, Valve>{};
    
    auto line = string{};
    auto data_regex = regex(R"(Valve (\w{2}) has flow rate=(\d+); tunnels? leads? to valves? (.+))");
    auto targets_regex = regex(R"([A-Z]{2})");

    while(getline(cin, line)) { 
        auto match = std::smatch{};
        if (std::regex_match(line, match, data_regex)) {
            std::cout << "Match: " << match[1].str() << " flow rate: " << match[2].str() << " path to: " << match[3].str() << std::endl;

            auto valve = Valve{};
            
            valve.flowRate = std::atoi(match[2].str().c_str());
            auto targets = match[3].str();
            auto targets_begin = std::sregex_iterator{targets.cbegin(), targets.cend(), targets_regex};
            auto targets_end = std::sregex_iterator{};
            
            for (auto iter = targets_begin; iter != targets_end; ++iter) {
                valve.tunnels.push_back(iter->str());
            }
            graph[match[1].str()] = valve;
        } else {
            std::cerr << "Regex match failure" << std::endl;
        }
    }
    return graph;
}

auto encodeGraphNodes(std::unordered_map<std::string, Valve> const& inputGraph) -> std::unordered_map<std::string, int> {
    auto nodeIds = std::vector<std::string>(inputGraph.size(), "");
    auto i = 0;
    for(auto const& entry : inputGraph) {
        nodeIds[i] = entry.first;
        ++i;
    }
    ranges::sort(nodeIds);
    auto encoding = std::unordered_map<std::string, int>{};
    for(auto i = 0; i < nodeIds.size(); ++i) {
        encoding[nodeIds[i]] = i;
    }
    return encoding;
}

struct Node {
    int flowRate;
    std::unordered_map<int, int> edges; // nodeId => weight
};

auto toStandardEdgeList(std::unordered_map<std::string, Valve> const inputGraph, 
    std::unordered_map<std::string, int> const encoding) -> std::vector<Node> {
        auto edgeList = std::vector<Node>(encoding.size());
        for (auto const& entry : inputGraph) {
            auto nodeId = encoding.at(entry.first);
            auto node = Node{};
            node.flowRate = entry.second.flowRate;
            for (auto const& neighbor : entry.second.tunnels) { 
                node.edges.insert({encoding.at(neighbor), 1});
            }
            std::swap(edgeList[nodeId], node);
        }

        return edgeList;
}

auto floydWarshall(std::vector<Node> const edgeList) -> std::vector<std::vector<int>> {
    constexpr auto INF = std::numeric_limits<int>::max();
    
    auto distanceMatrix = std::vector<std::vector<int>>(edgeList.size(), std::vector<int>(edgeList.size(), INF));
    
    for(auto i = 0; i < edgeList.size(); ++i) {
        for(auto j = 0; j < edgeList.size(); ++j) {
            if (i == j) {
                distanceMatrix[i][j] = 0;
            } else if (edgeList[i].edges.count(j) != 0) {
                distanceMatrix[i][j] = edgeList[i].edges.at(j);
            }
        }
    }

    for(auto k = 0; k < edgeList.size(); ++k) {
        for(auto i = 0; i < edgeList.size(); ++i) {
            for(auto j = 0; j < edgeList.size(); ++j) {
                if(distanceMatrix[i][k] == INF || distanceMatrix[k][j] == INF) {
                    continue;
                }
                distanceMatrix[i][j] = std::min(distanceMatrix[i][j], distanceMatrix[i][k] + distanceMatrix[k][j]);
            }
        }
    }

    return distanceMatrix;
}

auto compressEdgeList(std::vector<Node> const& edgeList) -> std::vector<Node> {
    // Floyd-Warshall
    auto shortestPaths = floydWarshall(edgeList);
    // filter out non positive valves
    auto positiveFlowValves = std::set<int>{};
    for (auto i = 0; i < edgeList.size(); ++i) {
        if (edgeList[i].flowRate > 0) {
            positiveFlowValves.insert(i);
        }
    }
    // adding node with zero as that is the starting node, regardless of its valves value
    positiveFlowValves.insert(0);

    auto compressedEdgeList = std::vector<Node>{};
    for (auto i = 0; i < edgeList.size(); ++i) {
        if (positiveFlowValves.count(i) == 0) {
            continue;
        }
        auto node = Node{edgeList[i].flowRate, std::unordered_map<int, int>{}};
        for (auto j = 0; j < edgeList.size(); ++j) {
            if(i == j || positiveFlowValves.count(j) == 0) {
                continue;
            }
            node.edges.insert({j, shortestPaths[i][j]});
        }
        compressedEdgeList.push_back(std::move(node));
    }
    return compressedEdgeList;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto originalGraph = parseGraph();
    auto nodeEncoding = encodeGraphNodes(originalGraph);
    auto originalEdgeList = toStandardEdgeList(originalGraph, nodeEncoding);
    auto edgeList = compressEdgeList(originalEdgeList);

    for(auto i = 0; i < edgeList.size(); ++i) {
        std::cout << "Node: " << i << " flow: " << edgeList[i].flowRate <<  std::endl;
        for(auto const& edge : edgeList[i].edges) {
            std::cout << '\t' << "To: " << edge.first << " in: " << edge.second << std::endl;
        }
    }

    return 0;
}
