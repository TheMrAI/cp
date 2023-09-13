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

auto toStandardEdgeMap(std::unordered_map<std::string, Valve> const inputGraph, 
    std::unordered_map<std::string, int> const encoding) -> std::unordered_map<int, Node> {
        auto edgeMap = std::unordered_map<int, Node>{};
        for (auto const& entry : inputGraph) {
            auto nodeId = encoding.at(entry.first);
            auto node = Node{};
            node.flowRate = entry.second.flowRate;
            for (auto const& neighbor : entry.second.tunnels) { 
                node.edges.insert({encoding.at(neighbor), 1});
            }
            std::swap(edgeMap[nodeId], node);
        }

        return edgeMap;
}

auto floydWarshall(std::unordered_map<int, Node> const edgeMap) -> std::vector<std::vector<int>> {
    constexpr auto INF = std::numeric_limits<int>::max();
    
    auto distanceMatrix = std::vector<std::vector<int>>(edgeMap.size(), std::vector<int>(edgeMap.size(), INF));
    
    for(auto const& sourceNode : edgeMap) {
        distanceMatrix[sourceNode.first][sourceNode.first] = 0;
        for(auto const& targetNode : sourceNode.second.edges) {
            distanceMatrix[sourceNode.first][targetNode.first] = targetNode.second;
        }
    }

    for(auto k = 0; k < edgeMap.size(); ++k) {
        for(auto i = 0; i < edgeMap.size(); ++i) {
            for(auto j = 0; j < edgeMap.size(); ++j) {
                if(distanceMatrix[i][k] == INF || distanceMatrix[k][j] == INF) {
                    continue;
                }
                distanceMatrix[i][j] = std::min(distanceMatrix[i][j], distanceMatrix[i][k] + distanceMatrix[k][j]);
            }
        }
    }

    return distanceMatrix;
}

auto compressEdgeMap(std::unordered_map<int, Node> const& edgeMap) -> std::unordered_map<int, Node> {
    // Floyd-Warshall
    auto shortestPaths = floydWarshall(edgeMap);
    
    auto compressedEdgeMap = std::unordered_map<int, Node>{};
    //add all positive flow rate nodes
    for(auto const& node : edgeMap) {
        if (node.second.flowRate > 0) {
            compressedEdgeMap[node.first] = node.second;
            // clear all connections, they will get reconnected later
            compressedEdgeMap[node.first].edges = std::unordered_map<int, int>{};
        }
    }
    //connect all positive flow nodes to all other
    for(auto & node : compressedEdgeMap) {
        for(auto const& targetNode : compressedEdgeMap) {
            if (node.first == targetNode.first) {
                continue;
            }
            node.second.edges[targetNode.first] = shortestPaths[node.first][targetNode.first];
        }
    }
    //add zero id node, because that is the start even if it has zero flow
    compressedEdgeMap[0] = edgeMap.at(0);
    compressedEdgeMap[0].edges = std::unordered_map<int, int>{};
    //connect to all other nodes
    for(auto const& targetNode : compressedEdgeMap) {
        if(targetNode.first == 0) {
            continue;
        }
        compressedEdgeMap[0].edges[targetNode.first] = shortestPaths[0][targetNode.first];
    }
    
    return compressedEdgeMap;
}

auto dfs(std::unordered_map<int, Node> const& edgeMap, int nodeId, int remainingTime, int perMinuteRelease, uint seen) -> int {
    auto thisNodeVisited = (seen & (1u << nodeId)) != 0;
    if (thisNodeVisited || remainingTime <= 1 || (seen == (1u << edgeMap.size()) - 1)) {
        return perMinuteRelease * remainingTime;
    }

    auto released = 0;
    // open current valve if it is worth opening
    if (edgeMap.at(nodeId).flowRate != 0) {
        --remainingTime;
        released = perMinuteRelease;
    }
    // mark node as seen
    seen = seen | (1u << nodeId);
    perMinuteRelease += edgeMap.at(nodeId).flowRate;
    
    auto maximumThatCanBeReleased = 0;
    for(auto const& targetNode : edgeMap.at(nodeId).edges) {
        auto travelTime = targetNode.second;
        auto releasedWhileTraversing = travelTime * perMinuteRelease;
        maximumThatCanBeReleased = std::max(maximumThatCanBeReleased, 
                                     releasedWhileTraversing + dfs(edgeMap, targetNode.first, remainingTime - travelTime, perMinuteRelease, seen));
    }
    return released + maximumThatCanBeReleased;
}

auto main() -> int {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    auto originalGraph = parseGraph();
    auto nodeEncoding = encodeGraphNodes(originalGraph);
    auto originalEdgeMap = toStandardEdgeMap(originalGraph, nodeEncoding);
    auto edgeMap = compressEdgeMap(originalEdgeMap);

    for(auto const& node : edgeMap) {
        std::cout << "Node: " << node.first << " with flow rate: " << node.second.flowRate << std::endl;
        for(auto const& n : node.second.edges) {
            std::cout << '\t' << n.first << " in: " << n.second << std::endl;
        }
    }

    auto totalReleased = dfs(edgeMap, 0, 30, 0, 0u);
    std::cout << "Total pressure released: " << totalReleased << std::endl;

    return 0;
}
